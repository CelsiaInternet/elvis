package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
	_ "github.com/lib/pq"
)

func connect() {
	driver := envar.EnvarStr("", "DB_DRIVE")
	host := envar.EnvarStr("", "DB_HOST")
	port := envar.EnvarInt(5432, "DB_PORT")
	dbname := envar.EnvarStr("", "DB_NAME")
	user := envar.EnvarStr("", "DB_USER")
	password := envar.EnvarStr("", "DB_PASSWORD")
	application_name := envar.EnvarStr("elvis", "DB_APPLICATION_NAME")

	if driver == "" {
		console.FatalF(msg.ERR_ENV_REQUIRED, "DB_DRIVE")
	}

	if host == "" {
		console.FatalF(msg.ERR_ENV_REQUIRED, "DB_HOST")
	}

	if dbname == "" {
		console.FatalF(msg.ERR_ENV_REQUIRED, "DB_NAME")
	}

	if user == "" {
		console.FatalF(msg.ERR_ENV_REQUIRED, "DB_USER")
	}

	if password == "" {
		console.FatalF(msg.ERR_ENV_REQUIRED, "DB_PASSWORD")
	}

	var connect *sql.DB
	var connectStr string
	var err error
	connect, connectStr, err = Connected(driver, host, port, dbname, user, password, application_name)
	if err != nil {
		console.Fatal(err)
	}

	err = connect.Ping()
	if err != nil {
		tmp, _, err := Connected(driver, host, port, "postgres", user, password, application_name)
		if err != nil {
			console.Fatal(err)
		}

		err = CreateDatabase(tmp, dbname)
		if err != nil {
			console.Fatal(err)
		}
		defer tmp.Close()

		connect, connectStr, err = Connected(driver, host, port, dbname, user, password, application_name)
		if err != nil {
			console.Fatal(err)
		}
	}

	if conn == nil {
		conn = &Conn{
			Db: []*Db{},
		}
	}

	idx := len(conn.Db)
	db := &Db{
		Index:      idx,
		Driver:     driver,
		Host:       host,
		Port:       port,
		Dbname:     dbname,
		User:       user,
		Connection: connectStr,
		Db:         connect,
	}

	conn.Db = append(conn.Db, db)
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
