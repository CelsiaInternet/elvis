package service

import (
	"fmt"
	"strconv"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/utility"
)

// Type message
type TpMessage int

const (
	TpTransactional TpMessage = iota
	TpComercial
)

var (
	packageName string
)

func (tp TpMessage) String() string {
	switch tp {
	case TpTransactional:
		return "transactional"
	case TpComercial:
		return "comercial"
	default:
		return "unknown"
	}
}

/**
* SetPackageName
* @param name string
**/
func SetPackageName(name string) {
	packageName = name
}

/**
* GetId
* @param client_id, kind, description string
* @response string
**/
func GetId(client_id, kind, description string) string {
	now := utility.Now()
	result := utility.UUID()
	data := et.Json{
		"created_at":  now,
		"service_id":  result,
		"client_id":   client_id,
		"kind":        kind,
		"description": description,
	}
	event.Work("service/client", data)
	cache.SetH(result, data)

	return result
}

/**
* SetStatus
* @param serviceId string, status et.Json
**/
func SetStatus(serviceId string, status et.Json) {
	event.Work("service/status", et.Json{
		"service_id": serviceId,
		"status":     status,
	})

	cache.SetH(serviceId, status)
}

/**
* GetStatus
* @param serviceId string
* @response et.Json, error
**/
func GetStatus(serviceId string) (et.Json, error) {
	return cache.GetJson(serviceId)
}

/**
* ServiceId
* @param serviceId string, packageName, path string, params ...interface{}
* @return string
**/
func ServiceId(serviceId string, path string, context ...interface{}) string {
	if serviceId == "" || serviceId == "new" {
		serviceId = utility.UUID()
	}

	event.Work("service/trace", et.Json{
		"created_at":   utility.Now(),
		"service_id":   serviceId,
		"package_name": packageName,
		"path":         path,
		"context":      context,
	})

	return serviceId
}

/**
* CalcularDV
* @param nit string
* @return int, error
**/
func CalcularDV(nit string) (int, error) {
	pesos := []int{71, 67, 59, 53, 47, 43, 41, 37, 29, 23, 19, 17, 13, 7, 3}
	suma := 0
	nitLen := len(nit)

	if nitLen > len(pesos) {
		return 0, fmt.Errorf("RUT demasiado largo")
	}

	for i := 0; i < nitLen; i++ {
		digito, err := strconv.Atoi(string(nit[nitLen-1-i]))
		if err != nil {
			return 0, fmt.Errorf("RUT inválido: %v", err)
		}
		suma += digito * pesos[len(pesos)-1-i]
	}

	residuo := suma % 11
	var dv int
	switch residuo {
	case 0:
		dv = 0
	case 1:
		dv = 1
	default:
		dv = 11 - residuo
	}

	return dv, nil
}
