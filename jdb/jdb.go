package jdb

import (
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
)

/**
* Load
* @return *Conn, error
**/
func Load() (*DB, error) {
	conn, err := ConnectTo(et.Json{
		"driver":           envar.GetStr("", "DB_DRIVER"),
		"host":             envar.GetStr("", "DB_HOST"),
		"port":             envar.GetInt(5432, "DB_PORT"),
		"dbname":           envar.GetStr("", "DB_NAME"),
		"user":             envar.GetStr("", "DB_USER"),
		"password":         envar.GetStr("", "DB_PASSWORD"),
		"application_name": envar.GetStr("elvis", "DB_APPLICATION_NAME"),
	})
	if err != nil {
		return nil, err
	}
	conn.UseCore = envar.GetBool(true, "USE_CORE")

	if !conn.UseCore {
		return conn, nil
	}

	err = InitCore(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
