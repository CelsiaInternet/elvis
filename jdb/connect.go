package jdb

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/envar"
	. "github.com/cgalvisleon/elvis/msg"
	_ "github.com/joho/godotenv/autoload"
)

func connect() {
	driver := EnvarStr("", "DB_DRIVE")
	host := EnvarStr("", "DB_HOST")
	port := EnvarInt(5432, "DB_PORT")
	dbname := EnvarStr("", "DB_NAME")
	user := EnvarStr("", "DB_USER")
	password := EnvarStr("", "DB_PASSWORD")

	if driver == "" {
		console.FatalF(ERR_ENV_REQUIRED, "DB_DRIVE")
	}

	if host == "" {
		console.FatalF(ERR_ENV_REQUIRED, "DB_HOST")
	}

	if dbname == "" {
		console.FatalF(ERR_ENV_REQUIRED, "DB_NAME")
	}

	if user == "" {
		console.FatalF(ERR_ENV_REQUIRED, "DB_USER")
	}

	if password == "" {
		console.FatalF(ERR_ENV_REQUIRED, "DB_PASSWORD")
	}

	_, err := Connected(driver, host, port, dbname, user, password)
	if err != nil {
		console.Fatal(err)
	}

	return
}

func Connected(driver, host string, port int, dbname, user, password string) (int, error) {
	url := ""
	switch driver {
	case Postgres:
		url = fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable`, driver, user, password, host, port, dbname)
	case Mysql:
		url = fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s`, user, password, host, port, dbname)
	case Sqlserver:
		url = fmt.Sprintf(`server=%s;user id=%s;password=%s;port=%d;database=%s;`, host, user, password, port, dbname)
	case Firebird:
		url = fmt.Sprintf(`%s/%s@%s;`, user, password, host)
	default:
		panic(NOT_SELECT_DRIVE)
	}

	sqlDB, err := sql.Open(driver, url)
	if err != nil {
		return -1, console.Error(err)
	}

	console.LogKF(driver, "Connected host:%s:%d", host, port)

	if conn == nil {
		conn = &Conn{
			Db: []*Db{},
		}
	}

	idx := len(conn.Db)
	db := &Db{
		Index:  idx,
		Driver: driver,
		Host:   host,
		Port:   port,
		Dbname: dbname,
		User:   user,
		URL:    url,
		Db:     sqlDB,
	}

	conn.Db = append(conn.Db, db)

	return idx, nil
}
