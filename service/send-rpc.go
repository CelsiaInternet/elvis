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
* @return et.Item, error
**/
func sendRpc(method, eventName string, data et.Json) (et.Item, error) {
	pkg := envar.GetStrOs(DEFAULT_SEND_RPC, "SEND_RPC")
	return dispatchRpc(pkg, method, eventName, data)
}

/**
* otpRpc: Genera la OTP por RPC en otp (paquete OTP_RPC, por defecto Service-otp.Services).
* @param method string, eventName string, data et.Json
* @return et.Item, error
**/
func otpRpc(method, eventName string, data et.Json) (et.Item, error) {
	pkg := envar.GetStrOs(DEFAULT_OTP_RPC, "OTP_RPC")
	return dispatchRpc(pkg, method, eventName, data)
}

/**
* dispatchRpc: Llama <pkg>.<method> por RPC; si no resuelve el método publica event.Work.
* @param pkg string, method string, eventName string, data et.Json
* @return et.Item, error
**/
func dispatchRpc(pkg, method, eventName string, data et.Json) (et.Item, error) {
	serviceId := data.Str("service_id")
	rpcMethod := strs.Format(`%s.%s`, pkg, method)
	solver, err := rpcSolve(rpcMethod)
	if err != nil {
		console.AlertF(`%s: %s; se publica event.Work %s`, rpcMethod, err.Error(), eventName)
		result := workPublish(eventName, data)
		result["service_id"] = serviceId
		return et.Item{Ok: true, Result: result}, nil
	}

	return rpcCall(solver, data)
}
