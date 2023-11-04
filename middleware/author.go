package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/cgalvisleon/elvis/claim"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/response"
	. "github.com/cgalvisleon/elvis/utilities"
)

func tokenFromAuthorization(authorization string) (string, error) {
	if authorization == "" {
		return "", console.ErrorM("Autorization is required")
	}

	if !strings.HasPrefix(authorization, "Bearer") {
		return "", console.ErrorM("Invalid autorization format")
	}

	l := strings.Split(authorization, " ")
	if len(l) != 2 {
		return "", console.ErrorM("Invalid autorization format")
	}

	return l[1], nil
}

func GetAuthorization(w http.ResponseWriter, r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")
	result, err := tokenFromAuthorization(authorization)
	if err != nil {
		return "", err
	}

	return result, nil
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString, err := GetAuthorization(w, r)
		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		c, err := claim.GetFromToken(ctx, tokenString)
		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		ctx = context.WithValue(ctx, "clientId", c.ID)
		ctx = context.WithValue(ctx, "app", c.App)
		ctx = context.WithValue(ctx, "name", c.Name)
		ctx = context.WithValue(ctx, "kind", c.Kind)
		ctx = context.WithValue(ctx, "username", c.Username)
		ctx = context.WithValue(ctx, "token", tokenString)

		now := Now()
		hostName, _ := os.Hostname()
		data := Json{
			"clientId":  c.ID,
			"last_use":  now,
			"host_name": hostName,
		}

		event.Action("telemetry.token.last_use", data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
