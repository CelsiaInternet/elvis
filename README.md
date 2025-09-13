# 🎸 Elvis - Framework para Microservicios en Go

[![Go Version](https://img.shields.io/github/go-mod/go-version/celsiainternet/elvis?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/github/license/celsiainternet/elvis?style=flat-square)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/celsiainternet/elvis?style=flat-square&logo=github)](https://github.com/celsiainternet/elvis/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/celsiainternet/elvis?style=flat-square)](https://goreportcard.com/report/github.com/celsiainternet/elvis)
[![GitHub Stars](https://img.shields.io/github/stars/celsiainternet/elvis?style=flat-square&logo=github)](https://github.com/celsiainternet/elvis/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/celsiainternet/elvis?style=flat-square&logo=github)](https://github.com/celsiainternet/elvis/issues)
[![Documentation](https://img.shields.io/badge/docs-godoc-blue?style=flat-square&logo=go)](https://pkg.go.dev/github.com/celsiainternet/elvis)

> 🚀 **Framework moderno y robusto para el desarrollo de microservicios escalables en Go**

<div align="center">

![Elvis Logo](https://via.placeholder.com/200x100/1e40af/ffffff?text=🎸+ELVIS)

**[📚 Documentación](https://pkg.go.dev/github.com/celsiainternet/elvis)** •
**[🚀 Quick Start](#-quick-start)** •
**[📖 Ejemplos](https://github.com/celsiainternet/elvis/tree/main/examples)** •
**[🐛 Issues](https://github.com/celsiainternet/elvis/issues)** •
**[💬 Discusiones](https://github.com/celsiainternet/elvis/discussions)**

</div>

## 📑 Tabla de Contenidos

- [📖 Descripción](#-descripción)
- [Requisitos Previos](#requisitos-previos)
- [🚀 Instalación Rápida](#-instalación-rápida)
- [⚙️ Configuración de Desarrollo](#️-configuración-de-desarrollo)
- [🚀 Quick Start](#-quick-start)
- [Características Principales](#características-principales)
- [📁 Estructura del Proyecto](#-estructura-del-proyecto)
- [🔧 Comandos Disponibles](#-comandos-disponibles)
- [Configuración de Variables de Entorno](#configuración-de-variables-de-entorno)
- [💡 FAQ y Mejores Prácticas](#-faq-y-mejores-prácticas)
- [🤝 Contribución](#-contribución)
- [📄 Licencia](#-licencia)
- [Releases](#releases)
- [Versionamiento](#versionamiento)

## 📖 Descripción

Elvis es un framework moderno y robusto diseñado para facilitar el desarrollo de microservicios en Go. Proporciona un conjunto completo de herramientas y estructuras para crear aplicaciones escalables con capacidades avanzadas de:

- 🔄 **Comunicación en tiempo real** (WebSocket)
- 🛡️ **Sistema de resiliencia** y recuperación automática
- 📅 **Tareas programadas** (Crontab)
- 🗄️ **Base de datos avanzada** con triggers automáticos
- 🔄 **Sistema de eventos** distribuidos
- 💾 **Cache inteligente** multi-backend
- 🔐 **Middleware de seguridad** integrado
- 📊 **Telemetría y monitoreo** en tiempo real

## Requisitos Previos

- Go 1.23 o superior
- Git
- PostgreSQL (para base de datos)
- Redis (para cache y comunicación en tiempo real)
- NATS (para eventos)

## 🚀 Instalación Rápida

### 1. Inicializar el Proyecto

```bash
go mod init github.com/tu-usuario/tu-proyecto
```

### 2. Instalar Elvis

```bash
go get github.com/celsiainternet/elvis@v1.1.131
```

### 3. Crear Proyecto con Elvis

```bash
go run github.com/celsiainternet/elvis/cmd/create go
```

### 4. Configurar Variables de Entorno

Copia el archivo `.env.example` a `.env` y ajusta los valores según tu entorno:

```bash
cp .env.example .env
```

## ⚙️ Configuración de Desarrollo

### IDE Configuration (Cursor/VSCode)

El proyecto incluye configuración optimizada para Cursor y VSCode:

- **`.vscode/settings.json`**: Configuración del workspace con staticcheck deshabilitado
- **`staticcheck.conf`**: Configuración específica de staticcheck para evitar warnings molestos
- **Linting**: ST1020 y otras reglas de documentación están deshabilitadas para mayor comodidad

### Ejecutar en Modo Desarrollo

```bash
# Creacion de elementos
gofmt -w . && go run ./cmd/create go

# Creacion de elementos de jdb
gofmt -w . && go run ./cmd/jdb go

# Cliente WorkFlow
gofmt -w . && go run ./cmd/flow go
```

## Uso

### Creación del Proyecto

Para crear un nuevo proyecto con Elvis, ejecuta el siguiente comando:

```bash
go run github.com/celsiainternet/elvis/cmd/create go
```

Este comando generará:

- Estructura base del proyecto
- Microservicios iniciales
- Modelos de datos
- Configuraciones necesarias

### Ejecutar el Proyecto

Para ejecutar el proyecto:

```bash
gofmt -w . && go run ./cmd/test -port 3400 -rpc 4400
gofmt -w . && go run ./cmd/resilence
```

Donde:

- `-port`: Puerto para el servidor HTTP (default: 3400)
- `-rpc`: Puerto para el servidor gRPC (default: 4400)

## Características Principales

### 🔄 Comunicación en Tiempo Real (WebSocket)

Elvis incluye un sistema completo de WebSocket para comunicación en tiempo real:

```go
// Servidor WebSocket
hub := ws.ServerHttp(3300, "username", "password")

// Cliente WebSocket
client, err := ws.Login(&ws.ClientConfig{
    ClientId:  "client-1",
    Name:      "TestClient",
    Url:       "ws://localhost:3300/ws",
    Reconnect: 3,
})

// Suscribirse a canales
client.Subscribe("notifications", func(msg ws.Message) {
    fmt.Println("Mensaje recibido:", msg.Data)
})

// Publicar mensajes
client.Publish("notifications", map[string]interface{}{
    "message": "Hola mundo",
})
```

### 🛡️ Sistema de Resiliencia

Manejo robusto de errores y recuperación automática:

```go
// Configurar resiliencia
resilience.SetNotifyType(resilience.TpNotifyEmail)
resilience.SetContactNumbers([]string{"+573160479724"})

// Agregar transacción con reintentos automáticos
transaction := resilience.Add("email-send", "Enviar email de confirmación", sendEmail, userEmail, content)
```

### 📅 Tareas Programadas (Crontab)

Sistema de tareas programadas integrado:

```go
// Crear tarea programada
jobs := crontab.New()
jobs.AddJob("backup-daily", "Backup diario", "0 2 * * *", "backup-channel", map[string]interface{}{
    "type": "daily",
    "path": "/backup",
})

// Iniciar tareas
jobs.Start()
```

### 🗄️ Base de Datos Avanzada

Sistema de base de datos con triggers automáticos y sincronización:

```go
// Definir modelo con triggers
model := linq.NewModel(db, "users", "Usuarios", 1)
model.DefineColum("_id", "", "VARCHAR(80)", "-1")
model.DefineColum("name", "", "VARCHAR(250)", "")
model.DefineColum("email", "", "VARCHAR(250)", "")

// Configurar triggers
model.Trigger(linq.BeforeInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
    // Lógica antes de insertar
    return nil
})

model.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
    // Lógica después de insertar
    return nil
})
```

### 🔄 Sistema de Eventos

Gestión de eventos distribuidos:

```go
// Publicar evento
event.Publish("user.created", map[string]interface{}{
    "user_id": "123",
    "email":   "user@example.com",
})

// Suscribirse a eventos
event.Subscribe("user.created", func(msg event.EvenMessage) {
    fmt.Println("Usuario creado:", msg.Data)
})

// Trabajos distribuidos
work := event.Work("email.send", map[string]interface{}{
    "to":      "user@example.com",
    "subject": "Bienvenido",
})
```

### 💾 Cache Inteligente

Sistema de cache con múltiples backends:

```go
// Configurar cache
cache.Load()

// Operaciones de cache
cache.Set("key", "value", 3600)
value, err := cache.Get("key")
cache.Delete("key")

// Cache hash
cache.SetH("user:123", map[string]interface{}{
    "name":  "Juan",
    "email": "juan@example.com",
})
```

### 🔐 Middleware de Seguridad

Middleware integrado para autenticación y autorización:

```go
// Middleware de autenticación
r.Use(middleware.Authentication)

// Middleware de autorización
r.Use(middleware.Authorization)

// Middleware de CORS
r.Use(middleware.CORS)

// Middleware de logging
r.Use(middleware.Logger)
```

### 📊 Telemetría y Monitoreo

Sistema de telemetría integrado:

```go
// Enviar telemetría
realtime.Telemetry(map[string]interface{}{
    "service": "user-service",
    "method":  "POST",
    "duration": 150,
    "status":  "success",
})

// Logging estructurado
logs.Log("user-service", "Usuario creado exitosamente")
logs.Alert(errors.New("Error de conexión"))
```

## 📁 Estructura del Proyecto

```
elvis/
├── 📂 cmd/                    # Comandos ejecutables
│   ├── cmd/                   # Servidor principal
│   ├── create-go/            # Generador de proyectos
│   ├── daemon/               # Daemon del sistema
│   ├── resilence/            # Servidor de resiliencia
│   ├── rpc/                  # Servidores RPC (cliente/servidor)
│   └── ws/                   # Servidor WebSocket
├── 📂 cache/                 # Sistema de cache
├── 📂 claim/                 # Manejo de claims/tokens
├── 📂 config/                # Configuración global
├── 📂 console/               # Utilidades de consola
├── 📂 create/                # Generación de código
│   └── template/             # Plantillas de código
├── 📂 crontab/               # Tareas programadas
├── 📂 envar/                 # Variables de entorno
├── 📂 et/                    # Tipos y utilidades
├── 📂 event/                 # Sistema de eventos
├── 📂 file/                  # Manejo de archivos
├── 📂 jdb/                   # Database abstraction layer
├── 📂 jrpc/                  # JSON-RPC implementation
├── 📂 linq/                  # Query builder
├── 📂 logs/                  # Sistema de logging
├── 📂 mem/                   # Cache en memoria
├── 📂 middleware/            # Middleware HTTP
├── 📂 msg/                   # Mensajería
├── 📂 race/                  # Control de concurrencia
├── 📂 realtime/              # Comunicación en tiempo real
├── 📂 resilience/            # Sistema de resiliencia
├── 📂 response/              # Manejo de respuestas HTTP
├── 📂 router/                # Enrutamiento HTTP
├── 📂 service/               # Servicios base
├── 📂 stdrout/               # Salida estándar
├── 📂 strs/                  # Utilidades de strings
├── 📂 timezone/              # Manejo de zonas horarias
├── 📂 utility/               # Utilidades generales
├── 📂 ws/                    # WebSocket implementation
├── 📂 .vscode/               # Configuración IDE
│   ├── settings.json         # Configuración optimizada
│   └── launch.json           # Configuración debug
├── 📄 staticcheck.conf       # Configuración linting
├── 📄 go.mod                 # Dependencias Go
├── 📄 go.sum                 # Checksums dependencias
├── 📄 version.sh             # Script versionamiento
└── 📄 README.md              # Documentación
```

## 🚀 Quick Start

### Ejemplo Básico

```go
package main

import (
    "github.com/celsiainternet/elvis/router"
    "github.com/celsiainternet/elvis/middleware"
    "github.com/celsiainternet/elvis/response"
)

func main() {
    // Crear router
    r := router.New()

    // Agregar middleware
    r.Use(middleware.CORS)
    r.Use(middleware.Logger)

    // Definir rutas
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        response.JSON(w, r, 200, map[string]string{
            "status": "ok",
            "message": "Elvis está funcionando!",
        })
    })

    // Iniciar servidor
    r.Listen(":3400")
}
```

## 🔧 Comandos Disponibles

### Desarrollo Local

```bash
# Servidor principal con hot reload
gofmt -w . && go run ./cmd/test -port 3400 -rpc 4400

# Generar nuevo proyecto
go run github.com/celsiainternet/elvis/cmd/create go
```

### Servicios Adicionales

```bash
# Creaceon de elementos
gofmt -w . && go run ./cmd/create go

# Creacion de elementos de jdb
gofmt -w . && go run ./cmd/jdb go

# Cliente WorkFlow
gofmt -w . && go run ./cmd/flow

```

### Herramientas de Sistema

```bash
# Daemon del sistema
go run ./cmd/daemon --status
go run ./cmd/daemon --restart
go run ./cmd/daemon --stop

# Verificar versión
./version.sh
```

## Configuración de Variables de Entorno

```bash
# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=password

# Redis
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# NATS
NATS_URL=nats://localhost:4222

# WebSocket
WS_USERNAME=admin
WS_PASSWORD=secret
RT_URL=ws://localhost:3300/ws

# Resiliencia
RESILIENCE_ATTEMPTS=3
RESILIENCE_TIME_ATTEMPTS=30
```

## 💡 FAQ y Mejores Prácticas

### ¿Cómo deshabilitar warnings de linting?

El proyecto ya incluye configuración para deshabilitar warnings molestos como ST1020:

- Verifica que tu IDE esté usando la configuración en `.vscode/settings.json`
- El archivo `staticcheck.conf` excluye las reglas problemáticas
- Reinicia tu IDE después de clonar el proyecto

### ¿Cómo agregar un nuevo microservicio?

```bash
# Usar el generador incluido
go run github.com/celsiainternet/elvis/cmd/create-go create

# Seguir las convenciones de nombres
# - Servicios en cmd/nombre-servicio/
# - Modelos en internal/models/
# - Lógica de negocio en pkg/
```

### ¿Cómo manejar bases de datos?

```go
// Usar el sistema linq incluido
model := linq.NewModel(db, "table_name", "Display Name", 1)
model.DefineColum("id", "", "VARCHAR(80)", "-1")

// Los triggers se configuran automáticamente
model.Trigger(linq.BeforeInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
    // Tu lógica aquí
    return nil
})
```

### ¿Cómo configurar WebSocket?

```go
// Servidor
hub := ws.ServerHttp(3300, "username", "password")

// Cliente
client, err := ws.Login(&ws.ClientConfig{
    ClientId: "unique-id",
    Url:      "ws://localhost:3300/ws",
})
```

## 🤝 Contribución

¡Las contribuciones son bienvenidas! Para contribuir:

1. **Fork** el proyecto
2. **Crea** una rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. **Commit** tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. **Push** a la rama (`git push origin feature/nueva-funcionalidad`)
5. **Abre** un Pull Request

### Guías de Contribución

- Sigue las convenciones de Go (gofmt, golint)
- Agrega tests para nuevas funcionalidades
- Actualiza la documentación
- Usa conventional commits
- Asegúrate de que todos los tests pasen

### Reportar Bugs

Usa los [GitHub Issues](https://github.com/celsiainternet/elvis/issues) para reportar bugs:

- Describe el problema claramente
- Incluye pasos para reproducir
- Especifica tu versión de Go y OS
- Adjunta logs si es posible

## 📄 Licencia

Este proyecto está bajo la **Licencia MIT**. Ver el archivo [LICENSE](LICENSE) para más detalles.

### Resumen de la Licencia

- ✅ **Uso comercial** permitido
- ✅ **Modificación** permitida
- ✅ **Distribución** permitida
- ✅ **Uso privado** permitido
- ❌ **Sin garantía**
- ❌ **Sin responsabilidad**

---

## 👨‍💻 Autor

**César Galvis León**

- 📧 Email: [cesar@celsiainternet.com](mailto:cesar@celsiainternet.com)
- 🌐 Website: [celsiainternet.com](https://celsiainternet.com)
- 💼 LinkedIn: [César Galvis León](https://linkedin.com/in/cesargalvisleon)

---

**⭐ Si te gusta Elvis, ¡no olvides darle una estrella al repositorio!**

---

_Desarrollado con ❤️ en Colombia_

## Releases

### v1.1.131

- **Sistema de Comunicación en Tiempo Real**

  - Implementación completa de WebSocket con hub centralizado
  - Soporte para canales y colas de mensajes
  - Cliente WebSocket con reconexión automática
  - Adaptadores para Redis y WebSocket distribuido
  - Sistema de suscripciones y publicaciones

- **Sistema de Resiliencia Avanzado**

  - Manejo de transacciones con reintentos automáticos
  - Notificaciones por SMS, Email y WhatsApp
  - Persistencia en cache y memoria
  - Configuración de intentos y tiempos de espera
  - Monitoreo de estado de transacciones

- **Tareas Programadas (Crontab)**

  - Nuevo generador de tareas programadas
  - Soporte para expresiones cron estándar
  - Integración con el sistema de eventos
  - Manejo de errores en tareas programadas
  - Persistencia de configuración

- **Mejoras en Base de Datos**

  - Triggers automáticos para sincronización
  - Sistema de series automáticas
  - Reciclaje de registros eliminados
  - Notificaciones PostgreSQL nativas
  - Funciones SQL optimizadas

- **Sistema de Eventos Distribuidos**

  - Publicación y suscripción de eventos
  - Trabajos distribuidos con estados
  - Colas de mensajes con balanceo de carga
  - Integración con NATS
  - Telemetría y logging automático

- **Mejoras en la Generación de Código**

  - Nuevas plantillas para WebSocket
  - Plantillas para tareas programadas
  - Optimización en la generación de modelos
  - Mejora en la documentación generada
  - Soporte para Docker multi-stage

- **Correcciones y Optimizaciones**
  - Ajustes en el manejo de memoria
  - Mejora en la concurrencia
  - Actualización de dependencias
  - Corrección de bugs en WebSocket
  - Optimización de rendimiento

### v1.1.2

- Mejoras en la generación de microservicios
  - Optimización del rendimiento en la creación de modelos
  - Corrección de bugs en la inicialización del proyecto
  - Mejora en la gestión de dependencias
  - Actualización de la documentación

### v1.1.1

- Agregado soporte para configuración de puertos personalizados
  - Nuevo flag `-port` para servidor HTTP
  - Nuevo flag `-rpc` para servidor gRPC
- Mejoras en la documentación
  - Guía de instalación actualizada
  - Ejemplos de uso mejorados
- Actualización de dependencias
  - Go 1.21.0
  - gRPC v1.58.0

### v1.1.0

- Implementación de generador de microservicios
  - Soporte para múltiples servicios
  - Configuración automática de endpoints
- Soporte para modelos de datos
  - Generación de estructuras Go
  - Validación de datos
- Integración con gRPC
  - Servicios bidireccionales
  - Streaming de datos
- Estructura base del proyecto
  - Organización de directorios
  - Archivos de configuración

### v1.0.1

- Correcciones de bugs
  - Solución de problemas de concurrencia
  - Mejora en el manejo de errores
- Optimizaciones de rendimiento
  - Reducción de uso de memoria
  - Mejora en tiempos de respuesta

### v1.0.0

- Lanzamiento inicial del framework
  - Generador de proyectos básico
  - Configuración inicial de Go modules
  - Estructura base del proyecto
  - Documentación inicial

### v0.9.0

- Versión beta
  - Pruebas de concepto
  - Feedback inicial de usuarios
  - Ajustes basados en pruebas

### v0.8.0

- Versión alpha
  - Desarrollo inicial
  - Características básicas implementadas
  - Pruebas internas

## Versionamiento

El proyecto utiliza un script `version.sh` para manejar el versionamiento de manera consistente. Este script automatiza el proceso de actualización de versiones siguiendo el estándar de [Semantic Versioning](https://semver.org/).

### Uso del Script de Versionamiento

```bash
# Para crear una nueva versión
./version.sh [major|minor|patch]

# Ejemplos:
./version.sh patch  # Incrementa la versión de parche (1.1.2 -> 1.1.3)
./version.sh minor  # Incrementa la versión menor (1.1.2 -> 1.2.0)
./version.sh major  # Incrementa la versión mayor (1.1.2 -> 2.0.0)
```

### Funcionalidades del Script

El script `version.sh` realiza las siguientes acciones:

1. Actualiza el número de versión en:

   - Archivo `go.mod`
   - Archivo `VERSION`
   - Tags de Git

2. Crea un nuevo tag de Git con el formato `vX.Y.Z`

3. Genera un commit con el mensaje "Release vX.Y.Z"

### Convenciones de Versionamiento

- **MAJOR**: Incrementa cuando hay cambios incompatibles en la API
- **MINOR**: Incrementa cuando se agregan funcionalidades de manera compatible
- **PATCH**: Incrementa cuando se corrigen bugs de manera compatible

### Ejemplo de Flujo de Trabajo

```bash
# 1. Realizar cambios en el código
# 2. Probar los cambios
# 3. Actualizar la versión
./version.sh patch

# 4. Verificar los cambios
git status
git log --oneline

# 5. Push de los cambios y tags
git push origin main --tags
```

### Notas Importantes

- Asegúrate de que todos los cambios estén commiteados antes de ejecutar el script
- El script debe tener permisos de ejecución (`chmod +x version.sh`)
- Siempre verifica los cambios generados por el script antes de hacer push

application/x-www-form-urlencoded
