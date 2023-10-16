package router

import (
	"net/http"

	"github.com/cgalvisleon/elvis/middleware"
	"github.com/go-chi/chi"
)

const (
	Get     = "GET"
	Post    = "POST"
	Put     = "PUT"
	Patch   = "PATCH"
	Delete  = "DELETE"
	Head    = "HEAD"
	Options = "OPTIONS"
)

func ProtectRoute(r *chi.Mux, method, pattern string, h http.HandlerFunc) *chi.Mux {
	switch method {
	case "GET":
		r.With(middleware.Authorization).Get(pattern, h)
	case "POST":
		r.With(middleware.Authorization).Post(pattern, h)
	case "PUT":
		r.With(middleware.Authorization).Put(pattern, h)
	case "PATCH":
		r.With(middleware.Authorization).Patch(pattern, h)
	case "DELETE":
		r.With(middleware.Authorization).Delete(pattern, h)
	case "HEAD":
		r.With(middleware.Authorization).Head(pattern, h)
	case "OPTIONS":
		r.With(middleware.Authorization).Options(pattern, h)
	}

	return r
}

func PublicRoute(r *chi.Mux, method, pattern string, h http.HandlerFunc) *chi.Mux {
	switch method {
	case "GET":
		r.Get(pattern, h)
	case "POST":
		r.Post(pattern, h)
	case "PUT":
		r.Put(pattern, h)
	case "PATCH":
		r.Patch(pattern, h)
	case "DELETE":
		r.Delete(pattern, h)
	case "HEAD":
		r.Head(pattern, h)
	case "OPTIONS":
		r.Options(pattern, h)
	}

	return r
}
