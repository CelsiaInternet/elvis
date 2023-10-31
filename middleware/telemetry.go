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
	. "github.com/cgalvisleon/elvis/utilities"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shirou/gopsutil/mem"
)

type Request struct {
	Tag     string
	Day     int
	Hour    int
	Minute  int
	Seccond int
	Limit   int
}

func CallRequests(tag string) Request {
	return Request{
		Tag:     tag,
		Day:     cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/86400), 86400),
		Hour:    cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/3600), 3600),
		Minute:  cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/60), 60),
		Seccond: cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/2), 2),
		Limit:   envar.EnvarInt(1000, "REQUESTS_LIMIT"),
	}
}

func Telemetry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_id := NewId()

		if rctx := chi.RouteContext(ctx); rctx != nil {
			endPoint := r.URL.Path
			t1 := time.Now()
			hostName, _ := os.Hostname()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			headers := Json{}
			for key, val := range ww.Header() {
				headers[key] = val
			}
			var mTotal uint64
			var mUsed uint64
			var mFree uint64
			memory, err := mem.VirtualMemory()
			if err != nil {
				mFree = 0
				mTotal = 0
				mUsed = 0
			}
			mTotal = memory.Total
			mUsed = memory.Used
			mFree = memory.Free
			requests_host := CallRequests(hostName)
			requests_endpoint := CallRequests(endPoint)

			summary := Json{
				"_id":           _id,
				"datetime":      t1,
				"host_name":     hostName,
				"method":        rctx.RouteMethod,
				"endpoint":      endPoint,
				"status":        ww.Status(),
				"bytes_written": ww.BytesWritten(),
				"header":        headers,
				"since":         time.Since(t1),
				"memory": Json{
					"total": mTotal,
					"used":  mUsed,
					"free":  mFree,
				},
				"request_host": Json{
					"day":    requests_host.Day,
					"hour":   requests_host.Hour,
					"minute": requests_host.Minute,
					"second": requests_host.Seccond,
					"limit":  requests_host.Limit,
				},
				"requests_endpoint": Json{
					"day":    requests_endpoint.Day,
					"hour":   requests_endpoint.Hour,
					"minute": requests_endpoint.Minute,
					"second": requests_endpoint.Seccond,
					"limit":  requests_endpoint.Limit,
				},
			}
			event.EventPublish("telemetry", summary)

			if requests_host.Seccond >= requests_host.Limit {
				event.EventPublish("requests/overflow", summary)
			}
		}

		w.Header().Set("_id", _id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
