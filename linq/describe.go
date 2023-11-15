package linq

import (
	"github.com/cgalvisleon/elvis/jdb"
	e "github.com/cgalvisleon/elvis/json"
)

func Describe(db int, schema, model, filter string) e.Json {
	if len(model) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return e.Json{}
		}

		_model := _schema.Model(model)
		if _model == nil {
			return e.Json{}
		}

		result := _model.Describe()

		if len(filter) > 0 {
			return e.Json{
				filter: result.Get(filter),
			}
		}

		return result
	}

	if len(schema) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return e.Json{}
		}

		return _schema.Describe()
	}

	var describes []e.Json = []e.Json{}
	for _, schema := range schemas {
		describes = append(describes, schema.Describe())
	}

	result := jdb.DB(db).Describe()
	result.Set("schemas", describes)

	return result
}
