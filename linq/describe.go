package linq

import (
	"github.com/cgalvisleon/elvis/jdb"
	j "github.com/cgalvisleon/elvis/json"
)

func Describe(db int, schema, model, filter string) j.Json {
	if len(model) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return j.Json{}
		}

		_model := _schema.Model(model)
		if _model == nil {
			return j.Json{}
		}

		result := _model.Describe()

		if len(filter) > 0 {
			return j.Json{
				filter: result.Get(filter),
			}
		}

		return result
	}

	if len(schema) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return j.Json{}
		}

		return _schema.Describe()
	}

	var describes []j.Json = []j.Json{}
	for _, schema := range schemas {
		describes = append(describes, schema.Describe())
	}

	result := jdb.DB(db).Describe()
	result.Set("schemas", describes)

	return result
}
