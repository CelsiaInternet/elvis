package middleware

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	lg "github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/response"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type ContentLength struct {
	Header int
	Body   int
	Total  int
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
	SizeHeader int
	SizeBody   int
	SizeTotal  int
	Host       string
}

/**
* WriteHeader
* @params statusCode int
**/
func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
	totalHeader := 0
	for key, values := range rw.Header() {
		totalHeader += len(key)
		for _, value := range values {
			totalHeader += len(value) + len(": ") + len("\r\n")
		}
	}
	totalHeader += len("\r\n") * len(rw.Header())
	rw.SizeHeader = totalHeader
	rw.SizeTotal = rw.SizeHeader + rw.SizeBody
}

/**
* Write
* @params b []byte
**/
func (rw *ResponseWriterWrapper) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.SizeBody += size
	rw.SizeTotal = rw.SizeHeader + rw.SizeBody
	return size, err
}

/**
* ContentLength
* @return ContentLength
**/
func (rw *ResponseWriterWrapper) ContentLength() ContentLength {
	totalHeader := 0
	for key, values := range rw.Header() {
		totalHeader += len(key)
		for _, value := range values {
			totalHeader += len(value) + len(": ") + len("\r\n")
		}
	}
	totalHeader += len("\r\n") * len(rw.Header())
	rw.SizeHeader = totalHeader
	rw.SizeTotal = rw.SizeHeader + rw.SizeBody
	return ContentLength{
		Header: rw.SizeHeader,
		Body:   rw.SizeBody,
		Total:  rw.SizeTotal,
	}
}

/**
* headerLength
* @params res *http.Response
* @return int
**/
func headerLength(res *http.Response) int {
	result := 0
	for key, values := range res.Header {
		result += len(key)
		for _, value := range values {
			result += len(value) + len(": ") + len("\r\n")
		}
	}
	result += len("\r\n") * len(res.Header)

	return result
}

/**
* contentLength
* @params res *http.Response
* @return int
**/
func contentLength(res *http.Response) ContentLength {
	result := headerLength(res)

	return ContentLength{
		Header: result,
		Body:   int(res.ContentLength),
		Total:  result + int(res.ContentLength),
	}
}

type Metrics struct {
	ReqID            string
	TimeBegin        time.Time
	TimeEnd          time.Time
	TimeExec         time.Time
	SearchTime       time.Duration
	ResponseTime     time.Duration
	Downtime         time.Duration
	Latency          time.Duration
	StatusCode       int
	Status           string
	ContentLength    ContentLength
	Header           http.Header
	Host             string
	EndPoint         string
	Method           string
	Proto            string
	RemoteAddr       string
	HostName         string
	RequestsHost     Request
	RequestsEndpoint Request
	Scheme           string
	CPUUsage         float64
	MemoryTotal      uint64
	MemoeryUsage     uint64
	MmemoryFree      uint64
}

/**
* NewMetric
* @params r *http.Request
**/
func NewMetric(r *http.Request) *Metrics {
	result := &Metrics{}
	result.TimeBegin = time.Now()
	result.ReqID = utility.UUID()
	result.EndPoint = r.URL.Path
	result.Method = r.Method
	result.Proto = r.Proto
	result.RemoteAddr = r.Header.Get("X-Forwarded-For")
	if result.RemoteAddr == "" {
		result.RemoteAddr = r.Header.Get("X-Real-IP")
	}
	if result.RemoteAddr == "" {
		result.RemoteAddr = r.RemoteAddr
	} else {
		result.RemoteAddr = strs.Split(result.RemoteAddr, ",")[0]
	}
	result.HostName, _ = os.Hostname()
	result.RequestsHost = callRequests(result.HostName)
	result.RequestsEndpoint = callRequests(result.EndPoint)
	result.Scheme = "http"
	if r.TLS != nil {
		result.Scheme = "https"
	}

	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		result.CPUUsage = 0
	}
	result.CPUUsage = percentages[0]

	v, err := mem.VirtualMemory()
	if err != nil {
		result.MemoryTotal = 0
		result.MemoeryUsage = 0
		result.MmemoryFree = 0
	}
	result.MemoryTotal = v.Total
	result.MemoeryUsage = v.Used
	result.MmemoryFree = v.Free

	return result
}

