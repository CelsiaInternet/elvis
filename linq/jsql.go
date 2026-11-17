package linq

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/celsiainternet/elvis/et"
)

const (
	QUERY_DEFAULT_ROWS         = 100
	QUERY_MAX_ROWS             = 1000
	MSG_QUERY_COLUMN_NOT_VALID = "Columna no valida en la consulta (%s)"
	SEARCH_MAX_FIELDS          = 5
	SEARCH_MAX_WORDS           = 5
	SEARCH_MAX_LENGTH          = 100
	MSG_SEARCH_LENGTH          = "La busqueda admite maximo %d caracteres"
	MSG_SEARCH_FIELDS          = "La busqueda requiere entre 1 y %d campos (search.fields)"
	MSG_SEARCH_FIELD_NOT_VALID = "Campo no valido en la busqueda (%s)"
)

// Comodines de LIKE escapados con \ (el escape por defecto de Postgres) para buscarlos como texto
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// Tipos de columna que no se buscan como texto
var searchSkipTypes = []string{"JSON", "JSONB", "BYTEA", "TSVECTOR"}

/**
* RunQuery: Arma la consulta sobre model a partir de query (contrato de jdb.RunQuery) y devuelve la pagina pedida con el total.
* query admite wheres, selects, order_by ({ asc: [], desc: [] }), page y rows (limit es alias de rows; 100 por defecto,
* maximo 1000) y search ({ value, fields }, ver querySearch). Las condiciones de wheres van entre parentesis si
* traen or; wheres (opcional) son condiciones fijas del llamador que se aplican con AND sobre ellas, y search se
* aplica con AND despues de todo. En wheres se ignora la columna desconocida u oculta; en selects, order_by y
* search.fields da error. No soporta isNull, notNull, notIn ni between.
* @param model *Model, query et.Json, wheres ...*Where
* @return et.List, error
**/
func RunQuery(model *Model, query et.Json, wheres ...*Where) (et.List, error) {
	ql, err := queryLinq(model, query, wheres...)
	if err != nil {
		return et.List{}, err
	}

	page := queryPage(query)
	rows := queryRows(query)
	return ql.List(page, rows)
}

/**
* queryLinq: Arma el Linq de RunQuery (selects, wheres fijos, wheres de query, search y order_by) sin ejecutarlo.
* @param model *Model, query et.Json, wheres ...*Where
* @return *Linq, error
**/
func queryLinq(model *Model, query et.Json, wheres ...*Where) (*Linq, error) {
	selectNames := query.ArrayStr("selects")
	selected, err := querySelects(model, selectNames)
	if err != nil {
		return nil, err
	}

	ql := model.Data(selected)
	for _, where := range wheres {
		ql.Where(where)
	}

	queryWheres := query.ArrayJson("wheres")
	conditions := queryConditions(model, queryWheres)
	ql = applyConditions(ql, conditions)

	search, err := querySearch(model, query.Json("search"))
	if err != nil {
		return nil, err
	}
	for _, where := range search {
		ql.And(where)
	}

	orderBy := query.Json("order_by")
	return queryOrderBy(model, ql, orderBy)
}

/**
* querySearch: Condiciones de search = { value, fields } (contrato de los /query de catalog): cada palabra de value
* (hasta SEARCH_MAX_WORDS, separadas por espacios) debe aparecer, sin distinguir mayusculas (ILIKE), en los campos
* unidos con un espacio (concat), asi "ana gomez" encuentra nombre "Ana" y apellido "Gomez". Cada palabra es una
* condicion AND, de modo que la busqueda no salta los wheres fijos. % y _ de value se buscan como texto. Sin value
* (o si search no es un objeto) no agrega nada. Error si value pasa de SEARCH_MAX_LENGTH, si fields esta vacio o pasa
* de SEARCH_MAX_FIELDS, o si un campo no existe, es oculto o es JSON/binario.
* @param model *Model, search et.Json
* @return []*Where, error
**/
func querySearch(model *Model, search et.Json) ([]*Where, error) {
	value := strings.TrimSpace(search.Str("value"))
	if value == "" {
		return nil, nil
	}
	if len([]rune(value)) > SEARCH_MAX_LENGTH {
		return nil, fmt.Errorf(MSG_SEARCH_LENGTH, SEARCH_MAX_LENGTH)
	}

	fields := search.ArrayStr("fields")
	if len(fields) == 0 || len(fields) > SEARCH_MAX_FIELDS {
		return nil, fmt.Errorf(MSG_SEARCH_FIELDS, SEARCH_MAX_FIELDS)
	}
	for _, name := range fields {
		col := queryColumn(model, name)
		if col == nil || (col.Tp != TpColumn && col.Tp != TpAtrib) || slices.Contains(searchSkipTypes, strings.ToUpper(col.Type)) {
			return nil, fmt.Errorf(MSG_SEARCH_FIELD_NOT_VALID, name)
		}
	}

	words := strings.Fields(value)
	if len(words) > SEARCH_MAX_WORDS {
		words = words[:SEARCH_MAX_WORDS]
	}

	result := []*Where{}
	for _, word := range words {
		col := queryConcat(model, fields)
		pattern, _ := queryValue("%" + likeEscaper.Replace(word) + "%")
		result = append(result, col.Like(pattern))
	}

	return result, nil
}

