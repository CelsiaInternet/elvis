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

- `claim.Claim` — JWT payload struct with user identity fields (ID, App, Name, Username, Device, ProjectId, ProfileTp, Tag)
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
| `create/v1`, `create/v2` | CLI scaffolding for new microservice projects                                                                                                                                                 |
| `cmd/create`, `cmd/jdb`  | CLI entry points (`cmd/crontab`, `cmd/flow`, `cmd/install`, `cmd/jql` are example/demo mains for the corresponding packages)                                                                  |

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
| `PIPE_HOST`                                               | jrpc              | —                   |
| `AUTHORIZATION_METHOD`                                    | router/middleware | —                   |
| `RESILIENCE_TOTAL_ATTEMPTS`                               | resilience        | `3`                 |
| `RESILIENCE_TIME_ATTEMPTS`                                | resilience        | `30` (seconds)      |
| `PRODUCTION`                                              | dt                | `true`              |
