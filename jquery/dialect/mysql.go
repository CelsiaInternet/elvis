package dialect

import (
	"strings"

	"github.com/celsiainternet/elvis/strs"
)

/**
* MySQL es el nombre bajo el cual MySQLDialect se registra en el
* registry (ver Register/Get).
**/
const MySQL = "mysql"

func init() {
	Register(MySQL, func() Dialect {
		return &MySQLDialect{}
	})
}

/**
* MySQLDialect implementa Dialect para MySQL.
**/
type MySQLDialect struct{}

/**
* Name
* @return string
**/
func (d *MySQLDialect) Name() string {
	return MySQL
}

/**
* QuoteIdent: MySQL usa backtick (`) como caracter de comillado de
* identificadores.
* @param name string
* @return string
**/
func (d *MySQLDialect) QuoteIdent(name string) string {
	parts := strings.Split(name, ".")
	for i, p := range parts {
		p = strings.Trim(p, "`")
		if p == "" || p == "*" {
			parts[i] = p
			continue
		}
		parts[i] = strs.Format("`%s`", p)
	}

	return strings.Join(parts, ".")
}

/**
* Like: MySQL no tiene ILIKE; LIKE es insensible a mayusculas por
* defecto con las collations mas comunes (p.ej. utf8mb4_general_ci).
* @return string
**/
func (d *MySQLDialect) Like() string {
	return "LIKE"
}

/**
* LimitOffset
* @param rows int, offset int
* @return string
**/
func (d *MySQLDialect) LimitOffset(rows, offset int) string {
	if rows <= 0 {
		return ""
	}

	if offset > 0 {
		return strs.Format(`LIMIT %d OFFSET %d`, rows, offset)
	}

	return strs.Format(`LIMIT %d`, rows)
}
