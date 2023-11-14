package response

import (
	"encoding/json"
	"net/http"
	"strings"

	j "github.com/cgalvisleon/elvis/json"
	"github.com/go-chi/chi"
)

type Result struct {
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result"`
}

func GetBody(r *http.Request) (j.Json, error) {
	var body j.Json
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return j.Json{}, err
	}
	defer r.Body.Close()

	return body, nil
}

func WriteResponse(w http.ResponseWriter, statusCode int, j []byte) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(j)

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

	j, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return WriteResponse(w, statusCode, j)
}

func ITEM(w http.ResponseWriter, r *http.Request, statusCode int, dt j.Item) error {
	if &dt == (&j.Item{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	j, err := json.Marshal(dt)
	if err != nil {
		return err
	}

	return WriteResponse(w, statusCode, j)
}

func ITEMS(w http.ResponseWriter, r *http.Request, statusCode int, dt j.Items) error {
	if &dt == (&j.Items{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	j, err := json.Marshal(dt)
	if err != nil {
		return err
	}

	return WriteResponse(w, statusCode, j)
}

func HTTPError(w http.ResponseWriter, r *http.Request, statusCode int, message string) error {
	msg := j.Json{
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

	j, err := json.Marshal(dt)
	if err != nil {
		return err
	}

	WriteResponse(w, statusCode, j)
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
