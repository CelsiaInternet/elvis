package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/core"
	"github.com/cgalvisleon/elvis/event"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Folders *linq.Model

func DefineFolders() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
	}

	if Folders != nil {
		return nil
	}

	Folders = linq.NewModel(SchemaModule, "FOLDERS", "Tabla de carpetas", 1)
	Folders.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Folders.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Folders.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	Folders.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Folders.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Folders.DefineColum("main_id", "", "VARCHAR(80)", "-1")
	Folders.DefineColum("name", "", "VARCHAR(250)", "")
	Folders.DefineColum("description", "", "VARCHAR(250)", "")
	Folders.DefineColum("_data", "", "JSONB", "{}")
	Folders.DefineColum("index", "", "INTEGER", 0)
	Folders.DefinePrimaryKey([]string{"_id"})
	Folders.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"main_id",
		"name",
		"index",
	})
	Folders.DefineForeignKey("module_id", Modules.Column("_id"))
	Folders.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *e.Json, data e.Json) error {
		id := new.Id()
		moduleId := new.Key("module_id")
		CheckProfileFolder(moduleId, "PROFILE.ADMIN", id, true)
		CheckProfileFolder(moduleId, "PROFILE.DEV", id, true)
		CheckProfileFolder(moduleId, "PROFILE.SUPORT", id, true)

		return nil
	})
	Folders.Trigger(linq.AfterUpdate, func(model *linq.Model, old, new *e.Json, data e.Json) error {
		event.Action("folder/update", *new)
		oldState := old.Key("_state")
		newState := old.Key("_state")
		if oldState != newState {
			event.Action("folder/state", *new)
		}

		return nil
	})
	Folders.Trigger(linq.AfterDelete, func(model *linq.Model, old, new *e.Json, data e.Json) error {
		event.Action("folder/delete", *old)

		return nil
	})

	return core.InitModel(Folders)
}

/**
*	Folder
*	Handler for CRUD data
 */
func GetFolderByName(moduleId, mainId, name string) (e.Item, error) {
	return Folders.Data().
		Where(Folders.Column("module_id").Eq(moduleId)).
		And(Folders.Column("main_id").Eq(mainId)).
		And(Folders.Column("name").Eq(name)).
		First()
}

func InitFolder(moduleId, mainId, id, name, description string, data e.Json) (e.Item, error) {
	if !utility.ValidId(moduleId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidId(mainId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "main_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return e.Item{}, err
	}

	if !module.Ok {
		return e.Item{}, console.ErrorM(msg.MODULE_NOT_FOUND)
	}

	current, err := GetFolderByName(moduleId, mainId, name)
	if err != nil {
		return e.Item{}, err
	}

	if current.Ok && current.Id() != id {
		return e.Item{
			Ok: current.Ok,
			Result: e.Json{
				"message": msg.RECORD_FOUND,
				"_id":     id,
			},
		}, nil
	}

	id = utility.GenId(id)
	data["module_id"] = moduleId
	data["main_id"] = mainId
	data["_id"] = id
	data["name"] = name
	data["description"] = description
	item, err := Folders.Upsert(data).
		Where(Folders.Column("_id").Eq(id)).
		And(Folders.Column("_state").Eq(utility.ACTIVE)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func UpSetFolder(moduleId, mainId, name, description string, data e.Json) (e.Item, error) {
	if !utility.ValidId(moduleId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidId(mainId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "main_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return e.Item{}, err
	}

	if !module.Ok {
		return e.Item{}, console.ErrorM(msg.MODULE_NOT_FOUND)
	}

	current, err := Folders.Data(Folders.Column("_id")).
		Where(Folders.Column("module_id").Eq(moduleId)).
		And(Folders.Column("main_id").Eq(mainId)).
		And(Folders.Column("name").Eq(name)).
		First()
	if err != nil {
		return e.Item{}, err
	}

	id := current.Id()
	id = utility.GenId(id)
	data["module_id"] = moduleId
	data["main_id"] = mainId
	data["_id"] = id
	data["name"] = name
	data["description"] = description
	item, err := Folders.Upsert(data).
		Where(Folders.Column("_id").Eq(id)).
		And(Folders.Column("_state").Eq(utility.ACTIVE)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func GetFolderById(id string) (e.Item, error) {
	return Folders.Data().
		Where(Folders.Column("_id").Eq(id)).
		First()
}

func StateFolder(id, state string) (e.Item, error) {
	if !utility.ValidId(state) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	item, err := Folders.Update(e.Json{
		"_state": state,
	}).
		Where(Folders.Column("_id").Eq(id)).
		And(Folders.Column("_state").Neg(state)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func DeleteFolder(id string) (e.Item, error) {
	item, err := Folders.Delete().
		Where(Folders.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func AllFolders(state, search string, page, rows int) (e.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Folders.Data().
			Where(Folders.Concat("NAME:", Folders.Column("name"), ":DESCRIPTION", Folders.Column("description"), ":DATA:", Folders.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Folders.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Folders.Data().
			Where(Folders.Column("_state").Neg(state)).
			OrderBy(Folders.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Folders.Data().
			Where(Folders.Column("_state").In("-1", state)).
			OrderBy(Folders.Column("name"), true).
			List(page, rows)
	} else {
		return Folders.Data().
			Where(Folders.Column("_state").Eq(state)).
			OrderBy(Folders.Column("name"), true).
			List(page, rows)
	}
}
