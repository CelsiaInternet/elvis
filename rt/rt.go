package rt

import (
	"net/http"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
	"github.com/celsiainternet/elvis/ws"
)

const ServiceName = "Real Time"

var conn *ws.Client

/**
* LoadFrom
* @return erro
**/
func Load() (*ws.Client, error) {
	if conn != nil {
		return conn, nil
	}

	result, _ := Connect()
	return result, nil
}

/**
* Connect
* @return *ws.Client, error
**/
func Connect() (*ws.Client, error) {
	if conn != nil {
		return conn, nil
	}

	url := envar.GetStr("", "RT_HOST")
	if url == "" {
		return nil, console.NewError(MSG_RT_HOST_REQUIRED)
	}

	token := envar.GetStr("", "RT_AUTH")
	client, err := ws.NewClient(&ws.ClientConfig{
		ClientId: utility.UUID(),
		Name:     envar.GetStr("RealTime", "RT_NAME"),
		Url:      url,
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
		Reconnect: envar.GetInt(3, "RT_RECONCECT"),
	})
	if err != nil {
		return nil, err
	}

	conn = client

	logs.Logf(ServiceName, `Connected host:%s`, url)

	return conn, nil
}

/**
* Close
**/
func Close() {
	if conn == nil {
		return
	}

	conn.Close()
}
