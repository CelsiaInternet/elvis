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
* query
* @param db *DB
* @param sql string
* @param args ...any
* @return *sql.Rows
* @return error
**/
func query(db *DB, sql string, args ...any) (*sql.Rows, error) {
	if db == nil {
		return nil, console.AlertF(msg.ERR_COMM)
	}

	isSelect := func(query string) bool {
		query = strings.TrimSpace(query)
		return strings.HasPrefix(strings.ToLower(query), "select")
	}

	if !isSelect(sql) {
		return nil, console.Alert("Query is not a SELECT statement")
	}

	rows, err := db.db.Query(sql, args...)
	if err != nil {
		return nil, console.AlertF(msg.ERR_SQL, err.Error(), sql)
	}

	return rows, nil
}

/**
* command
* @param db *DB
* @param id string
* @param sql string
* @param args ...any
* @return *sql.Rows
* @return error
**/
func command(db *DB, opt, id, sql string, args ...any) (*sql.Rows, error) {
	if db == nil {
		return nil, console.AlertF(msg.ERR_COMM)
	}

	query := SQLParse(sql, args...)
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, console.ErrorF(msg.ERR_SQL, err.Error(), sql)
	}

	go db.SetCommand(opt, id, query)

	return rows, nil
}

/**
* Command
* @param db *DB
* @param sql string
* @param args ...any
* @return et.Items
* @return error
**/
func (d *DB) Command(opt, id, sql string, args ...any) (et.Item, error) {
	if d.db == nil {
		return et.Item{}, console.AlertF(msg.ERR_COMM)
	}

	rows, err := command(d, opt, id, sql, args...)
	if err != nil {
		return et.Item{}, err
	}
	defer rows.Close()

	result := rowsItem(rows)

	return result, nil
}

/**
* Query
* @param db *DB
* @param sql string
* @param args ...any
* @return et.Items
* @return error
**/
func (d *DB) Query(sql string, args ...any) (et.Items, error) {
	if d.db == nil {
		return et.Items{}, console.AlertF(msg.ERR_COMM)
	}

	rows, err := query(d, sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	result := rowsItems(rows)

	return result, nil
}

/**
* QueryOne
* @param db *DB
* @param sql string
* @param args ...any
* @return et.Item
* @return error
**/
func (d *DB) QueryOne(sql string, args ...any) (et.Item, error) {
	if d.db == nil {
		return et.Item{}, console.AlertF(msg.ERR_COMM)
	}

	items, err := query(d, sql, args...)
	if err != nil {
		return et.Item{}, err
	}

	result := rowsItem(items)

	return result, nil
}

/**
* Source
* @param db *DB
* @param sourceField string
* @param sql string
* @param args ...any
* @return et.Items
* @return error
**/
func (d *DB) Source(sourceField string, sql string, args ...any) (et.Items, error) {
	if d.db == nil {
		return et.Items{}, console.AlertF(msg.ERR_COMM)
	}

	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return et.Items{}, console.ErrorF(msg.ERR_SQL, err.Error(), sql)
	}
	defer rows.Close()

	items := sourceItems(rows, sourceField)

	return items, nil
}

/**
* SourceOne
* @param db *DB
* @param sourceField string
* @param sql string
* @param args ...any
* @return et.Item
* @return error
**/
func (d *DB) SourceOne(sourceField string, sql string, args ...any) (et.Item, error) {
	if d.db == nil {
		return et.Item{}, console.AlertF(msg.ERR_COMM)
	}

	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return et.Item{}, err
	}
	defer rows.Close()

	item := sourceItem(rows, sourceField)

	return item, nil
}
