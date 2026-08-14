package dialect

import (
	"strings"

	"github.com/celsiainternet/elvis/strs"
)

/**
* SQLServer es el nombre bajo el cual SQLServerDialect se registra en
* el registry (ver Register/Get).
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
* QuoteIdent: SQL Server usa corchetes ([col]) como caracter de
* comillado de identificadores.
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
* Like: SQL Server no tiene ILIKE; la sensibilidad a mayusculas de
* LIKE depende de la collation de la columna/base de datos.
* @return string
**/
func (d *SQLServerDialect) Like() string {
	return "LIKE"
}

/**
* LimitOffset: SQL Server no soporta LIMIT/OFFSET; usa la clausula
* OFFSET/FETCH del estandar ANSI SQL:2008 (SQL Server 2012+), que
* requiere una clausula ORDER BY previa en la consulta.
* @param rows int, offset int
* @return string
**/
func (d *SQLServerDialect) LimitOffset(rows, offset int) string {
	if rows <= 0 {
		return ""
	}

	return strs.Format(`OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, offset, rows)
}