/**
* println
* @return et.Json
**/
func (m *Metrics) println() et.Json {
	w := lg.Color(lg.NMagenta, fmt.Sprintf(" [%s]: ", m.Method))
	lg.CW(w, lg.NCyan, fmt.Sprintf("%s %s", m.EndPoint, m.Proto))
	lg.CW(w, lg.NWhite, fmt.Sprintf(" from %s", m.RemoteAddr))
	if m.StatusCode >= 500 {
		lg.CW(w, lg.NRed, fmt.Sprintf(" - %s", m.Status))
	} else if m.StatusCode >= 400 {
		lg.CW(w, lg.NYellow, fmt.Sprintf(" - %s", m.Status))
	} else if m.StatusCode >= 300 {
		lg.CW(w, lg.NCyan, fmt.Sprintf(" - %s", m.Status))
	} else {
		lg.CW(w, lg.NGreen, fmt.Sprintf(" - %s", m.Status))
	}
	lg.CW(w, lg.NCyan, fmt.Sprintf(" Size: %v%s", m.ContentLength.Total, "KB"))
	lg.CW(w, lg.NWhite, " in ")
	limitLatency := time.Duration(envar.EnvarInt64(500, "LIMIT_LATENCY")) * time.Millisecond
	if m.Latency < limitLatency {
		lg.CW(w, lg.NGreen, "Latency:%s", m.Latency)
	} else if m.Latency < 5*time.Second {
		lg.CW(w, lg.NYellow, "Latency:%s", m.Latency)
	} else {
		lg.CW(w, lg.NRed, "Latency:%s", m.Latency)
	}
	lg.CW(w, lg.NWhite, " Response:%s", m.ResponseTime)
	lg.CW(w, lg.NRed, " Downtime:%s", m.Downtime)
	if m.RequestsHost.Seccond > m.RequestsHost.Limit {
		lg.CW(w, lg.NRed, " - Request S:%vM:%vH:%vL:%v", m.RequestsHost.Seccond, m.RequestsHost.Minute, m.RequestsHost.Hour, m.RequestsHost.Limit)
	} else {
		lg.CW(w, lg.NYellow, " - Request S:%vM:%vH:%vL:%v", m.RequestsHost.Seccond, m.RequestsHost.Minute, m.RequestsHost.Hour, m.RequestsHost.Limit)
	}
	lg.Println(w)

	result := et.Json{
		"reqID":         m.ReqID,
		"time_begin":    m.TimeBegin,
		"time_end":      m.TimeEnd,
		"time_exec":     m.TimeExec,
		"latency":       m.Latency,
		"search_time":   m.SearchTime,
		"response_time": m.ResponseTime,
		"host_name":     m.HostName,
		"remote_addr":   m.RemoteAddr,
		"request": et.Json{
			"end_point": m.EndPoint,
			"method":    m.Method,
			"status":    m.Status,
			"size": et.Json{
				"header": m.ContentLength.Header,
				"body":   m.ContentLength.Body,
			},
			"header": m.Header,
			"scheme": m.Scheme,
			"host":   m.Host,
		},
		"system": et.Json{
			"unity":        "MB",
			"total":        m.MemoryTotal / 1024 / 1024,
			"used":         m.MemoeryUsage / 1024 / 1024,
			"free":         m.MmemoryFree / 1024 / 1024,
			"percent_free": math.Floor(float64(m.MmemoryFree) / float64(m.MemoryTotal)),
			"cpu_usage":    m.CPUUsage,
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

	go event.Log("telemetry", result)

	if m.RequestsHost.Seccond > m.RequestsHost.Limit {
		go event.Log("requests/overflow", result)
	}

	return result
}

/**
* CallExecute
**/
func (m *Metrics) CallExecute() {
	m.SearchTime = time.Since(m.TimeBegin)
	m.TimeExec = time.Now()
}

/**
* Done
* @params res *http.Response
**/
func (m *Metrics) Done(res *http.Response) et.Json {
	m.TimeEnd = time.Now()
	m.ResponseTime = time.Since(m.TimeExec)
	m.Latency = time.Since(m.TimeBegin)
	m.Downtime = m.SearchTime - m.ResponseTime
	m.StatusCode = res.StatusCode
	m.Status = strs.Format(`%d %s`, res.StatusCode, http.StatusText(res.StatusCode))
	m.ContentLength = contentLength(res)
	m.Header = res.Header
	m.Host = res.Request.Host

	return m.println()
}

/**
* DoneRWW
* @params rw *ResponseWriterWrapper
* @params r *http.Request
* @return js.Json
**/
func (m *Metrics) DoneRWW(w *ResponseWriterWrapper, r *http.Request) et.Json {
	m.TimeEnd = time.Now()
	m.ResponseTime = time.Since(m.TimeExec)
	m.Latency = time.Since(m.TimeBegin)
	m.Downtime = m.Latency - m.ResponseTime
	m.StatusCode = w.StatusCode
	m.Status = strs.Format(`%d %s`, w.StatusCode, http.StatusText(w.StatusCode))
	m.ContentLength = w.ContentLength()
	m.Header = w.Header()
	m.Host = w.Host

	return m.println()
}

/**
* Unauthorized
* @params w http.ResponseWriter
* @params r *http.Request
**/
func (m *Metrics) Unauthorized(w http.ResponseWriter, r *http.Request) {
	rw := &ResponseWriterWrapper{ResponseWriter: w, StatusCode: http.StatusUnauthorized, Host: r.Host}
	m.CallExecute()
	response.HTTPError(rw, r, http.StatusUnauthorized, "401 Unauthorized")
	go m.DoneRWW(rw, r)
}

/**
* NotFound
* @params handler http.HandlerFunc
* @params w http.ResponseWriter
* @params r *http.Request
**/
func (m *Metrics) NotFound(handler http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	rw := &ResponseWriterWrapper{ResponseWriter: w, StatusCode: http.StatusNotFound, Host: r.Host}
	m.CallExecute()
	handler(rw, r)
	go m.DoneRWW(rw, r)
}

/**
* Handler
* @params handler http.HandlerFunc
* @params w http.ResponseWriter
* @params r *http.Request
**/
func (m *Metrics) Handler(handler http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	rw := &ResponseWriterWrapper{ResponseWriter: w, StatusCode: http.StatusOK, Host: r.Host}
	isWebSocket := r.Method == http.MethodGet &&
		r.Header.Get("Upgrade") == "websocket" &&
		r.Header.Get("Connection") == "Upgrade" &&
		r.Header.Get("Sec-WebSocket-Key") != ""
	if isWebSocket {
		handler(w, r)
	} else {
		handler(rw, r)
	}

	go m.DoneRWW(rw, r)
}
