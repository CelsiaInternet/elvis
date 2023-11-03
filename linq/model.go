package linq

import (
	"strings"

	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
)

const BeforeInsert = 1
const AfterInsert = 2
const BeforeUpdate = 3
const AfterUpdate = 4
const BeforeDelete = 5
const AfterDelete = 6

type Trigger func(model *Model, old, new *Json, data Json)

type Model struct {
	Db                 int
	Database           *Db
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
	AfterReferences    Rerences
	Version            int
}

func (c *Model) Describe() Json {
	var colums []Json = []Json{}
	for _, atrib := range c.Definition {
		colums = append(colums, atrib.Describe())
	}

	var foreignKey []Json = []Json{}
	for _, atrib := range c.ForeignKey {
		foreignKey = append(foreignKey, atrib.Describe())
	}

	var primaryKeys []string = []string{}
	for _, key := range c.PrimaryKeys {
		primaryKeys = append(primaryKeys, key)
	}

	var index []string = []string{}
	for _, key := range c.Index {
		index = append(index, key)
	}

	return Json{
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

func (c *Model) Model() Json {
	var result Json = Json{}
	for _, col := range c.Definition {
		if !ContainsInt([]int{TpColumn, TpAtrib, TpDetail}, col.Tp) {
			continue
		}

		if len(col.Atribs) > 0 {
			for _, atr := range col.Atribs {
				result.Set(atr.name, atr.Default)
			}

		} else if col.Type == "JSONB" {
			result.Set(col.name, Json{})
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
		Db:              schema.Db,
		Database:        schema.Database,
		Schema:          schema.Name,
		Name:            Append(Lowcase(schema.Name), Uppcase(table), "."),
		Description:     description,
		Table:           Uppcase(table),
		UseSync:         schema.UseSync,
		Version:         version,
		SourceField:     schema.SourceField,
		DateMakeField:   schema.DateMakeField,
		DateUpdateField: schema.DateUpdateField,
		IndexField:      schema.IndexField,
		CodeField:       schema.CodeField,
		ProjectField:    schema.ProjectField,
		StateField:      schema.StateField,
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
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = $1
		AND table_name = $2);`

	item, err := DBQueryOne(c.Db, sql, c.Schema, c.Table)
	if err != nil {
		return err
	}

	exists := item.Bool("exists")

	if !exists {
		sql = c.DDL()

		_, err := DBQDDL(c.Db, sql)
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
				index = append(index, column.DDLIndex())
			}
		}
	}

	if len(c.PrimaryKeys) > 0 {
		keys := Format(`PRIMARY KEY (%s)`, strings.Join(c.PrimaryKeys, ", "))
		fields = append(fields, Uppcase(keys))
	}

	for _, def := range index {
		result = Append(result, def, "\n")
	}

	_fields := ""
	for i, def := range fields {
		if i == 0 {
			def = Format("\n%s", def)
		}
		_fields = Append(_fields, def, ",\n")
	}

	str := Format(`CREATE TABLE IF NOT EXISTS %s(%s);`, c.Name, _fields)
	result = Append(str, result, "\n")

	c.Ddl = result

	return result
}

func (c *Model) DDLMigration() string {
	var fields []string

	table := c.Name
	c.Table = "NEW_TABLE"
	c.Name = Append(c.Schema, c.Table, ",")
	ddl := c.DDL()

	for _, column := range c.Definition {
		fields = append(fields, column.name)
	}

	insert := Format(`INSERT INTO %s(%s) SELECT %s FROM %s;`, c.Name, strings.Join(fields, ", "), strings.Join(fields, ", "), table)

	drop := Format(`DROP TABLE %s CASCADE;`, c.Name)

	alter := Format(`ALTER TABLE %s RENAME TO %s;`, c.Name, table)

	result := Format(`%s %s %s %s`, ddl, insert, drop, alter)

	return result
}

func (c *Model) DropDDL() string {
	return Format(`DROP TABLE IF EXISTS %s CASCADE;`, c.Name)
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

func (c *Model) References(reference Rerences) {
	c.AfterReferences = reference
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
	return Uppcase(c.Name)
}

func (c *Model) Low() string {
	return Lowcase(c.Name)
}

func (c *Model) ColIdx(name string) int {
	for i, item := range c.Definition {
		if item.Up() == Uppcase(name) {
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
		if Uppcase(item.Title) == Uppcase(name) {
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
		if Lowcase(item.name) == Lowcase(name) {
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
	result.name = Lowcase(name)
	source.Atribs = append(source.Atribs, result)

	return c
}

func (c *Model) DefineIndex(index []string) *Model {
	c.Index = index
	for _, key := range c.Index {
		col := c.Col(key)
		if col != nil {
			col.Indexed = true
		}
	}
	return c
}

func (c *Model) DefineHidden(hidden []string) *Model {
	for _, key := range hidden {
		col := c.Col(key)
		if col != nil {
			col.Hidden = true
		}
	}
	return c
}

func (c *Model) DefinePrimaryKey(keys []string) *Model {
	c.PrimaryKeys = keys
	for _, key := range c.PrimaryKeys {
		col := c.Col(key)
		if col != nil {
			col.Indexed = true
			col.Unique = true
			col.Required = true
			col.PrimaryKey = true
		}
	}

	return c
}

func (c *Model) DefineForeignKey(thisKey string, otherKey *Column) *Model {
	col := c.Col(thisKey)
	if col != nil {
		col.Indexed = true
		col.ForeignKey = true
	}
	c.ForeignKey = append(c.ForeignKey, NewForeignKey(thisKey, otherKey))

	return c
}

func (c *Model) DefineReference(thisKey, name, otherKey string, column *Column) *Model {
	if name == "" {
		name = thisKey
	}
	col := c.Col(name)
	if col == nil {
		col = NewColumn(c, name, "", "REFERENCE", Json{"_id": "", "name": ""})
		col.Tp = TpReference
	}
	col.Title = name
	col.Reference = &Reference{thisKey, name, otherKey, column}

	return c
}

func (c *Model) DefineField(name, description string, _default any, definition string) *Model {
	result := NewColumn(c, name, "", "FIELD", _default)
	result.Tp = TpField
	result.Definition = definition

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

func (c *Model) Select(sel ...any) *Linq {
	result := From(c)
	result.Select(sel...)

	return result
}

func (c *Model) Insert(data Json) *Linq {
	tp := TpSelect
	if c.UseSource {
		tp = TpData
	}

	result := NewLinq(tp, ActInsert, c)
	result.data = data

	return result
}

func (c *Model) Update(data Json) *Linq {
	tp := TpSelect
	if c.UseSource {
		tp = TpData
	}
	result := NewLinq(tp, ActUpdate, c)
	result.data = data

	return result
}

func (c *Model) Delete() *Linq {
	tp := TpSelect
	if c.UseSource {
		tp = TpData
	}
	result := NewLinq(tp, ActDelete, c)

	return result
}

func (c *Model) Upsert(data Json) *Linq {
	tp := TpSelect
	if c.UseSource {
		tp = TpData
	}
	result := NewLinq(tp, ActUpsert, c)
	result.data = data

	return result
}
