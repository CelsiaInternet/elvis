package jquery

import (
	"fmt"
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jquery/dialect"
)

/**
* JQueryBuilder guarda el estado necesario para traducir un et.Json
* con la forma:
*
*	{
*	  "from": "table",
*	  "select": ["id", "name"],
*	  "wheres": {...},
*	  "limit": {"page": 1, "rows": 100},
*	  "order_by": ["name"],
*	  "order_by_desc": ["age"]
*	}
*
* a una sentencia SQL SELECT para el jquery/dialect.Dialect indicado.
**/
type JQueryBuilder struct {
	Dialect     dialect.Dialect
	From        string
	Select      []string
	Wheres      et.Json
	Page        int
	Rows        int
	OrderBy     []string
	OrderByDesc []string
}

/**
* NewJQueryBuilder crea un JQueryBuilder para PostgreSQL a partir del
* query recibido.
* @param query et.Json
* @return *JQueryBuilder, error
**/
func NewJQueryBuilder(query et.Json) (*JQueryBuilder, error) {
	return NewJQueryBuilderWithDialect(query)
}

/**
* NewJQueryBuilderWithDialect crea un JQueryBuilder para el dialecto
* indicado en el atributo "dialect" del query (ver jquery/dialect:
* Register/Get, patron factory por motor). Si el query no trae
* "dialect", cae al dialecto por defecto de jquery (dialect.Postgres).
* Hoy solo "postgres" esta implementado; agregar sqlite/mysql/oracle/
* sqlserver es cuestion de sumar su propio archivo al paquete
* jquery/dialect implementando dialect.Dialect (ver el comentario de
* paquete en dialect/dialect.go), sin tocar este archivo.
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

	limit := query.Json("limit")

	return &JQueryBuilder{
		Dialect:     d,
		From:        from,
		Select:      query.ArrayStr("select"),
		Wheres:      query.Json("wheres"),
		Page:        limit.Int("page"),
		Rows:        limit.Int("rows"),
		OrderBy:     query.ArrayStr("order_by"),
		OrderByDesc: query.ArrayStr("order_by_desc"),
	}, nil
}

/**
* Build arma la sentencia SQL SELECT completa a partir del estado del
* builder.
* @return string, error
**/
func (b *JQueryBuilder) Build() (string, error) {
	var sql strings.Builder

	sql.WriteString(b.buildSelect())

	sql.WriteString(" FROM ")
	sql.WriteString(b.Dialect.QuoteIdent(b.From))

	whereClause, err := b.buildWhere()
	if err != nil {
		return "", err
	}
	if whereClause != "" {
		sql.WriteString(" WHERE ")
		sql.WriteString(whereClause)
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
* buildSelect
* @return string
**/
func (b *JQueryBuilder) buildSelect() string {
	if len(b.Select) == 0 {
		return "SELECT *"
	}

	cols := make([]string, len(b.Select))
	for i, c := range b.Select {
		cols[i] = b.Dialect.QuoteIdent(c)
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
* limitOffset calcula el OFFSET a partir de page/rows (page 1-based).
* Paginas <= 0 se tratan como pagina 1.
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
