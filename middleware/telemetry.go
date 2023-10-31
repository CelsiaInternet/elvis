package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var requests map[int64]int

func Telemetry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if rctx := chi.RouteContext(ctx); rctx != nil {
			endPoint := "/"
			n := len(rctx.RoutePatterns)
			if n > 0 {
				endPoint = fmt.Sprintf(`%s`, rctx.RoutePatterns[n-1])
			}
			t1 := time.Now()
			hostName, _ := os.Hostname()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			headers := Json{}
			for key, val := range ww.Header() {
				headers[key] = val
			}

			limit := envar.EnvarInt(1000, "REQUESTS_LIMIT")
			current := time.Now().Unix()
			_, exist := requests[current]
			if exist {
				requests[current]++
			} else {
				requests[current] = 1
			}

			summary := Json{
				"datetime":      t1,
				"host_name":     hostName,
				"method":        rctx.RouteMethod,
				"endpoint":      endPoint,
				"status":        ww.Status(),
				"bytes_written": ww.BytesWritten(),
				"header":        headers,
				"since":         time.Since(t1),
				"requests": requests[current],
				"limit": limit,
			}
			event.EventPublish("telemetry", summary)

			if requests[current] >= limit {
				event.EventPublish("requests/overflow", summary)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
