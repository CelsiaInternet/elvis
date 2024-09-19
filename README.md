# Elvis framework

# Public

```
go mod tidy &&
gofmt -w . &&
git update &&
git tag v1.0.9 &&
git tags

git push origin --tags

go build ./cmd/create-go
gofmt -w . && go run ./cmd/create-go
gofmt -w . && go run github.com/cgalvisleon/elvis/cmd/create-go create
gofmt -w . && go run ./cmd/gateway -port 3300 -rpc 4200

go run github.com/cgalvisleon/elvis/cmd/create-go create
go run github.com/cgalvisleon/elvis/cmd/apigateway

go build ./cmd/apigateway

go get -u github.com/cgalvisleon/elvis@v1.0.9
go get github.com/cgalvisleon/elvis@v1.0.9
```

# Build

```
docker system prune -a --volumes -f

docker build --no-cache -t apigateway -f ./cmd/apigateway/Dockerfile .
docker scout quickview local://apigateway:latest --org cgalvisleon

docker-compose -p apigateway -f ./cmd/apigateway/docker-compose.yml up -d
docker-compose -p apigateway -f ./cmd/apigateway/docker-compose.yml down
```
