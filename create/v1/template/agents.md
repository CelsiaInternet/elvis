Proyecto `$1`: microservicio en Go construido sobre [elvis](https://github.com/celsiainternet/elvis). Este archivo guía a los agentes (Claude Code, Codex y compatibles) que ayudan a crear y hacer crecer el proyecto.

## Reglas

1. **El código debe ser entendible, mantenible y escalable por un humano.** Sigue los patrones que ya existen en el servicio y usa lo que resuelve elvis antes de agregar piezas nuevas. Si una decisión no es obvia, explícala en la respuesta, no con código ingenioso.
2. **Los commits los hace el usuario.** No ejecutes `git commit` ni `git push`. Deja los cambios en el árbol de trabajo y resume qué cambió y cómo se verificó, para que el usuario los revise y los commitee.

## Comandos

- Verificar un cambio: `gofmt -l . && go build ./... && go vet ./...`
- Ejecutar: `go run ./cmd/$1` (lee `.env`; necesita PostgreSQL, Redis y NATS).
- Agregar microservicios, modelos o RPC: `go run github.com/celsiainternet/elvis/cmd/create go`. La opción **Modelo** registra las rutas y ejecuta `go mod tidy` sola.
- Actualizar las guías de elvis después de subir de versión: `go run github.com/celsiainternet/elvis/cmd/install`.
