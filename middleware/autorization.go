package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
)

type Store interface {
	Author(projectId, profileId, method, path string) (bool, error)
	RemoveAuthor(projectId, profileId, method, path string) error
	SetAuthor(projectId, profileId, method, path string) error
	SetPath(method, path string) error
}

var store Store

/**
* SetStore
* @param s Store
**/
func SetAuthorizationStore(s Store) {
	store = s
}

/**
* Authorization
* @param next http.Handler
* @return http.Handler
**/
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appName, ok := r.Header["AppName"]
		if !ok {
			response.PreconditionRequired(w, r, "AppName")
			return
		}

		if store != nil {
			response.InternalServerError(w, r, errors.New("Author store is not set"))
			return
		}

		token, err := GetAuthorization(w, r)
		if err != nil {
			response.Unauthorized(w, r)
			return
		}

		clm, err := claim.ValidToken(token)
		if err != nil {
			response.Unauthorized(w, r)
			return
		}

		if clm == nil {
			response.Unauthorized(w, r)
			return
		}

		ok, err = store.Author(clm.ProjectId, clm.ProfileId, r.Method, r.URL.Path)
		if err != nil {
			response.InternalServerError(w, r, err)
			return
		}

		if !ok {
			response.Forbidden(w, r)
			return
		}

		serviceId := utility.UUID()
		_, ok = r.Header["ServiceId"]
		if ok {
			serviceId = r.Header.Get("ServiceId")
		}

		ownerId := ""
		_, ok = r.Header["OwnerId"]
		if ok {
			ownerId = r.Header.Get("OwnerId")
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, claim.ServiceIdKey, serviceId)
		ctx = context.WithValue(ctx, claim.OwnerIdKey, ownerId)
		ctx = context.WithValue(ctx, claim.ClientIdKey, clm.ID)
		ctx = context.WithValue(ctx, claim.AppKey, clm.App)
		ctx = context.WithValue(ctx, claim.DeviceKey, clm.Device)
		ctx = context.WithValue(ctx, claim.NameKey, clm.Name)
		ctx = context.WithValue(ctx, claim.SubjectKey, clm.Subject)
		ctx = context.WithValue(ctx, claim.UsernameKey, clm.Username)
		ctx = context.WithValue(ctx, claim.ProjectIdKey, clm.ProjectId)
		ctx = context.WithValue(ctx, claim.ProfileIdKey, clm.ProfileId)
		ctx = context.WithValue(ctx, claim.AppNAmeKey, appName)
		ctx = context.WithValue(ctx, claim.TokenKey, token)

		now := timezone.Now()
		data := et.Json{
			"last_use":   now,
			"host_name":  hostName,
			"service_id": serviceId,
			"owner_id":   ownerId,
			"client_id":  clm.ID,
			"app":        clm.App,
			"device":     clm.Device,
			"name":       clm.Name,
			"subject":    clm.Subject,
			"username":   clm.Username,
			"project_id": clm.ProjectId,
			"profile_id": clm.ProfileId,
			"app_name":   appName,
			"token":      token,
		}

		PushTokenLastUse(data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
