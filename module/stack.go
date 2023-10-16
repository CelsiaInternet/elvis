package module

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	vnt "github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/linq"
	. "github.com/cgalvisleon/elvis/utilities"
)

var Stacks *Model

func DefineStacks() error {
	if err := DefineCoreSchema(); err != nil {
		return console.PanicE(err)
	}

	if Stacks != nil {
		return nil
	}

	Stacks = NewModel(SchemaCore, "STACK", "Tabla de colas", 1)
	Stacks.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Stacks.DefineColum("_state", "", "VARCHAR(80)", ACTIVE)
	Stacks.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Stacks.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	Stacks.DefineColum("app", "", "VARCHAR(80)", "")
	Stacks.DefineColum("event", "", "TEXT", "")
	Stacks.DefineColum("worker", "", "VARCHAR(80)", "-1")
	Stacks.DefineColum("_data", "", "JSONB", "{}")
	Stacks.DefineColum("index", "", "INTEGER", 0)
	Stacks.DefinePrimaryKey([]string{"_id"})
	Stacks.DefineIndex([]string{
		"date_make",
		"_state",
		"app",
		"event",
		"worker",
		"index",
	})
	Stacks.DefineForeignKey("project_id", Projects.Column("_id"))
	Stacks.Trigger(AfterInsert, func(model *Model, old, new *Json, data Json) {
		app := new.Key("app")
		event := new.Key("event")
		_data := new.Json("_data")
		vnt.Publish("event/stack", event, Json{
			"app":   app,
			"event": event,
			"data":  _data,
		})
	})
	Stacks.Trigger(AfterUpdate, func(model *Model, old, new *Json, data Json) {
		worker := new.Key("worker")
		if worker != "-1" {
			app := new.Key("app")
			event := new.Key("event")
			vnt.Publish("event/stack", Format(`%s/working`, event), Json{
				"app":    app,
				"event":  event,
				"worker": worker,
			})
		}
	})
	Stacks.Trigger(AfterDelete, func(model *Model, old, new *Json, data Json) {
		app := old.Key("app")
		event := old.Key("event")
		worker := old.Key("worker")
		vnt.Publish("event/stack", Format(`%s/finished`, event), Json{
			"app":    app,
			"event":  event,
			"worker": worker,
		})
	})

	if err := InitModel(Stacks); err != nil {
		return console.PanicE(err)
	}

	return nil
}

func GetStackById(id string) (Item, error) {
	return Stacks.Select().
		Where(Stacks.Col("_id").Eq(id)).
		First()
}

func SetStack(projectId, app, event string, data Json) (Item, error) {
	data["project_id"] = projectId
	data["app"] = app
	data["event"] = event
	item, err := Stacks.Insert(data).
		Command()
	if err != nil {
		return Item{}, console.Error(err)
	}

	return item, nil
}

func DeleteStack(id string) (Item, error) {
	return Stacks.Delete().
		Where(Stacks.Col("_id").Eq(id)).
		Command()
}

func WorkerStack(workerId, event string) (Item, error) {
	data := Json{}
	data.Set("_state", IN_PROCESS)
	data.Set("worker", workerId)
	return Stacks.Update(data).
		Where(Stacks.Col("_id").Eq(From(Stacks, "A").
			Where(Stacks.Col("_state").Eq(ACTIVE)).
			And(Stacks.Col("event").Eq(event)).
			OrderBy(Stacks.Col("index"), true).
			Select(Stacks.Col("_id")).SQL())).
		Command()
}

func AllStacks(projectId, search string, page, rows int) (List, error) {
	return Stacks.Select().
		Where(Stacks.Col("project_id").Eq(projectId)).
		And(Stacks.Concat("APP:", Stacks.Col("APP"), ":EVENT:", Stacks.Col("event"), ":DATA", Stacks.Col("_data")).Like("%"+search+"%")).
		OrderBy(Stacks.Column("index"), true).
		List(page, rows)
}
