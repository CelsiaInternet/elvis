package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/cgalvisleon/elvis/claim"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/response"
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

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString, err := claim.GetAuthorization(w, r)
		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, "401 Unauthorized")
			return
		}

		c, err := claim.GetFromToken(ctx, tokenString)
		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, "401 Unauthorized")
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
