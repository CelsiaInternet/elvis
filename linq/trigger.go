package linq

import (
	"github.com/cgalvisleon/elvis/event"
	j "github.com/cgalvisleon/elvis/json"
)

func beforeInsert(model *Model, old, new *j.Json, data j.Json) error {
	return nil
}

func afterInsert(model *Model, old, new *j.Json, data j.Json) error {
	event.Action("model/insert", j.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}

func beforeUpdate(model *Model, old, new *j.Json, data j.Json) error {
	return nil
}

func afterUpdate(model *Model, old, new *j.Json, data j.Json) error {
	event.Action("model/update", j.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}

func beforeDelete(model *Model, old, new *j.Json, data j.Json) error {
	return nil
}

func afterDelete(model *Model, old, new *j.Json, data j.Json) error {
	event.Action("model/delete", j.Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}
