package linq

import (
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
)

var schemas []*Schema = []*Schema{}

type Schema struct {
	Db              int
	Database        *Db
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
		Name:            Lowcase(name),
		Database:        DB(db),
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
		if Uppcase(item.Name) == Uppcase(name) {
			return item
		}
	}

	return nil
}

/**
*
**/
func (c *Schema) Describe() Json {
	var models []Json = []Json{}
	for _, model := range c.Models {
		models = append(models, model.Describe())
	}

	return Json{
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
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_namespace
		WHERE nspname = $1);`

	item, err := DBQueryOne(c.Db, sql, c.Name)
	if err != nil {
		return err
	}

	exists := item.Bool("exists")

	if !exists {
		sql := c.DDL()

		_, err := DBQDDL(c.Db, sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Schema) DDL() string {
	var result string

	if len(c.Name) > 0 {
		result = Format(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE SCHEMA IF NOT EXISTS "%s";`, c.Name)
	}

	return result
}

/**
*
**/
func (c *Schema) Model(name string) *Model {
	for _, item := range c.Models {
		if Uppcase(item.Name) == Uppcase(name) {
			return item
		}
	}

	return nil
}

func (c *Schema) SetUseSync(val bool) *Schema {
	c.UseSync = val

	return c
}
