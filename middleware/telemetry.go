package middleware

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
	"github.com/shirou/gopsutil/v3/mem"
)

var DefaultTelemetry func(next http.Handler) http.Handler

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
		Seccond: cache.More(fmt.Sprintf(`%s-%d`, tag, time.Now().Unix()/1), 1),
		Limit:   envar.EnvarInt(10, "REQUESTS_LIMIT"),
	}
}

func Telemetry(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_id := NewId()
		ctx := r.Context()
		endPoint := r.URL.Path
		method := r.Method
		t1 := time.Now()
		hostName, _ := os.Hostname()
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
		mFree = memory.Total - memory.Used
		pFree := float64(mFree) / float64(mTotal) * 100
		requests_host := CallRequests(hostName)
		requests_endpoint := CallRequests(endPoint)

		defer func() {
			summary := Json{
				"_id":       _id,
				"date_time": t1,
				"host_name": hostName,
				"method":    method,
				"endpoint":  endPoint,
				"status":    http.StatusOK,
				"since": Json{
					"value": time.Since(t1).Milliseconds(),
					"unity": "Milliseconds",
				},
				"memory": Json{
					"unity":        "MB",
					"total":        mTotal / 1024 / 1024,
					"used":         mUsed / 1024 / 1024,
					"free":         mFree / 1024 / 1024,
					"percent_free": math.Floor(pFree*100) / 100,
				},
				"request_host": Json{
					"host":   requests_host.Tag,
					"day":    requests_host.Day,
					"hour":   requests_host.Hour,
					"minute": requests_host.Minute,
					"second": requests_host.Seccond,
					"limit":  requests_host.Limit,
				},
				"requests_endpoint": Json{
					"endpoint": requests_endpoint.Tag,
					"day":      requests_endpoint.Day,
					"hour":     requests_endpoint.Hour,
					"minute":   requests_endpoint.Minute,
					"second":   requests_endpoint.Seccond,
					"limit":    requests_endpoint.Limit,
				},
			}
			go event.EventPublish("telemetry", summary)

			if requests_host.Seccond > requests_host.Limit {
				go event.EventPublish("requests/overflow", summary)
			}
		}()

		w.Header().Set("_id", _id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
