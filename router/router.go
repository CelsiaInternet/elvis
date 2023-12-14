package router

import (
	"net/http"

	"github.com/cgalvisleon/elvis/event"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/middleware"
	"github.com/go-chi/chi"
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

func PublicRoute(r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagepath, host string) *chi.Mux {
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

	event.Publish("router", "apimanager/upsert", e.Json{
		"kind":         "public",
		"method":       method,
		"path":         path,
		"package_name": packageName,
		"package_path": packagepath,
		"host":         host,
	})

	return r
}

func ProtectRoute(r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagepath, host string) *chi.Mux {
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

	event.Publish("router", "apimanager/upsert", e.Json{
		"kind":         "protect",
		"method":       method,
		"path":         path,
		"package_name": packageName,
		"package_path": packagepath,
		"host":         host,
	})

	return r
}
