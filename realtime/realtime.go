package realtime

import (
	"net/http"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/utility"
	"github.com/celsiainternet/elvis/ws"
)

const ServiceName = "Real Time"

var conn *ws.Client

/**
* Load
* @return erro
**/
func Load() (*ws.Client, error) {
	if conn != nil {
		return conn, nil
	}

	url := envar.GetStr("", "RT_HOST")
	if url == "" {
		return nil, console.NewError(MSG_RT_HOST_REQUIRED)
	}

	token := envar.GetStr("", "RT_AUTH")
	if token == "" {
		return nil, console.NewError(MSG_RT_AUTH_REQUIRED)
	}

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
