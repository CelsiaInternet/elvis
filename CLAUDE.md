# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

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
go build -o bin/jdb ./cmd/jdb

# Bump patch/minor/major version tag and push
./version.sh --request   # patch: v1.0.X
./version.sh --minor     # minor: v1.X.0
./version.sh --major     # major: vX.0.0
```

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
- Middleware stack in `middleware/`: `Autentication` (JWT validation), `Authorization` (permission check via jRPC), `Ephemeral` (short-lived tokens), `Cors`, `Logger`, `RequestId`, `Recoverer`, `Telemetry`
- `claim.ClientId(r)`, `claim.GetClient(r)` — extract identity from authenticated requests

### HTTP Router (`router/`, `response/`)

Built on **go-chi/chi v5**:

- Route registration helpers: `router.PublicRoute()`, `router.ProtectRoute()` (requires auth), `router.EphemeralRoute()`, `router.AuthorizationRoute()` (auth + permissions), `router.With()` (custom middleware)
- All routes automatically publish themselves to the API Gateway via NATS (`apigateway/set/resolve`)
- `response` package provides HTTP helpers: `ITEM`, `ITEMS`, `JSON`, `HTTPError`, `HTTPAlert`, `Unauthorized`, `Forbidden`, `Stream` (streaming paginated JSON)
- `response.GetBody(r)` — parses request body as `et.Json`; `response.GetQuery(r)` — query params; `response.GetParam(r, key)` — chi URL params

### RPC (`jrpc/`)

TCP-based RPC for inter-service calls using Go's `net/rpc`; Redis is used only to store package/solver registrations:

- `jrpc.Load(name)` — initialize package with service name; registers host/port from `RPC_HOST`/`RPC_PORT`
- `jrpc.Mount(services)` — registers a struct's exported methods as RPC endpoints; method keys are `<package>.<Struct>.<Method>` (exactly 3 dot-separated parts)
- `jrpc.Call()`, `CallJson()`, `CallItem()`, `CallItems()`, `CallList()`, `CallPermitios()` — typed call helpers that dispatch to the right TCP host via Redis-stored solver registry
- `PIPE_HOST` env var (`host:port`) overrides solver lookup and routes all RPC calls through a single proxy host
- Used by the authorization middleware to call `AUTHORIZATION_METHOD` env var

### Dependency Installer (`cmd/install/`)

`go run github.com/celsiainternet/elvis/cmd/install` is meant to be run **from a consuming project** (right after `go get github.com/celsiainternet/elvis`, per `README.md`) — it shells out to `go get <module>@<version>` for a hardcoded list of third-party packages (`cmd/install/main.go`'s `dependencies` slice), which adds them to *that project's* `go.mod`, not elvis's own. It exists because the microservice code scaffolded by `create/v1` (`cmd/main.go`, `pkg/*/config.go`, etc.) imports several packages elvis itself doesn't depend on — e.g. `github.com/dimiro1/banner`, `github.com/mattn/go-colorable`, `github.com/mattn/go-isatty` (used only by the generated `modelDbApi`/`modelApi`/`modelhRpc` templates for the startup banner/RPC bootstrap) — plus everything elvis's own `go.mod` requires (chi, redis, nats, jwt, cron, cors, cobra, promptui, etc.).

The pinned versions are **manually kept in sync with `go.mod` and can drift**: e.g. as of this writing `cmd/install` pins `lib/pq@v1.10.9`, `golang.org/x/crypto@v0.37.0`, `spf13/cobra@v1.9.1`, while elvis's own `go.mod` has moved on to `v1.12.3`, `v0.38.0`, `v1.10.2` respectively. When touching a dependency version in `go.mod`, check whether the same module appears in `cmd/install/main.go`'s list and update it too.

`cmd/install` also calls `agentsguide.Install()` (see below) as its last step, dropping/appending the agents framework guide into the consuming project's `CLAUDE.md`.

### Agents Framework Guide (`agentsguide/`)

`agentsguide/ELVIS_AGENTS.md` is a Spanish-language operational guide for coding agents that **create or refactor** backend projects on top of elvis — it documents the actually-verified API surface (the `linq` command/query split, real `middleware`/`claim` symbol names, etc.), a refactor checklist, and the minimum env var table. `agentsguide/guide.go` embeds it via `//go:embed` and exposes `agentsguide.Install()`, which writes it into `./CLAUDE.md` in the current directory — creating the file if missing, or appending the guide (once, gated by the `<!-- elvis:agents-guide:start -->` marker) if a `CLAUDE.md` already exists — so Claude Code auto-loads it in that project without any extra step.

`Install()` is wired into two places so it fires wherever a consuming project touches elvis: `cmd/install/main.go` (after the dependency loop — covers projects that already exist) and `create/v1` + `create/v2`'s `MkProject` (the "Project" scaffolding option — covers newly generated projects). It is **not** called from `MkMicroservice`, `MkMolue`, or `MkRpc` directly, only from `MkProject`, which calls `MkMicroservice` itself.

When editing `agentsguide/ELVIS_AGENTS.md`, keep it in sync with reality the same way this file is kept in sync — re-verify any function signature or constant it claims against the actual source before writing it down; it is explicitly written to be the trustworthy fallback when other generated docs have drifted.

### Project Scaffolding (`create/`, `cmd/create/`)

`cmd/create` is a Cobra CLI (`go run github.com/celsiainternet/elvis/cmd/create go`) that interactively scaffolds a new microservice project that consumes `elvis`. `create/v1/promps.go`'s `PrompCreate()` drives a `promptui` menu with four options, each backed by a `create/v1/hMicroservice.go` entry point:

- **Project** → `MkProject` — full new project: `MkMicroservice` + `README.md` + `.env` + `.gitignore`
- **Microservice** → `MkMicroservice` — `cmd/<name>/` (Dockerfile + `main.go`), `deployments/<name>/` (`local.yml` for docker-compose, plus `oke-template.yml` and `oke-statefulset-template.yml` Kubernetes manifests — Service+Deployment and Service+StatefulSet respectively, with `$ROLE`/`$NS`/`$PORT`/etc. left as literal placeholders for the CI/CD pipeline to substitute, not by `file.MakeFile`'s own `$1`/`$2`/... positional substitution), `internal/service/<name>/` (+ `v1/api.go`), `pkg/<name>/` (controller/handler/router/event/msg/config), `scripts/<name>.http`, `test/`
- **Modelo** → `MkMolue` → `MakeModel` — adds a `linq`-backed model + CRUD handler (`h<Model>.go`) into an existing `pkg/<name>`
- **Rpc** → `MkRpc` → `MakeRpc` — adds an RPC `rpc.go` stub into an existing `pkg/<name>`

Whether a `schema` argument is supplied determines the template family: non-empty `schema` generates a `linq.Model`-backed controller/handler with full CRUD + `linq` triggers wired to a DB schema (`modelDbController`/`modelDbHandler`/`modelDbRouter` in `create/v1/model.go`); an empty `schema` generates a bare stub controller/handler with no persistence (`modelController`/`modelHandler`/`modelRouter`).

Templates are Go string constants in `create/v1/model.go`, rendered by `file.MakeFile(folder, name, template, args...)`, which does positional `$1`/`$2`/... substitution (see `params()` in `file/file.go`) — **not** `text/template`. `MakeFile` is idempotent: it silently no-ops if the target file already exists, so re-running a generator never overwrites hand-edited output.

**Two generator versions exist, wired to two different CLI entry points**: `create/v1` (reworked folder layout in `internal/service/<name>/`) is used by `cmd/create/main.go`; `create/v2` (`internal/models/<name>/` + `internal/services/<name>/` instead of `internal/service/<name>/`, model code split out of the handler file, `router-<model>.go` naming) is used by **`cmd/jdb/main.go`** — despite its name, `cmd/jdb` is not a database CLI, it runs the exact same `promptui` scaffolding menu as `cmd/create` but against `create/v2`. Check which `cmd/`/`create/` pair an instruction actually intends before editing generator templates.

### Other Packages

| Package                  | Purpose                                                                                                                                                                                       |
| ------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `envar/`                 | Environment variable helpers (`GetStr`, `GetInt`, `GetBool`, `SetStr`, etc.); auto-loads `.env` via `godotenv`                                                                                |
| `logs/`                  | Structured logging with levels (Log, Logf, Alert, Debug, Panic)                                                                                                                               |
| `strs/`                  | String utilities (Uppcase, Lowcase, Format, Append, DaskSpace, etc.)                                                                                                                          |
| `utility/`               | General utilities: UUID, OTP, validation, crypto, password hashing, ID generation                                                                                                             |
| `config/`                | Application config loading                                                                                                                                                                    |
| `health/`                | Health check endpoint helpers                                                                                                                                                                 |
| `resilience/`            | Retry/resilience pattern; `resilience.Add(id, tag, description, tags, team, level, fn, fnArgs...)` wraps any function with automatic retries; env vars `RESILIENCE_TOTAL_ATTEMPTS` (default 3) and `RESILIENCE_TIME_ATTEMPTS` (seconds, default 30) |
| `workflow/`              | Multi-step workflow orchestration (`Flow`, `Step`, `FnContext`) with rollback support, conditional expressions, and configurable consistency (`strong`/`eventual`)                            |
| `instances/`             | Persistent service/workflow instance registry backed by a `linq` model in the database                                                                                                        |
| `request/`               | HTTP client utilities for outbound calls (GET, POST, PUT, DELETE with TLS support)                                                                                                            |
| `jtls/`                  | Self-signed TLS certificate generation (`jtls.Create(certFile, keyFile, expire)`) used by services that need mTLS                                                                             |
| `file/`                  | File system helpers: `MakeFolder`, `MakeFile`, `ReadFile`, `RemoveFile`, `ExistPath`, `ExtencionFile`                                                                                        |
| `race/`                  | Concurrency race helpers                                                                                                                                                                      |
| `dt/`                    | Step-based resilience counter (`dt.Resilience`) and Redis-backed object cache (`dt.Object`)                                                                                                   |
| `reg/`                   | ID registry helpers                                                                                                                                                                           |
| `service/`               | HTTP service client                                                                                                                                                                           |
| `console/`               | Low-level internal logging (used by other elvis packages; prefer `logs/` in application code)                                                                                                 |
| `timezone/`              | Timezone parsing and conversion helpers                                                                                                                                                       |
| `stdrout/`               | Standard output / terminal rendering                                                                                                                                                          |
| `crontab/`               | Cron job scheduling wrapper (`robfig/cron/v3`)                                                                                                                                                |
| `authorization/`         | Permission model backed by a `linq` model; loaded via `authorization.Load(db, schema)`; emits `event:set:authorization` / `event:del:authorization` on changes                              |
| `inbox/`                 | Per-user inbox/notification records backed by a `linq` model; `inbox.Load(db, schema)` then query via `GetInboxesById`, `GetInboxesByCode`, `GetInboxesByMy`                                 |
| `msg/`                   | Centralized Spanish-language message/error string constants (`MSG_*`, `ERR_*`, `RECORD_*`) shared across packages                                                                             |
| `queue/`                 | Generic in-process batching queue (`queue.Queue[T]`); groups `Push`ed items and dispatches to a handler on max batch size or timeout, whichever comes first                                  |
| `jquery/`                | Translates an `et.Json` query description (from/join/select/wheres/group_by/having/limit/order_by) into a SQL `SELECT`; dialect is pluggable via `jquery/dialect` (Register/Get factory), defaults to PostgreSQL — see the package doc comment in `jquery/jquery.go` for the full JSON query format |
| `xls/`                   | Excel read/write helpers on top of `xuri/excelize` (`xls.ReadXls`/`ReadXlsFile`/`ReadXlsMultipart`, `xls.NewXls(...).ToFile/.ToWriter/.ToHttp`)                                              |
| `agentsguide/`           | Embeds `ELVIS_AGENTS.md` (agent-facing operational guide) and exposes `Install()`, which writes/appends it into a consuming project's `CLAUDE.md` — see "Agents Framework Guide" above       |
| `create/v1`, `create/v2` | CLI scaffolding for new microservice projects — see "Project Scaffolding" above for which `cmd/` entry point drives which version                                                            |
| `cmd/create`, `cmd/jdb`  | CLI entry points for `create/v1` and `create/v2` respectively (`cmd/jdb` is misleadingly named — it is not a database tool); `cmd/crontab`, `cmd/flow`, `cmd/install`, `cmd/jql` are example/demo mains for the corresponding packages |

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
