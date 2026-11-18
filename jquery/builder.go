package jquery

import (
	"fmt"
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jquery/dialect"
)

/**
* JQueryBuilder: Estado para traducir un et.Json (ver el paquete jquery) a un SELECT de un dialecto.
**/
type JQueryBuilder struct {
	Dialect     dialect.Dialect
	From        string
	Joins       []Join
	Select      []string
	Wheres      et.Json
	GroupBy     []string
	Having      et.Json
	Page        int
	Rows        int
	OrderBy     []string
	OrderByDesc []string
}

/**
* NewJQueryBuilder: Crea un JQueryBuilder para PostgreSQL.
* @param query et.Json
* @return *JQueryBuilder, error
**/
func NewJQueryBuilder(query et.Json) (*JQueryBuilder, error) {
	return NewJQueryBuilderWithDialect(query)
}

/**
* NewJQueryBuilderWithDialect: Crea un JQueryBuilder para el "dialect" del query (por defecto Postgres).
* @param query et.Json
* @return *JQueryBuilder, error
**/
func NewJQueryBuilderWithDialect(query et.Json) (*JQueryBuilder, error) {
	dialectName := strings.TrimSpace(query.Str("dialect"))
	if dialectName == "" {
		dialectName = dialect.Postgres
	}

	d, err := dialect.Get(dialectName)
	if err != nil {
		return nil, err
	}

	from := strings.TrimSpace(query.Str("from"))
	if from == "" {
		return nil, fmt.Errorf(ERR_FROM_REQUIRED)
	}

	joins, err := parseJoins(query.Get("join"))
	if err != nil {
		return nil, err
	}

	limit := query.Json("limit")

	return &JQueryBuilder{
		Dialect:     d,
		From:        from,
		Joins:       joins,
		Select:      query.ArrayStr("select"),
		Wheres:      query.Json("wheres"),
		GroupBy:     query.ArrayStr("group_by"),
		Having:      query.Json("having"),
		Page:        limit.Int("page"),
		Rows:        limit.Int("rows"),
		OrderBy:     query.ArrayStr("order_by"),
		OrderByDesc: query.ArrayStr("order_by_desc"),
	}, nil
}

/**
* Build: Arma la sentencia SELECT completa.
* @return string, error
**/
func (b *JQueryBuilder) Build() (string, error) {
	var sql strings.Builder

	sql.WriteString(b.buildSelect())

	sql.WriteString(" FROM ")
	sql.WriteString(renderTableRef(b.Dialect, b.From))

	joinClause, err := buildJoins(b.Dialect, b.Joins)
	if err != nil {
		return "", err
	}
	sql.WriteString(joinClause)

	whereClause, err := b.buildWhere()
	if err != nil {
		return "", err
	}
	if whereClause != "" {
		sql.WriteString(" WHERE ")
		sql.WriteString(whereClause)
	}

	if groupClause := b.buildGroupBy(); groupClause != "" {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(groupClause)
	}

	havingClause, err := b.buildHaving()
	if err != nil {
		return "", err
	}
	if havingClause != "" {
		sql.WriteString(" HAVING ")
		sql.WriteString(havingClause)
	}

	if orderClause := b.buildOrderBy(); orderClause != "" {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(orderClause)
	}

	if limitClause := b.Dialect.LimitOffset(b.limitRows(), b.limitOffset()); limitClause != "" {
		sql.WriteString(" ")
		sql.WriteString(limitClause)
	}

	return sql.String(), nil
}

/**
* buildSelect: Arma el SELECT con columnas o agregaciones COUNT/MAX/MIN/SUM.
* @return string
**/
func (b *JQueryBuilder) buildSelect() string {
	if len(b.Select) == 0 {
		return "SELECT *"
	}

	cols := make([]string, len(b.Select))
	for i, c := range b.Select {
		cols[i] = renderExpr(b.Dialect, c)
	}

	return "SELECT " + strings.Join(cols, ", ")
}

/**
* buildWhere
* @return string, error
**/
func (b *JQueryBuilder) buildWhere() (string, error) {
	if len(b.Wheres) == 0 {
		return "", nil
	}

	return buildWheres(b.Dialect, b.Wheres)
}

/**
* buildGroupBy
* @return string
**/
func (b *JQueryBuilder) buildGroupBy() string {
	if len(b.GroupBy) == 0 {
		return ""
	}

	cols := make([]string, len(b.GroupBy))
	for i, c := range b.GroupBy {
		cols[i] = b.Dialect.QuoteIdent(c)
	}

	return strings.Join(cols, ", ")
}

/**
* buildHaving: Arma el HAVING con la misma forma and/or que WHERE.
* @return string, error
**/
func (b *JQueryBuilder) buildHaving() (string, error) {
	if len(b.Having) == 0 {
		return "", nil
	}

	return buildWheres(b.Dialect, b.Having)
}

/**
* buildOrderBy
* @return string
**/
func (b *JQueryBuilder) buildOrderBy() string {
	var parts []string

	for _, c := range b.OrderBy {
		parts = append(parts, b.Dialect.QuoteIdent(c)+" ASC")
	}

	for _, c := range b.OrderByDesc {
		parts = append(parts, b.Dialect.QuoteIdent(c)+" DESC")
	}

	return strings.Join(parts, ", ")
}

/**
* limitRows
* @return int
**/
func (b *JQueryBuilder) limitRows() int {
	return b.Rows
}

/**
* limitOffset: Calcula el OFFSET desde page/rows (page base 1).
* @return int
**/
func (b *JQueryBuilder) limitOffset() int {
	if b.Rows <= 0 {
		return 0
	}

	page := b.Page
	if page <= 0 {
		page = 1
	}

	return (page - 1) * b.Rows
}
