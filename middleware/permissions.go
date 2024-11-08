package middleware

import (
	"net/http"

	"github.com/celsiainternet/elvis/module"
	"github.com/celsiainternet/elvis/response"
)

/**
* PermissionsMiddleware
* @param next http.Handler
* @return http.Handler
**/
func PermissionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		project_id := r.Header.Get("project_id")
		profile_tp := r.Header.Get("profile_tp")
		model := r.Header.Get("model")
		permisions, err := module.GetPermissions(project_id, profile_tp, model)
		if err != nil {
			response.InternalServerError(w, r)
			return
		}

		ok := permisions.Method(r)
		if !ok {
			response.Forbidden(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
