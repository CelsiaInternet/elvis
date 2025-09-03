package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	lg "github.com/celsiainternet/elvis/stdrout"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
)

var (
	hostName, _ = os.Hostname()
	serviceName = "telemetry"
)

const (
	TELEMETRY                = "telemetry"
	TELEMETRY_LOG            = "telemetry:log"
	TELEMETRY_OVERFLOW       = "telemetry:overflow"
	TELEMETRY_TOKEN_LAST_USE = "telemetry:token:last_use"
	STATUS_PENDING           = "pending"
	STATUS_OPERATIVE         = "operative"
	STATUS_FAILED            = "failed"
)

type Result struct {
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result"`
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	Size       int
	StatusCode int
}

/**
* Write
* @params b []byte
**/
func (rw *ResponseWriterWrapper) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.Size += size
	return size, err
}

/**
* SetServiceName
* @params name string
**/
func SetServiceName(name string) {
	serviceName = name
}

type Metrics struct {
	TimeStamp    time.Time     `json:"timestamp"`
	ServiceName  string        `json:"service_name"`
	ServiceId    string        `json:"service_id"`
	RemoteAddr   string        `json:"remote_addr"`
	Token        string        `json:"token"`
	Scheme       string        `json:"scheme"`
	Host         string        `json:"host"`
	Method       string        `json:"method"`
	Path         string        `json:"path"`
	StatusCode   int           `json:"status_code"`
	ResponseSize int           `json:"response_size"`
	SearchTime   time.Duration `json:"search_time"`
	ResponseTime time.Duration `json:"response_time"`
	Latency      time.Duration `json:"latency"`
	key          string
	mark         time.Time
	metrics      Telemetry
}

/**
* ToJson
* @return et.Json
**/
func (m *Metrics) ToJson() et.Json {
	return et.Json{
		"timestamp":     strs.FormatDateTime("02/01/2006 03:04:05 PM", m.TimeStamp),
		"service_id":    m.ServiceId,
		"remote_addr":   m.RemoteAddr,
		"token":         m.Token,
		"scheme":        m.Scheme,
		"host":          m.Host,
		"method":        m.Method,
		"path":          m.Path,
		"status_code":   m.StatusCode,
		"search_time":   m.SearchTime,
		"response_time": m.ResponseTime,
		"latency":       m.Latency,
		"response_size": m.ResponseSize,
		"metric":        m.metrics.ToJson(),
	}
}

type Telemetry struct {
	TimeStamp         string
	ServiceName       string
	Method            string
	Path              string
	RequestsPerSecond int64
	RequestsPerMinute int64
	RequestsPerHour   int64
	RequestsPerDay    int64
	RequestsLimit     int64
}

/**
* ToJson
* @return et.Json
**/
func (m *Telemetry) ToJson() et.Json {
	return et.Json{
		"timestamp":           m.TimeStamp,
		"method":              m.Method,
		"path":                m.Path,
		"service_name":        m.ServiceName,
		"requests_per_second": m.RequestsPerSecond,
		"requests_per_minute": m.RequestsPerMinute,
		"requests_per_hour":   m.RequestsPerHour,
		"requests_per_day":    m.RequestsPerDay,
		"requests_limit":      m.RequestsLimit,
	}
}

/**
* PushTelemetry
* @param data et.Json
**/
func PushTelemetry(data et.Json) {
	go event.Publish(TELEMETRY, data)
}

/**
* PushTelemetryLog
* @param data string
**/
func PushTelemetryLog(data string) {
	go event.Publish(TELEMETRY_LOG, et.Json{
		"log": data,
	})
}

/**
* PushTelemetryOverflow
* @param data et.Json
**/
func PushTelemetryOverflow(data et.Json) {
	go event.Publish(TELEMETRY_OVERFLOW, data)
}

/**
* TokenLastUse
* @param data et.Json
**/
func PushTokenLastUse(data et.Json) {
	go event.Publish(TELEMETRY_TOKEN_LAST_USE, data)
}

