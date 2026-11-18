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
* asJson: Normaliza un valor (map o et.Json) a et.Json.
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
* asArray: Normaliza un valor a []any.
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
* sortedKeys: Claves de data en orden alfabético, para un SQL determinista.
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
* buildWheres: Arma recursivamente el WHERE desde columnas y grupos and/or.
* @param d dialect.Dialect, group et.Json
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
* buildColumnConditions: Una condición SQL por operador de la columna.
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
* renderValue: Renderiza un literal o, con {"col": ...}, una columna.
* @param d dialect.Dialect, val any
* @return string
**/
func renderValue(d dialect.Dialect, val any) string {
	if obj, ok := asJson(val); ok {
		if col := strings.TrimSpace(obj.Str("col")); col != "" {
			return d.QuoteIdent(col)
		}
	}

	return fmt.Sprintf("%v", et.Unquote(val))
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
		return fmt.Sprintf("%s %s %s", ident, symbol, renderValue(d, val)), nil

	case LIKE:
		return fmt.Sprintf("%s %s %s", ident, d.Like(), renderValue(d, val)), nil

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
			rendered[i] = renderValue(d, v)
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

		return fmt.Sprintf("%s %s %s AND %s", ident, keyword, renderValue(d, values[0]), renderValue(d, values[1])), nil

	default:
		return "", fmt.Errorf(ERR_OPERATOR_INVALID, string(op))
	}
}
