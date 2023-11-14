package linq

import (
	"github.com/cgalvisleon/elvis/event"
	js "github.com/cgalvisleon/elvis/json"
)

func beforeInsert(model *Model, old, new *js.Json, data js.Json) error {
	return nil
}

func afterInsert(model *Model, old, new *js.Json, data js.Json) error {
	event.Action("model/insert", js.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}

func beforeUpdate(model *Model, old, new *js.Json, data js.Json) error {
	return nil
}

func afterUpdate(model *Model, old, new *js.Json, data js.Json) error {
	event.Action("model/update", js.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}

func beforeDelete(model *Model, old, new *js.Json, data js.Json) error {
	return nil
}

func afterDelete(model *Model, old, new *js.Json, data js.Json) error {
	event.Action("model/delete", js.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}
