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

- [Descripción](#-descripción)
- [Requisitos Previos](#-requisitos-previos)
- [Instalación](#-instalación)
- [Quick Start](#-quick-start)
- [Tipos de Datos Centrales (et)](#-tipos-de-datos-centrales-et)
- [Base de Datos (jdb)](#-base-de-datos-jdb)
- [ORM / Query Builder (linq)](#-orm--query-builder-linq)
- [Cache (cache / mem)](#-cache-cache--mem)
- [Eventos (event)](#-eventos-event)
- [Autenticación y Autorización (claim / middleware)](#-autenticación-y-autorización-claim--middleware)
- [HTTP Router y Respuestas (router / response)](#-http-router-y-respuestas-router--response)
- [Resiliencia](#-resiliencia)
- [Workflows](#-workflows)
- [Variables de Entorno](#-variables-de-entorno)
- [Comandos de Desarrollo](#-comandos-de-desarrollo)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Eventos del Sistema](#-eventos-del-sistema)
- [Licencia](#-licencia)

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
- 🛠️ **CLI de scaffolding** para generar nuevos proyectos microservicio

---

## 📋 Requisitos Previos

- **Go 1.23** o superior
- **PostgreSQL** — base de datos principal
- **Redis** — cache y almacenamiento de tokens JWT
- **NATS** — mensajería distribuida (eventos y RPC entre servicios)

---

## 🚀 Instalación

```bash
# En el módulo de tu proyecto
go get github.com/celsiainternet/elvis@latest
go get github.com/celsiainternet/elvis@v1.1.223
go run github.com/celsiainternet/elvis/cmd/install
```

Para generar un nuevo microservicio con la estructura base de Elvis:

```bash
go run github.com/celsiainternet/elvis/cmd/create go
```

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
    r.Use(middleware.Cors)
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
}).One()

// UPDATE
item, err := modelo.Update(et.Json{"NOMBRE": "Juan Pablo"}).
    Where(modelo.Col(jdb.KEY).Eq(id)).
    One()

// UPSERT
item, err := modelo.Upsert(et.Json{
    jdb.KEY:  id,
    "NOMBRE": "Juan",
}).One()

// DELETE
item, err := modelo.Delete().
    Where(modelo.Col(jdb.KEY).Eq(id)).
    One()

// SELECT múltiple
items, err := modelo.Select().
    Where(modelo.Col("_STATE").Eq("0")).
    OrderBy(modelo.Col("NOMBRE"), true).
    All()

// SELECT paginado
items, err := modelo.Select().
    Where(modelo.Col("_STATE").Eq("0")).
    Page(1, 20).
    List()

// SELECT uno
item, err := modelo.Select().
    Where(modelo.Col(jdb.KEY).Eq(id)).
    One()
```

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

// Suscribirse
err = event.Subscribe("pedido.nuevo", func(msg event.EvenMessage) {
    fmt.Println("Pedido:", msg.Data)
})

// Publicar
event.Publish("pedido.nuevo", et.Json{"pedido_id": "123"})

// Stack: handler persistente que se re-registra tras reconexión
event.Stack("canal/reset", func(msg event.EvenMessage) {
    // re-sincronización
})
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

// Token con autorización (incluye projectId y profileTp)
token, err := claim.NewAuthorization(userId, "mi-app", "Juan", "juan@e.com", "web", projectId, profileTp, 8*time.Hour)

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
    profileTp := claim.ProfileTp(r)
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
r.Use(middleware.Cors)
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(middleware.RequestId)
r.Use(middleware.Telemetry)

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

// Con token efímero
router.EphemeralRoute(r, "GET", "/descarga/{id}", descargar, packageName, packagePath, host)

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
    response.HTTPError(w, r, http.StatusBadRequest, "msg") // {"ok": false, "result": {"message": "msg"}}
    response.HTTPAlert(w, r, "alerta")                     // 400
    response.Unauthorized(w, r)                            // 401
    response.Forbidden(w, r)                               // 403
    response.InternalServerError(w, r, err)                // 500

    // Streaming paginado (útil para exportaciones grandes)
    response.Stream(w, r, 100, func(page, rows int) (et.Items, error) {
        return modelo.Select().Page(page, rows).List()
    })
}
```

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

| Variable                    | Paquete       | Default     | Descripción                                          |
| --------------------------- | ------------- | ----------- | ---------------------------------------------------- |
| `DB_DRIVER`                 | jdb           | —           | `postgres`, `mysql` u `oracle`                       |
| `DB_HOST`                   | jdb           | —           | Host de la base de datos                             |
| `DB_PORT`                   | jdb           | `5432`      | Puerto de la base de datos                           |
| `DB_NAME`                   | jdb           | —           | Nombre de la base de datos                           |
| `DB_USER`                   | jdb           | —           | Usuario de la base de datos                          |
| `DB_PASSWORD`               | jdb           | —           | Contraseña de la base de datos                       |
| `DB_APPLICATION_NAME`       | jdb           | `elvis`     | Nombre de la aplicación en PostgreSQL                |
| `USE_CORE`                  | jdb           | `true`      | Inicializar tablas core (series, records, recycling) |
| `REDIS_HOST`                | cache         | —           | Host de Redis (ej. `localhost:6379`)                 |
| `REDIS_PASSWORD`            | cache         | —           | Contraseña de Redis                                  |
| `REDIS_DB`                  | cache         | `0`         | Número de base de datos Redis                        |
| `NATS_HOST`                 | event         | —           | URL de conexión NATS                                 |
| `NATS_USER`                 | event         | —           | Usuario NATS                                         |
| `NATS_PASSWORD`             | event         | —           | Contraseña NATS                                      |
| `SECRET`                    | claim         | `1977`      | Clave de firma JWT                                   |
| `HOST`                      | jrpc / router | `localhost` | Host del servicio actual                             |
| `PORT`                      | servicio      | `3400`      | Puerto HTTP                                          |
| `RPC_HOST`                  | jrpc          | `HOST`      | Host para RPC entre servicios                        |
| `RPC_PORT`                  | jrpc          | `4200`      | Puerto RPC                                           |
| `AUTHORIZATION_METHOD`      | router        | —           | Método RPC para verificar permisos                   |
| `RESILIENCE_TOTAL_ATTEMPTS` | resilience    | `3`         | Intentos totales por operación                       |
| `RESILIENCE_TIME_ATTEMPTS`  | resilience    | `30`        | Segundos entre reintentos                            |

---

## 🔧 Comandos de Desarrollo

```bash
# Compilar el módulo
go build ./...

# Ejecutar tests
go test ./...

# Ejecutar un test específico con verbose
go test ./paquete/... -run NombreTest -v

# Formatear código y ejecutar CLI de scaffolding
gofmt -w . && go run ./cmd/create go

# CLI para operaciones de base de datos
gofmt -w . && go run ./cmd/jdb go

# Actualizar dependencias
go mod tidy
```

---

## 📁 Estructura del Proyecto

```
elvis/
├── cache/          # Cliente Redis (Set, Get, Delete, Pub/Sub)
├── claim/          # JWT: generación, validación, invalidación
├── cmd/
│   ├── create/     # CLI scaffolding de microservicios
│   └── jdb/        # CLI de operaciones de base de datos
├── config/         # Carga de configuración
├── console/        # Logging interno de bajo nivel
├── create/
│   ├── v1/         # Generador de proyectos v1
│   └── v2/         # Generador de proyectos v2
├── crontab/        # Tareas programadas (cron)
├── dt/             # Objetos de transferencia de datos
├── envar/          # Helpers de variables de entorno
├── et/             # Tipos centrales: Json, Item, Items, List, Any
├── event/          # Eventos local (emitter) y distribuido (NATS)
├── file/           # Manejo de archivos
├── health/         # Health check helpers
├── jdb/            # Abstracción de base de datos (Postgres/MySQL/Oracle)
├── jrpc/           # RPC entre servicios vía Redis
├── linq/           # ORM / query builder
├── logs/           # Logging estructurado
├── mem/            # Cache in-memory con TTL
├── middleware/     # Middleware HTTP chi (auth, cors, logger, etc.)
├── msg/            # Mensajes de error compartidos
├── race/           # Helpers de concurrencia
├── reg/            # Registro de IDs
├── resilience/     # Reintentos automáticos
├── response/       # Helpers de respuesta HTTP
├── router/         # Registro de rutas chi con API Gateway
├── service/        # Cliente HTTP entre servicios
├── stdrout/        # Salida estándar / terminal
├── strs/           # Utilidades de strings
├── timezone/       # Manejo de zonas horarias
├── utility/        # Utilidades generales (UUID, OTP, crypto, etc.)
└── workflow/       # Orquestación de flujos multi-paso
```

---

## 🔔 Eventos del Sistema

Eventos internos que emite la librería y a los que se puede suscribir:

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
