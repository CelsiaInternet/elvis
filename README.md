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

- 🛡️ **Sistema de resiliencia** y recuperación automática
- 📅 **Tareas programadas** (Crontab)
- 🗄️ **Base de datos avanzada** con triggers automáticos
- 🔄 **Sistema de eventos** distribuidos
- 💾 **Cache inteligente** multi-backend
- 🔐 **Middleware de seguridad** integrado
- 📊 **Telemetría, logging y monitoreo**

## Requisitos Previos

- Go 1.23 o superior
- Git
- PostgreSQL (para base de datos)
- Redis (para cache)
- NATS (para eventos)

## 🚀 Instalación Rápida

### 1. Inicializar el Proyecto

```bash
go mod init github.com/tu-usuario/tu-proyecto
```

### 2. Instalar Elvis

```bash
go get github.com/celsiainternet/elvis@v1.1.196
go run github.com/celsiainternet/elvis/cmd/prepare
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

# Crontab
gofmt -w . && go run ./cmd/crontab go
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

- 🛡️ **Sistema de resiliencia** y recuperación automática

```go
// Configurar resiliencia
resilience.SetNotifyType(resilience.TpNotifyEmail)
resilience.SetContactNumbers([]string{"+573160479724"})

// Agregar transacción con reintentos automáticos
transaction := resilience.Add("email-send", "Enviar email de confirmación", sendEmail, userEmail, content)
```

- 📅 **Tareas programadas** (Crontab)

```go
// Configurar crontab
// Crear tarea programada
jobs := crontab.New()
jobs.AddJob("backup-daily", "Backup diario", "0 2 * * *", "backup-channel", map[string]interface{}{
    "type": "daily",
    "path": "/backup",
})

// Iniciar tareas
jobs.Start()
```

- 🗄️ **Base de datos avanzada** con triggers automáticos

```go
// Configurar base de datos
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

- 🔄 **Sistema de eventos** distribuidos

```go
// Configurar eventos
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

- 💾 **Cache inteligente** multi-backend

```go
// Configurar cache
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

- 🔐 **Middleware de seguridad** integrado

```go
// Configurar seguridad
// Middleware de autenticación
r.Use(middleware.Authentication)

// Middleware de autorización
r.Use(middleware.Authorization)

// Middleware de CORS
r.Use(middleware.CORS)

// Middleware de logging
r.Use(middleware.Logger)
```

- 📊 **Telemetría, logging y monitoreo**

```go
// Configurar telemetría
// Logging estructurado
logs.Log("user-service", "Usuario creado exitosamente")
logs.Alert(fmt.Errorf("Error de conexión"))
```

## 📁 Estructura del Proyecto

elvis/
├── 📂 cmd/ # Comandos ejecutables
│ ├── cmd/ # Servidor principal
│ ├── create-go/ # Generador de proyectos
│ ├── daemon/ # Daemon del sistema
│ ├── resilence/ # Servidor de resiliencia
│ └── rpc/ # Servidores RPC (cliente/servidor)
├── 📂 cache/ # Sistema de cache
├── 📂 claim/ # Manejo de claims/tokens
├── 📂 config/ # Configuración global
├── 📂 console/ # Utilidades de consola
├── 📂 create/ # Generación de código
│ └── template/ # Plantillas de código
├── 📂 crontab/ # Tareas programadas
├── 📂 envar/ # Variables de entorno
├── 📂 et/ # Tipos y utilidades
├── 📂 event/ # Sistema de eventos
├── 📂 file/ # Manejo de archivos
├── 📂 jdb/ # Database abstraction layer
├── 📂 jrpc/ # JSON-RPC implementation
├── 📂 linq/ # Query builder
├── 📂 logs/ # Sistema de logging
├── 📂 mem/ # Cache en memoria
├── 📂 middleware/ # Middleware HTTP
├── 📂 msg/ # Mensajería
├── 📂 race/ # Control de concurrencia
├── 📂 resilience/ # Sistema de resiliencia
├── 📂 response/ # Manejo de respuestas HTTP
├── 📂 router/ # Enrutamiento HTTP
├── 📂 service/ # Servicios base
├── 📂 stdrout/ # Salida estándar
├── 📂 strs/ # Utilidades de strings
├── 📂 timezone/ # Manejo de zonas horarias
├── 📂 utility/ # Utilidades generales
├── 📂 .vscode/ # Configuración IDE
│ ├── settings.json # Configuración optimizada
│ └── launch.json # Configuración debug
├── 📄 staticcheck.conf # Configuración linting
├── 📄 go.mod # Dependencias Go
├── 📄 go.sum # Checksums dependencias
├── 📄 version.sh # Script versionamiento
└── 📄 README.md # Documentación

## 🚀 Quick Start

```bash
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

workflow:f2334584-71f5-4be7-9c2c-0c3352bc9d50
