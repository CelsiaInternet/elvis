package gateway

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/middleware"
	"github.com/cgalvisleon/elvis/response"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
}

/**
* WriteHeader
* @params statusCode int
**/
func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

/**
* Write
* @params b []byte
**/
func (rw *ResponseWriterWrapper) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

/**
* version
* @params w http.ResponseWriter
* @params r *http.Request
**/
func version(w http.ResponseWriter, r *http.Request) {
	result := Version()
	response.JSON(w, r, http.StatusOK, result)
}

/**
* notFounder
* @params w http.ResponseWriter
* @params r *http.Request
**/
func notFounder(w http.ResponseWriter, r *http.Request) {
	result := et.Json{
		"message": "404 Not Found.",
		"route":   r.RequestURI,
	}
	response.JSON(w, r, http.StatusNotFound, result)
}

/**
* upsert
* @params w http.ResponseWriter
* @params r *http.Request
**/
func upsert(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	method := body.Str("method")
	path := body.Str("path")
	resolve := body.Str("resolve")
	kind := body.ValStr("HTTP", "kind")
	stage := body.ValStr("default", "stage")
	packageName := body.Str("package")

	conn.http.AddRoute(method, path, resolve, kind, stage, packageName)

	response.JSON(w, r, http.StatusOK, et.Json{
		"message": "Router added",
	})
}

/**
* getAll
* @params w http.ResponseWriter
* @params r *http.Request
**/
func getAll(w http.ResponseWriter, r *http.Request) {
	_pakages, err := et.Marshal(conn.http.pakages)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, _pakages)
}

/**
* handlerFn
* @params w http.ResponseWriter
* @params r *http.Request
**/
func handlerFn(w http.ResponseWriter, r *http.Request) {
	finalHandler := http.HandlerFunc(handlerExec)
	middleware.Authorization(finalHandler).ServeHTTP(w, r)
}

func handlerExec(w http.ResponseWriter, r *http.Request) {
	// Begin telemetry
	metric := NewMetric(r)

	// Get resolute
	resolute := GetResolute(r)

	// Call search time since begin
	metric.SearchTime = time.Since(metric.TimeBegin)
	metric.TimeExec = time.Now()

	handlerExec := func(handler http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
		rw := &ResponseWriterWrapper{ResponseWriter: w}
		handler(rw, r)
		metric.ContentLength = int64(rw.size)
		metric.summary(r)
	}

	// If not found
	if resolute.Resolve == nil || resolute.URL == "" {
		metric.Downtime = time.Since(metric.TimeBegin)
		metric.NotFount = true
		r.RequestURI = fmt.Sprintf(`%s://%s%s`, resolute.Scheme, resolute.Host, resolute.Path)
		handlerExec(conn.http.notFoundHandler, w, r)
		return
	}

	// If HandlerFunc is handler
	kind := resolute.Resolve.Route.Resolve.ValStr("HTTP", "kind")
	if kind == HANDLER {
		metric.Downtime = time.Since(metric.TimeBegin)
		handler := conn.http.handlers[resolute.Resolve.Route._id]
		if handler == nil {
			metric.NotFount = true
			handlerExec(conn.http.notFoundHandler, w, r)
			return
		}

		handlerExec(handler, w, r)
		return
	}

	// If REST is handler
	request, err := http.NewRequest(resolute.Method, resolute.URL, resolute.Body)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	metric.Downtime = time.Since(metric.TimeBegin)
	request.Header = resolute.Header
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadGateway, err.Error())
		return
	}

	defer func() {
		go metric.done(res)
		res.Body.Close()
	}()

	for key, value := range res.Header {
		w.Header().Set(key, value[0])
	}
	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
	}
}
