<!-- elvis:agents-guide:start -->
## Marco de trabajo: proyectos backend sobre Elvis

> Esta sección se agregó automáticamente al instalar `github.com/celsiainternet/elvis`
> (`go run github.com/celsiainternet/elvis/cmd/install`, o al generar el proyecto con
> `cmd/create`). Es una guía operativa para agentes (Claude Code u otros) que **crean**
> o **refactorizan** microservicios backend sobre esta librería. No la resumas ni la
> borres: futuras sesiones de agentes vuelven a leerla completa. Si el contenido de una
> nueva versión de Elvis cambia, vuelve a correr `cmd/install` para regenerarla — no la
> edites a mano salvo para añadir notas propias del proyecto fuera de las marcas
> `elvis:agents-guide:start` / `elvis:agents-guide:end`.

Aplica a cualquier servicio que importe `github.com/celsiainternet/elvis` en su
`go.mod`: APIs HTTP, workers de eventos, servicios RPC. No aplica a scripts sueltos.

### 1. Crear un proyecto nuevo

Preferir el generador interactivo a escribir el andamiaje a mano:

```bash
go run github.com/celsiainternet/elvis/cmd/create go
```

Opciones del menú: **Project** (microservicio completo: `cmd/`, `internal/`, `pkg/`,
`deployments/`, `scripts/`, `README.md`, `.env`, `.gitignore`) · **Microservice**
(agrega un servicio dentro de un repo existente) · **Modelo** (agrega un modelo `linq`
+ handler CRUD a un `pkg/<nombre>` ya creado) · **Rpc** (agrega un stub RPC a un
`pkg/<nombre>` ya creado).

**No uses `go run github.com/celsiainternet/elvis/cmd/jdb go`** a menos que el
repositorio ya use su layout. Pese al nombre, `cmd/jdb` no es una herramienta de base
de datos: ejecuta el mismo asistente que `cmd/create` pero contra un generador
alternativo (`create/v2`) cuya estructura de carpetas (`internal/models/<n>/`,
`internal/services/<n>/`, `router-<modelo>.go`) es **incompatible** con la de
`cmd/create` (`internal/service/<n>/`). Mezclar ambos layouts en el mismo repo produce
dos árboles `internal/` inconsistentes. Regla: si el repo ya tiene `internal/models/`,
sigue usando `cmd/jdb`; si no, usa siempre `cmd/create`. Ante la duda, un `find internal
-maxdepth 2 -type d` te dice cuál layout ya está en uso antes de generar nada.

Además, el código que genera `create/v2` (`cmd/jdb`) importa la librería independiente
`github.com/celsiainternet/jdb/jdb` en lugar de `github.com/celsiainternet/elvis/jdb`
(que es lo que usa `cmd/create`). Si el repo usa el layout v2, asegúrate de que su
`go.mod` requiera `github.com/celsiainternet/jdb`: `cmd/install` no lo agrega.

Después de generar (o al clonar un repo existente que use Elvis), corre:

```bash
go run github.com/celsiainternet/elvis/cmd/install
```

para fijar las versiones de las dependencias de terceros que el scaffolding importa
(chi, redis, nats, jwt, cobra, promptui...). Esa lista se mantiene a mano dentro de
Elvis y puede quedar desactualizada frente al `go.mod` real de la librería; si
`go build ./...` falla por una versión incompatible tras correrlo, resuelve el módulo
puntual con `go get <módulo>@latest` en vez de desconfiar del resto del scaffolding.

### 2. Layout esperado (generador `cmd/create` / `create/v1`)

```
cmd/<servicio>/              main.go + Dockerfile
internal/service/<servicio>/ service.go (arma chi.Router, cache/db/event.Load)
internal/service/<servicio>/v1/ api.go (monta pkg.Router en pkg.PackagePath)
pkg/<servicio>/               controller.go, router.go, event.go, msg.go, config.go
                               (+ model.go, schema.go, h<Modelo>.go, rpc.go si hay schema)
deployments/<servicio>/       local.yml (compose + labels de Traefik)
                               oke-template.yml (Service + Deployment k8s)
                               oke-statefulset-template.yml (Service + StatefulSet k8s)
scripts/<servicio>.http       peticiones de ejemplo
```

Un `pkg/<nombre>` generado **con** schema (BD) trae CRUD completo (`Insert`,
`UpSert`, `State`, `Delete`, `All`) más rutas REST estándar; generado **sin** schema
trae un stub de controller/handler/router vacío para lógica no persistida.

Los manifiestos `oke-template.yml`/`oke-statefulset-template.yml` usan placeholders
literales (`$ROLE`, `$NS`, `$PORT`, `$IMAGE`, `$REPLICAS`, etc.) que el generador **no**
sustituye — a diferencia de `local.yml`, que usa `$1`/`$2`/`$3` reemplazados por
`file.MakeFile` al momento de generarlo. Esos `$NOMBRE` quedan tal cual para que el
pipeline de CI/CD (o un `envsubst`/`kubectl` con `--dry-run` y variables de entorno) los
resuelva en tiempo de despliegue; no los confundas con los `$1`/`$2` posicionales ni
intentes rellenarlos a mano en el generador.

### 3. API verificada — no asumas firmas de memoria

Estas son las firmas reales verificadas contra el código fuente de Elvis; la
documentación generada (README/godoc) de versiones antiguas ha tenido errores en
varios de estos puntos, así que trátalos como la referencia autoritativa:

**Modelos `linq` — comandos de escritura vs. consultas de lectura.** Son dos familias
de métodos distintas, no intercambiables:

