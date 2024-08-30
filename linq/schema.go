package linq

import (
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/strs"
)

var (
	SourceField     string    = "_DATA"
	DateMakeField   string    = "DATE_MAKE"
	DateUpdateField string    = "DATE_UPDATE"
	SerieField      string    = "INDEX"
	CodeField       string    = "CODE"
	ProjectField    string    = "PROJECT_ID"
	StateField      string    = "_STATE"
	IdTFiled        string    = "_IDT"
	schemas         []*Schema = []*Schema{}
	models          []*Model  = []*Model{}
)

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

	setListener(db)
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
	_, err := c.db.Command(c.Define)
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
