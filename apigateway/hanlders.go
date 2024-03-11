package apigateway

import (
	"io"
	"net/http"

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
	resolute := NewResolute(r)

	if resolute.Resolve == nil {
		conn.http.notFoundHandler(w, r)
		return
	}

	kind := resolute.Resolve.Node.Resolve.ValStr("HTTP", "kind")
	if kind == "REST" {
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
		defer res.Body.Close()

		for key, value := range res.Header {
			w.Header().Set(key, value[0])
		}
		_, err = io.Copy(w, res.Body)
		if err != nil {
			response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	/*
		if kind == "WEBSOCKET" {
			// TODO

		}
	*/

	http.Redirect(w, r, resolute.URL, http.StatusSeeOther)
}

/*
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

func getAll(w http.ResponseWriter, r *http.Request) {
	_routes, err := et.Marshal(routes)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, _routes)
}
*/
