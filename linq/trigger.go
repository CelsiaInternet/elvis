package linq

import (
	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
)

func beforeInsert(model *Model, old, new *Json, data Json) {

}

func afterInsert(model *Model, old, new *Json, data Json) {
	event.EventPublish("model/insert", Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

}

func beforeUpdate(model *Model, old, new *Json, data Json) {
}

func afterUpdate(model *Model, old, new *Json, data Json) {
	event.EventPublish("model/update", Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

}

func beforeDelete(model *Model, old, new *Json, data Json) {

}
func afterDelete(model *Model, old, new *Json, data Json) {
	event.EventPublish("model/delete", Json{
		"table": model.Name,
		"old":   old,
		"new":   new,
	})

}
