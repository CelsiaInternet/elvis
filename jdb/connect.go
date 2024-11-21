package jdb

import (
	"database/sql"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	goOra "github.com/sijms/go-ora/v2"
)

const (
	Postgres = "postgres"
	Oracle   = "oracle"
	Mysql    = "mysql"
)

type DB struct {
	Description string
	Driver      string
	Host        string
	Port        int
	Dbname      string
	Connection  string
	UseCore     bool
	db          *sql.DB
	dm          *sql.DB
}

func (c *DB) Close() error {
	err := c.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *DB) Describe() et.Json {
	host := strs.Format(`%s:%d`, c.Host, c.Port)
	return et.Json{
		"name":        c.Dbname,
		"description": c.Description,
		"driver":      c.Driver,
		"host":        host,
	}
}

func ConnectTo(params et.Json) (*DB, error) {
	driver := params.Str("driver")
	if !utility.ValidStr(driver, 0, []string{""}) {
		return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "driver")
	}
	host := params.Str("host")
	if !utility.ValidStr(host, 0, []string{""}) {
		return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "host")
	}
	port := params.Int("port")
	if port == 0 {
		return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "port")
	}
	dbname := params.Str("dbname")

	var connStr string
	switch driver {
	case Postgres:
		user := params.Str("user")
		if !utility.ValidStr(user, 0, []string{""}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "user")
		}
		password := params.Str("password")
		if !utility.ValidStr(password, 4, []string{""}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "password")
		}
		application_name := params.Str("application_name")
		if !utility.ValidStr(password, 4, []string{""}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "application_name")
		}

		connStr = strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, user, password, host, port, dbname, application_name)
	case Oracle:
		user := params.Str("user")
		if !utility.ValidStr(user, 0, []string{""}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "user")
		}
		password := params.Str("password")
		if !utility.ValidStr(password, 4, []string{""}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "password")
		}
		service_name := params.Str("service_name")
		if !utility.ValidStr(service_name, 0, []string{""}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "service_name")
		}
		ssl := params.Str("ssl")
		if !utility.ValidIn(ssl, 0, []string{"TREU", "FALSE", "true", "false"}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "ssl (boolean)")
		}
		sslVerify := params.Str("ssl_verify")
		if !utility.ValidIn(sslVerify, 0, []string{"TREU", "FALSE", "true", "false"}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "ssl_verify (boolean)")
		}

		urlOptions := map[string]string{
			"ssl":        ssl,
			"ssl verify": sslVerify,
		}
		connStr = goOra.BuildUrl(host, port, service_name, user, password, urlOptions)
	case Mysql:
		user := params.Str("user")
		if !utility.ValidStr(user, 0, []string{""}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "user")
		}
		password := params.Str("password")
		if !utility.ValidStr(password, 4, []string{""}) {
			return nil, logs.Errorf(msg.MSG_ATRIB_REQUIRED, "password")
		}

		connStr = strs.Format(`%s:%s@tcp(%s:%d)/%s`, user, password, host, port, dbname)
	default:
		return nil, logs.Errorm(msg.NOT_SELECT_DRIVE)
	}

	db, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logs.Logf(driver, "Connected host:%s:%d", host, port)

	return &DB{
		Driver:     driver,
		Host:       host,
		Port:       port,
		Dbname:     dbname,
		Connection: connStr,
		db:         db,
		UseCore:    false,
	}, nil
}
