# Elvis framework

# Public

```
go mod tidy &&
gofmt -w . &&
git update &&
git tag v0.0.133 &&
git tags
git push origin --tags

go build ./cmd/create-go
gofmt -w . && go run ./cmd/create-go
gofmt -w . && go run github.com/cgalvisleon/elvis/cmd/create-go create

go run github.com/cgalvisleon/elvis/cmd/create-go create
go run github.com/cgalvisleon/elvis/cmd/apigateway

go get -u github.com/cgalvisleon/elvis@v0.0.133
go get github.com/cgalvisleon/elvis@v0.0.133
```
