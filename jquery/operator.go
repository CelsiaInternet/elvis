package jquery

/**
* Operator: operador soportado dentro de un bloque "wheres". El valor
* string es exactamente la clave que se espera recibir en el JSON del
* query (p.ej. `{"age": {"more": 45}}` usa el operador MORE).
**/
type Operator string

const (
	EQ          Operator = "eq"
	NEG         Operator = "neg"
	LESS        Operator = "less"
	LESS_EQ     Operator = "less_eq"
	MORE        Operator = "more"
	MORE_EQ     Operator = "more_eq"
	LIKE        Operator = "like"
	IN          Operator = "in"
	NOT_IN      Operator = "not_in"
	IS          Operator = "is"
	IS_NOT      Operator = "is_not"
	NULL        Operator = "null"
	NOT_NULL    Operator = "not_null"
	BETWEEN     Operator = "between"
	NOT_BETWEEN Operator = "not_between"
)

/**
* operatorSymbols: mapea los operadores de comparacion simple
* (columna <simbolo> valor) a su simbolo SQL. Los demas operadores
* (LIKE, IN/NOT_IN, NULL/NOT_NULL, BETWEEN/NOT_BETWEEN) tienen su
* propia forma de renderizado en buildCondition, porque no encajan en
* el patron "columna <simbolo> valor".
**/
var operatorSymbols = map[Operator]string{
	EQ:      "=",
	NEG:     "!=",
	LESS:    "<",
	LESS_EQ: "<=",
	MORE:    ">",
	MORE_EQ: ">=",
	IS:      "IS",
	IS_NOT:  "IS NOT",
}

/**
* IsValidOperator
* @param op Operator
* @return bool
**/
func IsValidOperator(op Operator) bool {
	switch op {
	case EQ, NEG, LESS, LESS_EQ, MORE, MORE_EQ, LIKE, IN, NOT_IN, IS, IS_NOT, NULL, NOT_NULL, BETWEEN, NOT_BETWEEN:
		return true
	default:
		return false
	}
}
