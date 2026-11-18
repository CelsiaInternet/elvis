package jquery

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/celsiainternet/elvis/jquery/dialect"
)

var aggregateExprRe = regexp.MustCompile(`(?i)^(count|max|min|sum)\(([^)]*)\)$`)

/**
* renderExpr: Traduce una columna o una agregación COUNT/MAX/MIN/SUM a SQL.
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
