package service

import (
	"errors"
	"testing"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jrpc"
)

/* fakeDispatch reemplaza el RPC y el worker; devuelve lo que se llamo */
func fakeDispatch(t *testing.T, solved bool) (*string, *string, *et.Json) {
	t.Helper()
	var method, published string
	var sent et.Json
	oldSolve, oldCall, oldWork := rpcSolve, rpcCall, workPublish
	t.Cleanup(func() { rpcSolve, rpcCall, workPublish = oldSolve, oldCall, oldWork })

	rpcSolve = func(m string) (*jrpc.Solver, error) {
		method = m
		if !solved {
			return nil, errors.New("sin solver")
		}
		return &jrpc.Solver{Method: m}, nil
	}
	rpcCall = func(solver *jrpc.Solver, args et.Json) (et.Item, error) {
		sent = args
		return et.Item{Ok: true, Result: et.Json{"_id": "err-1"}}, nil
	}
	workPublish = func(name string, data et.Json) et.Json {
		published = name
		sent = data
		return et.Json{"_id": "work-1", "event": name, "data": data}
	}
	return &method, &published, &sent
}

func supportReport() event.ErrorReport {
	return event.ErrorReport{App: "app.test", Service: "test", Message: "fallo", OwnerKind: "user", OwnerId: "u1"}
}

func TestRegisterErrorRpc(t *testing.T) {
	method, published, sent := fakeDispatch(t, true)
	t.Setenv("CATALOG_RPC", "")

	item, err := RegisterError(supportReport())
	if err != nil {
		t.Fatal(err)
	}
	if *method != "catalog.Services.RegisterError" || *published != "" {
		t.Errorf("metodo %q, publicado %q", *method, *published)
	}
	if item.Result.Str("_id") != "err-1" || sent.Str("message") != "fallo" || sent.Str("level") != event.ERROR_LEVEL_ERROR {
		t.Errorf("respuesta %v, enviado %v", item, *sent)
	}
}

func TestRegisterErrorWorkFallback(t *testing.T) {
	method, published, sent := fakeDispatch(t, false)
	t.Setenv("CATALOG_RPC", "otro.Services")

	item, err := RegisterError(supportReport())
	if err != nil {
		t.Fatal(err)
	}
	if *method != "otro.Services.RegisterError" || *published != EVENT_SUPPORT_ERROR {
		t.Errorf("metodo %q, publicado %q", *method, *published)
	}
	if !item.Ok || item.Result.Str("_id") != "work-1" || sent.Str("owner_id") != "u1" {
		t.Errorf("respaldo %v, enviado %v", item, *sent)
	}
}

func TestRegisterErrorInvalid(t *testing.T) {
	method, _, _ := fakeDispatch(t, true)
	report := supportReport()
	report.OwnerId = ""
	if _, err := RegisterError(report); err == nil {
		t.Error("acepto un reporte sin owner_id")
	}
	if *method != "" {
		t.Error("llamo al RPC con un reporte invalido")
	}
}
