package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

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

			requests_day := cache.More(fmt.Sprintf(`%d`, time.Now().Unix()/86400), 86400)
			requests_hour := cache.More(fmt.Sprintf(`%d`, time.Now().Unix()/3600), 3600)
			requests_minute := cache.More(fmt.Sprintf(`%d`, time.Now().Unix()/60), 60)
			requests_second := cache.More(fmt.Sprintf(`%d`, time.Now().Unix()/2), 2)
			limit := envar.EnvarInt(1000, "REQUESTS_LIMIT")

			summary := Json{
				"datetime":        t1,
				"host_name":       hostName,
				"method":          rctx.RouteMethod,
				"endpoint":        endPoint,
				"status":          ww.Status(),
				"bytes_written":   ww.BytesWritten(),
				"header":          headers,
				"since":           time.Since(t1),
				"requests_day":    requests_day,
				"requests_hour":   requests_hour,
				"requests_minute": requests_minute,
				"requests_second": requests_second,
				"limit":           limit,
			}
			event.EventPublish("telemetry", summary)

			if requests_second >= limit {
				event.EventPublish("requests/overflow", summary)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
