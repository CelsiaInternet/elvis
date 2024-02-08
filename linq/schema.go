package linq

import (
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/strs"
)

var schemas []*Schema = []*Schema{}

type Schema struct {
	Db              int
	Database        *jdb.Db
	Name            string
	Description     string
	UseSync         bool
	SourceField     string
	DateMakeField   string
	DateUpdateField string
	IndexField      string
	CodeField       string
	ProjectField    string
	StateField      string
	Models          []*Model
}

func NewSchema(db int, name string) *Schema {
	result := &Schema{
		Db:              db,
		Name:            strs.Lowcase(name),
		Database:        jdb.DB(db),
		UseSync:         true,
		SourceField:     "_DATA",
		DateMakeField:   "DATE_MAKE",
		DateUpdateField: "DATE_UPDATE",
		IndexField:      "INDEX",
		CodeField:       "CODE",
		ProjectField:    "PROJECT_ID",
		StateField:      "_STATE",
		Models:          []*Model{},
	}

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
		"name":            c.Name,
		"description":     c.Description,
		"database":        c.Database.Dbname,
		"source_field":    c.SourceField,
		"dateMakeField":   c.DateMakeField,
		"dateUpdateField": c.DateUpdateField,
		"indexField":      c.IndexField,
		"codeField":       c.CodeField,
		"projectField":    c.ProjectField,
		"models":          models,
	}
}

func (c *Schema) Init() error {
	_, err := jdb.CreateSchema(c.Db, c.Name)
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

func (c *Schema) SetUseSync(val bool) *Schema {
	c.UseSync = val

	return c
}
