package realtime

import (
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
func Load(name string) (*ws.Client, error) {
	if conn != nil {
		return conn, nil
	}

	url := envar.GetStr("", "RT_URL")
	if url == "" {
		return nil, console.NewError(MSG_RT_URL_REQUIRED)
	}

	client, err := ws.NewClient(&ws.ClientConfig{
		ClientId:  utility.UUID(),
		Name:      name,
		Url:       url,
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
