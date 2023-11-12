package linq

import (
	"reflect"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utility"
)

func FunctionDef(linq *Linq, col *Column) string {
	var result string

	switch args := col.Definition.(type) {
	case []interface{}:
		for _, arg := range args {
			def := ""
			switch v := arg.(type) {
			case Column:
				def = v.As(linq)
				def = Append(def, v.cast, "::")
			case *Column:
				def = v.As(linq)
				def = Append(def, v.cast, "::")
			case Col:
				def = v.from
				def = Append(def, v.Up(), ".")
				def = Append(def, v.cast, "::")
			case *Col:
				def = v.from
				def = Append(def, v.Up(), ".")
				def = Append(def, v.cast, "::")
			case string:
				def = Format(`%v`, Quoted(v))
			default:
				console.ErrorF(`FunctionDef:%s; value:%v`, reflect.TypeOf(v), v)
			}
			result = Append(result, def, ", ")
		}
	default:
		console.ErrorF(`FunctionDef:%s; value:%v`, reflect.TypeOf(args), args)
	}

	if len(result) > 0 {
		result = Format(`%s(%s)`, col.Function, result)
	}

	return result
}

func Concat(args ...any) *Column {
	result := &Column{
		Tp:         TpFunction,
		Definition: args,
		Function:   "CONCAT",
	}

	return result
}

/**
*
**/
func (c *Model) Concat(args ...any) *Column {
	return Concat(args...)
}

func (c *Model) Avg(args ...any) *Column {
	result := &Column{
		Tp:         TpFunction,
		Definition: args,
		Function:   "AVG",
	}

	return result
}

func (c *Model) Count(args ...any) *Column {
	result := &Column{
		Tp:         TpFunction,
		Definition: args,
		Function:   "COUNT",
	}

	return result
}

func (c *Model) Sum(args ...any) *Column {
	result := &Column{
		Tp:         TpFunction,
		Definition: args,
		Function:   "SUM",
	}

	return result
}

func (c *Model) Max(args ...any) *Column {
	result := &Column{
		Tp:         TpFunction,
		Definition: args,
		Function:   "MAX",
	}

	return result
}

func (c *Model) Min(args ...any) *Column {
	result := &Column{
		Tp:         TpFunction,
		Definition: args,
		Function:   "MIN",
	}

	return result
}
