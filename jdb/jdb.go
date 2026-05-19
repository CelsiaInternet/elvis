package jdb

import (
	"errors"
	"sync"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
)

const (
	KEY = "_id"
)

var (
	connections = map[string]*DB{}
	mu          sync.Mutex
)

/**
* LoadTo
* @param dbname string
* @return *DB, error
**/
func LoadTo(dbname string) (*DB, error) {
	if dbname == "" {
		return nil, errors.New("dbname is required")
	}

	mu.Lock()
	if conn, ok := connections[dbname]; ok {
		mu.Unlock()
		return conn, nil
	}
	mu.Unlock()

	conn, err := ConnectTo(et.Json{
		"driver":           envar.GetStr("", "DB_DRIVER"),
		"host":             envar.GetStr("", "DB_HOST"),
		"port":             envar.GetInt(5432, "DB_PORT"),
		"dbname":           dbname,
		"user":             envar.GetStr("", "DB_USER"),
		"password":         envar.GetStr("", "DB_PASSWORD"),
		"application_name": envar.GetStr("elvis", "DB_APPLICATION_NAME"),
	})
	if err != nil {
		return nil, err
	}

	conn.UseCore = envar.GetBool(true, "USE_CORE")
	if conn.UseCore {
		err = InitCore(conn)
		if err != nil {
			return nil, err
		}
	}

	mu.Lock()
	connections[dbname] = conn
	mu.Unlock()

	return conn, nil
}

/**
* Load
* @return *DB, error
**/
func Load() (*DB, error) {
	dbname := envar.GetStr("", "DB_NAME")
	return LoadTo(dbname)
}