/**
* queryPage: Pagina pedida en query; 1 si no viene o es menor.
* @param query et.Json
* @return int
**/
func queryPage(query et.Json) int {
	page := query.Int("page")
	if page < 1 {
		page = 1
	}

	return page
}

/**
* queryRows: Filas por pagina de query (rows o su alias limit), QUERY_DEFAULT_ROWS si no viene y como maximo QUERY_MAX_ROWS.
* @param query et.Json
* @return int
**/
func queryRows(query et.Json) int {
	rows := query.Int("rows")
	if rows < 1 {
		rows = query.Int("limit")
	}
	if rows < 1 {
		rows = QUERY_DEFAULT_ROWS
	}
	if rows > QUERY_MAX_ROWS {
		rows = QUERY_MAX_ROWS
	}

	return rows
}

/**
* queryColumn: Columna visible de model con ese nombre; nil si no existe o es oculta.
* @param model *Model, name string
* @return *Column
**/
func queryColumn(model *Model, name string) *Column {
	col := model.Column(name)
	if col == nil || col.Hidden {
		return nil
	}

	return col
}

/**
* querySelects: Columnas de selects; error si alguna no existe o es oculta. Sin selects devuelve una lista vacia (todas).
* @param model *Model, names []string
* @return []*Column, error
**/
func querySelects(model *Model, names []string) ([]*Column, error) {
	result := []*Column{}
	for _, name := range names {
		col := queryColumn(model, name)
		if col == nil {
			return nil, fmt.Errorf(MSG_QUERY_COLUMN_NOT_VALID, name)
		}

		result = append(result, col)
	}

	return result, nil
}

/**
* queryOrderBy: Aplica order_by ({ asc: [], desc: [] }) a ql; error si una columna no existe o es oculta.
* @param model *Model, ql *Linq, orderBy et.Json
* @return *Linq, error
**/
func queryOrderBy(model *Model, ql *Linq, orderBy et.Json) (*Linq, error) {
	for _, direction := range []string{"asc", "desc"} {
		names := orderBy.ArrayStr(direction)
		sorted := direction == "asc"
		for _, name := range names {
			col := queryColumn(model, name)
			if col == nil {
				return nil, fmt.Errorf(MSG_QUERY_COLUMN_NOT_VALID, name)
			}

			ql.OrderBy(col, sorted)
		}
	}

	return ql, nil
}

/**
* sortedKeys: Llaves de data en orden alfabetico, para que el SQL no dependa del orden del mapa.
* @param data et.Json
* @return []string
**/
func sortedKeys(data et.Json) []string {
	result := []string{}
	for key := range data {
		result = append(result, key)
	}
	sort.Strings(result)

	return result
}

/**
* queryConditions: Condiciones de wheres con su conector: { campo: {...} } va con AND y { and|or: { campo: {...} } }
* con ese conector (la primera condicion reconocida del grupo). Lo que no se reconoce se omite.
* @param model *Model, wheres []et.Json
* @return []*Where
**/
func queryConditions(model *Model, wheres []et.Json) []*Where {
	result := []*Where{}
	for _, where := range wheres {
		for _, column := range sortedKeys(where) {
			connector := strings.ToLower(column)
			if connector != "and" && connector != "or" {
				operators := where.Json(column)
				conditions := queryColumnConditions(model, "AND", column, operators)
				result = append(result, conditions...)
				continue
			}

			group := where.Json(column)
			for _, inner := range sortedKeys(group) {
				operators := group.Json(inner)
				conditions := queryColumnConditions(model, strings.ToUpper(connector), inner, operators)
				if len(conditions) > 0 {
					result = append(result, conditions...)
					break
				}
			}
		}
	}

	return result
}

