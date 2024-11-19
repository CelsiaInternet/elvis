package jdb

import (
	"database/sql"
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/strs"
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
		return nil, logs.Alertf(msg.NOT_CONNECT_DB)
	}

	rows, err := db.db.Query(sql, args...)
	if err != nil {
		return nil, logs.Alertf(msg.ERR_SQL, err.Error(), sql)
	}

	return rows, nil
}

/**
* exec
* @param db *DB
* @param id string
* @param sql string
* @param args ...any
* @return error
**/
func exec(db *DB, id, sql string, args ...any) error {
	if db == nil {
		return logs.Alertf(msg.NOT_CONNECT_DB)
	}

	query := SQLParse(sql, args...)
	_, err := db.db.Exec(query)
	if err != nil {
		return logs.Errorf(msg.ERR_SQL, err.Error(), sql)
	}

	err = db.upsertDDL(id, query)
	if err != nil {
		return err
	}

	return nil
}

/**
* Exec
* @param id string
* @param sql string
* @param args ...any
* @return et.Items
* @return error
**/
func (d *DB) Exec(id, sql string, args ...any) error {
	err := exec(d, id, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

/**
* QueryRows
* @param sql string
* @param args ...any
* @return *sql.Rows
* @return error
**/
func (d *DB) QueryRows(sql string, args ...any) (*sql.Rows, error) {
	if d.db == nil {
		return nil, logs.Alertf(msg.NOT_CONNECT_DB)
	}

	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return nil, logs.Alertf(msg.ERR_SQL, err.Error(), sql)
	}

	return rows, nil
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
		return et.Items{}, logs.Alertf(msg.NOT_CONNECT_DB)
	}

	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return et.Items{}, logs.Errorf(msg.ERR_SQL, err.Error(), sql)
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
		return et.Item{}, logs.Alertf(msg.NOT_CONNECT_DB)
	}

	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return et.Item{}, err
	}
	defer rows.Close()

	item := sourceItem(rows, sourceField)

	return item, nil
}
