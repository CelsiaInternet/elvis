/**
* Package jquery traduce un et.Json a una sentencia SQL SELECT.
*
* Soporta FROM (con alias opcional "tabla:alias"), JOIN/INNER JOIN/
* LEFT JOIN/RIGHT JOIN, SELECT (con columnas y agregaciones COUNT/MAX/
* MIN/SUM), WHERE (con conectores AND/OR anidables), GROUP BY, HAVING
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
*	  "from": "table:A",
*	  "dialect": "postgres",
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
*	  "select": ["A.id", "A.name", "A.age", "count(*)"],
*	  "wheres": {
*	    "A.name": {"eq": "cesar"},
*	    "and": [
*	      {"A.age": {"eq": 30}}
*	    ],
*	    "or": [
*	      {"A.age": {"more": 45}}
*	    ]
*	  },
*	  "group_by": ["A.name"],
*	  "having": {
*	    "A.name": {"eq": "cesar"},
*	    "and": [
*	      {"count(*)": {"eq": 30}}
*	    ]
*	  },
*	  "limit": {"page": 1, "rows": 100},
*	  "order_by": ["A.name"],
*	  "order_by_desc": ["A.age"]
*	}
*
* "join" acepta un unico objeto (como arriba) o un arreglo de objetos
* para varios joins, aplicados en orden; "type" es opcional
* ("join"/"inner"/"left"/"right", por defecto "join", equivalente a
* INNER JOIN). Ver JoinType/Join.
*
* Operadores soportados dentro de "wheres"/"having"/join.on (ver
* Operator): eq, neg, less, less_eq, more, more_eq, like, in, not_in,
* is, is_not, null, not_null, between, not_between. El valor de una
* condicion es un literal por defecto; para comparar contra otra
* columna (como en una clausula ON) use {"col": "identificador"} en
* vez de un literal (ver renderValue).
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
