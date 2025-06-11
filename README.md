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
go get github.com/celsiainternet/elvis@v1.1.2
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

### v1.1.2

- Mejoras en la generación de microservicios
- Optimización del rendimiento en la creación de modelos
- Corrección de bugs en la inicialización del proyecto

### v1.1.1

- Agregado soporte para configuración de puertos personalizados
- Mejoras en la documentación
- Actualización de dependencias

### v1.1.0

- Implementación de generador de microservicios
- Soporte para modelos de datos
- Integración con gRPC
- Estructura base del proyecto

### v1.0.0

- Lanzamiento inicial del framework
- Generador de proyectos básico
- Configuración inicial de Go modules
