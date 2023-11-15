package linq

import (
	"github.com/cgalvisleon/elvis/event"
	e "github.com/cgalvisleon/elvis/json"
)

func beforeInsert(model *Model, old, new *e.Json, data e.Json) error {
	return nil
}

func afterInsert(model *Model, old, new *e.Json, data e.Json) error {
	event.Action("model/insert", e.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}

func beforeUpdate(model *Model, old, new *e.Json, data e.Json) error {
	return nil
}

func afterUpdate(model *Model, old, new *e.Json, data e.Json) error {
	event.Action("model/update", e.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}

func beforeDelete(model *Model, old, new *e.Json, data e.Json) error {
	return nil
}

func afterDelete(model *Model, old, new *e.Json, data e.Json) error {
	event.Action("model/delete", e.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}
