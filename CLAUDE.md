# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build / static checks (both currently pass cleanly)
go build ./...
go vet ./...

# Run tests
go test ./...

# Run a single test
go test ./package/... -run TestName

# Run tests with verbose output
go test -v ./...

# Tidy dependencies
go mod tidy

# Build CLI tools
go build -o bin/create ./cmd/create

# Bump patch/minor/major version tag and push
./version.sh --request   # patch: v1.0.X
./version.sh --minor     # minor: v1.X.0
./version.sh --major     # major: vX.0.0
```

Test coverage is very thin: the only test files are `queue/example_test.go` and `jquery/test/*_test.go` (a separate package under `jquery/test/`, so run it with `go test ./jquery/...`, not `./jquery`). `create/v1/model_test.go` checks that every Go template of the generator renders to valid, gofmt-clean code and that no template embeds a JWT (`go test ./create/...`). For most packages, `go test ./...` only proves they compile; `go build`/`go vet` are the real check.

`version.sh` is outward-facing: it creates an annotated tag and immediately runs `git push origin --tags`. Only run it when the user asks for a release. Before tagging, it `sed`-replaces the current tag string with the new one in `README.md` (e.g. the `go get github.com/celsiainternet/elvis@vX.Y.Z` line), but it **does not commit** that edit. The tag therefore points at a commit whose README still shows the old version, and the working tree is left with a modified `README.md` that has to be committed separately.

## Code style

### Comments

All doc comments for functions, methods, and types must use this block style:

```go
/**
* FunctionName: Brief description.
* @param paramName type
* @return type
**/
```

- Use `@param` for each parameter and `@return` for the return value(s).
- Inline comments inside function bodies stay as `//`.
- Never use single-line `//` doc comments above a function or type declaration.

## Architecture Overview

**elvis** is a Go library (`github.com/celsiainternet/elvis`) providing infrastructure primitives for building microservices. It is not an application—it is a shared library consumed by other services.

### Core Data Types (`et/`)

The `et` package is the foundation used throughout the library:

- `et.Json` — `map[string]interface{}` with rich accessor methods (`.Str()`, `.Int()`, `.Bool()`, `.Key()`, etc.)
- `et.Item` — single result with `Ok bool` and `Result et.Json`
- `et.Items` — paginated result set with `Ok`, `Count`, `Result []et.Json`
- `et.List` — list with pagination metadata (rows, all, count, page, start, end)
- `et.Any` — generic value wrapper with typed conversion methods
- `et.MapBool` — `map[string]bool` used for permission maps; implements `ToString()`

### Database Layer (`jdb/`)

Multi-driver database abstraction supporting **PostgreSQL**, **MySQL**, and **Oracle**:

