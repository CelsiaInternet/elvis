package linq

import (
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
)

func Describe(db int, schema, model, filter string) Json {
	if len(model) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return Json{}
		}

		_model := _schema.Model(model)
		if _model == nil {
			return Json{}
		}

		result := _model.Describe()

		if len(filter) > 0 {
			return Json{
				filter: result.Get(filter),
			}
		}

		return result
	}

	if len(schema) > 0 {
		_schema := GetSchema(schema)
		if _schema == nil {
			return Json{}
		}

		return _schema.Describe()
	}

	var describes []Json = []Json{}
	for _, schema := range schemas {
		describes = append(describes, schema.Describe())
	}

	result := DB(db).Describe()
	result.Set("schemas", describes)

	return result
}
