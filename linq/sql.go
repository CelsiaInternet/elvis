package linq

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
)

func (c *Linq) Sql() string {
	if c.Act == ActInsert {
		c.PrepareInsert()
		return c.SqlInsert()
	} else if c.Act == ActUpdate {
		c.PrepareUpdate(c.data)
		return c.SqlUpdate()
	} else if c.Act == ActDelete {
		c.PrepareDelete(c.data)
		return c.SqlDelete()
	} else {
		return c.SqlSelect()
	}
}

func (c *Linq) SQL() SQL {
	c.Sql()
	return SQL{
		val: Format(`(%s)`, c.sql),
	}
}

/**
*
**/
func (c *Linq) SqlColumDef(cols ...*Column) string {
	var result string

	if c.Tp == TpData {
		data := ""
		json := ""

		for _, col := range cols {
			def := col.Def(c)

			if col.name == col.Model.SourceField {
				data = col.As(c)
			} else {
				json = Append(json, def, ",\n")
			}
		}

		json = Format(`jsonb_build_object(%s)`, json)
		return Append(data, json, "||")
	}

	for _, col := range cols {
		def := col.Def(c)

		result = Append(result, def, ",")
	}

	return result
}

func (c *Linq) SqlColums(cols ...*Column) string {
	var result string
	n := len(cols)

	if c.Tp == TpData && n == 0 {
		for _, from := range c.from {
			for _, col := range from.model.Definition {
				if col.Tp != TpAtrib {
					cols = append(cols, col)
				}
				if col.Tp == TpDetail {
					c.details = append(c.details, col)
				}
			}

			res := c.SqlColumDef(cols...)
			result = Append(result, res, ",")
		}

		result = Format(`%s AS %s`, result, c.from[0].model.SourceField)
	} else if c.Tp == TpData {
		result = c.SqlColumDef(cols...)

		result = Format(`%s AS %s`, result, c.from[0].model.SourceField)
	} else if n > 0 {
		result = c.SqlColumDef(cols...)
	} else {
		result = "*"
	}

	return result
}

/**
*
**/
func (c *Linq) SqlSelect() string {
	result := c.SqlColums(c._select...)

	c.sql = Format(`SELECT %s`, result)

	c.SqlFrom()

	c.SqlJoin()

	c.SqlWhere()

	c.SqlGroupBy()

	c.SqlOrderBy()

	return c.sql
}

func (c *Linq) SqlReturn() string {
	result := c.SqlColums(c._return...)

	if len(result) > 0 {
		result = Format(`RETURNING %s`, result)
	}

	c.sql = Append(c.sql, result, "\n")

	return result
}

func (c *Linq) SqlCurrent() string {
	var result string
	var cols []*Column
	model := c.from[0].model

	for _, col := range model.Definition {
		if col.Tp == TpColumn {
			cols = append(cols, col)
		}
	}

	n := len(cols)

	if n > 0 {
		result = c.SqlColumDef(cols...)
		if c.Tp == TpData {
			result = Format(`%s AS %s`, result, c.from[0].model.SourceField)
		}
	} else {
		result = "*"
	}

	c.sql = Format(`SELECT %s`, result)

	c.SqlFrom()

	c.SqlKey()

	if c.debug {
		console.Log(c.sql)
	}

	return c.sql
}

func (c *Linq) SqlCount() string {
	c.sql = "SELECT COUNT(*) AS COUNT"

	c.SqlFrom()

	c.SqlJoin()

	c.SqlWhere()

	c.SqlGroupBy()

	return c.sql
}

func (c *Linq) SqlFrom() string {
	var result string
	for _, from := range c.from {
		result = Append(result, from.NameAs(), ", ")
	}

	result = Format(`FROM %s`, result)

	c.sql = Append(c.sql, result, "\n")

	return result
}

func (c *Linq) SqlJoin() string {
	var result string
	for _, join := range c._join {
		where := join.where.Define(c).where
		def := Append(join.join.model.Name, join.join.as, " AS ")
		def = Format(`%s %s ON %s`, join.kind, def, where)
		result = Append(result, def, "\n")
	}

	c.sql = Append(c.sql, result, "\n")

	return result
}

func (c *Linq) SqlWhere() string {
	var result string
	var wh string
	for _, where := range c.where {
		def := where.Define(c)
		if len(result) == 0 {
			wh = def.where
		} else {
			wh = Append(def.connector, def.where, " ")
		}
		result = Append(result, wh, "\n")
	}

	if len(result) > 0 {
		result = Format(`WHERE %s`, result)
	}

	c.sql = Append(c.sql, result, "\n")

	return result
}

