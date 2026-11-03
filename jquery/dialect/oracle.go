package dialect

import (
	"strings"

	"github.com/celsiainternet/elvis/strs"
)

/**
* Oracle es el nombre bajo el cual OracleDialect se registra en el
* registry (ver Register/Get).
**/
const Oracle = "oracle"

func init() {
	Register(Oracle, func() Dialect {
		return &OracleDialect{}
	})
}

/**
* OracleDialect implementa Dialect para Oracle Database.
**/
type OracleDialect struct{}

/**
* Name
* @return string
**/
func (d *OracleDialect) Name() string {
	return Oracle
}

/**
* QuoteIdent
* @param name string
* @return string
**/
func (d *OracleDialect) QuoteIdent(name string) string {
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
* Like: Oracle no tiene ILIKE; LIKE es sensible a mayusculas.
* @return string
**/
func (d *OracleDialect) Like() string {
	return "LIKE"
}

/**
* LimitOffset: usa la clausula OFFSET/FETCH del estandar ANSI
* SQL:2008, soportada desde Oracle Database 12c.
* @param rows int, offset int
* @return string
**/
func (d *OracleDialect) LimitOffset(rows, offset int) string {
	if rows <= 0 {
		return ""
	}

	return strs.Format(`OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, offset, rows)
}
