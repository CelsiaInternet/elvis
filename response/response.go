package response

import (
	"encoding/json"
	"net/http"
	"strings"

	e "github.com/cgalvisleon/elvis/json"
	"github.com/go-chi/chi"
)

type Result struct {
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result"`
}

func GetBody(r *http.Request) (e.Json, error) {
	var body e.Json
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return e.Json{}, err
	}
	defer r.Body.Close()

	return body, nil
}

func WriteResponse(w http.ResponseWriter, statusCode int, e []byte) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(e)

	return nil
}

func JSON(w http.ResponseWriter, r *http.Request, statusCode int, dt interface{}) error {
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

	return WriteResponse(w, statusCode, e)
}

func ITEM(w http.ResponseWriter, r *http.Request, statusCode int, dt e.Item) error {
	if &dt == (&e.Item{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	e, err := json.Marshal(dt)
	if err != nil {
		return err
	}

	return WriteResponse(w, statusCode, e)
}

func ITEMS(w http.ResponseWriter, r *http.Request, statusCode int, dt e.Items) error {
	if &dt == (&e.Items{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	e, err := json.Marshal(dt)
	if err != nil {
		return err
	}

	return WriteResponse(w, statusCode, e)
}

func HTTPError(w http.ResponseWriter, r *http.Request, statusCode int, message string) error {
	msg := e.Json{
		"message": message,
	}

	return JSON(w, r, statusCode, msg)
}

func HTTPAlert(w http.ResponseWriter, r *http.Request, message string) error {
	return HTTPError(w, r, http.StatusBadRequest, message)
}

func Stream(w http.ResponseWriter, r *http.Request, statusCode int, dt interface{}) error {
	if dt == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	e, err := json.Marshal(dt)
	if err != nil {
		return err
	}

	WriteResponse(w, statusCode, e)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}

func HTTPApp(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
