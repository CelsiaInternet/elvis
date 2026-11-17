package linq

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/celsiainternet/elvis/et"
)

const (
	QUERY_DEFAULT_ROWS           = 100
	QUERY_MAX_ROWS               = 1000
	MSG_QUERY_COLUMN_NOT_VALID   = "Columna no valida en la consulta (%s)"
	MSG_QUERY_OPERATOR_NOT_VALID = "Operador no valido en la consulta (%s)"
	MSG_QUERY_VALUE_NOT_VALID    = "Valor no valido en la consulta (%s)"
	MSG_QUERY_KEY_NOT_VALID      = "Llave no valida en la consulta (%s)"
	SEARCH_MAX_FIELDS            = 5
	SEARCH_MAX_WORDS             = 5
	SEARCH_MAX_LENGTH            = 100
	MSG_SEARCH_LENGTH            = "La busqueda admite maximo %d caracteres"
	MSG_SEARCH_FIELDS            = "La busqueda requiere entre 1 y %d campos (search.fields)"
	MSG_SEARCH_FIELD_NOT_VALID   = "Campo no valido en la busqueda (%s)"
)

// Comodines de LIKE escapados con \ (el escape por defecto de Postgres) para buscarlos como texto
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// Operadores que admite wheres
var queryOperators = []string{"eq", "neg", "in", "like", "more", "less", "moreEq", "lessEq"}

// Llaves raiz que admite query
var queryRootKeys = []string{"wheres", "selects", "order_by", "search", "page", "rows", "limit"}

// Tipos de columna que no se buscan como texto
var searchSkipTypes = []string{"JSON", "JSONB", "BYTEA", "TSVECTOR"}

