package router

import (
	"net/http"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/middleware"
	"github.com/go-chi/chi/v5"
)

type TypeRoute int

const (
	HTTP TypeRoute = iota
	REST
)

const (
	Get         = "GET"
	Post        = "POST"
	Put         = "PUT"
	Patch       = "PATCH"
	Delete      = "DELETE"
	Head        = "HEAD"
	Options     = "OPTIONS"
	HandlerFunc = "HandlerFunc"
)

func SetApiGateway(method, path, resolve string, kind TypeRoute, packageName string) {
	event.Publish("gateway", "gateway/upsert", et.Json{
		"kind":    kind,
		"method":  method,
		"path":    path,
		"resolve": resolve,
		"package": packageName,
	})
}

func pushApiGateway(method, path, packagePath, host, packageName string) {
	path = packagePath + path
	resolve := host + path

	SetApiGateway(method, path, resolve, HTTP, packageName)
}

func PublicRoute(r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string) *chi.Mux {
	switch method {
	case "GET":
		r.Get(path, h)
	case "POST":
		r.Post(path, h)
	case "PUT":
		r.Put(path, h)
	case "PATCH":
		r.Patch(path, h)
	case "DELETE":
		r.Delete(path, h)
	case "HEAD":
		r.Head(path, h)
	case "OPTIONS":
		r.Options(path, h)
	case "HandlerFunc":
		r.HandleFunc(path, h)
	}

	pushApiGateway(method, path, packagePath, host, packageName)

	return r
}

func ProtectRoute(r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string) *chi.Mux {
	switch method {
	case "GET":
		r.With(middleware.Authorization).Get(path, h)
	case "POST":
		r.With(middleware.Authorization).Post(path, h)
	case "PUT":
		r.With(middleware.Authorization).Put(path, h)
	case "PATCH":
		r.With(middleware.Authorization).Patch(path, h)
	case "DELETE":
		r.With(middleware.Authorization).Delete(path, h)
	case "HEAD":
		r.With(middleware.Authorization).Head(path, h)
	case "OPTIONS":
		r.With(middleware.Authorization).Options(path, h)
	case "HandlerFunc":
		r.With(middleware.Authorization).HandleFunc(path, h)
	}

	pushApiGateway(method, path, packagePath, host, packageName)

	return r
}
