# Elvis framework

# Public

```
go mod tidy &&
gofmt -w . &&
git update &&
git tag v1.0.83 &&
git tags

git push origin --tags

go build ./cmd/create-go
gofmt -w . && go run ./cmd/create-go
gofmt -w . && go run github.com/celsiainternet/elvis/cmd/create-go create
gofmt -w . && go run ./cmd/rpc/server
gofmt -w . && go run ./cmd/rpc/client

gofmt -w . && go run ./cmd/gateway -port 3300 -rpc 4200

go run github.com/celsiainternet/elvis/cmd/create-go create
go run github.com/celsiainternet/elvis/cmd/apigateway

go build ./cmd/apigateway

go get -u github.com/celsiainternet/elvis@v1.0.83
go get github.com/celsiainternet/elvis@v1.0.83
```

# Build

```
docker system prune -a --volumes -f

docker build --no-cache -t apigateway -f ./cmd/apigateway/Dockerfile .
docker scout quickview local://apigateway:latest --org celsiainternet

docker-compose -p apigateway -f ./cmd/apigateway/docker-compose.yml up -d
docker-compose -p apigateway -f ./cmd/apigateway/docker-compose.yml down
```
