# Elvis - Framework para Microservicios en Go

## Descripción

Elvis es un framework diseñado para facilitar el desarrollo de microservicios en Go, proporcionando herramientas y estructuras para crear aplicaciones robustas y escalables con capacidades de comunicación en tiempo real, resiliencia y gestión de transacciones.

## Requisitos Previos

- Go 1.23 o superior
- Git
- PostgreSQL (para base de datos)
- Redis (para cache y comunicación en tiempo real)
- NATS (para eventos)

## Instalación

### Inicializar el Proyecto

```bash
go mod init github.com/test/api
```

### Instalar Dependencias

```bash
go get github.com/celsiainternet/elvis@v1.1.99
```

## Uso

### Creación del Proyecto

Para crear un nuevo proyecto con Elvis, ejecuta el siguiente comando:

```bash
go run github.com/celsiainternet/elvis/cmd/create-go create
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

## Estructura del Proyecto

```
.
├── cmd/
│   ├── test/
│   ├── ws/
│   ├── daemon/
│   └── resilence/
├── internal/
│   ├── models/
│   └── services/
├── pkg/
├── cache/
├── event/
├── ws/
├── realtime/
├── resilience/
├── crontab/
└── go.mod
```

## Comandos Disponibles

### Servidor WebSocket

```bash
go run ./cmd/ws -port 3300 -username admin -password secret
```

### Servidor de Resiliencia

```bash
go run ./cmd/resilence
```

### Daemon del Sistema

```bash
go run ./cmd/daemon --status
go run ./cmd/daemon --restart
```

## Configuración de Variables de Entorno

```bash
# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_NAME=elvis_db
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

## Contribución

Las contribuciones son bienvenidas. Por favor, lee nuestras guías de contribución antes de enviar un pull request.

## Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## Releases

### v1.1.99

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
