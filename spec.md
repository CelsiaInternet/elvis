# spec — elvis

## Objetivo

Librería Go (`github.com/celsiainternet/elvis`) con las primitivas de infraestructura de los microservicios de Celsia Internet: tipos comunes, base de datos y ORM, caché, eventos, autenticación, router HTTP, RPC entre servicios y envíos, más un generador de proyectos y un marco de trabajo para agentes. La consumen casi todos los servicios Go del workspace; dentro del `go.work` compilan contra esta copia y fuera de él contra el último tag (`v1.1.325`).

## Arquitectura y estado por capa

| Capa | Paquetes | Estado |
|------|----------|--------|
| Tipos y utilidades | `et`, `envar`, `strs`, `utility`, `msg`, `timezone`, `file` | Compila; sin pruebas propias; los usan todos los servicios. |
| Datos | `jdb`, `linq`, `jquery` | Compila; `linq/jsql_test.go` pasa; `jquery/test` tiene una prueba que falla; el resto sin pruebas. |
| Caché y eventos | `cache`, `mem`, `event`, `console`, `logs` | Compila; `event/error_test.go` pasa; sin pruebas contra Redis o NATS. |
| HTTP y auth | `router`, `response`, `middleware`, `claim`, `request` | Compila; `request/request_test.go` cubre el `Content-Type` por defecto. |
| RPC y envíos | `jrpc`, `service` | Compila; `service/send_test.go` pasa con un servidor net/rpc de prueba; el envío por RPC (F-37) no está en ningún tag. |
| Orquestación | `resilience`, `workflow`, `instances`, `crontab`, `queue` | Compila; solo `queue` tiene un ejemplo. |
| Generador y agentes | `create/v1`, `cmd/create`, `cmd/install`, `agentsguide` | `create/v1/model_test.go` pasa (plantillas válidas y con `gofmt`). |
| Publicación | `version.sh` | `bash -n` pasa; no se ejecutó. |

## Requisitos

| ID | Descripción | Peso | Estado |
|----|-------------|------|--------|
| B1 | Tipos `et` (`Json`, `Item`, `Items`, `List`, `Any`) con accesores tipados | 2 | 0.5 |
| B2 | `jdb`: conexión a PostgreSQL, MySQL y Oracle desde variables de entorno y tablas core (series, records, recycling) | 3 | 0.5 |
| B3 | `linq`: modelos, atributos JSONB, referencias, triggers y CRUD fluido | 3 | 0.5 |
| B4 | `linq.RunQuery`: contrato jsql validado (wheres, selects, order_by, search, paginación) y publicado desde `v1.1.324` | 3 | 1 |
| B5 | `jquery`: `et.Json` a `SELECT` con dialecto intercambiable | 2 | 0.5 |
| C1 | `cache` (Redis, pub/sub) y `mem` (memoria con TTL) | 2 | 0.5 |
| C2 | `event`: eventos locales y por NATS (`Subscribe`, `Queue`, `Stack` con re-registro) y `Work`/`WorkState` | 3 | 0.5 |
| C3 | `event.Error`: reporte validado al registro central de errores | 1 | 1 |
| C4 | `jrpc`: RPC TCP con registro en Redis, llamadas tipadas y `PIPE_HOST` | 3 | 0.5 |
| C5 | `service`: envíos y OTP por RPC con respaldo a `event.Work` (F-37) | 2 | 0.5 |
| H1 | `router`: rutas públicas, protegidas y con permisos publicadas en el API Gateway | 2 | 0.5 |
| H2 | `response`: helpers de respuesta y `Stream` | 1 | 0.5 |
| H3 | `claim` + `middleware`: JWT en Redis, autenticación, autorización, efímeros, CORS, logger, request id, recover | 3 | 0.5 |
| H4 | `request`: verbos, TLS, timeout, JSON, formulario y `multipart/form-data` | 2 | 0.5 |
| H5 | `request`: el timeout cancela la petición en curso (`context.WithTimeout`) | 1 | 0 |
| S1 | `secret`: AES-256-GCM con formato `v1:` compatible con catalog y send | 2 | 1 |
| S2 | `linq`: las condiciones `Column.Eq/In/Like/...` escapan la comilla simple | 2 | 0 |
| G1 | `cmd/create` + `create/v1`: proyecto, microservicio, modelo y RPC con plantillas que compilan contra el último tag | 3 | 1 |
| G2 | `cmd/install` y `agentsguide`: prepara el proyecto e instala el marco para agentes sin fijar versiones | 2 | 0.5 |
| O1 | `resilience` y `workflow`: reintentos y flujos multi-paso con rollback | 2 | 0.5 |
| O2 | `crontab`, `xls`, `queue`, `jtls` | 1 | 0.5 |
| Q1 | `go build ./...` y `go vet ./...` pasan | 1 | 1 |
| Q2 | `go test ./...` pasa completo | 2 | 0 |
| Q3 | Pruebas de `jdb`, `linq` (CRUD), `cache`, `event` NATS y `jrpc` contra servicios reales | 3 | 0 |
| D1 | README en inglés y en español con la misma estructura, instalación primero | 1 | 1 |
| D2 | `CLAUDE.md` alineado con el código | 1 | 1 |

