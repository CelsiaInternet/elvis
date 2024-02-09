package linq

import (
	"strings"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
)

const BeforeInsert = 1
const AfterInsert = 2
const BeforeUpdate = 3
const AfterUpdate = 4
const BeforeDelete = 5
const AfterDelete = 6

type Trigger func(model *Model, old, new *et.Json, data et.Json) error

type Model struct {
	Db                 int
	Database           *jdb.Db
	Name               string
	Description        string
	Schema             string
	Table              string
	Definition         []*Column
	PrimaryKeys        []string
	ForeignKey         []*Reference
	Index              []string
	SourceField        string
	DateMakeField      string
	DateUpdateField    string
	IndexField         string
	CodeField          string
	ProjectField       string
	StateField         string
	Ddl                string
	integrityAtrib     bool
	integrityReference bool
	UseState           bool
	UseSource          bool
	UseDateMake        bool
	UseDateUpdate      bool
	UseProject         bool
	UseIndex           bool
	UseSync            bool
	UseRecycle         bool
	BeforeInsert       []Trigger
	AfterInsert        []Trigger
	BeforeUpdate       []Trigger
	AfterUpdate        []Trigger
	BeforeDelete       []Trigger
	AfterDelete        []Trigger
	Version            int
}

func (c *Model) Driver() string {
	return c.Database.Driver
}

func (c *Model) Describe() et.Json {
	var colums []et.Json = []et.Json{}
	for _, atrib := range c.Definition {
		colums = append(colums, atrib.Describe())
	}

	var foreignKey []et.Json = []et.Json{}
	for _, atrib := range c.ForeignKey {
		foreignKey = append(foreignKey, atrib.Describe())
	}

	var primaryKeys []string = append([]string{}, c.PrimaryKeys...)
	var index []string = append([]string{}, c.Index...)

	return et.Json{
		"name":               c.Name,
		"description":        c.Description,
		"schema":             c.Schema,
		"table":              c.Table,
		"colums":             colums,
		"primaryKeys":        primaryKeys,
		"foreignKeys":        foreignKey,
		"index":              index,
		"sourceField":        c.SourceField,
		"dateMakeField":      c.DateMakeField,
		"dateUpdateField":    c.DateUpdateField,
		"indexField":         c.IndexField,
		"codeField":          c.CodeField,
		"projectField":       c.ProjectField,
		"integrityAtrib":     c.integrityAtrib,
		"integrityReference": c.integrityReference,
		"useState":           c.UseState,
		"useDateMake":        c.UseDateMake,
		"useDateUpdate":      c.UseDateUpdate,
		"useProject":         c.UseProject,
		"useReciclig":        c.UseRecycle,
		"model":              c.Model(),
	}
}

func (c *Model) Model() et.Json {
	var result et.Json = et.Json{}
	for _, col := range c.Definition {
		if !utility.ContainsInt([]int{TpColumn, TpAtrib, TpDetail}, col.Tp) {
			continue
		}

		if len(col.Atribs) > 0 {
			for _, atr := range col.Atribs {
				result.Set(atr.name, atr.Default)
			}
		} else if col.name == c.SourceField {
			continue
		} else if col.Type == "JSON" && col.Default == "[]" {
			result.Set(col.name, []et.Json{})
		} else if col.Type == "JSON" {
			result.Set(col.name, et.Json{})
		} else if col.Type == "JSONB" && col.Default == "[]" {
			result.Set(col.name, []et.Json{})
		} else if col.Type == "JSONB" {
			result.Set(col.name, et.Json{})
		} else if col.Type == "TIMESTAMP" && col.Default == "NOW()" {
			result.Set(col.name, utility.Now())
		} else if col.Type == "TIMESTAMP" && col.Default == "NULL" {
			result.Set(col.name, nil)
		} else if col.Type == "BOOLEAN" && col.Default == "TRUE" {
			result.Set(col.name, true)
		} else if col.Type == "BOOLEAN" && col.Default == "FALSE" {
			result.Set(col.name, false)
		} else {
			result.Set(col.name, col.Default)
		}
	}

	return result
}

/**
*
**/
func NewModel(schema *Schema, table, description string, version int) *Model {
	result := &Model{
		Db:                 schema.Db,
		Database:           schema.Database,
		Schema:             schema.Name,
		Name:               strs.Append(strs.Lowcase(schema.Name), strs.Uppcase(table), "."),
		Description:        description,
		Table:              strs.Uppcase(table),
		UseSync:            schema.UseSync,
		Version:            version,
		SourceField:        schema.SourceField,
		DateMakeField:      schema.DateMakeField,
		DateUpdateField:    schema.DateUpdateField,
		IndexField:         schema.IndexField,
		CodeField:          schema.CodeField,
		ProjectField:       schema.ProjectField,
		StateField:         schema.StateField,
		integrityReference: true,
	}

	result.BeforeInsert = append(result.BeforeInsert, beforeInsert)
	result.AfterInsert = append(result.AfterInsert, afterInsert)
	result.BeforeUpdate = append(result.BeforeUpdate, beforeUpdate)
	result.AfterUpdate = append(result.AfterUpdate, afterUpdate)
	result.BeforeDelete = append(result.BeforeDelete, beforeDelete)
	result.AfterDelete = append(result.AfterDelete, afterDelete)

	schema.Models = append(schema.Models, result)

	return result
}

