package linq

import (
	"strings"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/strs"
)

func (c *Model) DefineColum(name, description, _type string, _default any) *Model {
	name = strs.Uppcase(name)
	NewColumn(c, name, description, _type, _default)

	return c
}

func (c *Model) DefineAtrib(name, description, _type string, _default any) *Model {
	name = strs.Lowcase(name)
	source := NewColumn(c, SourceField.Upp(), "", "JSONB", "{}")
	result := NewColumn(c, name, description, _type, _default)
	result.Tp = TpAtrib
	result.Column = source
	source.Atribs = append(source.Atribs, result)

	return c
}

func (c *Model) DefineIndex(index []string) *Model {
	for _, name := range index {
		idx := c.ColIdx(name)
		if idx != -1 {
			c.Definition[idx].Indexed = true
			c.IndexAdd(name)
		}
	}

	return c
}

func (c *Model) DefineUniqueIndex(index []string) *Model {
	for _, name := range c.Index {
		col := c.Col(name)
		if col != nil {
			col.Unique = true
			col.Indexed = true
			c.IndexAdd(name)
		}
	}

	return c
}

func (c *Model) DefineHidden(hiddens []string) *Model {
	for _, key := range hiddens {
		col := c.Col(key)
		if col != nil {
			col.Hidden = true
		}
	}

	return c
}

func (c *Model) DefinePrimaryKey(keys []string) *Model {
	for _, name := range keys {
		col := c.Col(name)
		if col != nil {
			col.Required = true
			col.PrimaryKey = true
			c.PrimaryKeys = append(c.PrimaryKeys, name)
		}
	}

	return c
}

func (c *Model) DefineForeignKey(thisKey string, otherKey *Column) *Model {
	col := c.Col(thisKey)
	if col != nil {
		col.ForeignKey = true
		col.Reference = NewForeignKey(thisKey, otherKey)
		c.ForeignKey = append(c.ForeignKey, col.Reference)
		c.IndexAdd(thisKey)
		otherKey.ReferencesAdd(col)
	}

	return c
}

func (c *Model) DefineReference(thisKey, name, otherKey string, column *Column, showThisKey bool) *Model {
	if name == "" {
		name = thisKey
	}
	idxName := c.ColIdx(name)
	if idxName == -1 {
		col := NewColumn(c, name, "", "REFERENCE", et.Json{"_id": "", "name": ""})
		col.Tp = TpReference
		col.Title = name
		col.Reference = &Reference{thisKey, name, otherKey, column}
		idxThisKey := c.ColIdx(thisKey)
		if idxThisKey != -1 {
			c.Definition[idxThisKey].Hidden = !showThisKey
			c.Definition[idxThisKey].Indexed = true
			c.Definition[idxThisKey].Model.IndexAdd(c.Definition[idxThisKey].name)
			_otherKey := column.Model.Col(otherKey)
			if _otherKey != nil {
				_otherKey.ReferencesAdd(c.Definition[idxThisKey])
			}
		}
	}

	return c
}

func (c *Model) DefineCaption(thisKey, name, otherKey string, column *Column, _default any) *Model {
	if name == "" {
		name = thisKey
	}
	idx := c.ColIdx(name)
	if idx == -1 {
		col := NewColumn(c, name, "", "CAPTION", _default)
		col.Tp = TpCaption
		col.Title = name
		col.Reference = &Reference{thisKey, name, otherKey, column}
		idx := c.ColIdx(thisKey)
		if idx != -1 {
			c.Definition[idx].Indexed = true
			c.Definition[idx].Model.IndexAdd(c.Definition[idx].name)
			_otherKey := column.Model.Col(otherKey)
			if _otherKey != nil {
				_otherKey.ReferencesAdd(c.Definition[idx])
			}
		}
	}

	return c
}

func (c *Model) DefineField(name, description string, _default any, definition string) *Model {
	result := NewColumn(c, name, "", "FIELD", _default)
	result.Tp = TpField
	result.Definition = definition

	return c
}

func (c *Model) DefineRequired(names []string) *Model {
	for _, name := range names {
		list := strings.Split(name, ":")
		key := list[0]
		col := c.Col(key)
		if col != nil {
			col.Required = true
		}

		if len(list) > 1 {
			msg := list[1]
			if msg == "" {
				col.RequiredMsg = msg
			}
		} else {
			col.RequiredMsg = strs.Format(msg.MSG_ATRIB_REQUIRED, col.name)
		}
	}

	return c
}

/**
* DefineEventError
* @param fn Event
* @return *Model
**/
func (c *Model) DefineEventError(fn Event) *Model {
	c.EventError = fn

	return c
}

/**
* DefineEventInsert
* @param fn Event
* @return *Model
**/
func (c *Model) DefineEventInsert(fn Event) *Model {
	c.EventInsert = fn

	return c
}

/**
* DefineEventUpdate
* @param fn Event
* @return *Model
**/
func (c *Model) DefineEventUpdate(fn Event) *Model {
	c.EventUpdate = fn

	return c
}

/**
* DefineEventDelete
* @param fn Event
* @return *Model
**/
func (c *Model) DefineEventDelete(fn Event) *Model {
	c.EventDelete = fn

	return c
}
