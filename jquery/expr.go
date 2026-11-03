package jquery

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/celsiainternet/elvis/jquery/dialect"
)

var aggregateExprRe = regexp.MustCompile(`(?i)^(count|max|min|sum)\(([^)]*)\)$`)

/**
* renderExpr traduce una expresion de columna a SQL. Soporta columnas
* simples/calificadas ("name", "tabla.columna", via d.QuoteIdent) y
* llamadas a funciones de agregacion COUNT/MAX/MIN/SUM ("count(*)",
* "sum(price)"), usadas tanto en "select" como en las claves de
* columna dentro de "having". Una llamada sin argumento ("count()")
* se trata como COUNT(*); "*" nunca se cita.
* @param d dialect.Dialect, expr string
* @return string
**/
func renderExpr(d dialect.Dialect, expr string) string {
	expr = strings.TrimSpace(expr)

	match := aggregateExprRe.FindStringSubmatch(expr)
	if match == nil {
		return d.QuoteIdent(expr)
	}

	fn := strings.ToUpper(match[1])
	arg := strings.TrimSpace(match[2])

	if arg == "" || arg == "*" {
		return fmt.Sprintf("%s(*)", fn)
	}

	return fmt.Sprintf("%s(%s)", fn, d.QuoteIdent(arg))
}
