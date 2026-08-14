package dialect

import (
	"strings"

	"github.com/celsiainternet/elvis/strs"
)

/**
* Postgres es el nombre bajo el cual PostgresDialect se registra en
* el registry (ver Get/Register).
**/
const Postgres = "postgres"

func init() {
	Register(Postgres, func() Dialect {
		return &PostgresDialect{}
	})
}

/**
* PostgresDialect implementa Dialect para PostgreSQL.
**/
type PostgresDialect struct{}

/**
* Name
* @return string
**/
func (d *PostgresDialect) Name() string {
	return Postgres
}

/**
* QuoteIdent
* @param name string
* @return string
**/
func (d *PostgresDialect) QuoteIdent(name string) string {
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
* Like
* @return string
**/
func (d *PostgresDialect) Like() string {
	return "ILIKE"
}

/**
* LimitOffset
* @param rows int, offset int
* @return string
**/
func (d *PostgresDialect) LimitOffset(rows, offset int) string {
	if rows <= 0 {
		return ""
	}

	if offset > 0 {
		return strs.Format(`LIMIT %d OFFSET %d`, rows, offset)
	}

	return strs.Format(`LIMIT %d`, rows)
}