/**
* queryColumnConditions: Condiciones de una columna (o de concat: { fields: [], like: "%texto%" }, los campos unidos
* con un espacio). La primera lleva connector y las demas AND. Vacio si la columna no existe o es oculta.
* @param model *Model, connector, name string, operators et.Json
* @return []*Where
**/
func queryColumnConditions(model *Model, connector, name string, operators et.Json) []*Where {
	var col *Column
	if strings.EqualFold(name, "concat") {
		col = queryConcat(model, operators.ArrayStr("fields"))
	} else {
		col = queryColumn(model, name)
	}

	result := []*Where{}
	if col == nil {
		return result
	}

	for _, operator := range sortedKeys(operators) {
		value := operators[operator]
		where := queryOperator(col, operator, value)
		if where == nil {
			continue
		}

		where.connector = "AND"
		if len(result) == 0 {
			where.connector = connector
		}
		result = append(result, where)
	}

	return result
}

/**
* queryConcat: Columna CONCAT de los campos visibles de fields separados por un espacio; nil si ninguno existe.
* @param model *Model, fields []string
* @return *Column
**/
func queryConcat(model *Model, fields []string) *Column {
	args := []any{}
	for _, name := range fields {
		col := queryColumn(model, name)
		if col == nil {
			continue
		}

		if len(args) > 0 {
			args = append(args, " ")
		}
		args = append(args, col)
	}

	if len(args) == 0 {
		return nil
	}

	return model.Concat(args...)
}

/**
* queryOperator: Condicion de col para el operador (eq, neg, in, like, more, less, moreEq, lessEq); nil si el
* operador no se soporta o el valor no es texto, numero o booleano.
* @param col *Column, operator string, value any
* @return *Where
**/
func queryOperator(col *Column, operator string, value any) *Where {
	if operator == "in" {
		values := queryValues(value)
		if len(values) == 0 {
			return nil
		}

		return col.In(values...)
	}

	sqlValue, ok := queryValue(value)
	if !ok {
		return nil
	}

	switch operator {
	case "eq":
		return col.Eq(sqlValue)
	case "neg":
		return col.Neg(sqlValue)
	case "like":
		return col.Like(sqlValue)
	case "more":
		return col.More(sqlValue)
	case "less":
		return col.Less(sqlValue)
	case "moreEq":
		return col.MoreEq(sqlValue)
	case "lessEq":
		return col.LessEq(sqlValue)
	default:
		return nil
	}
}

/**
* queryValues: Valores validos de la lista de in (se omiten los que no son texto, numero o booleano).
* @param value any
* @return []any
**/
func queryValues(value any) []any {
	result := []any{}
	list, ok := value.([]any)
	if !ok {
		return result
	}

	for _, item := range list {
		sqlValue, ok := queryValue(item)
		if ok {
			result = append(result, sqlValue)
		}
	}

	return result
}

/**
* queryValue: Valor listo para el SQL: el texto va entre comillas simples con las comillas internas duplicadas (no pasa
* por et.Unquote, que no las escapa); numeros y booleanos van tal cual. ok=false para otros tipos.
* @param value any
* @return any, bool
**/
func queryValue(value any) (any, bool) {
	switch v := value.(type) {
	case string:
		escaped := strings.ReplaceAll(v, "'", "''")
		quoted := fmt.Sprintf(`'%s'`, escaped)
		return SQL{val: quoted}, true
	case float64, float32, int, int32, int64, bool:
		return v, true
	default:
		return nil, false
	}
}

/**
* applyConditions: Agrega conditions a ql. Si alguna es OR van todas en un grupo entre parentesis, para que el OR no
* salte las condiciones fijas ni las que van antes.
* @param ql *Linq, conditions []*Where
* @return *Linq
**/
func applyConditions(ql *Linq, conditions []*Where) *Linq {
	hasOr := false
	for _, condition := range conditions {
		if condition.connector == "OR" {
			hasOr = true
			break
		}
	}

	if !hasOr {
		for _, condition := range conditions {
			ql.And(condition)
		}

		return ql
	}

	group := ""
	for i, condition := range conditions {
		condition.linq = ql
		def := condition.Define(ql).where
		if i == 0 {
			group = def
			continue
		}

		group = fmt.Sprintf(`%s %s %s`, group, condition.connector, def)
	}

	text := fmt.Sprintf(`(%s)`, group)
	where := NewWhere(SQL{val: text}, "", SQL{})
	ql.Where(where)

	return ql
}
