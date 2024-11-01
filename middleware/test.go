package middleware

import (
	"context"
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/logs"
)

var app = claim.ContextKey("app")

func Test(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, app, "elvis")
		logs.Debug("middleware.Middleware next.ServeHTTP")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
