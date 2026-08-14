/**
* Package jquery traduce un et.Json a una sentencia SQL SELECT.
*
* Soporta FROM, SELECT (con columnas y agregaciones COUNT/MAX/MIN/
* SUM), WHERE (con conectores AND/OR anidables), GROUP BY, HAVING
* (misma forma recursiva AND/OR que WHERE, incluidas agregaciones en
* sus columnas), LIMIT/OFFSET (via limit.page + limit.rows) y ORDER BY
* (ASC/DESC). El dialecto por defecto es PostgreSQL; ver el paquete
* jquery/dialect (Dialect + registry Register/Get con patron factory)
* para agregar sqlite, mysql, oracle o sqlserver, o para seleccionar
* un dialecto distinto via el atributo "dialect" del query.
*
* Formato esperado del query:
*
*	{
*	  "from": "table",
*	  "dialect": "postgres",
*	  "select": ["id", "name", "age", "count(*)"],
*	  "wheres": {
*	    "name": {"eq": "cesar"},
*	    "and": [
*	      {"age": {"eq": 30}}
*	    ],
*	    "or": [
*	      {"age": {"more": 45}}
*	    ]
*	  },
*	  "group_by": ["name"],
*	  "having": {
*	    "name": {"eq": "cesar"},
*	    "and": [
*	      {"count(*)": {"eq": 30}}
*	    ]
*	  },
*	  "limit": {"page": 1, "rows": 100},
*	  "order_by": ["name"],
*	  "order_by_desc": ["age"]
*	}
*
* Operadores soportados dentro de "wheres"/"having" (ver Operator):
* eq, neg, less, less_eq, more, more_eq, like, in, not_in, is, is_not,
* null, not_null, between, not_between.
*
* Agregaciones soportadas en "select" y en las claves de columna de
* "having": count(col), max(col), min(col), sum(col) (case-insensitive).
* count(*) y count() son equivalentes; "*" nunca se cita.
**/
package jquery

import "github.com/celsiainternet/elvis/et"

/**
* JQuery traduce query a una sentencia SQL SELECT completa para
* PostgreSQL. Para usar otro dialecto (una vez implementado), use
* NewJQueryBuilderWithDialect directamente.
* @param query et.Json
* @return string, error
**/
func JQuery(query et.Json) (string, error) {
	builder, err := NewJQueryBuilder(query)
	if err != nil {
		return "", err
	}

	return builder.Build()
}