func (c *Linq) SqlKey() string {
	var result string
	var wh string
	for _, where := range c.where {
		col := c.Col(where.val1)
		if col.PrimaryKey {
			def := where.Define(c)
			if len(result) == 0 {
				wh = def.where
			} else {
				wh = Append(def.connector, def.where, " ")
			}
			result = Append(result, wh, "\n")
		}
	}

	if len(result) > 0 {
		result = Format(`WHERE %s`, result)
	}

	c.sql = Append(c.sql, result, "\n")

	return result
}

func (c *Linq) SqlGroupBy() string {
	var result string
	for _, col := range c.groupBy {
		def := col.As(c)
		result = Append(result, def, ", ")
	}

	if len(result) > 0 {
		result = Format(`GROUP BY %s`, result)
	}

	c.sql = Append(c.sql, result, "\n")

	return result
}

func (c *Linq) SqlOrderBy() string {
	var result string
	var group string
	for _, order := range c.orderBy {
		if order.sorted {
			group = Format(`%s ASC`, order.colum.As(c))
		} else {
			group = Format(`%s DESC`, order.colum.As(c))
		}

		result = Append(result, group, ", ")
	}

	if len(result) > 0 {
		result = Format(`ORDER BY %s`, result)
	}

	c.sql = Append(c.sql, result, "\n")

	return result
}

func (c *Linq) SqlLimit(limit int) string {
	c.SqlSelect()

	result := Format(`LIMIT %d;`, limit)

	c.sql = Append(c.sql, result, "\n")

	if c.debug {
		console.Log(c.sql)
	}

	return c.sql
}

func (c *Linq) SqlOffset(limit, offset int) string {
	c.SqlSelect()

	result := Format(`LIMIT %d OFFSET %d;`, limit, offset)

	c.sql = Append(c.sql, result, "\n")

	if c.debug {
		console.Log(c.sql)
	}

	return c.sql
}

func (c *Linq) SqlIndex() string {
	var result string
	var cols []*Column = []*Column{}
	from := c.from[0].model
	if from.UseIndex {
		col := from.Col(from.IndexField)
		cols = append(cols, col)
	} else {
		for _, key := range from.PrimaryKeys {
			col := from.Col(key)
			cols = append(cols, col)
		}
	}

	result = c.SqlColumDef(cols...)
	if c.Tp == TpData {
		result = Format(`%s AS %s`, result, c.from[0].model.SourceField)
	}

	if len(result) > 0 {
		result = Format(`RETURNING %s`, result)
	}

	c.sql = Append(c.sql, result, "\n")

	return result
}

/**
*
**/
func (c *Linq) SqlInsert() string {
	model := c.from[0].model
	var fields string
	var values string

	for key, val := range *c.new {
		field := Uppcase(key)
		value := Quoted(val)

		fields = Append(fields, field, ", ")
		values = Append(values, Format(`%v`, value), ", ")
	}

	c.sql = Format("INSERT INTO %s(%s)\nVALUES (%s)", model.Name, fields, values)

	c.SqlReturn()

	c.sql = Format(`%s;`, c.sql)

	if c.debug {
		console.Log(c.sql)
	}

	return c.sql
}

func (c *Linq) SqlUpdate() string {
	model := c.from[0].model
	var fieldValues string

	for key, val := range *c.new {
		field := Uppcase(key)
		value := Quoted(val)

		if model.UseSource && field == Uppcase(model.SourceField) {
			vals := Uppcase(model.SourceField)
			atribs := c.new.Json(Lowcase(field))

			for ak, av := range atribs {
				ak = Lowcase(ak)
				av = DoubleQuoted(av)

				vals = Format(`jsonb_set(%s, '{%s}', '%v', true)`, vals, ak, av)
			}
			value = vals
		}

		fieldValue := Format(`%s=%v`, field, value)
		fieldValues = Append(fieldValues, fieldValue, ",\n")
	}

	c.sql = Format(`UPDATE %s AS A SET %s`, model.Name, fieldValues)

	c.SqlWhere()

	c.SetAs(model, "A")

	c.SqlReturn()

	c.sql = Format(`%s;`, c.sql)

	if c.debug {
		console.Log(c.sql)
	}

	return c.sql
}

func (c *Linq) SqlDelete() string {
	model := c.from[0].model

	c.sql = Format(`DELETE FROM %s`, model.Name)

	c.SqlWhere()

	c.SqlIndex()

	c.sql = Format(`%s;`, c.sql)

	if c.debug {
		console.Log(c.sql)
	}

	return c.sql
}
