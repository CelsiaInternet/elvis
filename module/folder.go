package module

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/linq"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
)

var Folders *Model

func DefineFolders() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
	}

	if Folders != nil {
		return nil
	}

	Folders = NewModel(SchemaModule, "FOLDERS", "Tabla de carpetas", 1)
	Folders.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Folders.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Folders.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	Folders.DefineColum("_state", "", "VARCHAR(80)", ACTIVE)
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
	Folders.Trigger(AfterInsert, func(model *Model, old, new *Json, data Json) error {
		id := new.Id()
		moduleId := new.Key("module_id")
		CheckProfileFolder(moduleId, "PROFILE.ADMIN", id, true)
		CheckProfileFolder(moduleId, "PROFILE.DEV", id, true)
		CheckProfileFolder(moduleId, "PROFILE.SUPORT", id, true)

		return nil
	})
	Folders.Trigger(AfterUpdate, func(model *Model, old, new *Json, data Json) error {
		event.EventPublish("folder/update", *new)
		oldState := old.Key("_state")
		newState := old.Key("_state")
		if oldState != newState {
			event.EventPublish("folder/state", *new)
		}

		return nil
	})
	Folders.Trigger(AfterDelete, func(model *Model, old, new *Json, data Json) error {
		event.EventPublish("folder/delete", *old)

		return nil
	})

	return InitModel(Folders)
}

/**
*	Folder
*	Handler for CRUD data
 */
func GetFolderByName(moduleId, mainId, name string) (Item, error) {
	return Folders.Select().
		Where(Folders.Column("module_id").Eq(moduleId)).
		And(Folders.Column("main_id").Eq(mainId)).
		And(Folders.Column("name").Eq(name)).
		First()
}

func InitFolder(moduleId, mainId, id, name, description string, data Json) (Item, error) {
	if !ValidId(moduleId) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "module_id")
	}

	if !ValidId(mainId) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "main_id")
	}

	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "name")
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return Item{}, err
	}

	if !module.Ok {
		return Item{}, console.ErrorM(MODULE_NOT_FOUND)
	}

	current, err := GetFolderByName(moduleId, mainId, name)
	if err != nil {
		return Item{}, err
	}

	if current.Ok && current.Id() != id {
		return Item{
			Ok: current.Ok,
			Result: Json{
				"message": RECORD_FOUND,
				"_id":     id,
			},
		}, nil
	}

	id = GenId(id)
	data["module_id"] = moduleId
	data["main_id"] = mainId
	data["_id"] = id
	data["name"] = name
	data["description"] = description
	item, err := Folders.Upsert(data).
		Where(Folders.Column("_id").Eq(id)).
		And(Folders.Column("_state").Eq(ACTIVE)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func UpSetFolder(moduleId, mainId, name, description string, data Json) (Item, error) {
	if !ValidId(moduleId) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "module_id")
	}

	if !ValidId(mainId) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "main_id")
	}

	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "name")
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return Item{}, err
	}

	if !module.Ok {
		return Item{}, console.ErrorM(MODULE_NOT_FOUND)
	}

	current, err := Folders.Select(Folders.Column("_id")).
		Where(Folders.Column("module_id").Eq(moduleId)).
		And(Folders.Column("main_id").Eq(mainId)).
		And(Folders.Column("name").Eq(name)).
		First()
	if err != nil {
		return Item{}, err
	}

	id := current.Id()
	id = GenId(id)
	data["module_id"] = moduleId
	data["main_id"] = mainId
	data["_id"] = id
	data["name"] = name
	data["description"] = description
	item, err := Folders.Upsert(data).
		Where(Folders.Column("_id").Eq(id)).
		And(Folders.Column("_state").Eq(ACTIVE)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func GetFolderById(id string) (Item, error) {
	return Folders.Select().
		Where(Folders.Column("_id").Eq(id)).
		First()
}

func StateFolder(id, state string) (Item, error) {
	if !ValidId(state) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "state")
	}

	item, err := Folders.Upsert(Json{
		"_state": state,
	}).
		Where(Folders.Column("_id").Eq(id)).
		And(Folders.Column("_state").Neg(state)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func DeleteFolder(id string) (Item, error) {
	item, err := Folders.Delete().
		Where(Folders.Column("_id").Eq(id)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func AllFolders(state, search string, page, rows int) (List, error) {
	if state == "" {
		state = ACTIVE
	}

	auxState := state

	if auxState == "*" {
		state = FOR_DELETE

		return Folders.Select().
			Where(Folders.Column("_state").Neg(state)).
			And(Folders.Concat("NAME:", Folders.Column("name"), ":DESCRIPTION", Folders.Column("description"), ":DATA:", Folders.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Folders.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Folders.Select().
			Where(Folders.Column("_state").In("-1", state)).
			And(Folders.Concat("NAME:", Folders.Column("name"), ":DESCRIPTION", Folders.Column("description"), ":DATA:", Folders.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Folders.Column("name"), true).
			List(page, rows)
	} else {
		return Folders.Select().
			Where(Folders.Column("_state").Eq(state)).
			And(Folders.Concat("NAME:", Folders.Column("name"), ":DESCRIPTION", Folders.Column("description"), ":DATA:", Folders.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Folders.Column("name"), true).
			List(page, rows)
	}
}