- `jdb.DB` is the main connection struct wrapping `database/sql`
- `jdb.Load()` / `jdb.LoadTo(dbname)` — connect using env vars (`DB_DRIVER`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_APPLICATION_NAME`)
- `InitCore()` initializes three internal tables: series (auto-increment sequences), records (audit), and recycling (soft deletes)
- `USE_CORE=true` env var controls whether core tables are initialized on connect
- `jdb.NextSerie(db, tag)` / `jdb.NextCode(db, tag, prefix)` — generate sequential numbers and prefixed codes (e.g. `"USR000001"`)

### ORM / Query Builder (`linq/`)

LINQ-style query builder that sits on top of `jdb`:

- `linq.Schema` — maps to a PostgreSQL schema (calls `CREATE SCHEMA IF NOT EXISTS`)
- `linq.NewModel(schema, name, description, version)` — maps to a table; call `.DefineColum()`, `.DefinePrimaryKey()`, `.DefineIndex()` to configure, then `.Init()` to create the table in the DB
- `linq.Mutation(schema, name, description, version)` — like `NewModel` but marks the model write-only (no queries)
- `linq.Linq` — fluent query builder with `From()`, `Where()`, `And()`, `Or()`, `OrderBy()`, `GroupBy()`, `Returns()`
- Column types: `TpColumn` (real column), `TpAtrib` (JSONB sub-key), `TpReference` (foreign lookup), `TpCaption`, `TpDetail`, `TpFunction`, `TpClone`, `TpField`
- Default special fields: `_DATA` (JSONB source), `DATE_MAKE`, `DATE_UPDATE`, `INDEX` (series), `CODE`, `PROJECT_ID`, `_STATE`, `_IDT`
- Two query modes: `TpData` (returns JSONB-built object) vs standard row query
- Triggers: `linq.BeforeInsert`, `AfterInsert`, `BeforeUpdate`, `AfterUpdate`, `BeforeDelete`, `AfterDelete`

### Cache (`cache/`)

Redis client wrapper using `go-redis/v9`:

- `cache.Load()` — singleton connect using `REDIS_HOST`, `REDIS_PASSWORD`, `REDIS_DB`
- Supports pub/sub via `cache/pubsub.go`
- `cache.GenKey(parts...)` — builds cache keys

### In-Memory Cache (`mem/`)

Thread-safe in-memory store with TTL support, initialized automatically via `init()`. Used as a lightweight alternative to Redis.

### Event System (`event/`)

Dual-mode event system:

- **Local events** (in-process): `event.On(channel, handler)` / `event.Emit(channel, data)` via `EventEmiter`
- **Distributed events** (NATS): `event.Stack(channel, handler)` / `event.Publish(channel, data)` via NATS connection (`NATS_HOST`, `NATS_USER`, `NATS_PASSWORD`)
- `event.Stack` re-registers the handler automatically on reconnect; use it for reset/sync subscriptions

### Authentication & Authorization (`claim/`, `middleware/`)

- `claim.Claim` — JWT payload struct with user identity fields (ID, App, Name, Username, Device, ProjectId, ProfileId, Tag)
- `claim.NewToken()` / `claim.ValidToken()` — JWT generation and validation; tokens stored in Redis for invalidation
- `SECRET` env var is the JWT signing key (defaults to `"1977"`)
- Middleware stack in `middleware/`: `Autentication` (JWT validation; puts the claim, including `ProjectId`/`ProfileId`, in the context), `Authorization` (requires the `AppName` header, 428 without it, stored as a string; asks the `middleware.Store` set with `SetAuthorizationStore`, e.g. `jrpc.LoadAuthorization`), `Ephemeral` (short-lived tokens), `AllowAll(origins).Handler` (CORS via `rs/cors`; there is no `Cors`), `Logger`, `RequestID` (not `RequestId`), `Recoverer`. Telemetry is exposed as helpers (`NewMetric`, `PushTelemetry`), not as a chi middleware
- `claim.ClientId(r)`, `claim.GetClient(r)` — extract identity from authenticated requests

### HTTP Router (`router/`, `response/`)

Built on **go-chi/chi v5**:

- Route registration helpers: `router.PublicRoute()`, `router.ProtectRoute()` (requires auth), `router.AuthorizationRoute()` (auth + permissions), `router.With()` (custom middleware)
- All routes automatically publish themselves to the API Gateway via NATS (`apigateway/set/resolve`)
- `response` package provides HTTP helpers: `ITEM`, `ITEMS`, `JSON`, `HTTPError`, `HTTPAlert`, `Unauthorized`, `Forbidden`, `Stream` (streaming paginated JSON)
- `response.GetBody(r)` — parses request body as `et.Json`; `response.GetQuery(r)` — query params; `response.GetParam(r, key)` — chi URL params

### RPC (`jrpc/`)

TCP-based RPC for inter-service calls using Go's `net/rpc`; Redis is used only to store package/solver registrations:

- `jrpc.Load(name)` — initialize package with service name; registers host/port from `RPC_HOST`/`RPC_PORT`
- `jrpc.Mount(services)` — registers a struct's exported methods as RPC endpoints; method keys are `<package>.<Struct>.<Method>` (exactly 3 dot-separated parts)
- `jrpc.Call()`, `CallJson()`, `CallItem()`, `CallItems()`, `CallList()`, `CallPermitios()` — typed call helpers that dispatch to the right TCP host via Redis-stored solver registry
- `PIPE_HOST` env var (`host:port`) overrides solver lookup and routes all RPC calls through a single proxy host
- `jrpc.LoadAuthorization(method)` registers `AuthorStore` as the `middleware.Authorization` store: `Author` calls the RPC `method` (default `AUTHORIZATION_METHOD` env var, else `catalog.Services.Author`) with `project_id`, `profile_id`, `method`, `path` and reads `result.authorized`; answers are cached in `mem` for 10 s per replica. It lives in `jrpc` because `jrpc` already imports `middleware`

### Dependency Installer (`cmd/install/`)

`go run github.com/celsiainternet/elvis/cmd/install` is meant to be run **from a consuming project** (right after `go get github.com/celsiainternet/elvis`, per `README.md`) — it shells out to `go get <module>@<version>` for a hardcoded list of third-party packages (`cmd/install/main.go`'s `dependencies` slice), which adds them to _that project's_ `go.mod`, not elvis's own. It exists because the microservice code scaffolded by `create/v1` (`cmd/main.go`, `pkg/*/config.go`, etc.) imports several packages elvis itself doesn't depend on — e.g. `github.com/dimiro1/banner`, `github.com/mattn/go-colorable`, `github.com/mattn/go-isatty` (used only by the generated `modelDbApi`/`modelApi`/`modelhRpc` templates for the startup banner/RPC bootstrap) — plus everything elvis's own `go.mod` requires (chi, redis, nats, jwt, cron, cors, cobra, promptui, etc.).

The pinned versions are **manually kept in sync with `go.mod` and can drift**: e.g. as of this writing `cmd/install` pins `lib/pq@v1.10.9`, `golang.org/x/crypto@v0.37.0`, `spf13/cobra@v1.9.1`, while elvis's own `go.mod` has moved on to `v1.12.3`, `v0.38.0`, `v1.10.2` respectively. When touching a dependency version in `go.mod`, check whether the same module appears in `cmd/install/main.go`'s list and update it too.

`cmd/install` also calls `agentsguide.Install()` (see below) as its last step, installing the agents framework into the consuming project.

### Agents Framework Guide (`agentsguide/`)

`agentsguide/` holds the Spanish-language framework that tells coding agents how to write code on top of elvis (and jdb) so it stays readable, maintainable and scalable: `MARCO.md` (priority order, decision rules such as jdb vs `linq`, `console` vs `logs`, error convention, version upgrades), `ELVIS.md` ("for X use Y, not Z" catalog with verified signatures, env vars, scaffolding) and `REVISION.md` (review checklist). They are embedded with `//go:embed` in `guide.go`.

`agentsguide.Install()` writes them to `.claude/guides/elvis/` in the current directory and upserts a block delimited by `<!-- elvis:agents-guide:start <version> -->` / `<!-- elvis:agents-guide:end -->` in `CLAUDE.md` (as `@` imports, which Claude Code loads) and `AGENTS.md` (as a reading list, for Codex and other agents). Re-running it overwrites the guides and replaces the block, so it also migrates the old single-block installs; content outside the markers is never touched. The version comes from `runtime/debug.ReadBuildInfo` (`devel` in workspace mode). The generic `InstallGuide(name, modulePath, fs, names)` is exported so `jdb/agentsguide` installs `JDB.md` the same way under `.claude/guides/jdb/` with `jdb:` markers.

`Install()` is called from `cmd/install/main.go` (after the dependency loop; it also prints a hint to run `jdb/cmd/install` when the project's `go.mod` requires jdb) and from `create/v1`'s `MkProject`.

When editing these guides, re-verify every signature or constant against the source before writing it down: they are meant to be the trustworthy reference when other docs drift. Changing `InstallGuide`'s signature breaks `jdb/agentsguide`.

### Project Scaffolding (`create/`, `cmd/create/`)

`cmd/create` is a Cobra CLI (`go run github.com/celsiainternet/elvis/cmd/create go`) that interactively scaffolds a new microservice project that consumes `elvis`. `create/v1/promps.go`'s `PrompCreate()` drives a `promptui` menu with four options, each backed by a `create/v1/hMicroservice.go` entry point:

- **Project** → `MkProject` — full new project: `MkMicroservice` + `README.md` + `.env` + `.gitignore`
- **Microservice** → `MkMicroservice` — `cmd/<name>/` (Dockerfile + `main.go`), `deployments/<name>/` (`local.yml` for docker-compose, plus `oke-template.yml` and `oke-statefulset-template.yml` Kubernetes manifests — Service+Deployment and Service+StatefulSet respectively, with `$ROLE`/`$NS`/`$PORT`/etc. left as literal placeholders for the CI/CD pipeline to substitute, not by `file.MakeFile`'s own `$1`/`$2`/... positional substitution), `internal/service/<name>/` (+ `v1/api.go`), `pkg/<name>/` (controller/handler/router/event/msg/config), `scripts/<name>.http`, `test/`
- **Modelo** → `MkMolue` → `MakeModel` — adds a `linq`-backed model + CRUD handler (`h<Model>.go`) into an existing `pkg/<name>`
- **Modelo** and **Rpc** fail with a clear error if `pkg/<name>` does not exist (an orphan `rpc.go` would not compile: it uses `PackageName`/`HostName` from `router.go`)
- **Rpc** → `MkRpc(name)` → `MakeRpc(name)` — adds `rpc.go` (`Services` + `StartRpcServer` + `Version`) to an existing `pkg/<name>` if missing

Whether a `schema` argument is supplied determines the template family: non-empty `schema` generates a `linq.Model`-backed controller/handler with full CRUD + `linq` triggers wired to a DB schema (`modelDbController`/`modelDbHandler`/`modelDbRouter` in `create/v1/model.go`); an empty `schema` generates a bare stub controller/handler with no persistence (`modelController`/`modelHandler`/`modelRouter`).

Templates are Go string constants in `create/v1/model.go`, rendered by `file.MakeFile(folder, name, template, args...)`, which does positional `$1`/`$2`/... substitution (see `params()` in `file/file.go`) — **not** `text/template`. `MakeFile` is idempotent: it silently no-ops if the target file already exists, so re-running a generator never overwrites hand-edited output.

**Only `create/v1` exists here** (`create/v2` and `cmd/jdb` were removed). Projects that use the standalone `jdb` ORM are scaffolded with `go run github.com/celsiainternet/jdb/cmd/create go`, which lives in the `jdb` repo and generates the `internal/models` + `internal/services` layout. Project and Microservice must run where `go.mod` already exists: the module path comes from it (`utility.GoMod`), and the CLI stops with a hint to run `go mod init` otherwise.

**What the v1 templates guarantee** (verified by generating both variants and running them against Postgres/Redis/NATS):

- Every package gets `rpc.go`; `api.go` loads config → `cache.Load` → `event.Load` → `pkg.StartRpcServer()` (→ `jdb.Load` with schema), and `Close` also closes the DB. Each `h<Model>.go` with schema also exposes `Get<Model>ById` as RPC `<pkg>.Services.Get<Model>ById`.
- Models use the field names `linq` recognizes (`date_make`, `date_update`, `_state`, `project_id`, `index`, `_data`), so dates and serie are filled by `linq`'s own triggers; `IntegrityAtrib` is left at its default (false) so attributes that are not columns go into `_data` (and are searchable). `Insert<Model>` generates the id with `utility.GenKey` when it comes empty.
- `Controller.Init` returns an error and `Routes()` panics on it; `main.go` shuts down on SIGINT and SIGTERM with `http.Server.Shutdown`.
- Templates follow `agentsguide/ELVIS.md`: no deprecated `envar.Envar*`, errors with `fmt.Errorf(msg.X)`, `response.LIST` for `et.List`, `middleware.AllowAll`/`RequestID`.
- **Modelo** wires everything except a missing DB: `upsertModelInit` appends `Define<Model>(db)` to `initModels`, and `upsertRoutes` inserts the routes (`modelDbRoutes`/`modelRoutes` in `model.go`) into `Routes()` right before the `ctx := context.Background()` line that `router.go` is generated with. Both are idempotent. If that anchor line was edited away, `MkMolue` returns `false` and the CLI tells the user to copy the snippet at the bottom of `h<Model>.go`. A model with schema in a package created without schema (no `initModels(c.Db)` in `controller.go`) is refused before writing any file, with the three manual changes listed.
- `MkMicroservice` and `MkMolue` finish with `go mod tidy` (`GoModTidy` in `hExec.go`); a failure only prints a warning. Without a `replace`, tidy resolves the **published** elvis, so generated code must compile against the latest tag, not only against the local checkout.
- Imports of the project itself go in their own import group so `gofmt` output does not depend on the module name.

**Organization-specific strings are baked into the templates** (`create/v1/model.go`): the `local.yml`/`oke-*.yml` manifests hardcode `APP`/`COMPANY` = `Celsia Internet` and `URL_API_FAILURE` = `https://api.celsiainternet.com`; treat them as placeholders. `TOKEN_API` is read from `secret-$NS` and `scripts/<name>.http` leaves `@token=` empty; do not put real JWTs back in the templates (the test rejects them).

### Other Packages

| Package          | Purpose                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `envar/`         | Environment variable helpers (`GetStr`, `GetInt`, `GetBool`, `SetStr`, etc.); auto-loads `.env` via `godotenv`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `logs/`          | Local-only logging with levels (Log, Logf, Alert, Debug, Panic). Used inside elvis packages that `console` cannot be imported from (`console` imports `event`); application code uses `console`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| `strs/`          | String utilities (Uppcase, Lowcase, Format, Append, DaskSpace, etc.)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| `utility/`       | General utilities: UUID, OTP, validation, crypto, password hashing, ID generation                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| `config/`        | Application config loading                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| `health/`        | Health check endpoint helpers                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `resilience/`    | Retry/resilience pattern; `resilience.Add(id, tag, description, tags, team, level, fn, fnArgs...)` wraps any function with automatic retries; env vars `RESILIENCE_TOTAL_ATTEMPTS` (default 3) and `RESILIENCE_TIME_ATTEMPTS` (seconds, default 30)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `workflow/`      | Multi-step workflow orchestration (`Flow`, `Step`, `FnContext`) with rollback support, conditional expressions, and configurable consistency (`strong`/`eventual`)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `instances/`     | Persistent service/workflow instance registry backed by a `linq` model in the database                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `request/`       | HTTP client for outbound calls. Verbs `Get/Post/Put/Delete/Patch/Options` return `(*Body, Status)`; variants: `<Verb>WithTls(..., tlsConfig)`, `<Verb>WithTimeout(..., timeout, defaultValue)`, `<Verb>WithTlsTimeout(..., tlsConfig, timeout, defaultValue)`, plus generic `Http`/`HttpCtx`/`HttpWithTimeout`/`HttpCtxWithTimeout` (take `method`) and `NewTlsConfig(ca, cert, key)`. Everything funnels into `httpDo` (runs in a goroutine; `timeout=0` = wait forever). On timeout it returns `newBody(defaultValue)` + `Status{Code: 408, Message: "timeout"}`. On any other transport/parse error `*Body` is **nil** — check `Status.Ok` before calling `Body.ToJson()`/etc. (value receivers, nil panics). **Known bugs:** timeout only abandons the result (no `context.WithTimeout`, so the request lives until the 120-min client timeout), and the `multipart/form-data` branch writes to a nil `*bytes.Buffer` and panics inside the goroutine (uncatchable by callers/`Recoverer`) |
| `jtls/`          | Self-signed TLS certificate generation (`jtls.Create(certFile, keyFile, expire)`) used by services that need mTLS                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| `file/`          | File system helpers: `MakeFolder`, `MakeFile`, `ReadFile`, `RemoveFile`, `ExistPath`, `ExtencionFile`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `race/`          | Concurrency race helpers                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| `dt/`            | Step-based resilience counter (`dt.Resilience`) and Redis-backed object cache (`dt.Object`)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| `reg/`           | ID registry helpers                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `service/`       | HTTP service client                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `console/`       | Application logging: prints and also publishes each entry as a NATS `logs` event (disable with `CONSOLE_LOG_EVENTS=false`). This is what services, the generators and `jdb` use                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| `timezone/`      | Timezone parsing and conversion helpers                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| `stdrout/`       | Standard output / terminal rendering                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| `crontab/`       | Cron job scheduling wrapper (`robfig/cron/v3`)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `authorization/` | Permission model backed by a `linq` model; loaded via `authorization.Load(db, schema)`; emits `event:set:authorization` / `event:del:authorization` on changes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `inbox/`         | Per-user inbox/notification records backed by a `linq` model; `inbox.Load(db, schema)` then query via `GetInboxesById`, `GetInboxesByCode`, `GetInboxesByMy`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| `msg/`           | Centralized Spanish-language message/error string constants (`MSG_*`, `ERR_*`, `RECORD_*`) shared across packages                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| `queue/`         | Generic in-process batching queue (`queue.Queue[T]`); groups `Push`ed items and dispatches to a handler on max batch size or timeout, whichever comes first                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| `jquery/`        | Translates an `et.Json` query description (from/join/select/wheres/group_by/having/limit/order_by) into a SQL `SELECT`; dialect is pluggable via `jquery/dialect` (Register/Get factory), defaults to PostgreSQL — see the package doc comment in `jquery/jquery.go` for the full JSON query format                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `xls/`           | Excel read/write helpers on top of `xuri/excelize` (`xls.ReadXls`/`ReadXlsFile`/`ReadXlsMultipart`, `xls.NewXls(...).ToFile/.ToWriter/.ToHttp`)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| `agentsguide/` | Agents framework (`MARCO.md`, `ELVIS.md`, `REVISION.md`) and `Install()`/`InstallGuide()`, which install it into a consuming project — see "Agents Framework Guide" above |
| `create/v1` | CLI scaffolding for new microservice projects — see "Project Scaffolding" above |
| `cmd/create` | CLI entry point for `create/v1`; `cmd/crontab`, `cmd/flow`, `cmd/install`, `cmd/jql` are example/demo mains for the corresponding packages |

### Key Environment Variables

| Variable                                                  | Used By           | Default             |
| --------------------------------------------------------- | ----------------- | ------------------- |
| `DB_DRIVER`                                               | jdb               | —                   |
| `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` | jdb               | —                   |
| `DB_APPLICATION_NAME`                                     | jdb               | `elvis`             |
| `USE_CORE`                                                | jdb               | `true`              |
| `REDIS_HOST`, `REDIS_PASSWORD`, `REDIS_DB`                | cache             | —                   |
| `NATS_HOST`, `NATS_USER`, `NATS_PASSWORD`                 | event             | —                   |
| `SECRET`                                                  | claim             | `"1977"`            |
| `HOST`, `RPC_HOST`, `RPC_PORT`                            | jrpc              | `localhost`, `4200` |
| `PIPE_HOST`, `PIPE_PORT`                                  | jrpc              | —, `4200`           |
| `AUTHORIZATION_METHOD`                                    | router/middleware | —                   |
| `RESILIENCE_TOTAL_ATTEMPTS`                               | resilience        | `3`                 |
| `RESILIENCE_TIME_ATTEMPTS`                                | resilience        | `30` (seconds)      |
| `PRODUCTION`                                              | dt                | `true`              |
