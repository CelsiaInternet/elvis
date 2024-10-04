package router

import (
	"net/http"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/middleware"
	"github.com/cgalvisleon/elvis/utility"
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

type TpHeader int

const (
	TpKeepHeader TpHeader = iota
	TpJoinHeader
	TpReplaceHeader
)

/**
* IntToTpHeader
* @param tp int
* @return TpHeader
**/
func IntToTpHeader(tp int) TpHeader {
	switch tp {
	case 1:
		return TpJoinHeader
	case 2:
		return TpReplaceHeader
	default:
		return TpKeepHeader
	}
}

/**
* String
* @return string
**/
func (t TpHeader) String() string {
	switch t {
	case TpKeepHeader:
		return "Keep the resolve header"
	case TpJoinHeader:
		return "Join request header with the resolve header"
	case TpReplaceHeader:
		return "Replace resolve header with request header"
	default:
		return "Unknown"
	}
}

/**
* ToTpHeader
* @param str string
* @return TpHeader
**/
func ToTpHeader(tp int) TpHeader {
	switch tp {
	case 1:
		return TpJoinHeader
	case 2:
		return TpReplaceHeader
	default:
		return TpKeepHeader
	}
}

/**
* PushApiGateway
* @param method string
* @param path string
* @param resolve string
* @param host string
* @param packageName string
* @param private bool
**/
func PushApiGateway(method, path, resolve, host, packageName string, private bool) {
	event.Work("apigateway/http/resolve", et.Json{
		"kind":         HTTP,
		"method":       method,
		"path":         path,
		"resolve":      resolve,
		"package":      packageName,
		"tpHeader":     TpReplaceHeader,
		"private":      private,
		"package_name": packageName,
		"_id":          utility.UUID(),
	})
}

/**
* PopApiGatewayById
* @param id string
**/
func PopApiGatewayById(id string) {
	event.Work("apigateway/http/pop", et.Json{
		"_id": id,
	})
}

/**
* PushApiGateway
* @param method string
* @param path string
* @param packagePath string
* @param host string
* @param packageName string
* @param private bool
**/
func pushApiGateway(method, path, packagePath, host, packageName string, private bool) {
	path = packagePath + path
	resolve := host + path

	PushApiGateway(method, path, resolve, host, packageName, private)
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

	pushApiGateway(method, path, packagePath, host, packageName, false)

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

	pushApiGateway(method, path, packagePath, host, packageName, true)

	return r
}
