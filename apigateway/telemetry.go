package apigateway

import (
	"math"
	"net/http"
	"os"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
	"github.com/shirou/gopsutil/v3/mem"
)

var DefaultTelemetry func(next http.Handler) http.Handler

type Metrics struct {
	ReqID            string
	TimeBegin        time.Time
	TimeEnd          time.Time
	TimeExec         time.Time
	SearchTime       time.Duration
	ResponseTime     time.Duration
	TotalTime        time.Duration
	EndPoint         string
	Method           string
	RemoteAddr       string
	HostName         string
	Proto            string
	MTotal           uint64
	MUsed            uint64
	MFree            uint64
	PFree            float64
	RequestsHost     Request
	RequestsEndpoint Request
	Scheme           string
}

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
		Day:     cache.More(strs.Format(`%s-%d`, tag, time.Now().Unix()/86400), 86400),
		Hour:    cache.More(strs.Format(`%s-%d`, tag, time.Now().Unix()/3600), 3600),
		Minute:  cache.More(strs.Format(`%s-%d`, tag, time.Now().Unix()/60), 60),
		Seccond: cache.More(strs.Format(`%s-%d`, tag, time.Now().Unix()/1), 1),
		Limit:   envar.EnvarInt(400, "REQUESTS_LIMIT"),
	}
}

func telemetryNew(r *http.Request) *Metrics {
	result := &Metrics{}
	result.TimeBegin = time.Now()
	result.ReqID = utility.NewId()
	result.EndPoint = r.URL.Path
	result.Method = r.Method
	result.Proto = r.Proto
	result.RemoteAddr = r.RemoteAddr
	result.HostName, _ = os.Hostname()
	memory, err := mem.VirtualMemory()
	if err != nil {
		result.MFree = 0
		result.MTotal = 0
		result.MUsed = 0
		result.PFree = 0
	} else {
		result.MTotal = memory.Total
		result.MUsed = memory.Used
		result.MFree = memory.Total - memory.Used
		result.PFree = float64(result.MFree) / float64(result.MTotal) * 100
	}
	result.RequestsHost = CallRequests(result.HostName)
	result.RequestsEndpoint = CallRequests(result.EndPoint)
	result.Scheme = "http"
	if r.TLS != nil {
		result.Scheme = "https"
	}

	return result
}

func (m *Metrics) done(res *http.Response) et.Json {
	m.TimeEnd = time.Now()
	m.ResponseTime = time.Since(m.TimeExec)
	m.TotalTime = time.Since(m.TimeBegin)

	console.LogKF(m.Method, "%s://%s %s %v%s %s %v %s", m.Scheme, m.EndPoint, m.Proto, res.ContentLength, "KB", res.Status, m.TotalTime, "ms")

	result := et.Json{
		"reqID":         m.ReqID,
		"time_begin":    m.TimeBegin,
		"time_end":      m.TimeEnd,
		"time_exec":     m.TimeExec,
		"search_time":   m.SearchTime,
		"response_time": m.ResponseTime,
		"host_name":     m.HostName,
		"remote_addr":   m.RemoteAddr,
		"request": et.Json{
			"end_point": m.EndPoint,
			"method":    m.Method,
			"status":    res.Status,
			"bytes":     res.ContentLength,
			"header":    res.Header,
			"scheme":    m.Scheme,
			"host":      res.Request.Host,
		},
		"memory": et.Json{
			"unity":        "MB",
			"total":        m.MTotal / 1024 / 1024,
			"used":         m.MUsed / 1024 / 1024,
			"free":         m.MFree / 1024 / 1024,
			"percent_free": math.Floor(m.PFree*100) / 100,
		},
		"request_host": et.Json{
			"host":   m.RequestsHost.Tag,
			"day":    m.RequestsHost.Day,
			"hour":   m.RequestsHost.Hour,
			"minute": m.RequestsHost.Minute,
			"second": m.RequestsHost.Seccond,
			"limit":  m.RequestsHost.Limit,
		},
		"requests_endpoint": et.Json{
			"endpoint": m.RequestsEndpoint.Tag,
			"day":      m.RequestsEndpoint.Day,
			"hour":     m.RequestsEndpoint.Hour,
			"minute":   m.RequestsEndpoint.Minute,
			"second":   m.RequestsEndpoint.Seccond,
			"limit":    m.RequestsEndpoint.Limit,
		},
	}

	go event.Action("telemetry", result)

	if m.RequestsHost.Seccond > m.RequestsHost.Limit {
		go event.Action("requests/overflow", result)
	}

	return result
}