/**
* RunQuery: Arma la consulta sobre model a partir de query (contrato de jdb.RunQuery) y devuelve la pagina pedida con el total.
* query admite wheres, selects, order_by ({ asc: [], desc: [] }), page y rows (limit es alias de rows; 100 por defecto,
* maximo 1000) y search ({ value, fields }, ver querySearch). Las condiciones de wheres van entre parentesis si
* traen or; wheres (opcional) son condiciones fijas del llamador que se aplican con AND sobre ellas, y search se
* aplica con AND despues de todo. Da error (sin consultar la base) ante una llave raiz desconocida
* (MSG_QUERY_KEY_NOT_VALID), una columna desconocida u oculta en selects, order_by o wheres (tambien en grupos and/or
* y en concat; MSG_QUERY_COLUMN_NOT_VALID), un operador no soportado (MSG_QUERY_OPERATOR_NOT_VALID) o un valor de
* wheres que no es texto, numero ni booleano, o un in vacio (MSG_QUERY_VALUE_NOT_VALID), y en search ver querySearch.
* Los wheres fijos del llamador no se validan. No soporta isNull, notNull, notIn ni between.
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
* Antes valida las llaves raiz (queryKeys).
* @param model *Model, query et.Json, wheres ...*Where
* @return *Linq, error
**/
func queryLinq(model *Model, query et.Json, wheres ...*Where) (*Linq, error) {
	if err := queryKeys(query); err != nil {
		return nil, err
	}

	selectNames := query.ArrayStr("selects")
	selected, err := querySelects(model, selectNames)
	if err != nil {
		return nil, err
	}

	ql := model.Data(selected)
	for _, where := range wheres {
		ql.Where(where)
	}

	queryWheres, err := queryWheres(query)
	if err != nil {
		return nil, err
	}
	conditions, err := queryConditions(model, queryWheres)
	if err != nil {
		return nil, err
	}
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
* queryWheres: Lista de wheres de query como objetos; error si wheres no es una lista o trae algo que no es un objeto.
* @param query et.Json
* @return []et.Json, error
**/
func queryWheres(query et.Json) ([]et.Json, error) {
	value, ok := query["wheres"]
	if !ok || value == nil {
		return []et.Json{}, nil
	}

	list, ok := queryList(value)
	if !ok {
		return nil, fmt.Errorf(MSG_QUERY_VALUE_NOT_VALID, "wheres")
	}

	result := []et.Json{}
	for _, item := range list {
		where, ok := queryObject(item)
		if !ok {
			return nil, fmt.Errorf(MSG_QUERY_VALUE_NOT_VALID, "wheres")
		}
		result = append(result, where)
	}

	return result, nil
}

/**
* queryConditions: Condiciones de wheres con su conector: { campo: {...} } va con AND y { and|or: { campo: {...}, ... } }
* con ese conector en cada campo del grupo (en orden alfabetico). Error si un campo no existe o es oculto, si un grupo
* no es un objeto o esta vacio, o si un operador o un valor no son validos (ver queryColumnConditions).
* @param model *Model, wheres []et.Json
* @return []*Where, error
**/
func queryConditions(model *Model, wheres []et.Json) ([]*Where, error) {
	result := []*Where{}
	for _, where := range wheres {
		for _, column := range sortedKeys(where) {
			connector := strings.ToLower(column)
			if connector != "and" && connector != "or" {
				conditions, err := queryColumnConditions(model, "AND", column, where[column])
				if err != nil {
					return nil, err
				}
				result = append(result, conditions...)
				continue
			}

			group, ok := queryObject(where[column])
			if !ok || len(group) == 0 {
				return nil, fmt.Errorf(MSG_QUERY_VALUE_NOT_VALID, column)
			}

			for _, inner := range sortedKeys(group) {
				conditions, err := queryColumnConditions(model, strings.ToUpper(connector), inner, group[inner])
				if err != nil {
					return nil, err
				}
				result = append(result, conditions...)
			}
		}
	}

	return result, nil
}

/**
* queryColumnConditions: Condiciones de una columna (o de concat: { fields: [], like: "%texto%" }, los campos unidos
* con un espacio). La primera lleva connector y las demas AND. Error si la columna (o un campo de concat) no existe o
* es oculta, si value no es un objeto con al menos un operador, si un operador no se soporta o si su valor no es valido.
* @param model *Model, connector, name string, value any
* @return []*Where, error
**/
func queryColumnConditions(model *Model, connector, name string, value any) ([]*Where, error) {
	operators, ok := queryObject(value)
	if !ok {
		return nil, fmt.Errorf(MSG_QUERY_VALUE_NOT_VALID, name)
	}

	var col *Column
	if strings.EqualFold(name, "concat") {
		fields, ok := queryList(operators["fields"])
		if !ok || len(fields) == 0 {
			return nil, fmt.Errorf(MSG_QUERY_VALUE_NOT_VALID, name+".fields")
		}

		names := []string{}
		for _, field := range fields {
			fieldName, ok := field.(string)
			if !ok || queryColumn(model, fieldName) == nil {
				return nil, fmt.Errorf(MSG_QUERY_COLUMN_NOT_VALID, fmt.Sprint(field))
			}
			names = append(names, fieldName)
		}

		col = queryConcat(model, names)
		operators = operators.Clone()
		delete(operators, "fields")
	} else {
		col = queryColumn(model, name)
	}

	if col == nil {
		return nil, fmt.Errorf(MSG_QUERY_COLUMN_NOT_VALID, name)
	}

	if len(operators) == 0 {
		return nil, fmt.Errorf(MSG_QUERY_VALUE_NOT_VALID, name)
	}

	result := []*Where{}
	for _, operator := range sortedKeys(operators) {
		where, err := queryOperator(col, name, operator, operators[operator])
		if err != nil {
			return nil, err
		}

		where.connector = "AND"
		if len(result) == 0 {
			where.connector = connector
		}
		result = append(result, where)
	}

	return result, nil
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
* queryOperator: Condicion de col para el operador (eq, neg, in, like, more, less, moreEq, lessEq). Error
* MSG_QUERY_OPERATOR_NOT_VALID si el operador no se soporta y MSG_QUERY_VALUE_NOT_VALID si el valor no es texto,
* numero o booleano (en in, si no es una lista, esta vacia o algun elemento no es valido).
* @param col *Column, name, operator string, value any
* @return *Where, error
**/
func queryOperator(col *Column, name, operator string, value any) (*Where, error) {
	if !slices.Contains(queryOperators, operator) {
		return nil, fmt.Errorf(MSG_QUERY_OPERATOR_NOT_VALID, name+"."+operator)
	}

	if operator == "in" {
		values, ok := queryValues(value)
		if !ok {
			return nil, fmt.Errorf(MSG_QUERY_VALUE_NOT_VALID, name+"."+operator)
		}

		return col.In(values...), nil
	}

	sqlValue, ok := queryValue(value)
	if !ok {
		return nil, fmt.Errorf(MSG_QUERY_VALUE_NOT_VALID, name+"."+operator)
	}

	switch operator {
	case "eq":
		return col.Eq(sqlValue), nil
	case "neg":
		return col.Neg(sqlValue), nil
	case "like":
		return col.Like(sqlValue), nil
	case "more":
		return col.More(sqlValue), nil
	case "less":
		return col.Less(sqlValue), nil
	case "moreEq":
		return col.MoreEq(sqlValue), nil
	default:
		return col.LessEq(sqlValue), nil
	}
}

/**
* queryValues: Valores de la lista de in listos para el SQL; ok=false si value no es una lista, esta vacia o algun
* elemento no es texto, numero o booleano.
* @param value any
* @return []any, bool
**/
func queryValues(value any) ([]any, bool) {
	list, ok := queryList(value)
	if !ok || len(list) == 0 {
		return nil, false
	}

	result := []any{}
	for _, item := range list {
		sqlValue, ok := queryValue(item)
		if !ok {
			return nil, false
		}
		result = append(result, sqlValue)
	}

	return result, true
}

/**
* queryList: Elementos de value si es una lista (slice o arreglo de cualquier tipo, como []any o []string); ok=false si no lo es.
* @param value any
* @return []any, bool
**/
func queryList(value any) ([]any, bool) {
	if list, ok := value.([]any); ok {
		return list, true
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() || (rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array) {
		return nil, false
	}

	result := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}

	return result, true
}

/**
* queryObject: value como objeto (et.Json o map[string]any); ok=false si no lo es.
* @param value any
* @return et.Json, bool
**/
func queryObject(value any) (et.Json, bool) {
	switch v := value.(type) {
	case et.Json:
		return v, v != nil
	case map[string]any:
		return et.Json(v), v != nil
	default:
		return nil, false
	}
}

/**
* queryKeys: Error MSG_QUERY_KEY_NOT_VALID si query trae una llave raiz distinta de queryRootKeys.
* @param query et.Json
* @return error
**/
func queryKeys(query et.Json) error {
	for _, key := range sortedKeys(query) {
		if !slices.Contains(queryRootKeys, key) {
			return fmt.Errorf(MSG_QUERY_KEY_NOT_VALID, key)
		}
	}

	return nil
}

/**
* queryValue: Valor listo para el SQL: el texto va entre comillas simples con las comillas internas duplicadas (no pasa
* por et.Unquote, que no las escapa); numeros y booleanos van tal cual. ok=false para otros tipos (null, objetos, listas).
* @param value any
* @return any, bool
**/
func queryValue(value any) (any, bool) {
	switch v := value.(type) {
	case string:
		escaped := strings.ReplaceAll(v, "'", "''")
		quoted := fmt.Sprintf(`'%s'`, escaped)
		return SQL{val: quoted}, true
	case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, bool:
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