- Escritura — `model.Insert(data)`, `.Update(data)`, `.Upsert(data)`, `.Delete()` →
  encadenan `.Command()` (`et.Items`) o `.CommandOne()` (`et.Item`, primer resultado).
- Lectura — `model.Select(...)` / `model.Data(...)` → encadenan `.First()` (un
  registro, `et.Item`), `.All()` (`et.Items`), `.Page(page, rows)` (`et.Items`
  paginado, sin conteo total) o `.List(page, rows)` (`et.List` con `Count`/`All`
  además del `Result` de la página).

  **No existe `.One()` en ningún caso**, y `.Page(...)` ya ejecuta la consulta y
  devuelve `(et.Items, error)` — no se le puede encadenar `.List()` encima.

**Middleware HTTP (`middleware/`):**

- CORS: `middleware.AllowAll(origenes []string).Handler` (envuelve `rs/cors`;
  `nil`/slice vacío = permitir cualquier origen). **No existe** `middleware.Cors`.
- Request ID: `middleware.RequestID` (así, con `ID` en mayúsculas). **No existe**
  `middleware.RequestId`.
- Telemetría: no hay un middleware `r.Use(middleware.Telemetry)`. La telemetría se
  expone como helpers (`middleware.NewMetric`, `middleware.PushTelemetry`, etc.), no
  como middleware de chi.
- Sí existen tal cual: `middleware.Autentication`, `middleware.Authorization`,
  `middleware.Ephemeral`, `middleware.Logger`, `middleware.Recoverer`.

**JWT (`claim/`):** el campo de perfil del claim es `claim.Claim.ProfileId` /
`claim.ProfileId(r)`. **No existe** `ProfileTp`.

**Routing (`router/`):** usa siempre `router.PublicRoute`, `ProtectRoute`,
`AuthorizationRoute`, `EphemeralRoute` o `With` en vez de registrar rutas
directamente en `chi.Mux` — estas funciones publican la ruta al API Gateway vía NATS
(`apigateway/set/resolve`) automáticamente; registrar con `r.Get`/`r.Post` a mano
salta ese registro y la ruta queda invisible para el gateway.

### 4. Checklist para refactorizar un backend existente hacia Elvis

- [ ] Confirmar el driver de BD: solo `postgres`, `mysql` u `oracle` (`jdb.Postgres`,
      `jdb.Mysql`, `jdb.Oracle`) — no hay soporte de conexión a SQLite ni SQL Server.
- [ ] Sustituir SQL manual por `linq.Model` cuando la tabla tenga CRUD estándar;
      conservar `db.Query`/`db.Command`/`db.QueryOne` (jdb crudo) solo para consultas
      que `linq` no modele bien.
- [ ] Sustituir respuestas HTTP ad-hoc por `response.JSON/ITEM/ITEMS/HTTPError/
      HTTPAlert/Unauthorized/Forbidden/Stream`.
- [ ] Sustituir routing manual de chi por `router.PublicRoute/ProtectRoute/
      AuthorizationRoute/EphemeralRoute` (ver punto 3).
- [ ] Sustituir JWT casero por `claim.NewToken`/`claim.ValidToken` +
      `middleware.Autentication` (+ `middleware.Authorization` si hay permisos).
- [ ] Si el proyecto emite/consume eventos, migrar a `event.Publish` +
      `event.Subscribe`/`event.Queue`/`event.Stack` (usa `Stack` solo para
      suscripciones de reset/resincronización, se re-registra solo tras reconectar).
- [ ] Revisar nombres de variables de entorno contra la tabla del punto 5 antes de
      inventar una nueva — puede que ya exista un equivalente en Elvis.
- [ ] Correr `go build ./...` y `go test ./...` después de cada paso del refactor,
      no solo al final.

### 5. Variables de entorno mínimas

| Variable | Paquete | Default |
|---|---|---|
| `DB_DRIVER`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` | jdb | — |
| `DB_APPLICATION_NAME` | jdb | `elvis` |
| `USE_CORE` | jdb | `true` (inicializa tablas series/records/recycling) |
| `REDIS_HOST`, `REDIS_PASSWORD`, `REDIS_DB` | cache | — |
| `NATS_HOST`, `NATS_USER`, `NATS_PASSWORD` | event | — |
| `SECRET` | claim | `1977` (clave de firma JWT) |
| `HOST`, `RPC_HOST`, `RPC_PORT` | jrpc/router | `localhost`, `4200` |
| `AUTHORIZATION_METHOD` | router/middleware | — |
| `RESILIENCE_TOTAL_ATTEMPTS`, `RESILIENCE_TIME_ATTEMPTS` | resilience | `3`, `30`s |

### 6. Antes de confiar en algo que no está aquí

Esta guía cubre solo los puntos donde ya se detectaron inconsistencias entre la
documentación y el código. Para todo lo demás (workflows, resiliencia, eventos de
sistema, `jquery`, `xls`, `authorization`, `inbox`...), o si trabajas contra una
versión de Elvis distinta a la que generó este archivo, verifica el símbolo real en
el módulo instalado antes de usarlo:

```bash
grep -rn "func <Nombre>" "$(go env GOMODCACHE)/github.com/celsiainternet/elvis@$(go list -m -f '{{.Version}}' github.com/celsiainternet/elvis)"
```

o revisa `CLAUDE.md`/`README.md` en el repo fuente (`github.com/celsiainternet/elvis`)
correspondientes a la versión fijada en tu `go.mod`.
<!-- elvis:agents-guide:end -->
