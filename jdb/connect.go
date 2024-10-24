package jdb

import (
	"database/sql"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/strs"
	_ "github.com/lib/pq"
)

const Postgres = "postgres"

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

func connect() (*DB, error) {
	driver := envar.GetStr("", "DB_DRIVE")
	host := envar.GetStr("", "DB_HOST")
	port := envar.GetInt(5432, "DB_PORT")
	dbname := envar.GetStr("", "DB_NAME")
	user := envar.GetStr("", "DB_USER")
	password := envar.GetStr("", "DB_PASSWORD")
	application_name := envar.GetStr("elvis", "DB_APPLICATION_NAME")

	if driver == "" {
		return nil, console.PanicF(msg.ERR_ENV_REQUIRED, "DB_DRIVE")
	}

	if host == "" {
		return nil, console.PanicF(msg.ERR_ENV_REQUIRED, "DB_HOST")
	}

	if dbname == "" {
		return nil, console.PanicF(msg.ERR_ENV_REQUIRED, "DB_NAME")
	}

	if user == "" {
		return nil, console.PanicF(msg.ERR_ENV_REQUIRED, "DB_USER")
	}

	if password == "" {
		return nil, console.PanicF(msg.ERR_ENV_REQUIRED, "DB_PASSWORD")
	}

	db, err := ConnectTo(driver, host, port, dbname, user, password, application_name)
	if err != nil {
		return nil, err
	}

	db.UseCore = true

	return db, nil
}

func ConnectTo(driver, host string, port int, dbname, user, password, application_name string) (*DB, error) {
	var connStr string
	switch driver {
	case Postgres:
		connStr = strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, user, password, host, port, dbname, application_name)
	default:
		panic(msg.NOT_SELECT_DRIVE)
	}

	db, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, console.Alert(err.Error())
	}

	console.LogKF(driver, "Connected host:%s:%d", host, port)

	return &DB{
		Driver:     driver,
		Host:       host,
		Port:       port,
		Dbname:     dbname,
		Connection: connStr,
		db:         db,
	}, nil
}