Notas de estado:

- **B4**: `linq/jsql_test.go` (SQL, sin base) pasa; `linq/jsql.go` está en los tags `v1.1.324` y `v1.1.325`.
- **C3, S1, G1**: cubiertos por `event/error_test.go`, `secret/secret_test.go` y `create/v1/model_test.go`, que pasan.
- **C5**: `service/send_test.go` pasa con dobles y un servidor net/rpc local; no se llamó a send ni a otp reales. `service/send-rpc.go` y el cambio de `jrpc/handler.go` están sin commit y sin tag, así que los servicios fuera del `go.work` siguen enviando por `event.Work`.
- **H4**: `request/request_test.go` cubre el `Content-Type` por defecto (JSON y formulario); `multipart/form-data` no tiene prueba.
- **H5**: el timeout solo abandona el resultado; la petición sigue viva hasta el timeout del cliente (120 min).
- **S2**: el camino normal pasa los valores por `et.Unquote`, que no escapa `'`; `RunQuery` sí escapa.
- **Q2**: falla `TestJQuery_MultipleColumnsAreAnded` (`jquery/test/jquery_test.go:84`: `'and' debe ser un arreglo de objetos condicion`).
- **B1–B3, B5, C1, C2, C4, H1–H3, G2, O1, O2**: compilan y los usan los servicios en producción, pero no tienen pruebas automáticas ni se probaron en esta revisión.

## Brecha

`brecha = (total − avance) / total × 100 %`

| Bloque | Peso | Avance | Brecha (pts) |
|---|---|---|---|
| B | 13 | 8 | 5 |
| C | 11 | 6 | 5 |
| H | 9 | 4 | 5 |
| S | 4 | 2 | 2 |
| G | 5 | 4 | 1 |
| O | 3 | 1,5 | 1,5 |
| Q | 6 | 1 | 5 |
| D | 2 | 2 | 0 |
| **Total** | **53** | **28,5** | **24,5** |

Brecha actual: 46,2 % (24,5 de 53 pts)

## Orden sugerido

1. Q2: corregir `jquery` o su prueba para que `go test ./...` pase.
2. C5: commit y tag de elvis con el envío por RPC y subir los pins de los servicios (F-37, F-26).
3. S2: escapar la comilla simple en las condiciones de `linq`.
4. H5: cancelar la petición de `request` al vencer el timeout.
5. Q3: pruebas de integración de `jdb`/`linq`, `cache`, `event` y `jrpc` contra Postgres, Redis y NATS locales.

## Verificación (2026-11-18)

- `go build ./...` y `go vet ./...` pasan.
- `go test ./...`: pasan `create/v1`, `event`, `linq`, `queue`, `request`, `secret` y `service`; falla `jquery/test` (Q2).
- `bash -n version.sh` pasa. No se ejecutó `version.sh`.
- No se probó nada contra Postgres, Redis ni NATS.

## Historial

| Fecha | Cambio | Brecha |
|-------|--------|--------|
| 2026-11-18 | Spec inicial (F-39): 26 requisitos en 8 bloques, estado según build, vet y pruebas | 46,2 % |
