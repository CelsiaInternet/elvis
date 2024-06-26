package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
	_ "github.com/lib/pq"
)

func connect() (*Db, error) {
	driver := envar.EnvarStr("", "DB_DRIVE")
	host := envar.EnvarStr("", "DB_HOST")
	port := envar.EnvarInt(5432, "DB_PORT")
	dbname := envar.EnvarStr("", "DB_NAME")
	user := envar.EnvarStr("", "DB_USER")
	password := envar.EnvarStr("", "DB_PASSWORD")
	application_name := envar.EnvarStr("elvis", "DB_APPLICATION_NAME")

	if driver == "" {
		return nil, console.AlertF(msg.ERR_ENV_REQUIRED, "DB_DRIVE")
	}

	if host == "" {
		return nil, console.AlertF(msg.ERR_ENV_REQUIRED, "DB_HOST")
	}

	if dbname == "" {
		return nil, console.AlertF(msg.ERR_ENV_REQUIRED, "DB_NAME")
	}

	if user == "" {
		return nil, console.AlertF(msg.ERR_ENV_REQUIRED, "DB_USER")
	}

	if password == "" {
		return nil, console.AlertF(msg.ERR_ENV_REQUIRED, "DB_PASSWORD")
	}

	var connect *sql.DB
	var connectStr string
	var err error
	connect, connectStr, err = Connected(driver, host, port, dbname, user, password, application_name)
	if err != nil {
		logs.Fatal(err)
	}

	return &Db{
		Index:      0,
		Driver:     driver,
		Host:       host,
		Port:       port,
		Dbname:     dbname,
		User:       user,
		Connection: connectStr,
		Db:         connect,
	}, nil
}

func Connected(driver, host string, port int, dbname, user, password, application_name string) (*sql.DB, string, error) {
	var connStr string
	switch driver {
	case Postgres:
		connStr = strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, user, password, host, port, dbname, application_name)
	case Mysql:
		connStr = strs.Format(`%s:%s@tcp(%s:%d)/%s`, user, password, host, port, dbname)
	case Sqlserver:
		connStr = strs.Format(`server=%s;user id=%s;password=%s;port=%d;database=%s;`, host, user, password, port, dbname)
	case Firebird:
		connStr = strs.Format(`%s/%s@%s;`, user, password, host)
	default:
		panic(msg.NOT_SELECT_DRIVE)
	}

	result, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, "", console.Alert(err.Error())
	}

	console.LogKF(driver, "Connected host:%s:%d", host, port)

	return result, connStr, nil
}
