package linq

import (
	"github.com/cgalvisleon/elvis/jdb"
	ej "github.com/cgalvisleon/elvis/json"
)

func Describe(db int, schema, model, filter string) ej.Json {
	if len(model) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return ej.Json{}
		}

		_model := _schema.Model(model)
		if _model == nil {
			return ej.Json{}
		}

		result := _model.Describe()

		if len(filter) > 0 {
			return ej.Json{
				filter: result.Get(filter),
			}
		}

		return result
	}

	if len(schema) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return ej.Json{}
		}

		return _schema.Describe()
	}

	var describes []ej.Json = []ej.Json{}
	for _, schema := range schemas {
		describes = append(describes, schema.Describe())
	}

	result := jdb.DB(db).Describe()
	result.Set("schemas", describes)

	return result
}
