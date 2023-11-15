package linq

import (
	"strings"

	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
)

type Where struct {
	linq      *Linq
	connector string
	where     string
	val1      any
	operator  string
	val2      any
}

func (c *Where) Str1() string {
	result := ""
	switch v := c.val1.(type) {
	case []any:
		for _, vl := range v {
			def := c.Def(vl)
			result = utility.Append(result, def, ",")
		}
		result = utility.Format(`(%s)`, result)
	default:
		result = c.Def(v)
	}

	return result
}

func (c *Where) Str2() string {
	result := ""
	switch v := c.val2.(type) {
	case []any:
		for _, vl := range v {
			def := c.Def(vl)
			result = utility.Append(result, def, ",")
		}
		result = utility.Format(`(%s)`, result)
	default:
		result = c.Def(v)
	}

	return result
}

func StrToCols(str string) []string {
	str = utility.ReplaceAll(str, []string{" "}, "")
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
		return utility.Append(as, v.cast, "::")
	case *Column:
		as := v.As(c.linq)
		return utility.Append(as, v.cast, "::")
	case Col:
		as := v.from
		as = utility.Append(as, v.name, ".")
		return utility.Append(as, v.cast, "::")
	case *Col:
		as := v.from
		as = utility.Append(as, v.name, ".")
		return utility.Append(as, v.cast, "::")
	case SQL:
		return utility.Format(`%v`, v.val)
	default:
		return utility.Format(`%v`, e.Quoted(v))
	}
}

func (c *Where) Define(linq *Linq) *Where {
	var where string

	result := c.Str1()
	where = utility.Format(`%s %s`, result, c.operator)
	result = c.Str2()
	where = utility.Format(`%s %s`, where, result)

	c.where = where

	return c
}

func NewWhere(val1 any, operator string, val2 any) *Where {
	return &Where{val1: val1, operator: operator, val2: val2}
}
