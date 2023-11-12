package jdb

import (
	"database/sql"

	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utility"
)

const Postgres = "postgres"
const Mysql = "mysql"
const Sqlserver = "sqlserver"
const Firebird = "firebird"

type Db struct {
	Index       int
	Description string
	Driver      string
	Host        string
	Port        int
	Dbname      string
	User        string
	URL         string
	token       string
	Db          *sql.DB
}

func (c *Db) Close() error {
	err := c.Db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Db) Describe() Json {
	host := Format(`%s:%d`, c.Host, c.Port)
	return Json{
		"name":        c.Dbname,
		"description": c.Description,
		"driver":      c.Driver,
		"host":        host,
	}
}
