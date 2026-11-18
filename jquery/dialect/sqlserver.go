package dialect

import (
	"strings"

	"github.com/celsiainternet/elvis/strs"
)

/**
* SQLServer: Nombre con que se registra SQLServerDialect.
**/
const SQLServer = "sqlserver"

func init() {
	Register(SQLServer, func() Dialect {
		return &SQLServerDialect{}
	})
}

/**
* SQLServerDialect implementa Dialect para Microsoft SQL Server.
**/
type SQLServerDialect struct{}

/**
* Name
* @return string
**/
func (d *SQLServerDialect) Name() string {
	return SQLServer
}

/**
* QuoteIdent: Cita el identificador con corchetes.
* @param name string
* @return string
**/
func (d *SQLServerDialect) QuoteIdent(name string) string {
	parts := strings.Split(name, ".")
	for i, p := range parts {
		p = strings.Trim(p, "[]")
		if p == "" || p == "*" {
			parts[i] = p
			continue
		}
		parts[i] = strs.Format("[%s]", p)
	}

	return strings.Join(parts, ".")
}

/**
* Like: LIKE, porque SQL Server no tiene ILIKE.
* @return string
**/
func (d *SQLServerDialect) Like() string {
	return "LIKE"
}

/**
* LimitOffset: Usa OFFSET/FETCH (SQL Server 2012+, requiere ORDER BY).
* @param rows int, offset int
* @return string
**/
func (d *SQLServerDialect) LimitOffset(rows, offset int) string {
	if rows <= 0 {
		return ""
	}

	return strs.Format(`OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, offset, rows)
}
