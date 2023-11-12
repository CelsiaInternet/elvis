package linq

import (
	"github.com/cgalvisleon/elvis/jdb"
	js "github.com/cgalvisleon/elvis/json"
)

func Describe(db int, schema, model, filter string) js.Json {
	if len(model) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return js.Json{}
		}

		_model := _schema.Model(model)
		if _model == nil {
			return js.Json{}
		}

		result := _model.Describe()

		if len(filter) > 0 {
			return js.Json{
				filter: result.Get(filter),
			}
		}

		return result
	}

	if len(schema) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return js.Json{}
		}

		return _schema.Describe()
	}

	var describes []js.Json = []js.Json{}
	for _, schema := range schemas {
		describes = append(describes, schema.Describe())
	}

	result := jdb.DB(db).Describe()
	result.Set("schemas", describes)

	return result
}
