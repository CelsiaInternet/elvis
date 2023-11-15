package master

import (
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
)

func (c *Node) InsertValues(data e.Json) (fields, values string) {
	for k, v := range data {
		k = utility.Uppcase(k)
		v = e.Quoted(v)

		if len(fields) == 0 {
			fields = k
			values = utility.Format(`%v`, v)
		} else {
			fields = utility.Format(`%s, %s`, fields, k)
			values = utility.Format(`%s, %v`, values, v)
		}
	}

	return
}

func (c *Node) UpsertValues(data e.Json) (fields, values, fieldValue string) {
	for k, v := range data {
		k = utility.Uppcase(k)
		v = e.Quoted(v)

		if len(fieldValue) == 0 {
			fields = k
			values = utility.Format(`%v`, v)
			if k == "_IDT" {
				v = "-1"
			}
			fieldValue = utility.Format(`%s=%v`, k, v)
		} else {
			fields = utility.Format(`%s, %s`, fields, k)
			values = utility.Format(`%s, %v`, values, v)
			if k == "_IDT" {
				v = "-1"
			}
			fieldValue = utility.Append(fieldValue, utility.Format(`%s=%v`, k, v), ",\n")
		}
	}

	return
}

func (c *Node) SqlField(schema, table string, data e.Json) string {
	fields, _ := c.InsertValues(data)
	result := utility.Format(`INSERT INTO %s.%s (%s)`, utility.Lowcase(schema), utility.Uppcase(table), fields)
	result = utility.Append(result, "VALUES", "\n")
	return result
}

func (c *Node) ToSql(schema, table, idT string, data e.Json, action string) (string, bool) {
	var result string
	var ok bool
	if action == "INSERT" {
		_, values := c.InsertValues(data)
		result = utility.Format(`(%s)`, values)
	} else if action == "UPDATE" {
		fields, values, fieldValue := c.UpsertValues(data)
		result = utility.Format(`INSERT INTO %s.%s (%s)`, utility.Lowcase(schema), utility.Uppcase(table), fields)
		result = utility.Append(result, utility.Format(`VALUES (%s)`, values), "\n")
		result = utility.Append(result, "ON CONFLICT (_IDT) DO UPDATE SET", "\n")
		result = utility.Append(result, fieldValue, "\n")
		result = utility.Format(`%s;`, result)
	} else if action == "DELETE" {
		result = utility.Format(`DELETE FROM %s.%s WHERE _IDT=%s`, utility.Lowcase(schema), utility.Uppcase(table), idT)
	}

	return result, ok
}
