package jdb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
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
	channels    map[string]HandlerListend
	db          *sql.DB
	dm          *sql.DB
}

/**
* Close
* @return error
**/
func (c *DB) Close() error {
	err := c.db.Close()
	if err != nil {
		return err
	}

	return nil
}

/**
* Describe
* @return et.Json
**/
func (c *DB) Describe() et.Json {
	host := strs.Format(`%s:%d`, c.Host, c.Port)
	return et.Json{
		"name":        c.Dbname,
		"description": c.Description,
		"driver":      c.Driver,
		"host":        host,
	}
}

/**
* HealthCheck
* @return bool
**/
func (c *DB) HealthCheck() bool {
	err := c.db.Ping()
	if err != nil {
		return false
	}

	return true
}

/**
* connectTo
* @param driver, chain string
* @return *sql.DB, error
**/
func connectTo(driver, chain string) (*sql.DB, error) {
	db, err := sql.Open(driver, chain)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	maxOpen := envar.GetInt(25, "DB_POOL_MAX_OPEN")
	maxIdle := envar.GetInt(5, "DB_POOL_MAX_IDLE")
	connLifetime := envar.GetInt(900, "DB_POOL_CONN_LIFETIME")
	connIdleTime := envar.GetInt(300, "DB_POOL_CONN_IDLE_TIME")

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(connLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(connIdleTime) * time.Second)

	return db, nil
}

/**
* ExistDatabase
* @param db *DB, name string
* @return bool, error
**/
func existDatabase(db *sql.DB, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
	SELECT 1
	FROM pg_database
	WHERE UPPER(datname) = UPPER($1));`
	rows, err := db.Query(sql, name)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	items := rowsItems(rows)

	if items.Count == 0 {
		return false, nil
	}

	return items.Bool(0, "exists"), nil
}

/**
* CreateDatabase
* @param db *sql.DB, name string
* @return error
**/
func createDatabase(db *sql.DB, name string) error {
	exist, err := existDatabase(db, name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := fmt.Sprintf(`CREATE DATABASE %s;`, name)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	console.LogKF("postgres", `Database %s created`, name)

	return nil
}

/**
* ConnectTo
* @param params et.Json
* @return *DB, error
**/
func ConnectTo(params et.Json) (*DB, error) {
	driver := params.Str("driver")
	if !utility.ValidStr(driver, 0, []string{""}) {
		return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "driver")
	}
	host := params.Str("host")
	if !utility.ValidStr(host, 0, []string{""}) {
		return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "host")
	}
	port := params.Int("port")
	if port == 0 {
		return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "port")
	}
	dbname := params.Str("dbname")

	var connStr string
	switch driver {
	case Postgres:
		user := params.Str("user")
		if !utility.ValidStr(user, 0, []string{""}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "user")
		}
		password := params.Str("password")
		if !utility.ValidStr(password, 4, []string{""}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "password")
		}
		application_name := params.Str("application_name")
		if !utility.ValidStr(password, 4, []string{""}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "application_name")
		}

		connStr = strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, user, password, host, port, "postgres", application_name)
		db, err := connectTo(driver, connStr)
		if err != nil {
			return nil, err
		}

		exist, err := existDatabase(db, dbname)
		if err != nil {
			return nil, err
		}

		if !exist {
			err := createDatabase(db, dbname)
			if err != nil {
				return nil, err
			}
		}

		db.Close()
		connStr = strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, user, password, host, port, dbname, application_name)
	case Oracle:
		user := params.Str("user")
		if !utility.ValidStr(user, 0, []string{""}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "user")
		}
		password := params.Str("password")
		if !utility.ValidStr(password, 4, []string{""}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "password")
		}
		service_name := params.Str("service_name")
		if !utility.ValidStr(service_name, 0, []string{""}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "service_name")
		}
		ssl := params.Str("ssl")
		if !utility.ValidIn(ssl, 0, []string{"TREU", "FALSE", "true", "false"}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "ssl (boolean)")
		}
		sslVerify := params.Str("ssl_verify")
		if !utility.ValidIn(sslVerify, 0, []string{"TREU", "FALSE", "true", "false"}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "ssl_verify (boolean)")
		}

		urlOptions := map[string]string{
			"ssl":        ssl,
			"ssl verify": sslVerify,
		}
		connStr = goOra.BuildUrl(host, port, service_name, user, password, urlOptions)
	case Mysql:
		user := params.Str("user")
		if !utility.ValidStr(user, 0, []string{""}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "user")
		}
		password := params.Str("password")
		if !utility.ValidStr(password, 4, []string{""}) {
			return nil, logs.Errorf("ConnectTo", msg.MSG_ATRIB_REQUIRED, "password")
		}

		connStr = strs.Format(`%s:%s@tcp(%s:%d)/%s`, user, password, host, port, dbname)
	default:
		return nil, logs.Errorm("jdb", msg.NOT_SELECT_DRIVE)
	}

	db, err := connectTo(driver, connStr)
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
		UseCore:    false,
		channels:   make(map[string]HandlerListend),
		db:         db,
	}, nil
}
