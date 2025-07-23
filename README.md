# Elvis - Framework para Microservicios en Go

## Descripción

Elvis es un framework diseñado para facilitar el desarrollo de microservicios en Go, proporcionando herramientas y estructuras para crear aplicaciones robustas y escalables.

## Requisitos Previos

- Go 1.16 o superior
- Git

## Instalación

### Inicializar el Proyecto

```bash
go mod init github.com/test/api
```

### Instalar Dependencias

```bash
go get github.com/celsiainternet/elvis@v1.1.94
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

## Estructura del Proyecto

```
.
├── cmd/
│   └── test/
├── internal/
│   ├── models/
│   └── services/
├── pkg/
└── go.mod
```

## Contribución

Las contribuciones son bienvenidas. Por favor, lee nuestras guías de contribución antes de enviar un pull request.

## Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## Releases

### v1.1.94

- Implementación de sistema de tareas programadas (Crontab)
  - Nuevo generador de tareas programadas
  - Soporte para expresiones cron estándar
  - Integración con el sistema de logging
  - Manejo de errores en tareas programadas
- Mejoras en la generación de código
  - Nuevas plantillas para tareas programadas
  - Optimización en la generación de modelos
  - Mejora en la documentación generada
- Correcciones
  - Ajustes en el manejo de memoria
  - Actualización de dependencias

### v1.1.94

- Mejoras en la estabilidad del framework
  - Optimización del sistema de generación de código
  - Mejora en el manejo de dependencias circulares
  - Actualización de las plantillas de generación
- Nuevas características
  - Soporte para middleware personalizado
  - Integración mejorada con sistemas de logging
  - Nuevos helpers para validación de datos
- Correcciones
  - Solución de problemas de memoria en servicios largos
  - Mejora en el manejo de errores en la generación de modelos
  - Actualización de dependencias de seguridad

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
