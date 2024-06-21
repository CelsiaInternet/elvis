package jdb

import (
	"database/sql"
	"strings"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
)

/**
* Data Definition Language
**/

// SQLQuote quote SQL
func SQLQuote(sql string) string {
	sql = strings.TrimSpace(sql)

	result := strs.Replace(sql, `'`, `"`)
	result = strs.Trim(result)

	return result
}

// SQLDDL SQL Data Definition Language
func SQLDDL(sql string, args ...any) string {
	sql = strings.TrimSpace(sql)

	for i, arg := range args {
		old := strs.Format(`$%d`, i+1)
		new := strs.Format(`%v`, arg)
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

// SQLParse SQL Parse
func SQLParse(sql string, args ...any) string {
	for i := range args {
		old := strs.Format(`$%d`, i+1)
		new := strs.Format(`{$%d}`, i+1)
		sql = strings.ReplaceAll(sql, old, new)
	}

	for i, arg := range args {
		old := strs.Format(`{$%d}`, i+1)
		new := strs.Format(`%v`, et.Unquote(arg))
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

/**
* DBQDDL
**/

// DBQUERY database query
func DBQuery(db *sql.DB, sql string, args ...any) (et.Items, error) {
	if db == nil {
		return et.Items{}, console.AlertF(msg.ERR_COMM)
	}

	sql = SQLParse(sql, args...)
	rows, err := db.Query(sql)
	if err != nil {
		return et.Items{}, console.ErrorF(msg.ERR_SQL, err.Error(), sql)
	}
	defer rows.Close()

	items := rowsItems(rows)

	return items, nil
}

func DBQueryOne(db *sql.DB, sql string, args ...any) (et.Item, error) {
	if db == nil {
		return et.Item{}, console.AlertF(msg.ERR_COMM)
	}

	sql = SQLParse(sql, args...)
	items, err := DBQuery(db, sql, args...)
	if err != nil {
		return et.Item{}, err
	}

	if items.Count == 0 {
		return et.Item{
			Ok:     false,
			Result: et.Json{},
		}, nil
	}

	return et.Item{
		Ok:     items.Ok,
		Result: items.Result[0],
	}, nil
}

func IDXQuery(index int, sql string, args ...any) (et.Items, error) {
	if conn == nil || len(conn.Db) == 0 || conn.Db[index].Db == nil {
		return et.Items{}, console.AlertF(msg.ERR_COMM)
	}

	db := conn.Db[index].Db
	return DBQuery(db, sql, args...)
}

func IDXQueryOne(index int, sql string, args ...any) (et.Item, error) {
	if conn == nil || len(conn.Db) == 0 || conn.Db[index].Db == nil {
		return et.Item{}, console.AlertF(msg.ERR_COMM)
	}

	db := conn.Db[index].Db
	items, err := DBQuery(db, sql, args...)
	if err != nil {
		return et.Item{}, err
	}

	if items.Count == 0 {
		return et.Item{
			Ok:     false,
			Result: et.Json{},
		}, nil
	}

	return et.Item{
		Ok:     items.Ok,
		Result: items.Result[0],
	}, nil
}

func IDXQueryCount(index int, sql string, args ...any) int {
	item, err := IDXQueryOne(index, sql, args...)
	if err != nil {
		return -1
	}

	return item.Int("count")
}

func IDXQueryAtrib(index int, sql, atrib string, args ...any) (et.Items, error) {
	if conn == nil || len(conn.Db) == 0 || conn.Db[index].Db == nil {
		return et.Items{}, console.AlertF(msg.ERR_COMM)
	}

	sql = SQLParse(sql, args...)
	db := conn.Db[index].Db
	rows, err := db.Query(sql)
	if err != nil {
		return et.Items{}, console.ErrorF(msg.ERR_SQL, err.Error(), sql)
	}
	defer rows.Close()

	items := atribItems(rows, atrib)

	return items, nil
}

func IDXQueryAtribOne(index int, sql, atrib string, args ...any) (et.Item, error) {
	items, err := IDXQueryAtrib(index, sql, atrib, args...)
	if err != nil {
		return et.Item{}, err
	}

	if items.Count == 0 {
		return et.Item{
			Ok:     false,
			Result: et.Json{},
		}, nil
	}

	return et.Item{
		Ok:     items.Ok,
		Result: items.Result[0],
	}, nil
}

func IDXQueryData(index int, sql string, args ...any) (et.Items, error) {
	return IDXQueryAtrib(index, sql, "_data", args...)
}

func IDXQueryDataOne(index int, sql string, args ...any) (et.Item, error) {
	return IDXQueryAtribOne(index, sql, "_data", args...)
}

/**
* Query
**/
func QDDL(sql string, args ...any) (et.Items, error) {
	return IDXQuery(0, sql, args...)
}

func Query(sql string, args ...any) (et.Items, error) {
	return IDXQuery(0, sql, args...)
}

func QueryOne(sql string, args ...any) (et.Item, error) {
	return IDXQueryOne(0, sql, args...)
}

func QueryCount(sql string, args ...any) int {
	return IDXQueryCount(0, sql, args...)
}

func QueryAtrib(sql, atrib string, args ...any) (et.Items, error) {
	return IDXQueryAtrib(0, sql, atrib, args...)
}

func QueryAtribOne(sql, atrib string, args ...any) (et.Item, error) {
	return IDXQueryAtribOne(0, sql, atrib, args...)
}

func QueryData(sql string, args ...any) (et.Items, error) {
	return IDXQueryData(0, sql, args...)
}

func QueryDataOne(sql string, args ...any) (et.Item, error) {
	return IDXQueryDataOne(0, sql, args...)
}

/**
*
**/
func HttpQuery(sql string, args []any) (et.Items, error) {
	if !utility.ValidStr(sql, 0, []string{""}) {
		return et.Items{}, console.AlertF("SQL is empty")
	}

	return Query(sql, args...)
}
