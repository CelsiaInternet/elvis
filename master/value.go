package master

import (
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utility"
)

func (c *Node) InsertValues(data Json) (fields, values string) {
	for k, v := range data {
		k = Uppcase(k)
		v = Quoted(v)

		if len(fields) == 0 {
			fields = k
			values = Format(`%v`, v)
		} else {
			fields = Format(`%s, %s`, fields, k)
			values = Format(`%s, %v`, values, v)
		}
	}

	return
}

func (c *Node) UpsertValues(data Json) (fields, values, fieldValue string) {
	for k, v := range data {
		k = Uppcase(k)
		v = Quoted(v)

		if len(fieldValue) == 0 {
			fields = k
			values = Format(`%v`, v)
			if k == "_IDT" {
				v = "-1"
			}
			fieldValue = Format(`%s=%v`, k, v)
		} else {
			fields = Format(`%s, %s`, fields, k)
			values = Format(`%s, %v`, values, v)
			if k == "_IDT" {
				v = "-1"
			}
			fieldValue = Append(fieldValue, Format(`%s=%v`, k, v), ",\n")
		}
	}

	return
}

func (c *Node) SqlField(schema, table string, data Json) string {
	fields, _ := c.InsertValues(data)
	result := Format(`INSERT INTO %s.%s (%s)`, Lowcase(schema), Uppcase(table), fields)
	result = Append(result, "VALUES", "\n")
	return result
}

func (c *Node) ToSql(schema, table, idT string, data Json, action string) (string, bool) {
	var result string
	var ok bool
	if action == "INSERT" {
		_, values := c.InsertValues(data)
		result = Format(`(%s)`, values)
	} else if action == "UPDATE" {
		fields, values, fieldValue := c.UpsertValues(data)
		result = Format(`INSERT INTO %s.%s (%s)`, Lowcase(schema), Uppcase(table), fields)
		result = Append(result, Format(`VALUES (%s)`, values), "\n")
		result = Append(result, "ON CONFLICT (_IDT) DO UPDATE SET", "\n")
		result = Append(result, fieldValue, "\n")
		result = Format(`%s;`, result)
	} else if action == "DELETE" {
		result = Format(`DELETE FROM %s.%s WHERE _IDT=%s`, Lowcase(schema), Uppcase(table), idT)
	}

	return result, ok
}
