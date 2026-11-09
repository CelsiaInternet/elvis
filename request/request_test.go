package request

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celsiainternet/elvis/et"
)

/**
* echoServer: Servidor de prueba que devuelve el método, el Content-Type y el body recibidos.
* @param t *testing.T
* @return *httptest.Server
**/
func echoServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.NewEncoder(w).Encode(map[string]string{
			"method":      r.Method,
			"contentType": r.Header.Get("Content-Type"),
			"body":        string(body),
		})
	}))
	t.Cleanup(server.Close)
	return server
}

func TestGetWithoutContentType(t *testing.T) {
	server := echoServer(t)

	res, status := Get(server.URL, et.Json{"Authorization": "Bearer x"})
	if !status.Ok {
		t.Fatalf("GET sin Content-Type no llegó al servidor: %d %s", status.Code, status.Message)
	}

	got, err := res.ToJson()
	if err != nil {
		t.Fatal(err)
	}
	if got.Str("method") != http.MethodGet || got.Str("body") != "" {
		t.Errorf("respuesta: %v", got)
	}
}

func TestPostWithoutContentTypeSendsJson(t *testing.T) {
	server := echoServer(t)

	res, status := Post(server.URL, et.Json{}, et.Json{"name": "plan"})
	if !status.Ok {
		t.Fatalf("POST sin Content-Type no llegó al servidor: %d %s", status.Code, status.Message)
	}

	got, err := res.ToJson()
	if err != nil {
		t.Fatal(err)
	}
	if got.Str("contentType") != "application/json" {
		t.Errorf("Content-Type: got %q, want application/json", got.Str("contentType"))
	}
	if got.Str("body") != `{"name":"plan"}` {
		t.Errorf("body: got %q", got.Str("body"))
	}
}

func TestPostKeepsContentType(t *testing.T) {
	server := echoServer(t)

	header := et.Json{"Content-Type": "application/x-www-form-urlencoded"}
	res, status := Post(server.URL, header, et.Json{"name": "plan"})
	if !status.Ok {
		t.Fatalf("POST form: %d %s", status.Code, status.Message)
	}

	got, err := res.ToJson()
	if err != nil {
		t.Fatal(err)
	}
	if got.Str("contentType") != "application/x-www-form-urlencoded" || got.Str("body") != "name=plan" {
		t.Errorf("respuesta: %v", got)
	}
}
