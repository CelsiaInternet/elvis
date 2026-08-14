/**
* Package jquery traduce un et.Json a una sentencia SQL SELECT.
*
* Soporta FROM, SELECT, WHERE (con conectores AND/OR anidables),
* LIMIT/OFFSET (via limit.page + limit.rows) y ORDER BY (ASC/DESC).
* El dialecto por defecto es PostgreSQL; ver el paquete jquery/dialect
* (Dialect + registry Register/Get con patron factory) para agregar
* sqlite, mysql, oracle o sqlserver, o para seleccionar un dialecto
* distinto via el atributo "dialect" del query.
*
* Formato esperado del query:
*
*	{
*	  "from": "table",
*	  "dialect": "postgres",
*	  "select": ["id", "name", "age"],
*	  "wheres": {
*	    "name": {"eq": "cesar"},
*	    "and": [
*	      {"age": {"eq": 30}}
*	    ],
*	    "or": [
*	      {"age": {"more": 45}}
*	    ]
*	  },
*	  "limit": {"page": 1, "rows": 100},
*	  "order_by": ["name"],
*	  "order_by_desc": ["age"]
*	}
*
* Operadores soportados dentro de "wheres" (ver Operator):
* eq, neg, less, less_eq, more, more_eq, like, in, not_in, is, is_not,
* null, not_null, between, not_between.
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
