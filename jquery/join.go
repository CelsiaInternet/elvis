package jquery

import (
	"fmt"
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jquery/dialect"
)

/**
* JoinType: variante de JOIN soportada dentro de "join" (ver Join).
* El valor string es exactamente la clave que se espera recibir en el
* atributo "type" de cada join.
**/
type JoinType string

const (
	JoinTypeJoin  JoinType = "join"
	JoinTypeInner JoinType = "inner"
	JoinTypeLeft  JoinType = "left"
	JoinTypeRight JoinType = "right"
)

/**
* joinKeywords mapea cada JoinType a su palabra clave SQL.
**/
var joinKeywords = map[JoinType]string{
	JoinTypeJoin:  "JOIN",
	JoinTypeInner: "INNER JOIN",
	JoinTypeLeft:  "LEFT JOIN",
	JoinTypeRight: "RIGHT JOIN",
}

/**
* Join: una entrada de la clausula "join" del query, con la forma:
*
*	{
*	  "type": "left",
*	  "to": "roles:B",
*	  "on": {
*	    "B.user_id": {"eq": {"col": "A.id"}},
*	    "and": [
*	      {"B.role": {"neg": "admin"}}
*	    ]
*	  }
*	}
*
* "type" es opcional y por defecto es "join" (JOIN simple, equivalente
* a INNER JOIN). "to" acepta la misma sintaxis "tabla:alias" que
* "from" (ver renderTableRef). "on" tiene la misma forma recursiva
* and/or que "wheres"/"having" (ver buildWheres); sus valores aceptan
* referencias a columna via {"col": "identificador"} (ver renderValue),
* lo usual para condiciones columna-a-columna como B.user_id = A.id.
**/
type Join struct {
	Type JoinType
	To   string
	On   et.Json
}

/**
* parseJoins normaliza el atributo "join" del query (un unico objeto o
* un arreglo de objetos) a []Join. raw nil (atributo ausente) produce
* (nil, nil): el query simplemente no tiene joins.
* @param raw any
* @return []Join, error
**/
func parseJoins(raw any) ([]Join, error) {
	if raw == nil {
		return nil, nil
	}

	var items []any
	if arr, ok := asArray(raw); ok {
		items = arr
	} else if obj, ok := asJson(raw); ok {
		items = []any{obj}
	} else {
		return nil, fmt.Errorf(ERR_JOIN_INVALID)
	}

	joins := make([]Join, 0, len(items))
	for _, rawItem := range items {
		obj, ok := asJson(rawItem)
		if !ok {
			return nil, fmt.Errorf(ERR_JOIN_INVALID)
		}

		to := strings.TrimSpace(obj.Str("to"))
		if to == "" {
			return nil, fmt.Errorf(ERR_JOIN_TO_REQUIRED)
		}

		on := obj.Json("on")
		if len(on) == 0 {
			return nil, fmt.Errorf(ERR_JOIN_ON_REQUIRED, to)
		}

		typ := JoinType(strings.ToLower(strings.TrimSpace(obj.Str("type"))))
		if typ == "" {
			typ = JoinTypeJoin
		}
		if _, ok := joinKeywords[typ]; !ok {
			return nil, fmt.Errorf(ERR_JOIN_TYPE_INVALID, string(typ))
		}

		joins = append(joins, Join{Type: typ, To: to, On: on})
	}

	return joins, nil
}

/**
* renderTableRef renderiza una referencia de tabla con alias opcional,
* con la forma "tabla" o "tabla:alias" (p.ej. "users:A" produce
* "\"users\" AS \"A\""). "tabla" puede venir calificada por esquema
* (p.ej. "public.users:A").
* @param d dialect.Dialect, ref string
* @return string
**/
func renderTableRef(d dialect.Dialect, ref string) string {
	parts := strings.SplitN(strings.TrimSpace(ref), ":", 2)

	table := d.QuoteIdent(strings.TrimSpace(parts[0]))
	if len(parts) == 1 {
		return table
	}

	alias := strings.TrimSpace(parts[1])
	if alias == "" {
		return table
	}

	return table + " AS " + d.QuoteIdent(alias)
}

/**
* buildJoins arma las clausulas JOIN (una por cada elemento de Joins),
* en el orden en que fueron declaradas.
* @param d dialect.Dialect, joins []Join
* @return string, error
**/
func buildJoins(d dialect.Dialect, joins []Join) (string, error) {
	var sql strings.Builder

	for _, j := range joins {
		onClause, err := buildWheres(d, j.On)
		if err != nil {
			return "", err
		}

		sql.WriteString(" ")
		sql.WriteString(joinKeywords[j.Type])
		sql.WriteString(" ")
		sql.WriteString(renderTableRef(d, j.To))
		sql.WriteString(" ON ")
		sql.WriteString(onClause)
	}

	return sql.String(), nil
}
