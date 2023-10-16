package linq

import (
	"strings"

	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
)

type Where struct {
	linq      *Linq
	connector string
	where     string
	val1      any
	operator  string
	val2      any
}

func StrToColN(str string) []string {
	str = ReplaceAll(str, []string{" "}, "")
	cols := strings.Split(str, ",")

	return cols
}

func (c *Linq) Col(val any) *Column {
	switch v := val.(type) {
	case Column:
		return &v
	case *Column:
		return v
	default:
		return &Column{}
	}
}

func (c *Where) Def(val any) string {
	switch v := val.(type) {
	case Column:
		as := v.As(c.linq)
		return Append(as, v.cast, "::")
	case *Column:
		as := v.As(c.linq)
		return Append(as, v.cast, "::")
	case Col:
		as := v.from
		as = Append(as, v.name, ".")
		return Append(as, v.cast, "::")
	case *Col:
		as := v.from
		as = Append(as, v.name, ".")
		return Append(as, v.cast, "::")
	case SQL:
		return Format(`%v`, v.val)
	default:
		return Format(`%v`, Quoted(v))
	}
}

func (c *Where) Define(linq *Linq) *Where {
	var where string

	result := ""
	switch v := c.val1.(type) {
	case []any:
		for _, vl := range v {
			def := c.Def(vl)
			result = Append(result, def, ",")
		}
		result = Format(`(%s)`, result)
	default:
		result = c.Def(v)
	}
	where = Format(`%s %s`, result, c.operator)

	result = ""
	switch v := c.val2.(type) {
	case []any:
		for _, vl := range v {
			def := c.Def(vl)
			result = Append(result, def, ",")
		}
		result = Format(`(%s)`, result)
	default:
		result = c.Def(v)
	}
	where = Format(`%s %s`, where, result)

	c.where = where

	return c
}

func NewWhere(val1 any, operator string, val2 any) *Where {
	return &Where{val1: val1, operator: operator, val2: val2}
}
