package service

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jrpc"
	"github.com/celsiainternet/elvis/strs"
)

const (
	// Paquetes RPC por defecto de send y otp (strs.DaskSpace del PackageName de cada servicio + ".Services")
	DEFAULT_SEND_RPC = "Service-send.Services"
	DEFAULT_OTP_RPC  = "Service-otp.Services"
)

var (
	// Se reemplazan en las pruebas, que no tienen Redis, NATS ni send/otp
	rpcSolve = jrpc.GetSolver
	rpcCall  = func(solver *jrpc.Solver, args et.Json) (et.Item, error) {
		return jrpc.CallItemToHost(solver.Method, solver.Host, solver.Port, args)
	}
	workPublish = event.Work
)

/**
* sendRpc: Envia por RPC a send (paquete SEND_RPC, por defecto Service-send.Services).
* @param method string, eventName string, data et.Json
* @return et.Json
**/
func sendRpc(method, eventName string, data et.Json) et.Json {
	pkg := envar.GetStrOs(DEFAULT_SEND_RPC, "SEND_RPC")
	return dispatchRpc(pkg, method, eventName, data)
}

/**
* otpRpc: Genera la OTP por RPC en otp (paquete OTP_RPC, por defecto Service-otp.Services).
* @param method string, eventName string, data et.Json
* @return et.Json
**/
func otpRpc(method, eventName string, data et.Json) et.Json {
	pkg := envar.GetStrOs(DEFAULT_OTP_RPC, "OTP_RPC")
	return dispatchRpc(pkg, method, eventName, data)
}

/**
* dispatchRpc: Llama <pkg>.<method> con data y espera la respuesta (RPC_TIMEOUT). Si GetSolver no
* resuelve el metodo (no registrado, Redis sin cargar o nombre invalido) la llamada no se intento,
* asi que publica event.Work(eventName, data) como antes. Si la llamada falla despues de resolverse
* devuelve el error sin reenviar, para no duplicar mensajes.
* Retorno: el result del RPC mas ok, service_id y, si fallo, error con el texto; en el respaldo,
* el sobre del trabajo (_id, event, data, created_at) mas service_id.
* @param pkg string, method string, eventName string, data et.Json
* @return et.Json
**/
func dispatchRpc(pkg, method, eventName string, data et.Json) et.Json {
	serviceId := data.Str("service_id")
	rpcMethod := strs.Format(`%s.%s`, pkg, method)
	solver, err := rpcSolve(rpcMethod)
	if err != nil {
		console.AlertF(`%s: %s; se publica event.Work %s`, rpcMethod, err.Error(), eventName)
		result := workPublish(eventName, data)
		result["service_id"] = serviceId
		return result
	}

	item, err := rpcCall(solver, data)
	return rpcResult(item, err, serviceId)
}

/**
* rpcResult: Arma el retorno de los Send*: result del RPC, ok, service_id y error si fallo.
* @param item et.Item, err error, serviceId string
* @return et.Json
**/
func rpcResult(item et.Item, err error, serviceId string) et.Json {
	result := et.Json{}
	for k, v := range item.Result {
		result[k] = v
	}

	result["ok"] = item.Ok
	result["service_id"] = serviceId
	if err != nil {
		result["ok"] = false
		result["error"] = err.Error()
	}

	return result
}
