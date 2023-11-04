package linq

import (
	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
)

func beforeInsert(model *Model, old, new *Json, data Json) error {
	return nil
}

func afterInsert(model *Model, old, new *Json, data Json) error {
	event.EventPublish("model/insert", Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}

func beforeUpdate(model *Model, old, new *Json, data Json) error {
	return nil
}

func afterUpdate(model *Model, old, new *Json, data Json) error {
	event.EventPublish("model/update", Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}

func beforeDelete(model *Model, old, new *Json, data Json) error {
	return nil
}

func afterDelete(model *Model, old, new *Json, data Json) error {
	event.EventPublish("model/delete", Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

	return nil
}
