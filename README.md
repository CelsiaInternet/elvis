# 🎸 Elvis — Go Microservices Framework

[Español](README.es.md)

[![Go Version](https://img.shields.io/github/go-mod/go-version/celsiainternet/elvis?style=flat-square&logo=go)](https://golang.org)
[![Latest Release](https://img.shields.io/github/v/release/celsiainternet/elvis?style=flat-square&logo=github)](https://github.com/celsiainternet/elvis/releases)
[![Documentation](https://img.shields.io/badge/docs-godoc-blue?style=flat-square&logo=go)](https://pkg.go.dev/github.com/celsiainternet/elvis)

## 🚀 Installation

To start a project with Elvis, run these commands in a new folder:

```bash
go mod init github.com/<organizacion>/<proyecto>

go get github.com/celsiainternet/elvis@v1.1.331
go run github.com/celsiainternet/elvis/cmd/install
go run github.com/celsiainternet/elvis/cmd/create go
```

| Command                                           | What it does                                                                                                                                                                   |
| ------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `go mod init ...`                                 | Creates `go.mod`. The next two commands need it: they read the module path from it and stop if it does not exist.                                                              |
| `go get github.com/celsiainternet/elvis@v1.1.331` | Adds Elvis as a dependency and writes the `go.sum` entries the following `go run` commands need; without this step they fail with `missing go.sum entry`.                      |
| `go run .../cmd/install`                          | Prepares the project: runs `go mod tidy` when there are already Go packages and installs the agents framework (see the note below). It pins no versions of other dependencies. |
| `go run .../cmd/create go`                        | Opens an interactive menu to generate code: **Project** (full project), **Microservice**, **Modelo** (a `linq` model + CRUD in an existing package) or **Rpc**.                |

Requirements: **Go 1.23** or later and, to run the generated service, **PostgreSQL** (database), **Redis** (cache, JWT tokens and RPC registry) and **NATS** (events).

> Both `cmd/install` and the **Project** option of `cmd/create` install a framework for coding agents (Claude Code, Codex and compatible tools): they copy the guides to `.claude/guides/elvis/` and add a block to `CLAUDE.md` and `AGENTS.md` that references them. Running it again updates the guides to the installed Elvis version without touching the rest of those files. If the project uses jdb, `go run github.com/celsiainternet/jdb/cmd/install` adds its guide. See [`agentsguide/MARCO.md`](agentsguide/MARCO.md).

## 📖 What it is

**Elvis** (`github.com/celsiainternet/elvis`) is a Go library with the infrastructure primitives of the Celsia Internet microservices: database, ORM, cache, events, authentication, HTTP router and RPC between services. It is not an application: other services consume it, and its CLI generates new projects that already use it.

## ✨ Features

| Area                                            | What you get                                                                                                                                                       |
| ----------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Types (`et`)                                    | `et.Json`, `et.Item`, `et.Items`, `et.List` and `et.Any`, with typed accessors (`Str`, `Int`, `Bool`, …).                                                          |
| Database (`jdb`)                                | Connection to PostgreSQL, MySQL or Oracle from environment variables, direct queries and core tables (series, audit, recycle bin).                                 |
| ORM (`linq`)                                    | Models with columns, JSONB attributes, references and triggers; fluent CRUD and `linq.RunQuery` for the jsql contract of `POST /query`.                            |
| SQL from JSON (`jquery`)                        | Translates an `et.Json` (from/join/select/wheres/…) into a `SELECT`, with a pluggable dialect.                                                                     |
| Cache (`cache`, `mem`)                          | Redis with pub/sub and an in-memory cache with TTL.                                                                                                                |
| Events (`event`)                                | Local and NATS events (`Subscribe`, `Queue`, `Stack`), jobs with state (`Work`) and reports to the central error registry (`event.Error`).                         |
| Authentication (`claim`, `middleware`)          | JWT tokens stored in Redis (regular, with project/profile, and ephemeral) and middleware for authentication, authorization, CORS, logging, request id and recover. |
| Router (`router`, `response`)                   | Public, protected and permission-checked chi routes that publish themselves to the API Gateway; response helpers and paginated streaming.                          |
| RPC (`jrpc`)                                    | TCP RPC between services with a Redis registry, typed calls and an optional single proxy (`PIPE_HOST`).                                                            |
| Notifications (`service`)                       | SMS, WhatsApp, email and OTP through RPC to the send and otp services, falling back to `event.Work`.                                                               |
| HTTP client (`request`)                         | HTTP verbs with TLS, timeout, JSON, form and `multipart/form-data`.                                                                                                |
| Secrets (`secret`)                              | AES-256-GCM encryption of secrets a service stores.                                                                                                                |
| Resilience and flows (`resilience`, `workflow`) | Automatic retries and multi-step flows with rollback.                                                                                                              |
| Other                                           | Crontab, Excel (`xls`), batching queue (`queue`), TLS certificates (`jtls`), time zones, logging (`console`, `logs`).                                              |
| Generator (`cmd/create`)                        | Project, microservice, model or RPC, with docker-compose and Kubernetes manifests and the agents framework.                                                        |

## 🧭 Usage

### ⚡ Minimal service

**An HTTP service with one public and one protected route.** It connects to the database and Redis from environment variables and registers the routes, which are published to the API Gateway.

```go
package main

import (
    "net/http"

    "github.com/celsiainternet/elvis/cache"
    "github.com/celsiainternet/elvis/envar"
    "github.com/celsiainternet/elvis/et"
    "github.com/celsiainternet/elvis/jdb"
    "github.com/celsiainternet/elvis/logs"
    "github.com/celsiainternet/elvis/middleware"
    "github.com/celsiainternet/elvis/response"
    "github.com/celsiainternet/elvis/router"
    "github.com/go-chi/chi/v5"
)

func main() {
    // Conectar a la base de datos
    db, err := jdb.Load()
    if err != nil {
        logs.Alert(err)
        return
    }
    defer db.Close()

    // Conectar a Redis
    _, err = cache.Load()
    if err != nil {
        logs.Alert(err)
        return
    }
    defer cache.Close()

    // Crear router chi
    r := chi.NewRouter()
    r.Use(middleware.AllowAll(nil).Handler)
    r.Use(middleware.Logger)

    host        := envar.GetStr("localhost", "HOST")
    packagePath := "/api/v1"

    // Ruta pública
    router.PublicRoute(r, "GET", "/health", func(w http.ResponseWriter, req *http.Request) {
        response.JSON(w, req, http.StatusOK, et.Json{"status": "ok"})
    }, "mi-servicio", packagePath, host)

    // Ruta protegida (requiere JWT)
    router.ProtectRoute(r, "GET", "/me", func(w http.ResponseWriter, req *http.Request) {
        response.JSON(w, req, http.StatusOK, et.Json{"mensaje": "autenticado"})
    }, "mi-servicio", packagePath, host)

    addr := envar.GetStr(":3400", "ADDR")
    logs.Logf("HTTP", "Escuchando en %s", addr)
    http.ListenAndServe(addr, r)
}
```

### 🧱 Core data types (`et`)

The `et` package is the base of the whole library. It defines the common data types used across the API.

#### `et.Json`

`map[string]interface{}` with typed accessors:

```go
data := et.Json{
    "nombre": "Elvis",
    "edad":   42,
    "activo": true,
}

nombre := data.Str("nombre")   // "Elvis"
edad   := data.Int("edad")     // 42
activo := data.Bool("activo")  // true

// Escribir valores
data.Set("pais", "Colombia")

// Serializar a string JSON
str := data.ToString()
```

#### `et.Item`

Single result of a query:

```go
// Uso típico
item, err := modelo.QueryOne("SELECT * FROM tabla WHERE _id=$1", id)
if item.Ok {
    nombre := item.Result.Str("nombre")
}
```

#### `et.Items`

Paginated result of a query:

```go
// Acceder al primer elemento
first := items.First()

// Iterar
for _, item := range items.Result {
    fmt.Println(item.Str("nombre"))
}
```

#### `et.List`

Result with full pagination metadata: `Rows`, `All`, `Count`, `Page`, `Start`, `End`, `Result`.

### 🗄️ Database (`jdb`)

Multi-driver abstraction for **PostgreSQL**, **MySQL** and **Oracle**.

#### Connection

```go
import "github.com/celsiainternet/elvis/jdb"

// Carga usando variables de entorno
db, err := jdb.Load()

// Conexión a una base de datos específica
db, err := jdb.LoadTo("nombre_base_datos")

// Conexión manual con parámetros
db, err := jdb.ConnectTo(et.Json{
    "driver":           "postgres",
    "host":             "localhost",
    "port":             5432,
    "dbname":           "mi_db",
    "user":             "postgres",
    "password":         "secreto",
    "application_name": "mi-servicio",
})
```

#### Direct queries

```go
// Consulta múltiples filas
items, err := db.Query("SELECT * FROM usuarios WHERE activo = $1", true)

// Consulta una fila
item, err := db.QueryOne("SELECT * FROM usuarios WHERE _id = $1", id)

// Ejecutar DDL
err := db.Ddl("CREATE TABLE IF NOT EXISTS ejemplo (_id VARCHAR(80) PRIMARY KEY)")

// Ejecutar comando DML
items, err := db.Command("INSERT INTO tabla (_id, nombre) VALUES ($1, $2) RETURNING *", id, nombre)
```

#### Core tables

When `USE_CORE=true` (default), `jdb.Load()` initializes three internal tables:

- **series** — auto-increment sequences per tag
- **records** — change audit
- **recycling** — recycle bin (soft deletes)

```go
// Siguiente valor de serie numérica
siguiente := jdb.NextSerie(db, "mi.tabla")

// Siguiente código con prefijo
codigo := jdb.NextCode(db, "mi.tabla", "USR")
// Resultado: "USR000001"
```

### 🔍 ORM / Query builder (`linq`)

LINQ-style: define typed models and build queries fluently.

#### Define a schema and a model

```go
import (
    "github.com/celsiainternet/elvis/jdb"
    "github.com/celsiainternet/elvis/linq"
)

// Schema = esquema PostgreSQL
schema := linq.NewSchema(db, "public")

// Modelo = tabla dentro del schema
modelo := linq.NewModel(schema, "USUARIOS", "Tabla de usuarios", 1)

// Columnas reales
modelo.DefineColum(jdb.KEY,       "", "VARCHAR(80)",  "-1")    // llave primaria (_id)
modelo.DefineColum("NOMBRE",      "", "VARCHAR(250)", "")
modelo.DefineColum("EMAIL",       "", "VARCHAR(250)", "")
modelo.DefineColum("DATE_MAKE",   "", "TIMESTAMP",    "NOW()") // activa UseDateMake
modelo.DefineColum("DATE_UPDATE", "", "TIMESTAMP",    "NOW()") // activa UseDateUpdate
modelo.DefineColum("_STATE",      "", "VARCHAR(20)",  "0")     // activa UseState
modelo.DefineColum("_DATA",       "", "JSONB",        "{}")    // activa UseSource (modo JSONB)

// Atributos JSONB (sub-campos dentro de _DATA)
modelo.DefineAtrib("telefono", "Teléfono", "VARCHAR(20)",  "")
modelo.DefineAtrib("ciudad",   "Ciudad",   "VARCHAR(100)", "")

// Llave primaria
modelo.DefinePrimaryKey([]string{jdb.KEY})

// Índices
modelo.DefineIndex([]string{"EMAIL"})

// Campos requeridos (campo:mensaje_error)
modelo.DefineRequired([]string{"NOMBRE:El nombre es requerido", "EMAIL"})

// Crear tabla en base de datos (CREATE TABLE IF NOT EXISTS)
err := modelo.Init()
```

#### Triggers

```go
modelo.Trigger(linq.BeforeInsert, func(m *linq.Model, old, new *et.Json, data et.Json) error {
    new.Set("NOMBRE", strings.ToUpper(new.Str("NOMBRE")))
    return nil
})

modelo.Trigger(linq.AfterInsert, func(m *linq.Model, old, new *et.Json, data et.Json) error {
    event.Emit("usuario.creado", *new)
    return nil
})

// Constantes: linq.BeforeInsert, linq.AfterInsert,
//             linq.BeforeUpdate, linq.AfterUpdate,
//             linq.BeforeDelete, linq.AfterDelete
```

#### CRUD queries

```go
// INSERT
item, err := modelo.Insert(et.Json{
    jdb.KEY:  utility.UUID(),
    "NOMBRE": "Juan Pérez",
    "EMAIL":  "juan@ejemplo.com",
}).CommandOne()

// UPDATE
item, err := modelo.Update(et.Json{"NOMBRE": "Juan Pablo"}).
    Where(modelo.Col(jdb.KEY).Eq(id)).
    CommandOne()

// UPSERT
item, err := modelo.Upsert(et.Json{
    jdb.KEY:  id,
    "NOMBRE": "Juan",
}).CommandOne()

// DELETE
item, err := modelo.Delete().
    Where(modelo.Col(jdb.KEY).Eq(id)).
    CommandOne()

// SELECT múltiple
items, err := modelo.Select().
    Where(modelo.Col("_STATE").Eq("0")).
    OrderBy(modelo.Col("NOMBRE"), true).
    All()

// SELECT paginado (et.Items, sin conteo total)
items, err := modelo.Select().
    Where(modelo.Col("_STATE").Eq("0")).
    Page(1, 20)

// SELECT paginado con metadatos de paginación (et.List: rows, all, count, page...)
list, err := modelo.Select().
    Where(modelo.Col("_STATE").Eq("0")).
    List(1, 20)

// SELECT uno
item, err := modelo.Select().
    Where(modelo.Col(jdb.KEY).Eq(id)).
    First()
```

> Write commands (`Insert`, `Update`, `Upsert`, `Delete`) run with `.Command()` (`et.Items`) or `.CommandOne()` (`et.Item`, first result); queries (`Select`, `Data`) use `.All()`, `.Page(page, rows)`, `.List(page, rows)` or `.First()`, not `.One()`.

#### Available conditions

```go
col.Eq(val)       // =
col.Neg(val)      // !=
col.In(vals...)   // IN (...)
col.Like(val)     // ILIKE
col.More(val)     // >
col.Less(val)     // <
col.MoreEq(val)   // >=
col.LessEq(val)   // <=
col.Search(val)   // @@ (búsqueda full-text)
```

> These conditions pass the value through `et.Unquote`, which does not escape the single quote: do not pass them unvalidated client text. For queries that come from the client, use `linq.RunQuery`.

#### jsql query (`linq.RunQuery`)

`RunQuery(model, query, wheres...)` builds the query from an `et.Json` with the `jdb.RunQuery` contract and returns the page (`et.List`):

```go
query := et.Json{
    "wheres": []et.Json{
        {"owner_id": et.Json{"in": []string{"cc:1", "id_deal:2"}}},
        {"or": et.Json{"name": et.Json{"like": "%acta%"}}},
        {"concat": et.Json{"fields": []string{"name", "description"}, "like": "%juan%"}},
    },
    "search":   et.Json{"value": "ana gomez", "fields": []string{"name", "description"}},
    "selects":  []string{"_id", "name"},
    "order_by": et.Json{"desc": []string{"date_make"}},
    "page":     1,
    "rows":     50, // alias limit; 100 por defecto (QUERY_DEFAULT_ROWS), máximo 1000 (QUERY_MAX_ROWS)
}
projectWhere := modelo.Column("project_id").In("-1", projectId)
list, err := linq.RunQuery(modelo, query, projectWhere)
```

> Published since `v1.1.324`: a service built outside `go.work` needs that version or a later one in its `go.mod`.

**`query` format**

| Key        | Shape                                                         | What it does                                                                                              |
| ---------- | ------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `selects`  | `["_id", "name"]`                                             | Result columns. Empty or missing: all.                                                                    |
| `wheres`   | `[{ field: { operator: value } }, …]`                         | Conditions joined with AND. `{ "and": { field: {…} } }` or `{ "or": { field: {…} } }` use that connector. |
| `order_by` | `{ "asc": [...], "desc": [...] }`                             | The `asc` fields first, then the `desc` ones.                                                             |
| `search`   | `{ "value": "ana gomez", "fields": ["name", "description"] }` | Free-text search (the contract of catalog's `/query` routes), see below.                                  |
| `page`     | number                                                        | Requested page; 1 when missing or lower than 1.                                                           |
| `rows`     | number                                                        | Rows per page (alias `limit`). 100 by default (`QUERY_DEFAULT_ROWS`), at most 1000 (`QUERY_MAX_ROWS`).    |

- Operators: `eq` (`=`), `neg` (`!=`), `in` (list), `like` (`ILIKE`, case-insensitive; the client adds the `%`), `more` (`>`), `less` (`<`), `moreEq` (`>=`), `lessEq` (`<=`). Operator names are case-sensitive; `and`/`or` are not.
- A field can carry several operators (`{ "date_make": { "more": "2026-01-01", "less": "2026-12-31" } }`): the first one takes the connector and the others go with AND.
- An `and`/`or` group can hold several fields: each goes with the group's connector, in alphabetical order, and its extra operators with AND. `{ "or": { "a": { "eq": 1 }, "b": { "eq": 2 } } }` is `OR A = 1 OR B = 2`. The group must be an object with at least one field.
- `concat`: `{ "concat": { "fields": ["name", "description"], "like": "%juan%" } }` compares against the fields joined with a space. `fields` is required and every field must exist and be visible; besides `fields` it needs at least one operator.
- Fields can be columns or `_data` attributes (`DefineAtrib`).
- `wheres` (variadic) are the caller's fixed conditions, for example the project from the token: they go with AND, and the client's conditions are wrapped in parentheses when they contain `or`, so an `or` cannot bypass them.
- Strings are escaped (single quote doubled) without going through `et.Unquote`. Only strings, numbers and booleans are accepted; `in` takes any non-empty list (`[]any`, `[]string`, …) of those values.
- `search`: every word of `value` (split by spaces, at most 5: `SEARCH_MAX_WORDS`; the rest are ignored) must appear, case-insensitively, in the `fields` joined with a space: `CONCAT(A.NAME, ' ', A.DESCRIPTION)::TEXT ILIKE '%ana%' AND ... ILIKE '%gomez%'`. Each word is an AND condition placed **after** the fixed and the client's `wheres` (outside their parentheses), so the search cannot widen the scope. `%`, `_` and `\` in `value` are matched as text. Without `value`, or when `search` is not an object (for example jdb's plain-string search), it does not filter. `fields` takes visible columns and attributes (1 to 5: `SEARCH_MAX_FIELDS`), not JSON or binary. It uses no index (it is an `ILIKE` over a `CONCAT`).

**Response**: `et.List` (`rows`, `all` = total matching records, `count` = records in the page, `page`, `start`, `end`, `result`). Each item of `result` is the record object (`model.Data` query).

**Errors**

- `Llave no valida en la consulta (<key>)` (`MSG_QUERY_KEY_NOT_VALID`) when `query` has a root key other than `wheres`, `selects`, `order_by`, `search`, `page`, `rows` and `limit` (for example `project_id` or `folder`: the service that uses them removes them before calling `RunQuery`).
- `Columna no valida en la consulta (<field>)` (`MSG_QUERY_COLUMN_NOT_VALID`) when a field in `selects`, `order_by` or `wheres` (also inside an `and`/`or` group and in `concat.fields`) does not exist or is hidden (`DefineHidden`).
- `Operador no valido en la consulta (<field>.<operator>)` (`MSG_QUERY_OPERATOR_NOT_VALID`) when the operator is not supported (for example `isNull`, `notIn`, `between`, `search`).
- `Valor no valido en la consulta (<field>[.<operator>])` (`MSG_QUERY_VALUE_NOT_VALID`) when the value is not a string, number or boolean (`null`, object, list outside `in`), when `in` is not a list, is empty or has an invalid item, when a field has no object with at least one operator, when an `and`/`or` group is not an object or is empty, when `concat` has no `fields`, or when `wheres` is not a list of objects (`(wheres)`).
- In `search`: `MSG_SEARCH_LENGTH` when `value` is longer than 100 characters (`SEARCH_MAX_LENGTH`), `MSG_SEARCH_FIELDS` when `fields` is empty or has more than 5, and `MSG_SEARCH_FIELD_NOT_VALID` when a field does not exist, is hidden or is JSON/binary.
- The database error, when the query fails.

These errors come before the database is queried, so a malformed filter answers 400 instead of returning everything the fixed `wheres` allow. The caller's fixed (variadic) `wheres` are not validated. Without `value`, or when `search` is not an object, the search neither filters nor fails.

**Not supported** (they are in `jdb.RunQuery`): `isNull`, `notNull`, `notIn`, `between`, the `search` operator inside `wheres`, the root `search` key as plain text (jdb's `DefineSearch`), and JSON paths (`_data.field`). For "no value" in a text field use `{ "eq": "" }`.

**`POST /query` handler**

```go
func (s *Router) query(w http.ResponseWriter, r *http.Request) {
    body, err := response.GetBody(r)
    if err != nil {
        response.HTTPError(w, r, http.StatusBadRequest, err.Error())
        return
    }

    // El proyecto sale del token, nunca del body.
    projectId := claim.ProjectId(r)
    projectWhere := Documents.Column("project_id").Eq(projectId)
    list, err := linq.RunQuery(Documents, body, projectWhere)
    if err != nil {
        response.HTTPError(w, r, http.StatusBadRequest, err.Error())
        return
    }

    response.LIST(w, r, http.StatusOK, list) // {"ok": true, "rows", "all", "count", "page", "start", "end", "result"}
}
```

#### References between models

```go
// Llave foránea
modeloOrden.DefineForeignKey("USUARIO_ID", modeloUsuario.Col(jdb.KEY))

// Referencia (columna virtual que trae el nombre del otro modelo)
modeloOrden.DefineReference("USUARIO_ID", "USUARIO", jdb.KEY, modeloUsuario.Col("NOMBRE"), false)
```

### 💾 Cache (`cache` / `mem`)

#### Redis (`cache`)

```go
import (
    "time"
    "github.com/celsiainternet/elvis/cache"
)

// Conectar (usa REDIS_HOST, REDIS_PASSWORD, REDIS_DB)
_, err := cache.Load()
defer cache.Close()

// Operaciones básicas
err = cache.Set("clave", "valor", 30*time.Minute)
valor, err := cache.Get("clave", "default")
_, err = cache.Delete("clave")

// Verificar conexión
ok := cache.HealthCheck()
```

#### Memory (`mem`)

In-process cache with TTL, initialized automatically:

```go
import (
    "time"
    "github.com/celsiainternet/elvis/mem"
)

mem.Set("clave", "valor", 5*time.Minute)
valor, err := mem.Get("clave", "default")
mem.Del("clave")
mem.Clear("prefijo") // elimina todas las claves que contienen "prefijo"
```

### 🔄 Events (`event`)

The event system has two modes: **local** (in-process) and **distributed** (NATS).

#### Local events (in-process)

```go
import "github.com/celsiainternet/elvis/event"

// Registrar handler
event.On("usuario.creado", func(msg event.EvenMessage) {
    fmt.Println("Usuario:", msg.Data.Str("nombre"))
})

// Emitir evento
event.Emit("usuario.creado", et.Json{"nombre": "Juan"})
```

#### Distributed events (NATS)

```go
// Conectar a NATS (usa NATS_HOST, NATS_USER, NATS_PASSWORD)
conn, err := event.Load()
defer event.Close()

// Suscribirse (un consumidor por mensaje)
err = event.Subscribe("pedido.nuevo", func(msg event.EvenMessage) {
    fmt.Println("Pedido:", msg.Data)
})

// Queue: reparto de carga entre múltiples instancias del mismo servicio
err = event.Queue("pedido.nuevo", "mi-grupo", func(msg event.EvenMessage) {
    fmt.Println("Pedido:", msg.Data)
})

// Stack: Queue con grupo fijo ("stack"); se re-registra automáticamente tras reconexión
event.Stack("canal/reset", func(msg event.EvenMessage) {
    // re-sincronización
})

// Publicar
event.Publish("pedido.nuevo", et.Json{"pedido_id": "123"})
```

> The `Subscribe`, `Queue` and `Stack` handlers are protected with `recover()`. A panic inside a handler is logged as an error instead of killing the NATS goroutine.

#### Report errors to the central registry (`event.Error`)

Publishes an error on `EVENT_ERROR` (`"event:error"`). Catalog stores it in its central registry: it groups identical errors and the front end opens the ticket in ADA or GLPI. `Error` validates the report and returns that error; when it is valid it publishes in another goroutine, so it never blocks or fails the caller. Without a NATS connection it does nothing.

```go
err := event.Error(event.ErrorReport{
    App:         "app.vcore",                // obligatorio
    Service:     "vcore-api",
    Level:       event.ERROR_LEVEL_ERROR,    // error (por defecto), warning o critical
    Code:        "CRM-LEAD-SAVE",            // con code, los mismos errores se agrupan aunque cambie el mensaje
    Message:     "No se pudo guardar el lead", // obligatorio
    Detail:      et.Json{"request_id": "r-1"},
    OwnerKind:   "lead",                     // obligatorio: entidad afectada
    OwnerId:     "L-1234",                   // obligatorio
    UserId:      userId,                     // opcional: usuario que tuvo el error
    TicketKind:  "incidente",                // opcional
    TicketId:    "45871",                    // opcional; si viene, TicketOwner es obligatorio
    TicketOwner: event.TICKET_OWNER_GLPI,    // ADA o GLPI
})
```

`OccurredAt` is the current time when it is not sent. The channel's consumer reads the message with `event.ParseErrorReport(msg.Data)`, which normalizes and validates the same way.

#### Work and WorkState

Helpers to trace asynchronous jobs:

```go
// Emite EVENT_WORK ("event:worker") y el evento propio del canal
work := event.Work("factura.generar", et.Json{"factura_id": "F-001"})
workId := work.Str("_id")

// Actualizar estado de la tarea (emite EVENT_WORK_STATE)
event.WorkState(workId, event.WorkStatusProcessing, et.Json{})
event.WorkState(workId, event.WorkStatusCompleted, et.Json{"url": "..."})
// Estados: WorkStatusPending, WorkStatusAccepted, WorkStatusProcessing, WorkStatusCompleted, WorkStatusFailed
```

### 📨 Notifications and OTP (`service`)

`service.SendSms`, `SendWhatsapp`, `SendEmail`, `SendEmailByTemplate` and `SendEmailToByTemplate` call the send service through RPC (`jrpc`); `SendOtpSms`, `SendOtpWhatsapp` and `SendOtpEmail` call the otp service. The call is synchronous and uses `RPC_TIMEOUT`, so the caller gets the provider's result or error.

| Function                | RPC method                         | Fallback event           |
| ----------------------- | ---------------------------------- | ------------------------ |
| `SendSms`               | `<SEND_RPC>.SendSMS`               | `send/sms`               |
| `SendWhatsapp`          | `<SEND_RPC>.SendWhatsapp`          | `send/whatsapp`          |
| `SendEmail`             | `<SEND_RPC>.SendEmail`             | `send/email`             |
| `SendEmailByTemplate`   | `<SEND_RPC>.SendEmailByTemplate`   | `send/emailbytemplate`   |
| `SendEmailToByTemplate` | `<SEND_RPC>.SendEmailToByTemplate` | `send/emailtobytemplate` |
| `SendOtpSms`            | `<OTP_RPC>.GenerateOTPSMS`         | `generate/otp/sms`       |
| `SendOtpWhatsapp`       | `<OTP_RPC>.GenerateOTPWhatsapp`    | `generate/otp/whatsapp`  |
| `SendOtpEmail`          | `<OTP_RPC>.GenerateOTPEmail`       | `generate/otp/email`     |
| `RegisterError`         | `<CATALOG_RPC>.RegisterError`      | `support/error`          |

`SEND_RPC` defaults to `Service-send.Services`, `OTP_RPC` to `Service-otp.Services` and `CATALOG_RPC` to `catalog.Services`. `service.RegisterError(event.ErrorReport{...})` registers a support case in catalog's central error registry (Sistema › Soporte): it validates the report (`app`, `message`, `owner_kind`, `owner_id`), calls catalog by RPC and returns `et.Item` with the registered row, or, when the method does not resolve, the published `support/error` work, which catalog stores later. Unlike `event.Error` (fire and forget on `event:error`), the caller gets the result. The signature is unchanged: they return `et.Json` with the RPC `result` plus `ok`, `service_id` and, when the call failed, `error` with the text.

```go
result := service.SendSms(projectId, "", []string{"3001234567"}, "Hola {{name}}", []et.Json{{"name": "Ana"}}, service.TpTransactional, user)
if !result.Bool("ok") {
	console.Alert(result.Str("error"))
}
serviceId := result.Str("service_id") // id del envío en send
```

Fallback: when `jrpc.GetSolver` cannot resolve the method (send or otp not registered in Redis, Redis not loaded or an invalid name), the call is not attempted and `event.Work` is published as before; the return value is then the job envelope (`_id`, `event`, `data`, `created_at`) plus `service_id`. When the call fails after resolving, the error is returned and nothing is resent through NATS, to avoid duplicate messages. With `PIPE_HOST` the solver always resolves, so there is no fallback.

### 🔐 Authentication and authorization (`claim` / `middleware`)

#### Generate JWT tokens

```go
import (
    "time"
    "github.com/celsiainternet/elvis/claim"
)

// Token básico (almacenado en Redis para soporte de logout)
token, err := claim.NewToken(userId, "mi-app", "Juan", "juan@e.com", "web", 24*time.Hour)

// Token con autorización (incluye projectId y profileId)
token, err := claim.NewAuthorization(userId, "mi-app", "Juan", "juan@e.com", "web", projectId, profileId, 8*time.Hour)

// Token efímero (vida corta, requiere tag descriptivo)
token, err := claim.NewEphemeralToken(userId, "mi-app", "Juan", "juan@e.com", "web", "descarga-reporte", 15*time.Minute)
```

The `SECRET` environment variable (default `"1977"`) is the JWT signing key.

#### Validate and parse tokens

```go
// Parsear sin validar en Redis
c, err := claim.ParceToken(token)

// Validar (parsea + verifica que el token esté activo en Redis)
c, err := claim.ValidToken(token)

// Logout
err = claim.DeleteToken(app, device, userId)
err = claim.DeleteTokeByToken(token)
```

#### Read the client's data from the request

```go
import "github.com/celsiainternet/elvis/claim"

func handler(w http.ResponseWriter, r *http.Request) {
    clientId  := claim.ClientId(r)
    nombre    := claim.ClientName(r)
    username  := claim.Username(r)
    projectId := claim.ProjectId(r)
    profileId := claim.ProfileId(r)
    device    := claim.Device(r)
    tag       := claim.Tag(r)

    // O todo en un et.Json:
    cliente := claim.GetClient(r)
}
```

#### HTTP middleware

```go
import (
    "github.com/celsiainternet/elvis/middleware"
    "github.com/go-chi/chi/v5"
)

r := chi.NewRouter()
r.Use(middleware.AllowAll(nil).Handler) // CORS (origins vacío = permitir todos)
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(middleware.RequestID)

// Ruta con autenticación
r.With(middleware.Autentication).Get("/perfil", handler)

// Ruta con autenticación + verificación de permisos
r.With(middleware.Autentication).With(middleware.Authorization).Get("/admin", handler)

// Ruta con token efímero
r.With(middleware.Ephemeral).Get("/descarga", handler)
```

### 🌐 HTTP router and responses (`router` / `response`)

#### Route registration

Every registered route is published to the API Gateway through NATS.

```go
import (
    "github.com/celsiainternet/elvis/router"
    "github.com/go-chi/chi/v5"
)

r := chi.NewRouter()
packageName := "usuarios-service"
packagePath := "/api/v1/usuarios"
host        := "http://localhost:3400"

// Pública
router.PublicRoute(r, "GET",  "/",    listar,     packageName, packagePath, host)
router.PublicRoute(r, "POST", "/",    crear,      packageName, packagePath, host)

// Protegida (requiere JWT válido)
router.ProtectRoute(r, "GET",    "/{id}", obtener,    packageName, packagePath, host)
router.ProtectRoute(r, "PUT",    "/{id}", actualizar, packageName, packagePath, host)
router.ProtectRoute(r, "DELETE", "/{id}", eliminar,   packageName, packagePath, host)

// Con autenticación + verificación de permisos
router.AuthorizationRoute(r, "GET", "/admin", adminHandler, packageName, packagePath, host)

// Con middlewares personalizados
router.With(r, "POST", "/upload", []func(http.Handler)http.Handler{miMiddleware}, upload, packageName, packagePath, host)
```

#### Response helpers

```go
import (
    "net/http"
    "github.com/celsiainternet/elvis/response"
)

func handler(w http.ResponseWriter, r *http.Request) {
    body,   err := response.GetBody(r)    // et.Json del body
    params       := response.GetQuery(r)  // et.Json de query params
    id           := response.GetParam(r, "id") // URL param chi

    // Respuestas
    response.JSON(w, r, http.StatusOK, data)               // {"ok": true, "result": data}
    response.ITEM(w, r, http.StatusOK, item)               // et.Item
    response.ITEMS(w, r, http.StatusOK, items)             // et.Items
    response.LIST(w, r, http.StatusOK, list)               // et.List: {"ok": true, "rows", "all", "count", "page", "start", "end", "result"}
    response.HTTPError(w, r, http.StatusBadRequest, "msg") // {"ok": false, "result": {"message": "msg"}}
    response.HTTPAlert(w, r, "alerta")                     // 400
    response.Unauthorized(w, r)                            // 401
    response.Forbidden(w, r)                               // 403
    response.InternalServerError(w, r, err)                // 500

    // Streaming paginado (útil para exportaciones grandes)
    response.Stream(w, r, 100, func(page, rows int) (et.Items, error) {
        return modelo.Select().Page(page, rows)
    })
}
```

### 🌍 Outbound HTTP client (`request`)

Client to call other services. Every function returns `(*request.Body, request.Status)`.

```go
import (
    "time"

    "github.com/celsiainternet/elvis/et"
    "github.com/celsiainternet/elvis/request"
)

header := et.Json{"Content-Type": "application/json", "Authorization": "Bearer " + token}

body, status := request.Post("https://api.example.com/users", header, et.Json{"name": "Ana"})
if !status.Ok {
    return fmt.Errorf("%d %s", status.Code, status.Message)
}

user, err := body.ToJson() // también: ToItem, ToItems, ToArrayJson, ToString, ToInt, ToInt64, ToFloat, ToBool, ToTime
```

| Family                              | Signature (after `path`/`method`)                                     | Use                                                   |
| ----------------------------------- | --------------------------------------------------------------------- | ----------------------------------------------------- |
| `Get/Post/Put/Delete/Patch/Options` | `header[, body]`                                                      | Plain call                                            |
| `<Verb>WithTls`                     | `header[, body], tlsConfig`                                           | mTLS / own CA (`request.NewTlsConfig(ca, cert, key)`) |
| `<Verb>WithTimeout`                 | `header[, body], timeout, defaultValue`                               | With a time limit                                     |
| `<Verb>WithTlsTimeout`              | `header[, body], tlsConfig, timeout, defaultValue`                    | TLS + time limit                                      |
| `Http` / `HttpWithTimeout`          | `method, path, header, body, tlsConfig[, timeout, defaultValue]`      | Dynamic method                                        |
| `HttpCtx` / `HttpCtxWithTimeout`    | `ctx, method, path, header, body, tlsConfig[, timeout, defaultValue]` | With a `context.Context`                              |

Behavior to keep in mind:

- `Content-Type` decides how the body is serialized: `application/json` (the value when the header is missing), `application/x-www-form-urlencoded` or `multipart/form-data`. Any other type sends no body.
- `timeout == 0` waits with no limit. When the timeout expires it returns `Status{Ok: false, Code: 408, Message: "timeout"}` and a `Body` with `defaultValue`.
- On any other error (network, invalid method, invalid `Content-Type`) the `*Body` is **`nil`**: check `status.Ok` before calling `body.ToJson()` and friends.
- In `multipart/form-data`, `request.File{Name, Data}` (or `*request.File`) values are sent as files and the rest as text; the `Content-Type` with the boundary replaces the header's.
- Known limitation: the timeout abandons the result but does not cancel the request in flight (it lives until the client timeout, 120 min).

```go
// Con timeout: si tarda más de 3s devuelve 408 y el defaultValue
body, status := request.GetWithTimeout(url, header, 3*time.Second, []byte(`{}`))
```

```go
// Subir un archivo como multipart/form-data
header := et.Json{"Content-Type": "multipart/form-data"}
body, status := request.Post(url, header, et.Json{
    "folder": "actas",
    "file":   request.File{Name: "acta.pdf", Data: data},
})
```

### 🔒 Stored secrets (`secret`)

Encrypts and decrypts secrets a service stores in its database (for example a source database password) with **AES-256-GCM**. The stored text is `"v1:" + base64(nonce || ciphertext)`, with a random nonce per encryption, and the key is the `sha256` of the configured secret, so a key of any length works. Whether to store the secret is each service's call; the package only encrypts and decrypts. Use it instead of `utility.EncryptMessage` (CFB with `SECRET` used raw as the key, without detecting tampered text), which stays unchanged.

```go
// Llave de la primera variable no vacía: MIGRATION_SECRET o, si no está, SECRET
c, err := secret.FromEnv("MIGRATION_SECRET", "SECRET")
if err != nil {
    return err // secret.ErrKeyRequired
}

stored, err := c.Encrypt(password) // "v1:..."
plain, err := c.Decrypt(stored)    // secret.ErrInvalid si se alteró, la llave cambió o no tiene "v1:"

// Atajos sin crear el Cipher (sin nombres usa SECRET)
stored, err = secret.Encrypt(password, "MIGRATION_SECRET", "SECRET")
plain, err = secret.Decrypt(stored, "MIGRATION_SECRET", "SECRET")

// En tests, llave fija sin pasar por envar (envar guarda el primer valor que lee)
c, _ = secret.New("llave-de-pruebas")
```

- `secret.New(key)` / `secret.FromEnv(names...)` → `(*secret.Cipher, error)`; `secret.EnvKey(names...)` returns the value of the first non-empty variable.
- Errors: `secret.ErrKeyRequired` (empty key) and `secret.ErrInvalid` (no prefix, broken base64, tampered text or another key). When GCM cannot open the text the error is `*secret.InvalidError`: `errors.Is(err, secret.ErrInvalid)` is true and `errors.Unwrap(err)` gives the GCM error. Services map them to their own messages with `errors.Is`.
- Changing the key makes what is already stored unreadable: those secrets must be written again.

### 🛡️ Resilience

Automatic retries for operations that can fail.

```go
import (
    "github.com/celsiainternet/elvis/et"
    "github.com/celsiainternet/elvis/resilience"
)

// Registrar operación con reintentos
// Configurar con RESILIENCE_TOTAL_ATTEMPTS (default 3) y RESILIENCE_TIME_ATTEMPTS en segundos (default 30)
instance := resilience.Add(
    "",                           // id (vacío = generado automáticamente)
    "email-bienvenida",           // tag
    "Enviar email de bienvenida", // descripción
    et.Json{"user_id": "123"},    // metadata
    "growth",                     // equipo
    "critical",                   // nivel
    func(to, asunto string) error {
        return enviarEmail(to, asunto)
    },
    "juan@ejemplo.com", // arg 1
    "¡Bienvenido!",     // arg 2
)

// Con parámetros personalizados
instance = resilience.AddCustom(
    "",
    "pago",
    "Procesar pago",
    5,               // intentos totales
    10*time.Second,  // tiempo entre intentos
    et.Json{},
    "pagos", "high",
    procesarPago, datoPago,
)

// Control de instancias
resilience.Stop(instance.Id)
resilience.Restart(instance.Id)
```

### 🔁 Workflows

Multi-step process orchestration with rollback and conditional expressions.

```go
import (
    "time"
    "github.com/celsiainternet/elvis/et"
    "github.com/celsiainternet/elvis/workflow"
)

// Definir el flujo
flow := workflow.New(
    "onboarding-usuario",
    "v1",
    "Onboarding",
    "Proceso completo de alta de usuario",
    func(inst *workflow.Instance, ctx et.Json) (et.Json, error) {
        return et.Json{"iniciado": true}, nil
    },
    false,     // modo debug
    "sistema", // equipo
)

// Resiliencia del flujo
flow.Resilence(3, 30*time.Second, "growth", "critical")

// Agregar pasos
flow.Step("CrearUsuario", "Crear usuario en BD", func(inst *workflow.Instance, ctx et.Json) (et.Json, error) {
    return et.Json{"user_id": utility.UUID()}, nil
}, false) // false = continúa automáticamente

// Rollback del paso anterior
flow.Rollback(func(inst *workflow.Instance, ctx et.Json) (et.Json, error) {
    return et.Json{"revertido": true}, nil
})

flow.Step("EnviarEmail", "Enviar email de bienvenida", func(inst *workflow.Instance, ctx et.Json) (et.Json, error) {
    return et.Json{"email_enviado": true}, nil
}, true) // true = paso de parada (se detiene hasta reanudar)

// Ejecutar instancia
result, err := workflow.Run(
    "",                              // instanceId (vacío = generado automáticamente)
    "onboarding-usuario",            // tag del flow
    -1,                              // step (-1 = ejecutar siguiente paso)
    et.Json{"env": "prod"},          // contexto
    et.Json{"email": "juan@e.com"},  // datos de entrada
    "sistema",
)
```

### 🔔 System events

Internal events the library emits, which you can subscribe to:

**Event / Workers**
| Channel | Description |
|---|---|
| `event:worker` | A job was fired (`event.Work`) |
| `event:worker:state` | A job changed state (`event.WorkState`) |
| `event:source` | A model action was published (`event.Source`) |
| `event:subscribed` | A new subscription was registered |
| `event:requests:overflow` | Request overflow |

**Workflows**
| Event | Description |
|---|---|
| `workflow:set` | A flow was created or updated |
| `workflow:delete` | A flow was deleted |
| `workflow:status` | An instance changed state |
| `workflow:awaiting` | An instance is waiting |
| `workflow:results` | Results are available |

**Resilience**
| Event | Description |
|---|---|
| `resilience:status` | An instance changed state |
| `resilience:stop` | Stop an instance by id |
| `resilience:restart` | Restart an instance by id |
| `resilience:failed` | An instance used up all its attempts |

**Database (JDB)**
| Event | Description |
|---|---|
| `sql:error` | SQL query error |
| `sql:query` | SQL query run |
| `sql:definition` | DDL run |
| `sql:command` | DML command run |

**API Gateway (Router)**
| Channel | Description |
|---|---|
| `apigateway/set/resolve` | New route registered |
| `apigateway/delete/resolve` | Route removed |
| `apigateway/reset` | Request to resync routes |
| `apigateway/set/proxy` | New proxy registered |
| `apigateway/delete/proxy` | Proxy removed |

## ⚙️ Configuration

Environment variables Elvis reads (`envar` also loads a `.env` file):

| Variable                    | Package       | Default                 | Description                                                        |
| --------------------------- | ------------- | ----------------------- | ------------------------------------------------------------------ |
| `DB_DRIVER`                 | jdb           | —                       | `postgres`, `mysql` or `oracle`                                    |
| `DB_HOST`                   | jdb           | —                       | Database host                                                      |
| `DB_PORT`                   | jdb           | `5432`                  | Database port                                                      |
| `DB_NAME`                   | jdb           | —                       | Database name                                                      |
| `DB_USER`                   | jdb           | —                       | Database user                                                      |
| `DB_PASSWORD`               | jdb           | —                       | Database password                                                  |
| `DB_APPLICATION_NAME`       | jdb           | `elvis`                 | Application name in PostgreSQL                                     |
| `USE_CORE`                  | jdb           | `true`                  | Initialize the core tables (series, records, recycling)            |
| `REDIS_HOST`                | cache         | —                       | Redis host (e.g. `localhost:6379`)                                 |
| `REDIS_PASSWORD`            | cache         | —                       | Redis password                                                     |
| `REDIS_DB`                  | cache         | `0`                     | Redis database number                                              |
| `NATS_HOST`                 | event         | —                       | NATS connection URL                                                |
| `NATS_USER`                 | event         | —                       | NATS user                                                          |
| `NATS_PASSWORD`             | event         | —                       | NATS password                                                      |
| `SECRET`                    | claim         | `1977`                  | JWT signing key                                                    |
| `HOST`                      | jrpc / router | `localhost`             | Host of the current service                                        |
| `PORT`                      | service       | `3400`                  | HTTP port                                                          |
| `RPC_HOST`                  | jrpc          | `HOST`                  | Host for RPC between services                                      |
| `RPC_PORT`                  | jrpc          | `4200`                  | RPC port                                                           |
| `AUTHORIZATION_METHOD`      | router        | —                       | RPC method that checks permissions                                 |
| `RESILIENCE_TOTAL_ATTEMPTS` | resilience    | `3`                     | Total attempts per operation                                       |
| `RESILIENCE_TIME_ATTEMPTS`  | resilience    | `30`                    | Seconds between retries                                            |
| `PIPE_HOST`                 | jrpc          | —                       | Host of the single RPC proxy (combined with `PIPE_PORT`)           |
| `PIPE_PORT`                 | jrpc          | `4200`                  | Port of the single RPC proxy                                       |
| `RPC_TIMEOUT`               | jrpc          | `600`                   | Seconds to wait for an RPC call (notifications too)                |
| `SEND_RPC`                  | service       | `Service-send.Services` | send RPC package for `service.Send*`                               |
| `OTP_RPC`                   | service       | `Service-otp.Services`  | otp RPC package for `service.SendOtp*`                             |
| `STAGE`                     | event         | `local`                 | Environment prefix for the pipe channel (`pipe:<stage>:<channel>`) |
| `PRODUCTION`                | dt            | `true`                  | Enables Redis persistence of the `dt.Object` object cache          |

## 🔧 Development

```bash
# Build and check (the main check: there are few tests)
go build ./...
go vet ./...

# Tests
go test ./...
go test ./secret/ -run TestName -v   # one test
go test ./jquery/...                 # jquery's tests live in jquery/test/

# Generator from this checkout
gofmt -w . && go run ./cmd/create go

# Release a version: creates the tag and pushes it (only when releasing)
./version.sh --request   # v1.1.X
```

There are tests in `secret`, `linq` (SQL of `RunQuery`), `jquery/test`, `service` (RPC notifications), `event` (`ParseErrorReport`), `request`, `queue` and `create/v1` (the generator's templates compile and pass `gofmt`). For the other packages `go test ./...` only checks that they compile.

`version.sh` replaces the version in the READMEs' install command but does not commit that change.

For projects that use the jdb ORM (`internal/models` + `internal/services` layout) the generator lives in the jdb repo: `go run github.com/celsiainternet/jdb/cmd/create go`.

The documentation lives in three files: this `README.md` (English), `README.es.md` (the same structure in Spanish) and [`spec.md`](spec.md), with the requirements, their verified status and the gap. The context for agents is in `CLAUDE.md`.

### 📁 Structure

```
elvis/
├── agentsguide/    # Agents framework; Install() writes it into CLAUDE.md and AGENTS.md
├── authorization/  # Permission model (linq) with set/del events
├── cache/          # Redis client (Set, Get, Delete, Pub/Sub)
├── claim/          # JWT: generation, validation, invalidation
├── cmd/
│   ├── create/     # Generator CLI (uses create/v1)
│   ├── install/    # Prepares a consuming project (tidy + agents framework)
│   └── crontab/, flow/, jql/  # Examples of the matching packages
├── config/         # Configuration loading
├── console/        # Application logging: prints and publishes each entry as a `logs` NATS event
├── create/v1/      # Project generator (templates in model.go)
├── crontab/        # Scheduled jobs (cron)
├── dt/             # Step-based resilience counter and Redis-backed object cache
├── envar/          # Environment variables (loads .env)
├── et/             # Core types: Json, Item, Items, List, Any
├── event/          # Local (emitter) and distributed (NATS) events
├── file/           # File handling
├── health/         # Health check
├── inbox/          # Per-user notification inbox (linq)
├── instances/      # Persistent workflow/resilience instance registry in the DB
├── jdb/            # Database abstraction (Postgres/MySQL/Oracle)
├── jquery/         # Translates an et.Json (from/join/where/...) into SQL SELECT, pluggable dialect
├── jrpc/           # TCP RPC between services (net/rpc); Redis stores the method registry
├── jtls/           # Self-signed TLS certificates
├── linq/           # ORM / query builder
├── logs/           # Local logging, for elvis packages that cannot import console
├── mem/            # In-memory cache with TTL
├── middleware/     # chi HTTP middleware (auth, cors, logger, etc.)
├── msg/            # Shared messages
├── queue/          # In-process batching queue (queue.Queue[T])
├── race/           # Concurrency helpers
├── reg/            # ID registry
├── request/        # Outbound HTTP client
├── resilience/     # Automatic retries
├── response/       # HTTP response helpers
├── router/         # chi route registration with the API Gateway
├── secret/         # Encryption of stored secrets (AES-256-GCM, v1: format)
├── service/        # Notifications (send) and OTP (otp) through RPC, falling back to event.Work
├── stdrout/        # Standard output / terminal
├── strs/           # String utilities
├── timezone/       # Time zones
├── utility/        # General utilities (UUID, OTP, crypto, etc.)
├── workflow/       # Multi-step flow orchestration
└── xls/            # Excel read/write (excelize)
```

## 📄 License

Distributed under the MIT license. See [LICENSE](LICENSE) for details.
