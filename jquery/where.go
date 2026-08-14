package jquery

import (
	"fmt"
	"sort"
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jquery/dialect"
)

const (
	keyAnd = "and"
	keyOr  = "or"
)

/**
* asJson normaliza un valor crudo (proveniente de json.Unmarshal o de
* un et.Json construido a mano en Go) a et.Json. json.Unmarshal
* produce map[string]any; codigo Go que arma el query sin
* pasar por texto JSON normalmente produce et.Json directamente.
* @param val any
* @return et.Json, bool
**/
func asJson(val any) (et.Json, bool) {
	switch v := val.(type) {
	case et.Json:
		return v, true
	case map[string]any:
		return et.Json(v), true
	default:
		return nil, false
	}
}

/**
* asArray normaliza un valor crudo a []any, soportando los
* mismos dos origenes que asJson (JSON real vs et.Json/[]et.Json
* armado a mano).
* @param val any
* @return []any, bool
**/
func asArray(val any) ([]any, bool) {
	switch v := val.(type) {
	case []any:
		return v, true
	case []et.Json:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, true
	case []map[string]any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, true
	default:
		return nil, false
	}
}

/**
* sortedKeys devuelve las claves de un et.Json ordenadas
* alfabeticamente. El orden de iteracion de un map en Go no es
* determinista; como todas las condiciones de un mismo grupo se
* combinan con el mismo conector (AND, o AND entre los grupos "and"/
* "or"), el orden no afecta el significado de la consulta — pero sin
* ordenar, el mismo query.Json produciria un SQL con las condiciones
* en distinto orden en cada llamada, lo cual complica pruebas y el
* uso del SQL resultante como cache key.
* @param data et.Json
* @return []string
**/
func sortedKeys(data et.Json) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

/**
* buildWheres construye recursivamente la clausula WHERE (sin la
* palabra WHERE) a partir de un grupo de condiciones con la forma:
*
*	{
*	  "<columna>": {"<operador>": valor, ...},
*	  "and": [ {<grupo>}, {<grupo>}, ... ],
*	  "or":  [ {<grupo>}, {<grupo>}, ... ]
*	}
*
* Las columnas presentes directamente en el grupo se combinan entre si
* con AND. Cada elemento del arreglo "and"/"or" es a su vez un grupo
* del mismo tipo (recursivo); esos elementos se combinan entre si con
* AND u OR segun corresponda, y el resultado se agrega entre parentesis
* al resto del grupo padre cuando hay mas de un elemento.
* @param d dialect.Dialect
* @param group et.Json
* @return string, error
**/
func buildWheres(d dialect.Dialect, group et.Json) (string, error) {
	var parts []string

	for _, key := range sortedKeys(group) {
		switch key {
		case keyAnd, keyOr:
			items, ok := asArray(group.Get(key))
			if !ok {
				return "", fmt.Errorf(ERR_AND_OR_INVALID, key)
			}

			var subParts []string
			for _, rawItem := range items {
				itemJson, ok := asJson(rawItem)
				if !ok {
					return "", fmt.Errorf(ERR_AND_OR_INVALID, key)
				}

				part, err := buildWheres(d, itemJson)
				if err != nil {
					return "", err
				}
				if part != "" {
					subParts = append(subParts, part)
				}
			}

			if len(subParts) == 0 {
				continue
			}

			connector := " AND "
			if key == keyOr {
				connector = " OR "
			}

			joined := strings.Join(subParts, connector)
			if len(subParts) > 1 {
				joined = fmt.Sprintf("(%s)", joined)
			}

			parts = append(parts, joined)
		default:
			ops, ok := asJson(group.Get(key))
			if !ok || len(ops) == 0 {
				return "", fmt.Errorf(ERR_WHERE_INVALID, key)
			}

			condParts, err := buildColumnConditions(d, key, ops)
			if err != nil {
				return "", err
			}

			parts = append(parts, condParts...)
		}
	}

	return strings.Join(parts, " AND "), nil
}

/**
* buildColumnConditions arma una condicion SQL por cada operador
* definido para una columna (p.ej. {"more": 10, "less": 20} produce
* dos condiciones, combinadas luego con AND por buildWheres). column
* acepta tanto una columna simple/calificada como una llamada de
* agregacion COUNT/MAX/MIN/SUM (p.ej. "count(*)"), util cuando este
* grupo se usa para HAVING (ver renderExpr).
* @param d dialect.Dialect, column string, ops et.Json
* @return []string, error
**/
func buildColumnConditions(d dialect.Dialect, column string, ops et.Json) ([]string, error) {
	ident := renderExpr(d, column)

	var parts []string
	for _, key := range sortedKeys(ops) {
		op := Operator(key)
		if !IsValidOperator(op) {
			return nil, fmt.Errorf(ERR_OPERATOR_INVALID, key)
		}

		part, err := buildCondition(d, ident, op, ops.Get(key), column)
		if err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil
}

/**
* buildCondition renderiza una unica condicion "columna <operador> valor".
* @param d dialect.Dialect, ident string, op Operator, val any, column string
* @return string, error
**/
func buildCondition(d dialect.Dialect, ident string, op Operator, val any, column string) (string, error) {
	switch op {
	case EQ, NEG, LESS, LESS_EQ, MORE, MORE_EQ, IS, IS_NOT:
		symbol := operatorSymbols[op]
		return fmt.Sprintf("%s %s %v", ident, symbol, et.Unquote(val)), nil

	case LIKE:
		return fmt.Sprintf("%s %s %v", ident, d.Like(), et.Unquote(val)), nil

	case NULL:
		return fmt.Sprintf("%s IS NULL", ident), nil

	case NOT_NULL:
		return fmt.Sprintf("%s IS NOT NULL", ident), nil

	case IN, NOT_IN:
		values, ok := asArray(val)
		if !ok || len(values) == 0 {
			return "", fmt.Errorf(ERR_IN_VALUES, string(op), column)
		}

		rendered := make([]string, len(values))
		for i, v := range values {
			rendered[i] = fmt.Sprintf("%v", et.Unquote(v))
		}

		keyword := "IN"
		if op == NOT_IN {
			keyword = "NOT IN"
		}

		return fmt.Sprintf("%s %s (%s)", ident, keyword, strings.Join(rendered, ", ")), nil

	case BETWEEN, NOT_BETWEEN:
		values, ok := asArray(val)
		if !ok || len(values) != 2 {
			return "", fmt.Errorf(ERR_BETWEEN_VALUES, string(op), column)
		}

		keyword := "BETWEEN"
		if op == NOT_BETWEEN {
			keyword = "NOT BETWEEN"
		}

		return fmt.Sprintf("%s %s %v AND %v", ident, keyword, et.Unquote(values[0]), et.Unquote(values[1])), nil

	default:
		return "", fmt.Errorf(ERR_OPERATOR_INVALID, string(op))
	}
}