/**
* NewMetric
* @params r *http.Request
* @return *Metrics
**/
func NewMetric(r *http.Request) *Metrics {
	remoteAddr := r.RemoteAddr
	if remoteAddr == "" {
		remoteAddr = r.Header.Get("Origin")
	}
	if remoteAddr == "" {
		remoteAddr = r.Header.Get("X-Forwarded-For")
	}
	if remoteAddr == "" {
		remoteAddr = r.Header.Get("X-Real-IP")
	}
	if remoteAddr != "" {
		remoteAddr = strs.Split(remoteAddr, ",")[0]
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	serviceId := r.Header.Get("ServiceId")
	if serviceId == "" {
		serviceId = utility.UUID()
		r.Header.Set("ServiceId", serviceId)
	}

	token := r.Header.Get("Authorization")
	result := &Metrics{
		TimeStamp:   timezone.NowTime(),
		ServiceName: serviceName,
		ServiceId:   serviceId,
		RemoteAddr:  remoteAddr,
		Token:       token,
		Host:        hostName,
		Method:      r.Method,
		Path:        r.URL.Path,
		Scheme:      scheme,
		mark:        timezone.NowTime(),
		key:         strs.Format(`%s:%s`, r.Method, r.URL.Path),
	}

	return result
}

/**
* NewRpcMetric
* @params method string
* @return *Metrics
**/
func NewRpcMetric(method string) *Metrics {
	scheme := "rpc"
	result := &Metrics{
		TimeStamp:   timezone.NowTime(),
		ServiceName: serviceName,
		ServiceId:   utility.UUID(),
		Path:        method,
		Method:      strs.Uppcase(scheme),
		Scheme:      scheme,
		mark:        timezone.NowTime(),
		key:         strs.Format(`%s:%s`, strs.Uppcase(scheme), method),
	}

	return result
}

/**
* setRequest
* @params remove bool
**/
func (m *Metrics) setRequest(remove bool) {
	m.key = fmt.Sprintf(`%s:%s`, m.Method, m.Path)
	if remove {
		cache.LRem("telemetry:requests", m.key)
	} else {
		cache.LPush("telemetry:requests", m.key)
	}
}

/**
* SetPath
* @params val string
**/
func (m *Metrics) SetPath(val string) {
	if val == "" {
		return
	}

	m.Path = val
	m.setRequest(false)
}

/**
* CallSearchTime
**/
func (m *Metrics) CallSearchTime() {
	m.SearchTime = time.Since(m.mark)
	m.mark = timezone.NowTime()
}

/**
* CallResponseTime
**/
func (m *Metrics) CallResponseTime() {
	m.ResponseTime = time.Since(m.mark)
	m.mark = timezone.NowTime()
}

/**
* CallLatency
**/
func (m *Metrics) CallLatency() {
	m.Latency = time.Since(m.TimeStamp)
}

/**
* CallMetrics
* @return Telemetry
**/
func (m *Metrics) CallMetrics() Telemetry {
	timeNow := timezone.NowTime()
	date := timeNow.Format("2006-01-02")
	hour := timeNow.Format("2006-01-02-15")
	minute := timeNow.Format("2006-01-02-15:04")
	second := timeNow.Format("2006-01-02-15:04:05")

	return Telemetry{
		TimeStamp:         date,
		ServiceName:       serviceName,
		Method:            m.Method,
		Path:              m.Path,
		RequestsPerSecond: cache.Incr(cache.GenKey(m.key, second), 2*time.Second),
		RequestsPerMinute: cache.Incr(cache.GenKey(m.key, minute), 1*time.Minute+1*time.Second),
		RequestsPerHour:   cache.Incr(cache.GenKey(m.key, hour), 1*time.Hour+1*time.Second),
		RequestsPerDay:    cache.Incr(cache.GenKey(m.key, date), 24*time.Hour+1*time.Second),
		RequestsLimit:     envar.GetInt64(400, "LIMIT_REQUESTS"),
	}
}

/**
* println
* @return et.Json
**/
func (m *Metrics) println() et.Json {
	w := lg.Color(lg.NMagenta, " [%s] ", m.Method)
	lg.CW(w, lg.NCyan, "%s", m.Path)
	lg.CW(w, lg.NWhite, " from:%s", m.RemoteAddr)
	if m.StatusCode >= 500 {
		lg.CW(w, lg.NRed, " - %s", http.StatusText(m.StatusCode))
	} else if m.StatusCode >= 400 {
		lg.CW(w, lg.NYellow, " - %s", http.StatusText(m.StatusCode))
	} else if m.StatusCode >= 300 {
		lg.CW(w, lg.NCyan, " - %s", http.StatusText(m.StatusCode))
	} else {
		lg.CW(w, lg.NGreen, " - %s", http.StatusText(m.StatusCode))
	}
	size := float64(m.ResponseSize) / 1024
	lg.CW(w, lg.NCyan, ` Size:%.2f%s`, size, "KB")
	lg.CW(w, lg.NWhite, " in ")
	limitLatency := time.Duration(envar.GetInt64(1000, "LIMIT_LATENCY")) * time.Millisecond
	if m.Latency < limitLatency {
		lg.CW(w, lg.NGreen, " Latency:%s", m.Latency)
	} else if m.Latency < 5*time.Second {
		lg.CW(w, lg.NYellow, " Latency:%s", m.Latency)
	} else {
		lg.CW(w, lg.NRed, " Latency:%s", m.Latency)
	}
	lg.CW(w, lg.NWhite, " Response:%s", m.ResponseTime)
	m.metrics = m.CallMetrics()
	if m.metrics.RequestsPerSecond > m.metrics.RequestsLimit {
		lg.CW(w, lg.NRed, " - Request:S:%vM:%vH:%vD:%vL:%v", m.metrics.RequestsPerSecond, m.metrics.RequestsPerMinute, m.metrics.RequestsPerHour, m.metrics.RequestsPerDay, m.metrics.RequestsLimit)
	} else if m.metrics.RequestsPerSecond > int64(float64(m.metrics.RequestsLimit)*0.6) {
		lg.CW(w, lg.NYellow, " - Request:S:%vM:%vH:%vD:%vL:%v", m.metrics.RequestsPerSecond, m.metrics.RequestsPerMinute, m.metrics.RequestsPerHour, m.metrics.RequestsPerDay, m.metrics.RequestsLimit)
	} else {
		lg.CW(w, lg.NGreen, " - Request:S:%vM:%vH:%vD:%vL:%v", m.metrics.RequestsPerSecond, m.metrics.RequestsPerMinute, m.metrics.RequestsPerHour, m.metrics.RequestsPerDay, m.metrics.RequestsLimit)
	}
	lg.CW(w, lg.NMagenta, " [ServiceId]:%s", m.ServiceId)
	lg.Println(w)

	m.setRequest(true)
	PushTelemetryLog(w.String())

	return m.ToJson()
}

/**
* telemetry
* @return et.Json
**/
func (m *Metrics) telemetry() et.Json {
	result := m.ToJson()
	PushTelemetry(result)
	if m.metrics.RequestsPerSecond > m.metrics.RequestsLimit {
		PushTelemetryOverflow(m.metrics.ToJson())
	}

	return result
}

/**
* DoneHTTP
* @params rw *ResponseWriterWrapper
* @return et.Json
**/
func (m *Metrics) DoneHTTP(rw *ResponseWriterWrapper) et.Json {
	m.StatusCode = rw.StatusCode
	m.ResponseSize = rw.Size
	m.CallResponseTime()
	m.CallLatency()

	return m.println()
}

/**
* Done
* @params rw *ResponseWriterWrapper
**/
func (m *Metrics) DoneTelemetry(rw *ResponseWriterWrapper) et.Json {
	m.DoneHTTP(rw)
	return m.telemetry()
}

/**
* DoneRpc
* @params r et.Json
* @return et.Json
**/
func (m *Metrics) DoneRpc(r any) et.Json {
	str, ok := r.(string)
	if !ok {
		m.ResponseSize = 0
	} else {
		m.ResponseSize = len(str)
	}
	m.StatusCode = http.StatusOK
	m.CallResponseTime()
	m.CallLatency()
	m.println()

	return m.telemetry()
}

/**
* WriteResponse
* @params w http.ResponseWriter, r *http.Request, statusCode int, e []byte
**/
func (m *Metrics) WriteResponse(w http.ResponseWriter, r *http.Request, statusCode int, e []byte) error {
	rw := &ResponseWriterWrapper{ResponseWriter: w, StatusCode: statusCode}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(statusCode)
	rw.Write(e)

	m.DoneHTTP(rw)
	return nil
}

/**
* JSON
* @params w http.ResponseWriter, r *http.Request, statusCode int, dt interface{}
**/
func (m *Metrics) JSON(w http.ResponseWriter, r *http.Request, statusCode int, dt interface{}) error {
	if dt == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	result := Result{
		Ok:     http.StatusOK == statusCode,
		Result: dt,
	}

	e, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return m.WriteResponse(w, r, statusCode, e)
}

/**
* ITEM
* @params w http.ResponseWriter, r *http.Request, statusCode int, dt et.Item
**/
func (m *Metrics) ITEM(w http.ResponseWriter, r *http.Request, statusCode int, dt et.Item) error {
	if &dt == (&et.Item{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	e, err := json.Marshal(dt)
	if err != nil {
		return err
	}

	return m.WriteResponse(w, r, statusCode, e)
}

/**
* ITEMS
* @params w http.ResponseWriter, r *http.Request, statusCode int, dt et.Items
**/
func (m *Metrics) ITEMS(w http.ResponseWriter, r *http.Request, statusCode int, dt et.Items) error {
	if &dt == (&et.Items{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	e, err := json.Marshal(dt)
	if err != nil {
		return err
	}

	return m.WriteResponse(w, r, statusCode, e)
}

/**
* HTTPError
* @params w http.ResponseWriter, r *http.Request, statusCode int, message string
**/
func (m *Metrics) HTTPError(w http.ResponseWriter, r *http.Request, statusCode int, message string) error {
	msg := et.Json{
		"message": message,
	}

	return m.JSON(w, r, statusCode, msg)
}

/**
* Unauthorized
* @params w http.ResponseWriter, r *http.Request
**/
func (m *Metrics) Unauthorized(w http.ResponseWriter, r *http.Request) {
	m.HTTPError(w, r, http.StatusUnauthorized, "401 Unauthorized")
}
