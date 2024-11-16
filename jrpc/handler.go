package jrpc

import (
	"net/http"
	"net/rpc"

	"github.com/celsiainternet/elvis/et"
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
	routes, err := getRouters()
	if err != nil {
		return et.Items{}, err
	}

	for _, route := range routes {
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
* Call
* @param method string
* @param data et.Json
* @return et.Item
* @return error
**/
func Call(method string, data et.Json) (et.Item, error) {
	metric := middleware.NewRpcMetric(method)
	solver, err := GetSolver(method)
	if err != nil {
		return et.Item{}, err
	}

	if solver == nil {
		return et.Item{}, logs.NewErrorf(ERR_METHOD_NOT_FOUND, method)
	}

	address := strs.Format(`%s:%d`, solver.Host, solver.Port)
	metric.CallSearchTime()
	metric.ClientIP = address

	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return et.Item{}, err
	}
	defer client.Close()

	result := et.Item{}
	err = client.Call(method, data, &result)
	if err != nil {
		return et.Item{}, err
	}

	metric.DoneRpc(result)

	return result, nil
}

/**
* HttpCallRPC
* @param w http.ResponseWriter
* @param r *http.Request
**/
func HttpCallRPC(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	method := body.ValStr("", "method")
	data := body.Json("data")
	result, err := Call(method, data)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
	}

	response.JSON(w, r, http.StatusOK, result)
}
