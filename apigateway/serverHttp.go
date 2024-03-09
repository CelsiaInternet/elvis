package apigateway

import (
	"net/http"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/response"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/rs/cors"
)

type HttpServer struct {
	addr    string
	handler http.Handler
	mux     *http.ServeMux
}

func NewHttpServer() *HttpServer {
	// Create a new server
	mux := http.NewServeMux()

	// Handler router
	mux.HandleFunc("/version", version)
	mux.HandleFunc("/", handler)

	port := envar.EnvarInt(3300, "PORT")
	result := &HttpServer{
		addr:    strs.Format(":%d", port),
		handler: cors.AllowAll().Handler(mux),
		mux:     mux,
	}

	return result
}

func version(w http.ResponseWriter, r *http.Request) {
	result := Version()
	response.JSON(w, r, http.StatusOK, result)
}

func handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	proto := r.Proto
	path := r.URL.Path
	rawQuery := r.URL.RawQuery
	query := r.URL.Query()
	requestURI := r.RequestURI
	remoteAddr := r.RemoteAddr
	header := r.Header
	body := r.Body
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	resolve := GetResolve(method, path)

	if resolve != nil {
		url := strs.AppendStr(resolve.Resolve, rawQuery)
		http.Redirect(w, r, url, http.StatusSeeOther)

		/*
			request, err := http.NewRequest(method, url, body)
			if err != nil {
				response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
				return
			}

			request.Header = header
			client := &http.Client{}
			res, err := client.Do(request)
			if err != nil {
				response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			defer res.Body.Close()

			body, err := json.Marshal(res.Body)
			if err != nil {
				response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
				return
			}

			http.Redirect(w, r, url, res.StatusCode)

			w.Header().Set("Response Agent", PackageTitle)
			w.WriteHeader(res.StatusCode)
			w.Write(body)
		*/

		return
	}

	response.JSON(w, r, http.StatusOK, et.Json{
		"method":   method,
		"proto":    proto,
		"path":     path,
		"rawquery": rawQuery,
		"query":    query,
		"uri":      requestURI,
		"remote":   remoteAddr,
		"header":   header,
		"body":     body,
		"host":     host,
		"scheme":   scheme,
		"resolve":  resolve,
	})

	console.Log("handler")
}
