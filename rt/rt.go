package rt

import (
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/utility"
	"github.com/celsiainternet/elvis/ws"
)

var conn *ws.Client

/**
* LoadFrom
* @return erro
**/
func Load() error {
	if conn != nil {
		return nil
	}

	var err error
	params := &ws.ClientConfig{
		ClientId: utility.UUID(),
		Name:     envar.GetStr("Real Time", "RT_NAME"),
		Schema:   envar.GetStr("ws", "RT_SCHEME"),
		Host:     envar.GetStr("localhost", "RT_HOST"),
		Path:     envar.GetStr("/ws", "RT_PATH"),
	}
	conn, err = ws.NewClient(params)
	if err != nil {
		return err
	}

	return nil
}

/**
* Close
**/
func Close() {
	conn.Close()
}
