package jrpc

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"slices"
	"time"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jtls"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/middleware"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/strs"
)

/**
* GetRouters
* @return et.Items
* @return error
**/
func GetRouters() (et.Items, error) {
	var result = et.Items{Result: []et.Json{}}
	packages, err := getPackages()
	if err != nil {
		return et.Items{}, err
	}

	for _, route := range packages {
		_routes := []et.Json{}
		for k, v := range route.Solvers {
			_routes = append(_routes, et.Json{
				"method":  k,
				"inputs":  v.Inputs,
				"outputs": v.Output,
			})
		}

		result.Result = append(result.Result, et.Json{
			"packageName": route.Name,
			"host":        route.Host,
			"port":        route.Port,
			"count":       len(_routes),
			"routes":      _routes,
		})
		result.Ok = true
		result.Count++
	}

	return result, nil
}

/**
* call
* @param host string, port int, method string, args et.Json
* @return any, error
**/
func call(host string, port int, method string, args et.Json, result any) (*middleware.Metrics, error) {
	metric := middleware.NewRpcMetric(method)
	address := strs.Format(`%s:%d`, host, port)
	pipeHost := envar.GetStr("", "PIPE_HOST")
	pipePort := envar.GetInt(4200, "PIPE_PORT")
	pipeAddress := strs.Format(`%s:%d`, pipeHost, pipePort)
	if pipeAddress == address {
		token := envar.GetStr("", "PIPE_TOKEN")
		args.Set("Authorization", fmt.Sprintf(`Bearer %s`, token))

		pipePath := envar.GetStr("./.keys", "PIPE_PATH")
		conn, err := jtls.Deal(pipePath, host, port, 365*24*time.Hour)
		if err != nil {
			metric.DoneRpc(err.Error())
			return metric, err
		}
		defer conn.Close()

		client := rpc.NewClient(conn)
		defer client.Close()

		err = client.Call(method, args, result)
		if err != nil {
			metric.DoneRpc(err.Error())
			return metric, err
		}

		return metric, nil
	}

	timeOut := 10 * time.Second
	conn, err := net.DialTimeout(
		"tcp",
		address,
		timeOut,
	)

	if err != nil {
		metric.DoneRpc(err.Error())
		return metric, err
	}

	defer conn.Close()

	timeOutRead := time.Duration(envar.GetInt(600, "RPC_TIMEOUT")) * time.Second
	_ = conn.SetDeadline(
		time.Now().Add(timeOutRead),
	)

	client := rpc.NewClient(conn)
	defer client.Close()

	call := client.Go(
		method,
		args,
		result,
		make(chan *rpc.Call, 1),
	)

	select {
	case done := <-call.Done:

		if done.Error != nil {
			metric.DoneRpc(done.Error.Error())
			return metric, done.Error
		}

	case <-time.After(timeOutRead):
		err := errors.New("rpc timeout")
		metric.DoneRpc(err.Error())

		return metric, err
	}

	return metric, nil
}

/**
* Call
* @param method string, args et.Json, result et.JAny
* @return error
**/
func Call(method string, args et.Json) (any, error) {
	solver, err := GetSolver(method)
	if err != nil {
		return nil, err
	}

	if len(solver.Inputs) == 0 {
		// pipeHost path: sin metadata de tipos, decodificar como et.Item por defecto
		var result et.Item
		metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
		if err != nil {
			return nil, err
		}
		metric.DoneRpc(result.ToString())
		return result, nil
	}

	if len(solver.Inputs) != 3 {
		return nil, fmt.Errorf("invalid number of inputs")
	}

	tp := solver.Inputs[2]
	if tp == "et.Json" || tp == "*et.Json" {
		var result et.Json
		metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
		if err != nil {
			return nil, err
		}
		metric.DoneRpc(result.ToString())
		return result, nil
	} else if tp == "et.Item" || tp == "*et.Item" {
		var result et.Item
		metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
		if err != nil {
			return nil, err
		}
		metric.DoneRpc(result.ToString())
		return result, nil
	} else if tp == "et.Items" || tp == "*et.Items" {
		var result et.Items
		metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
		if err != nil {
			return nil, err
		}
		metric.DoneRpc(result.ToString())
		return result, nil
	} else if tp == "et.List" || tp == "*et.List" {
		var result et.List
		metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
		if err != nil {
			return nil, err
		}
		metric.DoneRpc(result.ToString())
		return result, nil
	} else if tp == "et.MapBool" || tp == "*et.MapBool" {
		var result et.MapBool
		metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
		if err != nil {
			return nil, err
		}
		metric.DoneRpc(result.ToString())
		return result, nil
	}

	return nil, fmt.Errorf("invalid type: %s", tp)
}

