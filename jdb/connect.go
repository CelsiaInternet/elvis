package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
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
	lastcomand  int64
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
	driver := envar.EnvarStr("", "DB_DRIVE")
	host := envar.EnvarStr("", "DB_HOST")
	port := envar.EnvarInt(5432, "DB_PORT")
	dbname := envar.EnvarStr("", "DB_NAME")
	user := envar.EnvarStr("", "DB_USER")
	password := envar.EnvarStr("", "DB_PASSWORD")
	application_name := envar.EnvarStr("elvis", "DB_APPLICATION_NAME")

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
