package linq

import (
	j "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
)

type Reference struct {
	Fkey      string
	Name      string
	Key       string
	Reference *Column
}

func (c *Reference) Describe() j.Json {
	return j.Json{
		"foreignKey": c.Fkey,
		"title":      c.Name,
		"key":        c.Key,
		"reference":  c.Reference.describe(),
	}
}

func (c *Reference) DDL() string {
	table := c.Reference.Model.Name
	return utility.Format(`REFERENCES %s(%s)`, table, c.Reference.Up())
}

func NewForeignKey(fKey string, reference *Column) *Reference {
	return &Reference{Fkey: fKey, Key: reference.name, Reference: reference}
}
