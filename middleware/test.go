package middleware

import (
	"context"
	"net/http"

	"github.com/cgalvisleon/elvis/logs"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "elvis/middleware context value " + k.name
}

func Test(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "app", "elvis")
		logs.Debug("middleware.Middleware next.ServeHTTP")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
