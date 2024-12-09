package linq

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/strs"
)

type DefaultField string

var (
	SourceField     DefaultField = "_DATA"
	DateMakeField   DefaultField = "DATE_MAKE"
	DateUpdateField DefaultField = "DATE_UPDATE"
	SerieField      DefaultField = "INDEX"
	CodeField       DefaultField = "CODE"
	ProjectField    DefaultField = "PROJECT_ID"
	StateField      DefaultField = "_STATE"
	IdTFiled        DefaultField = "_IDT"
	schemas         []*Schema    = []*Schema{}
	models          []*Model     = []*Model{}
)

func (s DefaultField) Low() string {
	return strs.Lowcase(string(s))
}

func (s DefaultField) Upp() string {
	return strs.Uppcase(string(s))
}

func (s DefaultField) Origin() string {
	return string(s)
}

type Schema struct {
	db          *jdb.DB
	Name        string
	Description string
	Define      string
	Models      []*Model
}

func NewSchema(db *jdb.DB, name string) *Schema {
	result := &Schema{
		db:     db,
		Name:   strs.Lowcase(name),
		Models: []*Model{},
	}

	SetListener(db)
	result.Init()
	schemas = append(schemas, result)

	return result
}

func GetSchema(name string) *Schema {
	for _, item := range schemas {
		if strs.Uppcase(item.Name) == strs.Uppcase(name) {
			return item
		}
	}

	return nil
}

/**
*
**/
func (c *Schema) Describe() et.Json {
	var models []et.Json = []et.Json{}
	for _, model := range c.Models {
		models = append(models, model.Describe())
	}

	return et.Json{
		"name":        c.Name,
		"description": c.Description,
		"models":      models,
	}
}

func (c *Schema) Init() error {
	c.Define = strs.Format(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE SCHEMA IF NOT EXISTS "%s";`, c.Name)
	id := strs.Format(`create-schema-%s`, c.Name)
	err := c.db.Exec(id, c.Define)
	if err != nil {
		return err
	}

	return nil
}

/**
*
**/
func (c *Schema) Model(name string) *Model {
	for _, item := range c.Models {
		if strs.Uppcase(item.Name) == strs.Uppcase(name) {
			return item
		}
	}

	return nil
}
