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
*	  "from": "table:A",
*	  "join": {
*	    "type": "left",
*	    "to": "roles:B",
*	    "on": {
*	      "B.user_id": {"eq": {"col": "A.id"}},
*	      "and": [
*	        {"B.role": {"neg": "admin"}}
*	      ]
*	    }
*	  },
*	  "select": ["A.id", "A.name", "count(*)"],
*	  "wheres": {...},
*	  "group_by": ["A.name"],
*	  "having": {...},
*	  "limit": {"page": 1, "rows": 100},
*	  "order_by": ["A.name"],
*	  "order_by_desc": ["A.age"]
*	}
*
* "from" y el "to" de cada join aceptan alias con la sintaxis
* "tabla:alias" (ver renderTableRef). "join" acepta un unico objeto o
* un arreglo de ellos (para varios joins); "type" es opcional
* ("join"/"inner"/"left"/"right", por defecto "join"); "on" tiene la
* misma forma recursiva and/or que "wheres" (ver Join/buildJoins).
*
* "select" y las claves de columna dentro de "having" aceptan tanto
* columnas simples/calificadas como llamadas a funciones de
* agregacion COUNT/MAX/MIN/SUM (p.ej. "count(*)", "sum(price)"); ver
* renderExpr. "having" tiene la misma forma recursiva and/or que
* "wheres" (ver buildWheres).
*
* Los valores de una condicion (en "wheres", "having" u "on") son
* literales por defecto; para comparar contra otra columna en vez de
* un literal, use {"col": "identificador"} (ver renderValue) — asi es
* como una clausula ON expresa "B.user_id = A.id".
*
* a una sentencia SQL SELECT para el jquery/dialect.Dialect indicado.
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
* Motores disponibles: dialect.Postgres, dialect.SQLite, dialect.MySQL,
* dialect.SQLServer y dialect.Oracle. Agregar uno nuevo es cuestion de
* sumar su propio archivo al paquete jquery/dialect implementando
* dialect.Dialect (ver el comentario de paquete en dialect/dialect.go),
* sin tocar este archivo.
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
* Build arma la sentencia SQL SELECT completa a partir del estado del
* builder.
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
* buildSelect. Cada elemento de Select puede ser una columna simple/
* calificada o una llamada de agregacion COUNT/MAX/MIN/SUM (ver
* renderExpr).
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
* buildHaving tiene la misma forma recursiva and/or que buildWhere;
* las claves de columna dentro de "having" tambien aceptan llamadas de
* agregacion (p.ej. "count(*)"), ver renderExpr/buildColumnConditions.
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
