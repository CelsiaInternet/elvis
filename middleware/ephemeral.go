package middleware

import (
	"context"
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
)

/**
* Ephemeral
* @param next http.Handler
**/
func Ephemeral(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetAuthorization(w, r)
		if err != nil {
			response.Unauthorized(w, r)
			return
		}

		clm, err := claim.ParceToken(token)
		if err != nil {
			response.Unauthorized(w, r)
			return
		}

		if clm.ExpiresAt < utility.NowTime().Unix() {
			response.Unauthorized(w, r)
			return
		}

		serviceId := utility.UUID()
		ctx := r.Context()
		ctx = context.WithValue(ctx, claim.ServiceIdKey, serviceId)
		ctx = context.WithValue(ctx, claim.ClientIdKey, clm.ID)
		ctx = context.WithValue(ctx, claim.AppKey, clm.App)
		ctx = context.WithValue(ctx, claim.DeviceKey, clm.Device)
		ctx = context.WithValue(ctx, claim.NameKey, clm.Name)
		ctx = context.WithValue(ctx, claim.SubjectKey, clm.Subject)
		ctx = context.WithValue(ctx, claim.UsernameKey, clm.Username)
		ctx = context.WithValue(ctx, claim.TokenKey, token)

		now := utility.Now()
		data := et.Json{
			"serviceId": serviceId,
			"clientId":  clm.ID,
			"last_use":  now,
			"host_name": hostName,
			"token":     token,
		}

		PushTokenLastUse(data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
