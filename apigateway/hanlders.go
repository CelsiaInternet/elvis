package apigateway

import (
	"io"
	"net/http"
	"time"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/response"
)

func version(w http.ResponseWriter, r *http.Request) {
	result := Version()
	response.JSON(w, r, http.StatusOK, result)
}

func notFounder(w http.ResponseWriter, r *http.Request) {
	response.HTTPError(w, r, http.StatusNotFound, "404 Not Found.")
}

func handlerFn(w http.ResponseWriter, r *http.Request) {
	// Begin telemetry
	telemetry := telemetryNew(r)

	// Get resolute
	resolute := GetResolute(r)

	// Call search time since begin
	telemetry.SearchTime = time.Since(telemetry.TimeBegin)
	telemetry.TimeExec = time.Now()

	if resolute.Resolve == nil {
		conn.http.notFoundHandler(w, r)
		return
	}

	kind := resolute.Resolve.Node.Resolve.ValStr("HTTP", "kind")
	if kind == "HANDLER" {
		handler := handlers[resolute.Resolve.Node._id]
		if handler == nil {
			response.HTTPError(w, r, http.StatusNotFound, "404 Not Found.")
			return
		}

		handler(w, r)
		return
	}

	/*
		if kind == "REST" {

		}

		if kind == "WEBSOCKET" {

		}
	*/

	// http.Redirect(w, r, resolute.URL, http.StatusSeeOther)
	request, err := http.NewRequest(resolute.Method, resolute.URL, resolute.Body)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	request.Header = resolute.Header
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer func() {
		telemetry.EndPoint = resolute.URL
		telemetry.done(res)
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

// Upsert a update or new route
func upsert(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	method := body.Str("method")
	path := body.Str("path")
	resolve := body.Str("resolve")
	kind := body.ValStr("HTTP", "kind")
	stage := body.ValStr("default", "stage")
	packageName := body.Str("package")

	AddRoute(method, path, resolve, kind, stage, packageName)

	response.JSON(w, r, http.StatusOK, et.Json{
		"message": "Router added",
	})
}

// Getall list of routes
func getAll(w http.ResponseWriter, r *http.Request) {
	_pakages, err := et.Marshal(pakages)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, _pakages)
}
