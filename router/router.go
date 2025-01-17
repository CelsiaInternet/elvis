package router

import (
	"errors"
	"net/http"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jrpc"
	"github.com/celsiainternet/elvis/middleware"
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

var router = make(map[string]et.Json)

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
* @params header et.Json
* @param tpHeader TpHeader
* @param private bool
* @param packageName string
**/
func PushApiGateway(id, method, path, resolve string, header et.Json, tpHeader TpHeader, excludeHeader []string, private bool, packageName string) {
	router[id] = et.Json{
		"_id":            id,
		"kind":           HTTP,
		"method":         method,
		"path":           path,
		"resolve":        resolve,
		"header":         header,
		"tp_header":      tpHeader,
		"exclude_header": excludeHeader,
		"private":        private,
		"package_name":   packageName,
	}

	event.Publish("apigateway/set/resolve", router[id])
}

/**
* DeleteApiGatewayById
* @param id string
**/
func DeleteApiGatewayById(id, method, path string) {
	delete(router, id)
	event.Publish("apigateway/delete/resolve", et.Json{
		"_id":    id,
		"method": method,
		"path":   path,
	})
}

/**
* GetRoutes
* @return map[string]et.Json
**/
func GetRoutes() map[string]et.Json {
	return router
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
	id := cache.GenKey(method, path)

	PushApiGateway(id, method, path, resolve, et.Json{}, TpReplaceHeader, []string{}, private, packageName)
}

/**
* PublicRoute
* @param r *chi.Mux
* @param method string
* @param path string
* @param h http.HandlerFunc
* @param packageName string
* @param packagePath string
* @param host string
* @return *chi.Mux
**/
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

/**
* ProtectRoute
* @param r *chi.Mux
* @param method string
* @param path string
* @param h http.HandlerFunc
* @param packageName string
* @param packagePath string
* @param host string
* @return *chi.Mux
**/
func ProtectRoute(r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string) *chi.Mux {
	switch method {
	case "GET":
		r.With(middleware.Autentication).Get(path, h)
	case "POST":
		r.With(middleware.Autentication).Post(path, h)
	case "PUT":
		r.With(middleware.Autentication).Put(path, h)
	case "PATCH":
		r.With(middleware.Autentication).Patch(path, h)
	case "DELETE":
		r.With(middleware.Autentication).Delete(path, h)
	case "HEAD":
		r.With(middleware.Autentication).Head(path, h)
	case "OPTIONS":
		r.With(middleware.Autentication).Options(path, h)
	case "HandlerFunc":
		r.With(middleware.Autentication).HandleFunc(path, h)
	}

	pushApiGateway(method, path, packagePath, host, packageName, true)

	return r
}

/**
* AuthorizationRoute
* @param r *chi.Mux
* @param method string
* @param path string
* @param h http.HandlerFunc
* @param packageName string
* @param packagePath string
* @param host string
* @return *chi.Mux
**/
func AuthorizationRoute(r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string) *chi.Mux {
	switch method {
	case "GET":
		r.With(middleware.Autentication).With(middleware.Authorization).Get(path, h)
	case "POST":
		r.With(middleware.Autentication).With(middleware.Authorization).Post(path, h)
	case "PUT":
		r.With(middleware.Autentication).With(middleware.Authorization).Put(path, h)
	case "PATCH":
		r.With(middleware.Autentication).With(middleware.Authorization).Patch(path, h)
	case "DELETE":
		r.With(middleware.Autentication).With(middleware.Authorization).Delete(path, h)
	case "HEAD":
		r.With(middleware.Autentication).With(middleware.Authorization).Head(path, h)
	case "OPTIONS":
		r.With(middleware.Autentication).With(middleware.Authorization).Options(path, h)
	case "HandlerFunc":
		r.With(middleware.Autentication).With(middleware.Authorization).HandleFunc(path, h)
	}

	pushApiGateway(method, path, packagePath, host, packageName, true)

	return r
}

func With(r *chi.Mux, method, path string, middlewares []func(http.Handler) http.Handler, h http.HandlerFunc, packageName, packagePath, host string) *chi.Mux {
	switch method {
	case "GET":
		r.With(middlewares...).Get(path, h)
	case "POST":
		r.With(middlewares...).Post(path, h)
	case "PUT":
		r.With(middlewares...).Put(path, h)
	case "PATCH":
		r.With(middlewares...).Patch(path, h)
	case "DELETE":
		r.With(middlewares...).Delete(path, h)
	case "HEAD":
		r.With(middlewares...).Head(path, h)
	case "OPTIONS":
		r.With(middlewares...).Options(path, h)
	case "HandlerFunc":
		r.With(middlewares...).HandleFunc(path, h)
	}

	pushApiGateway(method, path, packagePath, host, packageName, true)

	return r
}

func authorization(profile et.Json) (map[string]bool, error) {
	method := envar.GetStr("Module.Services.GetPermissions", "AUTHORIZATION_METHOD")
	if method == "" {
		return map[string]bool{}, errors.New("authorization method not found")
	}

	result, err := jrpc.CallPermitios(method, profile)
	if err != nil {
		return map[string]bool{}, err
	}

	return result, nil
}

func init() {
	middleware.SetAuthorizationFunc(authorization)
}
