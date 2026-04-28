package jdb

import (
	"context"
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
		return "control"
	case "COMMIT", "ROLLBACK", "SAVEPOINT", "SET":
		return "transaction"
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
* queryContext
* @param ctx context.Context, sql string, args ...any
* @return *sql.Rows, error
**/
func (s *DB) queryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	if s == nil {
		return nil, logs.Alertf(msg.NOT_CONNECT_DB)
	}

	rows, err := s.db.QueryContext(ctx, sql, args...)
	if err != nil {
		event.Publish(EVENT_SQL_ERROR, et.Json{
			"db_name": s.Dbname,
			"sql":     sql,
			"args":    args,
			"error":   err.Error(),
		})
		sql = SQLParse(sql, args...)
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
* DdlContext
* @param ctx context.Context, sql string, args ...any
* @return error
**/
func (d *DB) DdlContext(ctx context.Context, sql string, args ...any) error {
	_, err := d.queryContext(ctx, sql, args...)
	return err
}

func (d *DB) Ddl(sql string, args ...any) error {
	return d.DdlContext(context.Background(), sql, args...)
}

/**
* QueryContext
* @param ctx context.Context, sql string, args ...any
* @return et.Items, error
**/
func (d *DB) QueryContext(ctx context.Context, sql string, args ...any) (et.Items, error) {
	rows, err := d.queryContext(ctx, sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return rowsItems(rows), nil
}

func (d *DB) Query(sql string, args ...any) (et.Items, error) {
	return d.QueryContext(context.Background(), sql, args...)
}

/**
* QueryOneContext
* @param ctx context.Context, sql string, args ...any
* @return et.Item, error
**/
func (d *DB) QueryOneContext(ctx context.Context, sql string, args ...any) (et.Item, error) {
	result, err := d.QueryContext(ctx, sql, args...)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

func (d *DB) QueryOne(sql string, args ...any) (et.Item, error) {
	return d.QueryOneContext(context.Background(), sql, args...)
}

/**
* SourceContext
* @param ctx context.Context, sourceField string, sql string, args ...any
* @return et.Items, error
**/
func (d *DB) SourceContext(ctx context.Context, sourceField string, sql string, args ...any) (et.Items, error) {
	rows, err := d.queryContext(ctx, sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return sourceItems(rows, sourceField), nil
}

func (d *DB) Source(sourceField string, sql string, args ...any) (et.Items, error) {
	return d.SourceContext(context.Background(), sourceField, sql, args...)
}

/**
* CommandContext
* @param ctx context.Context, sql string, args ...any
* @return et.Items, error
**/
func (d *DB) CommandContext(ctx context.Context, sql string, args ...any) (et.Items, error) {
	rows, err := d.queryContext(ctx, sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return rowsItems(rows), nil
}

func (d *DB) Command(sql string, args ...any) (et.Items, error) {
	return d.CommandContext(context.Background(), sql, args...)
}

/**
* CommandSourceContext
* @param ctx context.Context, sourceField string, sql string, args ...any
* @return et.Items, error
**/
func (d *DB) CommandSourceContext(ctx context.Context, sourceField string, sql string, args ...any) (et.Items, error) {
	rows, err := d.queryContext(ctx, sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return sourceItems(rows, sourceField), nil
}

func (d *DB) CommandSource(sourceField string, sql string, args ...any) (et.Items, error) {
	return d.CommandSourceContext(context.Background(), sourceField, sql, args...)
}

/**
* BulckContext
* @param ctx context.Context, sql string, args ...any
* @return error
**/
func (d *DB) BulckContext(ctx context.Context, sql string, args ...any) error {
	_, err := d.queryContext(ctx, sql, args...)
	return err
}

func (d *DB) Bulck(sql string, args ...any) error {
	return d.BulckContext(context.Background(), sql, args...)
}
