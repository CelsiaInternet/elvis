package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
)

/**
* tokenFromAuthorization
* @param authorization string
* @return string
* @return error
**/
func tokenFromAuthorization(authorization string) (string, error) {
	if authorization == "" {
		return "", console.Alert("Autorization is required")
	}

	if !strings.HasPrefix(authorization, "Bearer") {
		return "", console.Alert("Invalid autorization format")
	}

	l := strings.Split(authorization, " ")
	if len(l) != 2 {
		return "", console.Alert("Invalid autorization format")
	}

	return l[1], nil
}

/**
* GetAuthorization
* @param w http.ResponseWriter
* @param r *http.Request
* @return string
* @return error
**/
func GetAuthorization(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		return cookie.Value, nil
	}

	authorization := r.Header.Get("Authorization")
	result, err := tokenFromAuthorization(authorization)
	if err != nil {
		return "", console.AlertE(err)
	}

	return result, nil
}

/**
* Authorization
* @param next http.Handler
**/
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetAuthorization(w, r)
		if err != nil {
			response.Unauthorized(w, r)
			return
		}

		c, err := claim.GetFromToken(token)
		if err != nil {
			response.Unauthorized(w, r)
			return
		}

		if c == nil {
			response.Unauthorized(w, r)
			return
		}

		serviceId := utility.UUID()
		ctx := r.Context()
		ctx = context.WithValue(ctx, claim.ServiceIdKey, serviceId)
		ctx = context.WithValue(ctx, claim.ClientIdKey, c.Id)
		ctx = context.WithValue(ctx, claim.AppKey, c.App)
		ctx = context.WithValue(ctx, claim.NameKey, c.Name)
		ctx = context.WithValue(ctx, claim.KindKey, c.Kind)
		ctx = context.WithValue(ctx, claim.UsernameKey, c.Username)
		ctx = context.WithValue(ctx, claim.TokenKey, token)

		now := utility.Now()
		data := et.Json{
			"serviceId": serviceId,
			"clientId":  c.Id,
			"last_use":  now,
			"host_name": hostName,
			"token":     token,
		}

		go event.TokenLastUse(data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
