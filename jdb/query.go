package jdb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/strs"
)

func TipoSQL(query string) string {
	q := strings.TrimSpace(strings.ToUpper(query))

	parts := strings.Fields(q)
	if len(parts) == 0 {
		return "DESCONOCIDO"
	}

	cmd := parts[0]

	switch cmd {
	case "SELECT":
		return "query"
	case "INSERT", "UPDATE", "DELETE", "MERGE":
		return "command"
	case "CREATE", "ALTER", "DROP", "TRUNCATE":
		return "definition"
	case "GRANT", "REVOKE":
		return "definition"
	case "COMMIT", "ROLLBACK", "SAVEPOINT", "SET":
		return "definition"
	default:
		return "desconocido"
	}
}

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
* @param sql string, args ...any
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
* @param sql string, args ...any
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
* @param db *DB, sql string, args ...any
* @return *sql.Rows, error
**/
func (s *DB) query(sql string, args ...any) (*sql.Rows, error) {
	if s == nil {
		return nil, logs.Alertf(msg.NOT_CONNECT_DB)
	}

	rows, err := s.db.Query(sql, args...)
	if err != nil {
		event.Publish(EVENT_SQL_ERROR, et.Json{
			"db_name": s.Dbname,
			"sql":     sql,
			"args":    args,
			"error":   err.Error(),
		})
		return nil, fmt.Errorf(msg.ERR_SQL, err.Error(), sql)
	}

	tp := TipoSQL(sql)
	event.Publish(fmt.Sprintf("sql:%s", tp), et.Json{
		"db_name": s.Dbname,
		"sql":     sql,
		"args":    args,
	})

	return rows, nil
}

/**
* Ddl
* @param sql string, args ...any
* @return error
**/
func (d *DB) Ddl(sql string, args ...any) error {
	_, err := d.query(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

/**
* Query
* @param sql string, args ...any
* @return et.Items, error
**/
func (d *DB) Query(sql string, args ...any) (et.Items, error) {
	rows, err := d.query(sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return rowsItems(rows), nil
}

/**
* Source
* @param sourceField string, sql string, args ...any
* @return et.Items, error
**/
func (d *DB) Source(sourceField string, sql string, args ...any) (et.Items, error) {
	rows, err := d.query(sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return sourceItems(rows, sourceField), nil
}

/**
* Command
* @param sql string, args ...any
* @return et.Items, error
**/
func (d *DB) Command(sql string, args ...any) (et.Items, error) {
	rows, err := d.query(sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return rowsItems(rows), nil
}

/**
* CommandSource
* @param sourceField string, sql string, args ...any
* @return et.Items, error
**/
func (d *DB) CommandSource(sourceField string, sql string, args ...any) (et.Items, error) {
	rows, err := d.query(sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return sourceItems(rows, sourceField), nil
}

/**
* Bulck
* @param sql string, args ...any
* @return error
**/
func (d *DB) Bulck(sql string, args ...any) error {
	_, err := d.query(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
