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

func (c contextKey) String(ctx context.Context, def string) string {
	val := ctx.Value(c)
	result, ok := val.(string)
	if !ok {
		return def
	}

	return result
}

const (
	ServiceIDKey contextKey = "serviceId"
	ClientIDKey  contextKey = "clientId"
	AppKey       contextKey = "app"
	NameKey      contextKey = "name"
	KindKey      contextKey = "kind"
	UsernameKey  contextKey = "username"
	TokenKey     contextKey = "token"
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
		ctx = context.WithValue(ctx, ServiceIDKey, serviceId)
		ctx = context.WithValue(ctx, ClientIDKey, c.ID)
		ctx = context.WithValue(ctx, AppKey, c.App)
		ctx = context.WithValue(ctx, NameKey, c.Name)
		ctx = context.WithValue(ctx, KindKey, c.Kind)
		ctx = context.WithValue(ctx, UsernameKey, c.Username)
		ctx = context.WithValue(ctx, TokenKey, tokenString)

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
