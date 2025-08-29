package router

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jrpc"
	"github.com/celsiainternet/elvis/middleware"
	"github.com/celsiainternet/elvis/strs"
	"github.com/go-chi/chi/v5"
)

type TypeRoute int

const (
	HTTP TypeRoute = iota
	REST
	PROXY
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
	APIGATEWAY  = "apigateway"
)

var (
	APIGATEWAY_SET_RESOLVE    = fmt.Sprintf("%s/set/resolve", APIGATEWAY)
	APIGATEWAY_DELETE_RESOLVE = fmt.Sprintf("%s/delete/resolve", APIGATEWAY)
	APIGATEWAY_RESET          = fmt.Sprintf("%s/reset", APIGATEWAY)
	APIGATEWAY_SET_PROXY      = fmt.Sprintf("%s/set/proxy", APIGATEWAY)
	APIGATEWAY_DELETE_PROXY   = fmt.Sprintf("%s/delete/proxy", APIGATEWAY)
)

type TpHeader int

const (
	TpKeepHeader TpHeader = iota
	TpJoinHeader
	TpReplaceHeader
)

type TpRouter struct {
	name   string
	routes map[string]et.Json
}

var router *TpRouter

func initRouter(name string) {
	if router == nil {
		router = &TpRouter{
			name:   name,
			routes: make(map[string]et.Json),
		}

		resetApigateway()
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
* @param id, method, path, resolve string, header et.Json, tpHeader TpHeader, excludeHeader []string, private bool, packageName string
**/
func PushApiGateway(id, method, path, resolve string, header et.Json, tpHeader TpHeader, excludeHeader []string, private bool, packageName string) {
	initRouter(packageName)
	router.routes[id] = et.Json{
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

	event.Publish(APIGATEWAY_SET_RESOLVE, router.routes[id])
}

/**
* DeleteApiGatewayById
* @param id, method, path string
**/
func DeleteApiGatewayById(id, method, path string) {
	delete(router.routes, id)

	event.Publish(APIGATEWAY_DELETE_RESOLVE, et.Json{
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
	return router.routes
}

/**
* pushApiGateway
* @param method, path, packagePath, host, packageName string, private bool
**/
func pushApiGateway(method, path, packagePath, host, packageName string, private bool) {
	id := cache.GenKey(method, path)
	path = strings.ReplaceAll(packagePath+path, "//", "/")
	resolve := host + path
	PushApiGateway(id, method, path, resolve, et.Json{}, TpReplaceHeader, []string{}, private, packageName)
}

/**
* resetApigateway
**/
func resetApigateway() {
	event.Stack(APIGATEWAY_RESET, func(m event.EvenMessage) {
		for _, r := range router.routes {
			event.Publish(APIGATEWAY_SET_RESOLVE, r)
		}
	})

	channel := strs.Format(`%s/%s`, APIGATEWAY_RESET, router.name)
	event.Stack(channel, func(m event.EvenMessage) {
		for _, r := range router.routes {
			event.Publish(APIGATEWAY_SET_RESOLVE, r)
		}
	})
}

/**
* PublicRoute
* @param r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string
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
* @param r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string
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
* EphemeralRoute
* @param r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string
* @return *chi.Mux
**/
func EphemeralRoute(r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string) *chi.Mux {
	switch method {
	case "GET":
		r.With(middleware.Ephemeral).Get(path, h)
	case "POST":
		r.With(middleware.Ephemeral).Post(path, h)
	case "PUT":
		r.With(middleware.Ephemeral).Put(path, h)
	case "PATCH":
		r.With(middleware.Ephemeral).Patch(path, h)
	case "DELETE":
		r.With(middleware.Ephemeral).Delete(path, h)
	case "HEAD":
		r.With(middleware.Ephemeral).Head(path, h)
	case "OPTIONS":
		r.With(middleware.Ephemeral).Options(path, h)
	case "HandlerFunc":
		r.With(middleware.Ephemeral).HandleFunc(path, h)
	}

	pushApiGateway(method, path, packagePath, host, packageName, true)

	return r
}

/**
* AuthorizationRoute
* @param r *chi.Mux, method, path string, h http.HandlerFunc, packageName, packagePath, host string
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

/**
* With
* @param r *chi.Mux, method, path string, middlewares []func(http.Handler) http.Handler, h http.HandlerFunc, packageName, packagePath, host string
* @return *chi.Mux
**/
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

/**
* authorization
* @param profile et.Json
* @return map[string]bool, error
**/
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