/**
* CallJsonToHost
* @param method string, host string, port int, args et.Json
* @return et.Json, error
**/
func CallJsonToHost(method, host string, port int, args et.Json) (et.Json, error) {
	var result et.Json
	metric, err := call(host, port, method, args, &result)
	if err != nil {
		return result, err
	}

	metric.DoneRpc(result.ToString())

	return result, nil
}

/**
* CallJson
* @param method string, args et.Json
* @return et.Json, error
**/
func CallJson(method string, args et.Json) (et.Json, error) {
	var result et.Json
	solver, err := GetSolver(method)
	if err != nil {
		return result, err
	}

	metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
	if err != nil {
		return result, err
	}

	metric.DoneRpc(result.ToString())

	return result, nil
}

/**
* CallItem
* @param method string, args et.Json
* @return et.Item, error
**/
func CallItem(method string, args et.Json) (et.Item, error) {
	var result et.Item
	solver, err := GetSolver(method)
	if err != nil {
		return result, err
	}

	metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
	if err != nil {
		return result, err
	}

	metric.DoneRpc(result.ToString())

	return result, nil
}

/**
* CallItems
* @param method string, args et.Json
* @return et.Items, error
**/
func CallItems(method string, args et.Json) (et.Items, error) {
	var result et.Items
	solver, err := GetSolver(method)
	if err != nil {
		return result, err
	}

	metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
	if err != nil {
		return result, err
	}

	metric.DoneRpc(result.ToString())

	return result, nil
}

/**
* CallList
* @param method string, args et.Json
* @return et.List, error
**/
func CallList(method string, args et.Json) (et.List, error) {
	var result et.List
	solver, err := GetSolver(method)
	if err != nil {
		return result, err
	}

	metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
	if err != nil {
		return result, err
	}

	metric.DoneRpc(result.ToString())

	return result, nil
}

/**
* CallPermitios
* @param method string, args et.Json
* @return map[string]bool, error
**/
func CallPermitios(method string, args et.Json) (map[string]bool, error) {
	var result et.MapBool
	solver, err := GetSolver(method)
	if err != nil {
		return result, err
	}

	metric, err := call(solver.Host, solver.Port, solver.Method, args, &result)
	if err != nil {
		return result, err
	}

	metric.DoneRpc(result.ToString())

	return result, nil
}

/**
* DeleteRouters
* @param host string, packageName string
* @return et.Item, error
**/
func DeleteRouters(host, packageName string) (et.Item, error) {
	packages, err := getPackages()
	if err != nil {
		return et.Item{}, err
	}

	idx := slices.IndexFunc(packages, func(e *Package) bool { return e.Host == host && e.Name == packageName })
	if idx == -1 {
		return et.Item{}, logs.Errorm("jrpc", MSG_PACKAGE_NOT_FOUND)
	} else {
		packages = append(packages[:idx], packages[idx+1:]...)
	}

	err = setPackages(packages)
	if err != nil {
		return et.Item{}, err
	}

	return et.Item{
		Ok: true,
		Result: et.Json{
			"message": MSG_PACKAGE_DELETE,
		},
	}, nil
}

/**
* HttpCallRPC
* @param w http.ResponseWriter
* @param r *http.Request
**/
func HttpCallRPC(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	method := body.ValStr("", "method")
	args := body.Json("args")
	result, err := CallItem(method, args)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(map[string]bool{})
	gob.Register(map[string]string{})
	gob.Register(map[string]int{})
	gob.Register(et.Json{})
	gob.Register([]et.Json{})
	gob.Register(et.Item{})
	gob.Register(et.Items{})
	gob.Register(et.List{})
}
