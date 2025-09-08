package linq

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
)

const BeforeInsert = 1
const AfterInsert = 2
const BeforeUpdate = 3
const AfterUpdate = 4
const BeforeDelete = 5
const AfterDelete = 6

type Trigger func(model *Model, old, new *et.Json, data et.Json) error

type Listener func(data et.Json)

type Model struct {
	db                 *jdb.DB
	Name               string
	Description        string
	Define             string
	Functions          string
	Schema             *Schema
	Table              string
	Definition         []*Column
	PrimaryKeys        []string
	ForeignKey         []*Reference
	Index              []string
	SourceField        *Column
	DateMakeField      *Column
	DateUpdateField    *Column
	SerieField         *Column
	CodeField          *Column
	ProjectField       *Column
	StateField         *Column
	IdTFiled           *Column
	Ddl                string
	DdlIndex           string
	integrityAtrib     bool
	integrityReference bool
	indexeSource       bool
	UseDateMake        bool
	UseDateUpdate      bool
	UseState           bool
	UseProject         bool
	UseSource          bool
	UseSerie           bool
	BeforeInsert       []Trigger
	AfterInsert        []Trigger
	BeforeUpdate       []Trigger
	AfterUpdate        []Trigger
	BeforeDelete       []Trigger
	AfterDelete        []Trigger
	OnListener         Listener
	EventError         Event
	EventInsert        Event
	EventUpdate        Event
	EventDelete        Event
	Version            int
	mutation           bool
}

func NewModel(schema *Schema, name, description string, version int) *Model {
	name = strs.Uppcase(name)
	table := strs.Append(schema.Name, name, ".")
	result := &Model{
		db:                 schema.db,
		Schema:             schema,
		Name:               name,
		Description:        description,
		Table:              table,
		Version:            version,
		integrityReference: true,
		indexeSource:       true,
		mutation:           false,
	}

	result.BeforeInsert = append(result.BeforeInsert, beforeInsert)
	result.AfterInsert = append(result.AfterInsert, afterInsert)
	result.BeforeUpdate = append(result.BeforeUpdate, beforeUpdate)
	result.AfterUpdate = append(result.AfterUpdate, afterUpdate)
	result.BeforeDelete = append(result.BeforeDelete, beforeDelete)
	result.AfterDelete = append(result.AfterDelete, afterDelete)

	schema.Models = append(schema.Models, result)
	models = append(models, result)

	return result
}

func Mutation(schema *Schema, name, description string, version int) *Model {
	result := NewModel(schema, name, description, version)

	return result
}

func Table(schema, name string) *Model {
	schema = strs.Uppcase(schema)
	name = strs.Uppcase(name)

	for _, model := range models {
		if model.Schema.Name == schema && model.Name == name {
			return model
		}
	}

	return nil
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
		"serieField":         c.SerieField,
		"codeField":          c.CodeField,
		"projectField":       c.ProjectField,
		"integrityAtrib":     c.integrityAtrib,
		"integrityReference": c.integrityReference,
		"indexeSource":       c.indexeSource,
		"useDateMake":        c.UseDateMake,
		"useDateUpdate":      c.UseDateUpdate,
		"useState":           c.UseState,
		"useProject":         c.UseProject,
		"useSerie":           c.UseSerie,
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
		} else if col.Up() == SourceField.Upp() {
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
* Db
* @return *jdb.DB
**/
func (c *Model) Db() *jdb.DB {
	return c.db
}

/**
* DDL
**/
func (c *Model) Init() error {
	c.Define = c.DDL()
	c.Functions = c.DDLFunction()

	exists, err := jdb.ExistTable(c.db, c.Schema.Name, c.Name)
	if err != nil {
		return err
	}

	if exists {
		err = c.db.Ddl(c.Functions)
		if err != nil {
			return err
		}

		return nil
	}

	err = c.db.Ddl(c.Define)
	if err != nil {
		return err
	}

	return nil
}

func (c *Model) DDL() string {
	return ddlTable(c)
}

func (c *Model) DDLFunction() string {
	return ddlFunctions(c)
}

func (c *Model) DDLMigration() string {
	return ddlMigration(c)
}

func (c *Model) DropDDL() string {
	return strs.Format(`DROP TABLE IF EXISTS %s CASCADE;`, c.Table)
}

func (c *Model) Trigger(event int, trigger Trigger) {
	switch event {
	case BeforeInsert:
		c.BeforeInsert = append(c.BeforeInsert, trigger)
	case AfterInsert:
		c.AfterInsert = append(c.AfterInsert, trigger)
	case BeforeUpdate:
		c.BeforeUpdate = append(c.BeforeUpdate, trigger)
	case AfterUpdate:
		c.AfterUpdate = append(c.AfterUpdate, trigger)
	case BeforeDelete:
		c.BeforeDelete = append(c.BeforeDelete, trigger)
	case AfterDelete:
		c.AfterDelete = append(c.BeforeDelete, trigger)
	}
}

func (c *Model) Details(name, description string, _default any, details Details) {
	col := NewColumn(c, name, "", "DETAIL", _default)
	col.Tp = TpDetail
	col.Hidden = true
	col.Details = details
}

/**
* NextCode
* @param tag string
* @param prefix string
* @return string
**/
func (c *Model) NexCode(tag, prefix string) string {
	return jdb.NextCode(c.db, tag, prefix)
}

/**
* NextSerie
* @param tag string
* @return int64
**/
func (c *Model) NextSerie(tag string) int64 {
	return jdb.NextSerie(c.db, tag)
}

/**
* Up
* @return string
**/
func (c *Model) Up() string {
	return strs.Uppcase(c.Table)
}

/**
* Low
* @return string
**/
func (c *Model) Low() string {
	return strs.Lowcase(c.Table)
}

func (c *Model) ColIdx(name string) int {
	for i, item := range c.Definition {
		if item.Up() == strs.Uppcase(name) {
			return i
		}
	}

	return -1
}

func (c *Model) Column(name string) *Column {
	idx := c.ColIdx(name)

	if idx == -1 {
		return nil
	}

	return c.Definition[idx]
}

func (c *Model) Col(name string) *Column {
	return c.Column(name)
}

func (c *Model) C(names string) *Column {
	return c.Column(names)
}

func (c *Model) As(as string) *FRom {
	return &FRom{
		model: c,
		as:    as,
	}
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
	source := c.Col(SourceField.Upp())
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

func (c *Model) A(names string) *Column {
	return c.Atrib(names)
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

func (c *Model) IntegrityAtrib(ok bool) *Model {
	c.integrityAtrib = ok

	return c
}

func (c *Model) IntegrityReference(ok bool) *Model {
	c.integrityReference = ok

	return c
}

func (c *Model) IndexSource(ok bool) *Model {
	c.indexeSource = ok

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

func (c *Model) Query(sql string, args ...any) (et.Items, error) {
	return c.db.Query(sql, args...)
}
