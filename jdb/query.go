package jdb

import (
	"database/sql"
	"strings"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
)

/**
* SQLQuote
* @param sql string
* @return string
**/
func SQLQuote(sql string) string {
	sql = strings.TrimSpace(sql)

	result := strs.Replace(sql, `'`, `"`)
	result = strs.Trim(result)

	return result
}

/**
* SQLDDL
* @param sql string
* @param args ...any
* @return string
**/
func SQLDDL(sql string, args ...any) string {
	sql = strings.TrimSpace(sql)

	for i, arg := range args {
		old := strs.Format(`$%d`, i+1)
		new := strs.Format(`%v`, arg)
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

/**
* SQLParse
* @param sql string
* @param args ...any
* @return string
**/
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
* Query
* @param db *sql.DB
* @param sql string
* @param args ...any
* @return et.Items
* @return error
**/
func Query(db *sql.DB, sql string, args ...any) (et.Items, error) {
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

/**
* QueryOne
* @param db *sql.DB
* @param sql string
* @param args ...any
* @return et.Item
* @return error
**/
func QueryOne(db *sql.DB, sql string, args ...any) (et.Item, error) {
	if db == nil {
		return et.Item{}, console.AlertF(msg.ERR_COMM)
	}

	sql = SQLParse(sql, args...)
	rows, err := db.Query(sql)
	if err != nil {
		return et.Item{}, err
	}
	defer rows.Close()

	item := rowsItem(rows)

	return item, nil
}

/**
* Source
* @param db *sql.DB
* @param sourceField string
* @param sql string
* @param args ...any
* @return et.Items
* @return error
**/
func Source(db *sql.DB, sourceField string, sql string, args ...any) (et.Items, error) {
	if db == nil {
		return et.Items{}, console.AlertF(msg.ERR_COMM)
	}

	sql = SQLParse(sql, args...)
	rows, err := db.Query(sql)
	if err != nil {
		return et.Items{}, console.ErrorF(msg.ERR_SQL, err.Error(), sql)
	}
	defer rows.Close()

	items := sourceItems(rows, sourceField)

	return items, nil
}

/**
* SourceOne
* @param db *sql.DB
* @param sourceField string
* @param sql string
* @param args ...any
* @return et.Item
* @return error
**/
func SourceOne(db *sql.DB, sourceField string, sql string, args ...any) (et.Item, error) {
	if db == nil {
		return et.Item{}, console.AlertF(msg.ERR_COMM)
	}

	sql = SQLParse(sql, args...)
	rows, err := db.Query(sql)
	if err != nil {
		return et.Item{}, err
	}
	defer rows.Close()

	item := sourceItem(rows, sourceField)

	return item, nil
}
