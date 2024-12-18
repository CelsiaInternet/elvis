package linq

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/strs"
)

type Reference struct {
	Fkey      string
	Name      string
	Key       string
	Reference *Column
}

func (c *Reference) Describe() et.Json {
	return et.Json{
		"foreignKey": c.Fkey,
		"title":      c.Name,
		"key":        c.Key,
		"reference":  c.Reference.describe(),
	}
}

func (c *Reference) DDL() string {
	return strs.Format(`REFERENCES %s(%s) ON DELETE CASCADE`, c.Reference.Model.Table, c.Reference.Up())
}

func NewForeignKey(fKey string, reference *Column) *Reference {
	return &Reference{Fkey: fKey, Key: reference.name, Reference: reference}
}
