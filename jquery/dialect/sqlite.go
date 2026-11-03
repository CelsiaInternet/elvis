package dialect

import (
	"strings"

	"github.com/celsiainternet/elvis/strs"
)

/**
* SQLite es el nombre bajo el cual SQLiteDialect se registra en el
* registry (ver Register/Get).
**/
const SQLite = "sqlite"

func init() {
	Register(SQLite, func() Dialect {
		return &SQLiteDialect{}
	})
}

/**
* SQLiteDialect implementa Dialect para SQLite.
**/
type SQLiteDialect struct{}

/**
* Name
* @return string
**/
func (d *SQLiteDialect) Name() string {
	return SQLite
}

/**
* QuoteIdent
* @param name string
* @return string
**/
func (d *SQLiteDialect) QuoteIdent(name string) string {
	parts := strings.Split(name, ".")
	for i, p := range parts {
		p = strings.Trim(p, `"`)
		if p == "" || p == "*" {
			parts[i] = p
			continue
		}
		parts[i] = strs.Format(`"%s"`, p)
	}

	return strings.Join(parts, ".")
}

/**
* Like: SQLite no tiene ILIKE; LIKE ya es insensible a mayusculas
* para ASCII por defecto.
* @return string
**/
func (d *SQLiteDialect) Like() string {
	return "LIKE"
}

/**
* LimitOffset
* @param rows int, offset int
* @return string
**/
func (d *SQLiteDialect) LimitOffset(rows, offset int) string {
	if rows <= 0 {
		return ""
	}

	if offset > 0 {
		return strs.Format(`LIMIT %d OFFSET %d`, rows, offset)
	}

	return strs.Format(`LIMIT %d`, rows)
}
