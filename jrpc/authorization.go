package jrpc

import (
	"fmt"
	"strconv"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/mem"
	"github.com/celsiainternet/elvis/middleware"
)

const (
	// Metodo RPC por defecto que resuelve los permisos (catalog)
	DEFAULT_AUTHORIZATION_METHOD = "catalog.Services.Author"
	// Segundos que se guarda en memoria cada respuesta
	AUTHORIZATION_CACHE_SECONDS = 10
)

/**
* AuthorStore: Implementa middleware.Store consultando los permisos por RPC.
**/
type AuthorStore struct {
	Method string
}

/**
* NewAuthorStore
* @param method string: metodo RPC; si es "" usa AUTHORIZATION_METHOD o DEFAULT_AUTHORIZATION_METHOD
* @return *AuthorStore
**/
func NewAuthorStore(method string) *AuthorStore {
	if method == "" {
		method = envar.GetStr(DEFAULT_AUTHORIZATION_METHOD, "AUTHORIZATION_METHOD")
	}

	return &AuthorStore{
		Method: method,
	}
}

/**
* LoadAuthorization: Registra un AuthorStore por RPC como store del middleware Authorization.
* @param method string
* @return *AuthorStore
**/
func LoadAuthorization(method string) *AuthorStore {
	store := NewAuthorStore(method)
	middleware.SetAuthorizationStore(store)

	return store
}

/**
* Author: Pregunta por RPC si el perfil tiene permiso en el proyecto (team) sobre el metodo y la ruta.
* Las respuestas se guardan en memoria AUTHORIZATION_CACHE_SECONDS, asi que un cambio de permisos
* puede tardar ese tiempo en verse en cada replica.
* @param projectId, profileId, method, path string
* @return bool, error
**/
func (s *AuthorStore) Author(projectId, profileId, method, path string) (bool, error) {
	if projectId == "" || profileId == "" {
		return false, nil
	}

	key := fmt.Sprintf("author:%s:%s:%s:%s", projectId, profileId, method, path)
	if cached, err := mem.Get(key, ""); err == nil {
		return strconv.ParseBool(cached)
	}

	result, err := CallItem(s.Method, et.Json{
		"project_id": projectId,
		"profile_id": profileId,
		"method":     method,
		"path":       path,
	})
	if err != nil {
		return false, err
	}

	if !result.Ok {
		return false, fmt.Errorf(MSG_INVALID_RESPONSE, s.Method)
	}

	ok := result.Bool("authorized")
	mem.Set(key, strconv.FormatBool(ok), AUTHORIZATION_CACHE_SECONDS)

	return ok, nil
}
