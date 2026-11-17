# 🎸 Elvis — Framework para Microservicios en Go

[![Go Version](https://img.shields.io/github/go-mod/go-version/celsiainternet/elvis?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/github/license/celsiainternet/elvis?style=flat-square)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/celsiainternet/elvis?style=flat-square&logo=github)](https://github.com/celsiainternet/elvis/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/celsiainternet/elvis?style=flat-square)](https://goreportcard.com/report/github.com/celsiainternet/elvis)
[![GitHub Stars](https://img.shields.io/github/stars/celsiainternet/elvis?style=flat-square&logo=github)](https://github.com/celsiainternet/elvis/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/celsiainternet/elvis?style=flat-square&logo=github)](https://github.com/celsiainternet/elvis/issues)
[![Documentation](https://img.shields.io/badge/docs-godoc-blue?style=flat-square&logo=go)](https://pkg.go.dev/github.com/celsiainternet/elvis)

> 🚀 **Librería de infraestructura para construir microservicios escalables en Go**

<div align="center">

**[📚 Documentación](https://pkg.go.dev/github.com/celsiainternet/elvis)** •
**[🚀 Quick Start](#-quick-start)** •
**[🐛 Issues](https://github.com/celsiainternet/elvis/issues)** •
**[💬 Discusiones](https://github.com/celsiainternet/elvis/discussions)**

</div>

---

## 📑 Tabla de Contenidos

- [Instalación](#-instalación)
- [Descripción](#-descripción)
- [Requisitos Previos](#-requisitos-previos)
- [Quick Start](#-quick-start)
- [Tipos de Datos Centrales (et)](#-tipos-de-datos-centrales-et)
- [Base de Datos (jdb)](#-base-de-datos-jdb)
- [ORM / Query Builder (linq)](#-orm--query-builder-linq)
- [Cache (cache / mem)](#-cache-cache--mem)
- [Eventos (event)](#-eventos-event)
- [Autenticación y Autorización (claim / middleware)](#-autenticación-y-autorización-claim--middleware)
- [HTTP Router y Respuestas (router / response)](#-http-router-y-respuestas-router--response)
- [Cliente HTTP saliente (request)](#-cliente-http-saliente-request)
- [Secretos guardados (secret)](#-secretos-guardados-secret)
- [Resiliencia](#-resiliencia)
- [Workflows](#-workflows)
- [Variables de Entorno](#-variables-de-entorno)
- [Comandos de Desarrollo](#-comandos-de-desarrollo)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Eventos del Sistema](#-eventos-del-sistema)
- [Licencia](#-licencia)

---

## 🚀 Instalación

Para empezar un proyecto con Elvis, ejecuta estos comandos en una carpeta nueva:

```bash
go mod init github.com/<organizacion>/<proyecto>

go get github.com/celsiainternet/elvis@v1.1.324
go run github.com/celsiainternet/elvis/cmd/install
go run github.com/celsiainternet/elvis/cmd/create go
```

| Comando | Qué hace |
| ------- | -------- |
| `go mod init ...` | Crea el `go.mod`. Los dos comandos siguientes lo necesitan: toman de ahí la ruta del módulo y se detienen si no existe. |
| `go get github.com/celsiainternet/elvis@v1.1.324` | Agrega Elvis como dependencia y deja en `go.sum` las entradas que necesitan los `go run` siguientes; sin este paso fallan con `missing go.sum entry`. |
| `go run .../cmd/install` | Prepara el proyecto: ejecuta `go mod tidy` si ya hay paquetes Go e instala el marco de trabajo para agentes (ver la nota de abajo). No fija versiones de otras dependencias. |
| `go run .../cmd/create go` | Abre un menú interactivo para generar código: **Project** (proyecto completo), **Microservice**, **Modelo** (modelo `linq` + CRUD en un paquete existente) o **Rpc**. |

> Tanto `cmd/install` como la opción **Project** de `cmd/create`
> instalan automáticamente un marco de trabajo para agentes (Claude Code, Codex
> y compatibles): copian las guías en `.claude/guides/elvis/` y agregan un
> bloque en `CLAUDE.md` y `AGENTS.md` que las referencia. Volver a ejecutarlo
> actualiza las guías a la versión instalada de Elvis sin tocar el resto de
> esos archivos. Si el proyecto usa jdb, `go run github.com/celsiainternet/jdb/cmd/install`
> agrega su guía. Ver [`agentsguide/MARCO.md`](agentsguide/MARCO.md).

---

## 📖 Descripción

**Elvis** es una librería Go (`github.com/celsiainternet/elvis`) que provee primitivas de infraestructura para construir microservicios. No es una aplicación en sí misma, sino una librería compartida que otros servicios consumen.

Incluye:

- 🗄️ **Abstracción de base de datos** multi-driver (PostgreSQL, MySQL, Oracle)
- 🔍 **Query builder** estilo LINQ con soporte a JSONB
- 💾 **Cache** multi-backend (Redis + memoria)
- 🔄 **Sistema de eventos** local (in-process) y distribuido (NATS)
- 🔐 **Autenticación JWT** con invalidación en Redis
- 🛡️ **Middleware HTTP** integrado (auth, CORS, logging, telemetría)
- 🔁 **Sistema de resiliencia** con reintentos automáticos
- 📋 **Workflows** con pasos, rollback y expresiones condicionales
- 📅 **Tareas programadas** (Crontab)
- 🛠️ **CLI de scaffolding** para generar nuevos proyectos microservicio, incluyendo
  `deployments/<servicio>/` con `local.yml` (docker-compose) y manifiestos de
  Kubernetes (`oke-template.yml` para `Deployment`, `oke-statefulset-template.yml`
  para `StatefulSet`)

---

## 📋 Requisitos Previos

- **Go 1.23** o superior
- **PostgreSQL** — base de datos principal
- **Redis** — cache y almacenamiento de tokens JWT
- **NATS** — mensajería distribuida (eventos y RPC entre servicios)

---

## ⚡ Quick Start

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

---

## 🧱 Tipos de Datos Centrales (`et`)

El paquete `et` es la base de toda la librería. Define los tipos de dato comunes usados en toda la API.

### `et.Json`

`map[string]interface{}` con métodos de acceso tipados:

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

### `et.Item`

Resultado unitario de una consulta:

```go
// Uso típico
item, err := modelo.QueryOne("SELECT * FROM tabla WHERE _id=$1", id)
if item.Ok {
    nombre := item.Result.Str("nombre")
}
```

### `et.Items`

Resultado paginado de una consulta:

```go
// Acceder al primer elemento
first := items.First()

// Iterar
for _, item := range items.Result {
    fmt.Println(item.Str("nombre"))
}
```

### `et.List`

Resultado con metadatos completos de paginación: `Rows`, `All`, `Count`, `Page`, `Start`, `End`, `Result`.

---

## 🗄️ Base de Datos (`jdb`)

Abstracción multi-driver que soporta **PostgreSQL**, **MySQL** y **Oracle**.

### Conexión

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

### Consultas directas

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

### Core Tables

Cuando `USE_CORE=true` (default), `jdb.Load()` inicializa tres tablas internas:

- **series** — secuencias auto-incrementales por tag
- **records** — auditoría de cambios
- **recycling** — papelera (soft deletes)

```go
// Siguiente valor de serie numérica
siguiente := jdb.NextSerie(db, "mi.tabla")

// Siguiente código con prefijo
codigo := jdb.NextCode(db, "mi.tabla", "USR")
// Resultado: "USR000001"
```

---

## 🔍 ORM / Query Builder (`linq`)

Basado en LINQ, permite definir modelos tipados y construir consultas de forma fluida.

### Definir Schema y Modelo

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

### Triggers

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

### Consultas CRUD

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

> Los comandos de escritura (`Insert`, `Update`, `Upsert`, `Delete`) se ejecutan con `.Command()` (`et.Items`) o `.CommandOne()` (`et.Item`, primer resultado); las consultas (`Select`, `Data`) usan `.All()`, `.Page(page, rows)`, `.List(page, rows)` o `.First()`, no `.One()`.

### Condiciones disponibles

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

> Estas condiciones pasan el valor por `et.Unquote`, que no escapa la comilla simple: no les pases texto del cliente sin validar. Para consultas que llegan del cliente usa `linq.RunQuery`.

### Consulta jsql (`linq.RunQuery`)

`RunQuery(model, query, wheres...)` arma la consulta desde un `et.Json` con el contrato de `jdb.RunQuery` y devuelve la página (`et.List`):

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

> Sin publicar: `linq/jsql.go` no está en ningún tag todavía (el último es `v1.1.324`). Un servicio que compila fuera del `go.work` necesita el tag que lo incluya.

**Formato de `query`**

| Llave | Forma | Qué hace |
|-------|-------|----------|
| `selects` | `["_id", "name"]` | Columnas del resultado. Vacío o ausente: todas. |
| `wheres` | `[{ campo: { operador: valor } }, …]` | Condiciones con AND. `{ "and": { campo: {…} } }` u `{ "or": { campo: {…} } }` usan ese conector. |
| `order_by` | `{ "asc": [...], "desc": [...] }` | Primero los `asc`, luego los `desc`. |
| `search` | `{ "value": "ana gomez", "fields": ["name", "description"] }` | Búsqueda de texto libre (contrato de los `/query` de catalog), ver abajo. |
| `page` | número | Página pedida; 1 si no viene o es menor que 1. |
| `rows` | número | Filas por página (alias `limit`). 100 por defecto (`QUERY_DEFAULT_ROWS`), máximo 1000 (`QUERY_MAX_ROWS`). |

- Operadores: `eq` (`=`), `neg` (`!=`), `in` (lista), `like` (`ILIKE`, sin distinguir mayúsculas; los `%` los pone el cliente), `more` (`>`), `less` (`<`), `moreEq` (`>=`), `lessEq` (`<=`). Los nombres de operador distinguen mayúsculas; `and`/`or` no.
- Un campo puede llevar varios operadores (`{ "date_make": { "more": "2026-01-01", "less": "2026-12-31" } }`): el primero lleva el conector y los demás van con AND.
- Un grupo `and`/`or` puede traer varios campos: cada uno va con el conector del grupo, en orden alfabético, y sus operadores adicionales con AND. `{ "or": { "a": { "eq": 1 }, "b": { "eq": 2 } } }` es `OR A = 1 OR B = 2`. El grupo debe ser un objeto con al menos un campo.
- `concat`: `{ "concat": { "fields": ["name", "description"], "like": "%juan%" } }` compara contra los campos unidos con un espacio. `fields` es obligatorio y todos deben existir y ser visibles; además de `fields` lleva al menos un operador.
- Los campos pueden ser columnas o atributos de `_data` (`DefineAtrib`).
- `wheres` (variádico) son condiciones fijas del llamador, por ejemplo el proyecto del token: van con AND y las del cliente van entre paréntesis si traen `or`, así un `or` no las salta.
- Los textos se escapan (comilla simple duplicada) sin pasar por `et.Unquote`. Solo se aceptan textos, números y booleanos; `in` acepta cualquier lista (`[]any`, `[]string`, …) no vacía de esos valores.
- `search`: cada palabra de `value` (separadas por espacios, máximo 5: `SEARCH_MAX_WORDS`; las demás se ignoran) debe aparecer, sin distinguir mayúsculas, en los `fields` unidos con un espacio: `CONCAT(A.NAME, ' ', A.DESCRIPTION)::TEXT ILIKE '%ana%' AND ... ILIKE '%gomez%'`. Cada palabra es una condición AND que va **después** de los `wheres` fijos y de los del cliente (fuera de su paréntesis), así la búsqueda no salta el alcance. `%`, `_` y `\` de `value` se buscan como texto. Sin `value`, o si `search` no es un objeto (por ejemplo el texto suelto de `jdb.RunQuery`), no filtra. `fields` admite columnas y atributos visibles (1 a 5: `SEARCH_MAX_FIELDS`), no JSON ni binarios. No usa índice (es `ILIKE` sobre un `CONCAT`).

**Respuesta**: `et.List` (`rows`, `all` = total de registros que cumplen, `count` = registros de la página, `page`, `start`, `end`, `result`). Cada elemento de `result` es el objeto del registro (consulta `model.Data`).

**Errores**

- `Llave no valida en la consulta (<llave>)` (`MSG_QUERY_KEY_NOT_VALID`) si `query` trae una llave raíz distinta de `wheres`, `selects`, `order_by`, `search`, `page`, `rows` y `limit` (por ejemplo `project_id` o `folder`: el servicio que las use las saca antes de llamar a `RunQuery`).
- `Columna no valida en la consulta (<campo>)` (`MSG_QUERY_COLUMN_NOT_VALID`) si un campo de `selects`, `order_by` o `wheres` (también dentro de un grupo `and`/`or` y en `concat.fields`) no existe o es oculto (`DefineHidden`).
- `Operador no valido en la consulta (<campo>.<operador>)` (`MSG_QUERY_OPERATOR_NOT_VALID`) si el operador no es uno de los soportados (por ejemplo `isNull`, `notIn`, `between`, `search`).
- `Valor no valido en la consulta (<campo>[.<operador>])` (`MSG_QUERY_VALUE_NOT_VALID`) si el valor no es texto, número ni booleano (`null`, objeto, lista fuera de `in`), si `in` no es una lista, está vacío o trae un elemento no válido, si un campo no trae un objeto con al menos un operador, si un grupo `and`/`or` no es un objeto o está vacío, si `concat` no trae `fields`, o si `wheres` no es una lista de objetos (`(wheres)`).
- En `search`: `MSG_SEARCH_LENGTH` si `value` pasa de 100 caracteres (`SEARCH_MAX_LENGTH`), `MSG_SEARCH_FIELDS` si `fields` está vacío o trae más de 5, y `MSG_SEARCH_FIELD_NOT_VALID` si un campo no existe, es oculto o es JSON/binario.
- El error de la base, si la consulta falla.

Estos errores salen antes de consultar la base, así que un filtro mal escrito responde 400 en vez de devolver todo lo que permiten los `wheres` fijos. Los `wheres` fijos del llamador (variádicos) no se validan. Sin `value`, o si `search` no es un objeto, la búsqueda sigue sin filtrar ni dar error.

**No soporta** (sí están en `jdb.RunQuery`): `isNull`, `notNull`, `notIn`, `between`, el operador `search` dentro de `wheres` ni la llave raíz `search` como texto (`DefineSearch` de jdb), ni rutas JSON (`_data.campo`). Para "sin valor" en un texto usa `{ "eq": "" }`.

**Handler `POST /query`**

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

### Referencias entre modelos

```go
// Llave foránea
modeloOrden.DefineForeignKey("USUARIO_ID", modeloUsuario.Col(jdb.KEY))

// Referencia (columna virtual que trae el nombre del otro modelo)
modeloOrden.DefineReference("USUARIO_ID", "USUARIO", jdb.KEY, modeloUsuario.Col("NOMBRE"), false)
```

---

## 💾 Cache (`cache` / `mem`)

### Redis (`cache`)

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

### Memoria (`mem`)

Cache in-process con TTL, inicializado automáticamente:

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

---

## 🔄 Eventos (`event`)

El sistema de eventos tiene dos modos: **local** (in-process) y **distribuido** (NATS).

### Eventos locales (in-process)

```go
import "github.com/celsiainternet/elvis/event"

// Registrar handler
event.On("usuario.creado", func(msg event.EvenMessage) {
    fmt.Println("Usuario:", msg.Data.Str("nombre"))
})

// Emitir evento
event.Emit("usuario.creado", et.Json{"nombre": "Juan"})
```

### Eventos distribuidos (NATS)

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

> Los handlers de `Subscribe`, `Queue` y `Stack` están protegidos con `recover()`. Un panic dentro del handler se loguea como error en lugar de matar el goroutine de NATS.

### Reportar errores al registro central (`event.Error`)

Publica un error en `EVENT_ERROR` (`"event:error"`). Catalog lo guarda en su registro central: agrupa los iguales y el front abre el ticket en ADA o GLPI. `Error` valida el reporte y devuelve ese error; si es válido publica en otra goroutine, así que nunca bloquea ni hace fallar a quien reporta. Sin conexión a NATS no hace nada.

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

`OccurredAt` es la hora actual si no se envía. Quien consume el canal lee el mensaje con `event.ParseErrorReport(msg.Data)`, que normaliza y valida igual.

### Work y WorkState

Utilidades para trazabilidad de tareas asíncronas:

```go
// Emite EVENT_WORK ("event:worker") y el evento propio del canal
work := event.Work("factura.generar", et.Json{"factura_id": "F-001"})
workId := work.Str("_id")

// Actualizar estado de la tarea (emite EVENT_WORK_STATE)
event.WorkState(workId, event.WorkStatusProcessing, et.Json{})
event.WorkState(workId, event.WorkStatusCompleted, et.Json{"url": "..."})
// Estados: WorkStatusPending, WorkStatusAccepted, WorkStatusProcessing, WorkStatusCompleted, WorkStatusFailed
```

---

## 🔐 Autenticación y Autorización (`claim` / `middleware`)

### Generar tokens JWT

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

La variable de entorno `SECRET` (default `"1977"`) es la clave de firma JWT.

### Validar y parsear tokens

```go
// Parsear sin validar en Redis
c, err := claim.ParceToken(token)

// Validar (parsea + verifica que el token esté activo en Redis)
c, err := claim.ValidToken(token)

// Logout
err = claim.DeleteToken(app, device, userId)
err = claim.DeleteTokeByToken(token)
```

### Leer datos del cliente desde el request

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

### Middleware HTTP

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

---

## 🌐 HTTP Router y Respuestas (`router` / `response`)

### Registro de rutas

Todas las rutas registradas se publican automáticamente al API Gateway vía NATS.

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

### Helpers de respuesta

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

---

## 🌍 Cliente HTTP saliente (`request`)

Cliente para llamar a otros servicios. Todas las funciones devuelven `(*request.Body, request.Status)`.

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

| Familia                             | Firma (después del `path`/`method`)                                   | Uso                                                      |
| ----------------------------------- | --------------------------------------------------------------------- | -------------------------------------------------------- |
| `Get/Post/Put/Delete/Patch/Options` | `header[, body]`                                                      | Llamada simple                                           |
| `<Verbo>WithTls`                    | `header[, body], tlsConfig`                                           | mTLS / CA propia (`request.NewTlsConfig(ca, cert, key)`) |
| `<Verbo>WithTimeout`                | `header[, body], timeout, defaultValue`                               | Con límite de tiempo                                     |
| `<Verbo>WithTlsTimeout`             | `header[, body], tlsConfig, timeout, defaultValue`                    | TLS + límite de tiempo                                   |
| `Http` / `HttpWithTimeout`          | `method, path, header, body, tlsConfig[, timeout, defaultValue]`      | Método dinámico                                          |
| `HttpCtx` / `HttpCtxWithTimeout`    | `ctx, method, path, header, body, tlsConfig[, timeout, defaultValue]` | Con `context.Context`                                    |

Comportamiento a tener en cuenta:

- `Content-Type` en el header es obligatorio y decide cómo se serializa el body (`application/json`, `application/x-www-form-urlencoded`; `multipart/form-data` **no funciona** todavía, ver más abajo).
- `timeout == 0` espera sin límite. Si vence el timeout se devuelve `Status{Ok: false, Code: 408, Message: "timeout"}` y un `Body` con `defaultValue`.
- En cualquier otro error (red, método inválido, `Content-Type` inválido) el `*Body` es **`nil`**: revisa `status.Ok` antes de llamar a `body.ToJson()` y compañía.
- Limitaciones conocidas: el timeout abandona el resultado pero no cancela la petición en curso (sigue viva hasta el timeout del cliente, 120 min), y `multipart/form-data` produce un panic (buffer nil dentro de una goroutine).

```go
// Con timeout: si tarda más de 3s devuelve 408 y el defaultValue
body, status := request.GetWithTimeout(url, header, 3*time.Second, []byte(`{}`))
```

---

## 🔒 Secretos guardados (`secret`)

Cifra y descifra secretos que un servicio guarda en su base (por ejemplo la contraseña de una base origen) con **AES-256-GCM**. El texto guardado es `"v1:" + base64(nonce || cifrado)`, con un nonce aleatorio en cada cifrado, y la llave es `sha256` del secreto configurado, así que acepta una llave de cualquier largo. Guardar o no el secreto lo decide cada servicio; el paquete solo cifra y descifra. Úsalo en lugar de `utility.EncryptMessage` (CFB con `SECRET` tal cual como llave, sin detectar un texto alterado), que se mantiene sin cambios.

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

- `secret.New(key)` / `secret.FromEnv(names...)` → `(*secret.Cipher, error)`; `secret.EnvKey(names...)` devuelve el valor de la primera variable no vacía.
- Errores: `secret.ErrKeyRequired` (llave vacía) y `secret.ErrInvalid` (sin prefijo, base64 roto, texto alterado u otra llave). Cuando GCM no abre el texto el error es `*secret.InvalidError`: `errors.Is(err, secret.ErrInvalid)` es verdadero y `errors.Unwrap(err)` da el error de GCM. Los servicios los traducen a sus propios mensajes con `errors.Is`.
- Cambiar la llave deja ilegible lo ya guardado: hay que volver a escribir esos secretos.

---

## 🛡️ Resiliencia

Sistema de reintentos automáticos para operaciones que pueden fallar.

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

---

## 🔁 Workflows

Orquestación de procesos multi-paso con soporte a rollback y expresiones condicionales.

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

---

## ⚙️ Variables de Entorno

| Variable                    | Paquete       | Default     | Descripción                                                     |
| --------------------------- | ------------- | ----------- | --------------------------------------------------------------- |
| `DB_DRIVER`                 | jdb           | —           | `postgres`, `mysql` u `oracle`                                  |
| `DB_HOST`                   | jdb           | —           | Host de la base de datos                                        |
| `DB_PORT`                   | jdb           | `5432`      | Puerto de la base de datos                                      |
| `DB_NAME`                   | jdb           | —           | Nombre de la base de datos                                      |
| `DB_USER`                   | jdb           | —           | Usuario de la base de datos                                     |
| `DB_PASSWORD`               | jdb           | —           | Contraseña de la base de datos                                  |
| `DB_APPLICATION_NAME`       | jdb           | `elvis`     | Nombre de la aplicación en PostgreSQL                           |
| `USE_CORE`                  | jdb           | `true`      | Inicializar tablas core (series, records, recycling)            |
| `REDIS_HOST`                | cache         | —           | Host de Redis (ej. `localhost:6379`)                            |
| `REDIS_PASSWORD`            | cache         | —           | Contraseña de Redis                                             |
| `REDIS_DB`                  | cache         | `0`         | Número de base de datos Redis                                   |
| `NATS_HOST`                 | event         | —           | URL de conexión NATS                                            |
| `NATS_USER`                 | event         | —           | Usuario NATS                                                    |
| `NATS_PASSWORD`             | event         | —           | Contraseña NATS                                                 |
| `SECRET`                    | claim         | `1977`      | Clave de firma JWT                                              |
| `HOST`                      | jrpc / router | `localhost` | Host del servicio actual                                        |
| `PORT`                      | servicio      | `3400`      | Puerto HTTP                                                     |
| `RPC_HOST`                  | jrpc          | `HOST`      | Host para RPC entre servicios                                   |
| `RPC_PORT`                  | jrpc          | `4200`      | Puerto RPC                                                      |
| `AUTHORIZATION_METHOD`      | router        | —           | Método RPC para verificar permisos                              |
| `RESILIENCE_TOTAL_ATTEMPTS` | resilience    | `3`         | Intentos totales por operación                                  |
| `RESILIENCE_TIME_ATTEMPTS`  | resilience    | `30`        | Segundos entre reintentos                                       |
| `PIPE_HOST`                 | jrpc          | —           | Host del proxy RPC único (se combina con `PIPE_PORT`)           |
| `PIPE_PORT`                 | jrpc          | `4200`      | Puerto del proxy RPC único                                      |
| `STAGE`                     | event         | `local`     | Prefijo de entorno para canal pipe (`pipe:<stage>:<canal>`)     |
| `PRODUCTION`                | dt            | `true`      | Habilita persistencia en Redis del cache de objetos `dt.Object` |

---

## 🔧 Comandos de Desarrollo

```bash
# Compilar el módulo
go build ./...

# Ejecutar tests
go test ./...

# Ejecutar un test específico con verbose
go test ./paquete/... -run NombreTest -v

# Formatear código y ejecutar CLI de scaffolding (create/v1)
gofmt -w . && go run ./cmd/create go

# Para proyectos que usan el ORM jdb (layout internal/models + internal/services),
# el generador vive en el repo jdb
go run github.com/celsiainternet/jdb/cmd/create go

# Actualizar dependencias
go mod tidy
```

---

## 📁 Estructura del Proyecto

```
elvis/
├── agentsguide/    # Marco de trabajo para agentes; Install() lo escribe en CLAUDE.md
├── authorization/  # Modelo de permisos (linq) con eventos set/del
├── cache/          # Cliente Redis (Set, Get, Delete, Pub/Sub)
├── claim/          # JWT: generación, validación, invalidación
├── cmd/
│   └── create/     # CLI scaffolding de microservicios (usa create/v1)
├── config/         # Carga de configuración
├── console/        # Logging interno de bajo nivel
├── create/
│   └── v1/         # Generador de proyectos (usado por cmd/create)
├── crontab/        # Tareas programadas (cron)
├── dt/             # Contador de pasos por resiliencia y cache de objetos respaldado en Redis
├── envar/          # Helpers de variables de entorno
├── et/             # Tipos centrales: Json, Item, Items, List, Any
├── event/          # Eventos local (emitter) y distribuido (NATS)
├── file/           # Manejo de archivos
├── health/         # Health check helpers
├── inbox/          # Bandeja de notificaciones por usuario (linq)
├── instances/      # Registro persistente de instancias de workflow/resiliencia en BD
├── jdb/            # Abstracción de base de datos (Postgres/MySQL/Oracle)
├── jquery/         # Traduce un et.Json (from/join/where/...) a SQL SELECT, dialecto pluggable
├── jrpc/           # RPC entre servicios vía Redis
├── jtls/           # Generación de certificados TLS auto-firmados
├── linq/           # ORM / query builder
├── logs/           # Logging estructurado
├── mem/            # Cache in-memory con TTL
├── middleware/     # Middleware HTTP chi (auth, cors, logger, etc.)
├── msg/            # Mensajes de error compartidos
├── queue/          # Cola de batching en proceso (queue.Queue[T])
├── race/           # Helpers de concurrencia
├── reg/            # Registro de IDs
├── request/        # Cliente HTTP saliente (GET/POST/PUT/DELETE/PATCH/OPTIONS, TLS y timeout)
├── resilience/     # Reintentos automáticos
├── response/       # Helpers de respuesta HTTP
├── router/         # Registro de rutas chi con API Gateway
├── secret/         # Cifrado de secretos guardados (AES-256-GCM, formato v1:)
├── service/        # Cliente HTTP entre servicios
├── stdrout/        # Salida estándar / terminal
├── strs/           # Utilidades de strings
├── timezone/       # Manejo de zonas horarias
├── utility/        # Utilidades generales (UUID, OTP, crypto, etc.)
├── workflow/       # Orquestación de flujos multi-paso
└── xls/            # Lectura/escritura de Excel (excelize)
```

---

## 🔔 Eventos del Sistema

Eventos internos que emite la librería y a los que se puede suscribir:

**Event / Workers**
| Canal | Descripción |
|---|---|
| `event:worker` | Se disparó un trabajo (`event.Work`) |
| `event:worker:state` | Cambio de estado de un trabajo (`event.WorkState`) |
| `event:source` | Acción de modelo publicada (`event.Source`) |
| `event:subscribed` | Nueva suscripción registrada |
| `event:requests:overflow` | Desbordamiento de solicitudes |

**Workflows**
| Evento | Descripción |
|---|---|
| `workflow:set` | Se creó o actualizó un flujo |
| `workflow:delete` | Se eliminó un flujo |
| `workflow:status` | Cambio de estado de una instancia |
| `workflow:awaiting` | Instancia en espera |
| `workflow:results` | Resultados disponibles |

**Resiliencia**
| Evento | Descripción |
|---|---|
| `resilience:status` | Cambio de estado de una instancia |
| `resilience:stop` | Detener instancia por id |
| `resilience:restart` | Reiniciar instancia por id |
| `resilience:failed` | Instancia agotó todos los intentos |

**Base de Datos (JDB)**
| Evento | Descripción |
|---|---|
| `sql:error` | Error en consulta SQL |
| `sql:query` | Consulta SQL ejecutada |
| `sql:definition` | DDL ejecutado |
| `sql:command` | Comando DML ejecutado |

**API Gateway (Router)**
| Canal | Descripción |
|---|---|
| `apigateway/set/resolve` | Nueva ruta registrada |
| `apigateway/delete/resolve` | Ruta eliminada |
| `apigateway/reset` | Solicitud de re-sincronización de rutas |
| `apigateway/set/proxy` | Nuevo proxy registrado |
| `apigateway/delete/proxy` | Proxy eliminado |

---

## 📄 Licencia

Distribuido bajo la licencia MIT. Ver [LICENSE](LICENSE) para más detalles.

## Test

```
go test ./queue/... -run Example_queue -v
```
