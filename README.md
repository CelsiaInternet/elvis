# Elvis framework

# Public

```
git push origin --tags

go build ./cmd/create-go
gofmt -w . && go run ./cmd/create-go
gofmt -w . && go run github.com/celsiainternet/elvis/cmd/create-go create
gofmt -w . && go run ./cmd/rpc/server
gofmt -w . && go run ./cmd/rpc/client

gofmt -w . && go run ./cmd/ws -port 3000 -username test -password test
gofmt -w . && go run ./cmd/ws -port 3100 -username test -password test

go run github.com/celsiainternet/elvis/cmd/create-go create
go run github.com/celsiainternet/elvis/cmd/apigateway

go build ./cmd/apigateway

go get -u github.com/celsiainternet/elvis@v1.1.68
go get github.com/celsiainternet/elvis@v1.1.68
```

## Create project

go mod init github.com/apimanager/api

## Dependencias

```
go get github.com/celsiainternet/elvis@v1.1.68
```

## Crear projecto, microservicios, modelos

```
go run github.com/celsiainternet/elvis/cmd/create-go create
```

## Version

```
go mod tidy &&
gofmt -w . &&
git update &&
git tag v1.1.68 &&
git push origin --tags

```

## Run

```
gofmt -w . && go run ./cmd/ws -port 3300 -rpc 4200
gofmt -w . && go run ./cmd/ws -port 3500 -rpc 4500
gofmt -w . && go run --race ./cmd/ws -port 3500
gofmt -w . && go run --race ./cmd/ws -port 3600

gofmt -w . && go run --race ./cmd/daemon -port 3600 -rpc 4600

gofmt -w . && go build --race -a -v -o ./daemon ./cmd/daemon

gofmt -w . && go build -race ./cmd/ws
```

# Build

```
docker system prune -a --volumes -f

docker build --no-cache -t apigateway -f ./cmd/apigateway/Dockerfile .
docker scout quickview local://apigateway:latest --org celsiainternet

docker-compose -p apigateway -f ./cmd/apigateway/docker-compose.yml up -d
docker-compose -p apigateway -f ./cmd/apigateway/docker-compose.yml down
```
