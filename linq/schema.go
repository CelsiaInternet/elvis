package linq

import (
	"github.com/cgalvisleon/elvis/jdb"
	js "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
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
		Name:            utility.Lowcase(name),
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
		if utility.Uppcase(item.Name) == utility.Uppcase(name) {
			return item
		}
	}

	return nil
}

/**
*
**/
func (c *Schema) Describe() js.Json {
	var models []js.Json = []js.Json{}
	for _, model := range c.Models {
		models = append(models, model.Describe())
	}

	return js.Json{
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

	item, err := jdb.DBQueryOne(c.Db, sql, c.Name)
	if err != nil {
		return err
	}

	exists := item.Bool("exists")

	if !exists {
		sql := c.DDL()

		_, err := jdb.DBQDDL(c.Db, sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Schema) DDL() string {
	var result string

	if len(c.Name) > 0 {
		result = utility.Format(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE SCHEMA IF NOT EXISTS "%s";`, c.Name)
	}

	return result
}

/**
*
**/
func (c *Schema) Model(name string) *Model {
	for _, item := range c.Models {
		if utility.Uppcase(item.Name) == utility.Uppcase(name) {
			return item
		}
	}

	return nil
}

func (c *Schema) SetUseSync(val bool) *Schema {
	c.UseSync = val

	return c
}
