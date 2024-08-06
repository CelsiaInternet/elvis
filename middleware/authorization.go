package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/cgalvisleon/elvis/claim"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/utility"
)

type contextKey string

const (
	serviceIDKey contextKey = "serviceId"
	clientIDKey  contextKey = "clientId"
	appKey       contextKey = "app"
	nameKey      contextKey = "name"
	kindKey      contextKey = "kind"
	usernameKey  contextKey = "username"
	tokenKey     contextKey = "token"
)

/**
* tokenFromAuthorization
* @param authorization string
* @return string
* @return error
**/
func tokenFromAuthorization(authorization string) (string, error) {
	if authorization == "" {
		return "", console.NewError("Autorization is required")
	}

	if !strings.HasPrefix(authorization, "Bearer") {
		return "", console.NewError("Invalid autorization format")
	}

	l := strings.Split(authorization, " ")
	if len(l) != 2 {
		return "", console.NewError("Invalid autorization format")
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
	authorization := r.Header.Get("Authorization")
	result, err := tokenFromAuthorization(authorization)
	if err != nil {
		return "", err
	}

	return result, nil
}

/**
* Authorization
* @param next http.Handler
**/
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metric := NewMetric(r)
		ctx := r.Context()
		tokenString, err := GetAuthorization(w, r)
		if err != nil {
			metric.Unauthorized(w, r)
			return
		}

		c, err := claim.GetFromToken(ctx, tokenString)
		if err != nil {
			metric.Unauthorized(w, r)
			return
		}

		serviceId := utility.UUID()
		ctx = context.WithValue(ctx, serviceIDKey, serviceId)
		ctx = context.WithValue(ctx, clientIDKey, c.ID)
		ctx = context.WithValue(ctx, appKey, c.App)
		ctx = context.WithValue(ctx, nameKey, c.Name)
		ctx = context.WithValue(ctx, kindKey, c.Kind)
		ctx = context.WithValue(ctx, usernameKey, c.Username)
		ctx = context.WithValue(ctx, tokenKey, tokenString)

		now := utility.Now()
		hostName, _ := os.Hostname()
		data := et.Json{
			"serviceId": serviceId,
			"clientId":  c.ID,
			"last_use":  now,
			"host_name": hostName,
			"token":     tokenString,
		}

		go event.Log("telemetry.token.last_use", data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
