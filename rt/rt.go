package rt

import (
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/utility"
	"github.com/celsiainternet/elvis/ws"
)

var conn *ws.ClientWS

/**
* LoadFrom
* @return erro
**/
func Load() error {
	if conn != nil {
		return nil
	}

	var err error
	name := envar.GetStr("Real Time", "RT_NAME")
	host := envar.GetStr("localhost", "RT_HOST")
	schema := envar.GetStr("ws", "RT_SCHEME")
	path := envar.GetStr("/ws", "RT_PATH")
	conn, err = ws.NewClientWS(utility.UUID(), name, schema, host, path)
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