/**
* DDL
**/
func (c *Model) Init() error {
	exists, err := jdb.ExistTable(c.Db, c.Schema, c.Table)
	if err != nil {
		return err
	}

	if !exists {
		sql := c.DDL()

		_, err := jdb.DBQDDL(c.Db, sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Model) SetUseSync(val bool) *Model {
	c.UseSync = val

	return c
}

func (c *Model) DDL() string {
	var result string
	var fields []string
	var index []string

	for _, column := range c.Definition {
		if column.Tp == TpColumn {
			fields = append(fields, column.DDL())
			if column.Indexed {
				if column.Unique {
					index = append(index, column.DDLUniqueIndex())
				} else {
					index = append(index, column.DDLIndex())
				}
			}
		}
	}

	if len(c.PrimaryKeys) > 0 {
		keys := strs.Format(`PRIMARY KEY (%s)`, strings.Join(c.PrimaryKeys, ", "))
		fields = append(fields, strs.Uppcase(keys))
	}

	for _, def := range index {
		result = strs.Append(result, def, "\n")
	}

	_fields := ""
	for i, def := range fields {
		if i == 0 {
			def = strs.Format("\n%s", def)
		}
		_fields = strs.Append(_fields, def, ",\n")
	}

	str := strs.Format(`CREATE TABLE IF NOT EXISTS %s(%s);`, c.Name, _fields)
	result = strs.Append(str, result, "\n")

	c.Ddl = result

	return result
}

func (c *Model) DDLMigration() string {
	var fields []string

	table := c.Name
	c.Table = "NEW_TABLE"
	c.Name = strs.Append(c.Schema, c.Table, ",")
	ddl := c.DDL()

	for _, column := range c.Definition {
		fields = append(fields, column.name)
	}

	insert := strs.Format(`INSERT INTO %s(%s) SELECT %s FROM %s;`, c.Name, strings.Join(fields, ", "), strings.Join(fields, ", "), table)

	drop := strs.Format(`DROP TABLE %s CASCADE;`, c.Name)

	alter := strs.Format(`ALTER TABLE %s RENAME TO %s;`, c.Name, table)

	result := strs.Format(`%s %s %s %s`, ddl, insert, drop, alter)

	return result
}

func (c *Model) DropDDL() string {
	return strs.Format(`DROP TABLE IF EXISTS %s CASCADE;`, c.Name)
}

/**
*
**/
func (c *Model) Trigger(event int, trigger Trigger) {
	if event == BeforeInsert {
		c.BeforeInsert = append(c.BeforeInsert, trigger)
	} else if event == AfterInsert {
		c.AfterInsert = append(c.AfterInsert, trigger)
	} else if event == BeforeUpdate {
		c.BeforeUpdate = append(c.BeforeUpdate, trigger)
	} else if event == AfterUpdate {
		c.AfterUpdate = append(c.AfterUpdate, trigger)
	} else if event == BeforeDelete {
		c.BeforeDelete = append(c.BeforeDelete, trigger)
	} else if event == AfterDelete {
		c.AfterDelete = append(c.BeforeDelete, trigger)
	}
}

func (c *Model) Details(name, description string, _default any, details Details) {
	col := NewColumn(c, name, "", "DETAIL", _default)
	col.Tp = TpDetail
	col.Details = details
}

/**
*
**/
func (c *Model) Up() string {
	return strs.Uppcase(c.Name)
}

func (c *Model) Low() string {
	return strs.Lowcase(c.Name)
}

func (c *Model) ColIdx(name string) int {
	for i, item := range c.Definition {
		if item.Up() == strs.Uppcase(name) {
			return i
		}
	}

	return -1
}

func (c *Model) Col(name string) *Column {
	idx := c.ColIdx(name)
	if idx == -1 && !c.integrityAtrib {
		return NewVirtualAtrib(c, name, "", "text", "")
	} else if idx == -1 {
		return nil
	}

	return c.Definition[idx]
}

func (c *Model) As(as string) *FRom {
	return &FRom{
		model: c,
		as:    as,
	}
}

func (c *Model) Column(name string) *Column {
	return c.Col(name)
}

func (c *Model) TitleIdx(name string) int {
	for i, item := range c.Definition {
		if strs.Uppcase(item.Title) == strs.Uppcase(name) {
			return i
		}
	}

	return -1
}

func (c *Model) AtribIdx(name string) int {
	source := c.Col(c.SourceField)
	if source == nil {
		return -1
	}

	for i, item := range source.Atribs {
		if strs.Lowcase(item.name) == strs.Lowcase(name) {
			return i
		}
	}

	return -1
}

func (c *Model) Atrib(name string) *Column {
	idx := c.ColIdx(name)
	if idx == -1 {
		return nil
	}

	return c.Definition[idx]
}

func (c *Model) IndexIdx(name string) int {
	for i, _name := range c.Index {
		if strs.Uppcase(_name) == strs.Uppcase(name) {
			return i
		}
	}

	return -1
}

func (c *Model) IndexAdd(name string) int {
	idx := c.IndexIdx(name)
	if idx == -1 {
		c.Index = append(c.Index, name)
		idx = len(c.Index) - 1
	}

	return idx
}

func (c *Model) All() []*Column {
	result := c.Definition

	return result
}

/**
*
**/
func (c *Model) DefineColum(name, description, _type string, _default any) *Model {
	NewColumn(c, name, description, _type, _default)

	return c
}

func (c *Model) DefineAtrib(name, description, _type string, _default any) *Model {
	source := c.Col(c.SourceField)
	result := NewColumn(c, name, description, _type, _default)
	result.Tp = TpAtrib
	result.Column = source
	result.name = strs.Lowcase(name)
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
			col.Indexed = true
			col.Unique = true
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
	for _, name := range c.PrimaryKeys {
		col := c.Col(name)
		if col != nil {
			col.Unique = true
			col.Required = true
			col.PrimaryKey = true
			c.PrimaryKeys = append(c.PrimaryKeys, name)
			c.IndexAdd(name)
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

func (c *Model) DefineReference(thisKey, name, otherKey string, column *Column) *Model {
	if name == "" {
		name = thisKey
	}
	idx := c.ColIdx(name)
	if idx == -1 {
		col := NewColumn(c, name, "", "REFERENCE", et.Json{"_id": "", "name": ""})
		col.Tp = TpReference
		col.Title = name
		col.Reference = &Reference{thisKey, name, otherKey, column}
		idx := c.ColIdx(thisKey)
		if idx != -1 {
			c.Definition[idx].ReferenceKey = true
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

func (c *Model) IntegrityAtrib(ok bool) *Model {
	c.integrityAtrib = ok

	return c
}

func (c *Model) IntegrityReference(ok bool) *Model {
	c.integrityReference = ok

	return c
}

/**
*
**/
func (c *Model) From() *Linq {
	return From(c)
}

func (c *Model) Data(sel ...any) *Linq {
	result := From(c)
	if !c.UseSource {
		result.Select(sel...)
	} else {
		result.Data(sel...)
	}

	return result
}

func (c *Model) Select(sel ...any) *Linq {
	result := From(c)
	result.Select(sel...)

	return result
}

/**
*
**/
func (c *Model) Insert(data et.Json) *Linq {
	tp := TpRow
	if c.UseSource {
		tp = TpData
	}

	result := NewLinq(ActInsert, c)
	result.SetTp(tp)
	result.data = data

	return result
}

func (c *Model) Update(data et.Json) *Linq {
	tp := TpRow
	if c.UseSource {
		tp = TpData
	}

	result := NewLinq(ActUpdate, c)
	result.SetTp(tp)
	result.data = data

	return result
}

func (c *Model) Delete() *Linq {
	tp := TpRow
	if c.UseSource {
		tp = TpData
	}

	result := NewLinq(ActDelete, c)
	result.SetTp(tp)

	return result
}

func (c *Model) Upsert(data et.Json) *Linq {
	tp := TpRow
	if c.UseSource {
		tp = TpData
	}

	result := NewLinq(ActUpsert, c)
	result.SetTp(tp)
	result.data = data

	return result
}

/**
* Row
**/
func (c *Model) InsertRow(data et.Json) *Linq {
	tp := TpRow

	result := NewLinq(ActInsert, c)
	result.SetTp(tp)
	result.data = data

	return result
}

func (c *Model) UpdateRow(data et.Json) *Linq {
	tp := TpRow

	result := NewLinq(ActUpdate, c)
	result.SetTp(tp)
	result.data = data

	return result
}

func (c *Model) DeleteRow() *Linq {
	tp := TpRow

	result := NewLinq(ActDelete, c)
	result.SetTp(tp)

	return result
}

func (c *Model) UpsertRow(data et.Json) *Linq {
	tp := TpRow

	result := NewLinq(ActUpsert, c)
	result.SetTp(tp)
	result.data = data

	return result
}
