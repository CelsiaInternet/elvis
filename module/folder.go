package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Folders *linq.Model

func DefineFolders(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return console.Panic(err)
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
	Folders.DefineForeignKey("module_id", Modules.Col("_id"))
	Folders.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"main_id",
		"name",
		"index",
	})

	if err := Folders.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* GetFolderById
* @param id string
* @return et.Item, error
**/
func GetFolderById(id string) (et.Item, error) {
	return Folders.Data().
		Where(Folders.Column("_id").Eq(id)).
		First()
}

/**
* GetFolderByName
* @param moduleId string
* @param mainId string
* @param name string
* @return et.Item, error
**/
func GetFolderByName(moduleId, mainId, name string) (et.Item, error) {
	return Folders.Data().
		Where(Folders.Column("module_id").Eq(moduleId)).
		And(Folders.Column("main_id").Eq(mainId)).
		And(Folders.Column("name").Eq(name)).
		First()
}

/**
* InitFolder
* @param moduleId string
* @param mainId string
* @param id string
* @param name string
* @param description string
* @param data et.Json
* @return et.Item, error
**/
func InitFolder(moduleId, mainId, id, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidId(moduleId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidId(mainId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "main_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	id = utility.GenKey(id)
	current, err := GetFolderById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		data["module_id"] = moduleId
		data["main_id"] = mainId
		data["_id"] = id
		data["name"] = name
		data["description"] = description
		item, err := Folders.Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		if item.Ok {
			CheckProfileFolder(moduleId, "PROFILE.ADMIN", id, true)
			CheckProfileFolder(moduleId, "PROFILE.DEV", id, true)
			CheckProfileFolder(moduleId, "PROFILE.SUPORT", id, true)
			CheckModuleFolder(moduleId, id, true)
		}

		return item, nil
	}

	return current, nil
}

/**
* UpSetFolder
* @param moduleId string
* @param mainId string
* @param name string
* @param description string
* @param data et.Json
* @return et.Item, error
**/
func UpSetFolder(moduleId, mainId, id, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidId(moduleId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidId(mainId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "main_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	id = utility.GenKey(id)
	current, err := GetFolderById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		data["module_id"] = moduleId
		data["main_id"] = mainId
		data["_id"] = id
		data["name"] = name
		data["description"] = description
		item, err := Folders.Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		if item.Ok {
			CheckProfileFolder(moduleId, "PROFILE.ADMIN", id, true)
			CheckProfileFolder(moduleId, "PROFILE.DEV", id, true)
			CheckProfileFolder(moduleId, "PROFILE.SUPORT", id, true)
			CheckModuleFolder(moduleId, id, true)
		}

		return item, nil
	}

	if current.State() == utility.OF_SYSTEM {
		return et.Item{}, console.Alert(msg.RECORD_IS_SYSTEM)
	} else if current.State() == utility.FOR_DELETE {
		return et.Item{}, console.Alert(msg.RECORD_DELETE)
	} else if current.State() != utility.ACTIVE {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_ACTIVE, current.State())
	}

	delete(data, "module_id")
	data["main_id"] = mainId
	data["_id"] = id
	data["name"] = name
	data["description"] = description
	return Folders.Update(data).
		Where(Folders.Column("_id").Eq(id)).
		CommandOne()
}

/**
* StateFolder
* @param id string
* @param state string
* @return et.Item, error
**/
func StateFolder(id, state string) (et.Item, error) {
	if !utility.ValidId(id) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "id")
	}

	if !utility.ValidStr(state, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := GetFolderById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, console.Alert(msg.RECORD_NOT_FOUND)
	}

	if current.State() == utility.OF_SYSTEM {
		return et.Item{}, console.Alert(msg.RECORD_IS_SYSTEM)
	} else if current.State() == utility.FOR_DELETE {
		return et.Item{}, console.Alert(msg.RECORD_DELETE)
	} else if current.State() == state {
		return et.Item{}, console.Alert(msg.RECORD_NOT_CHANGE)
	}

	result, err := Folders.Update(et.Json{
		"_state": state,
	}).
		Where(Folders.Column("_id").Eq(id)).
		And(Folders.Column("_state").Neg(state)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	if result.Ok {
		event.Work("folder/state", result.Result)
	}

	return result, nil
}

/**
* DeleteFolder
* @param id string
* @return et.Item, error
**/
func DeleteFolder(id string) (et.Item, error) {
	return StateFolder(id, utility.FOR_DELETE)
}

/**
* AllFolders
* @param state string
* @param search string
* @param page int
* @param rows int
* @return et.List, error
**/
func AllFolders(state, search string, page, rows int) (et.List, error) {
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

/**
* DefaultFolderUsers
* @param moduleId string
* @return error
**/
func DefaultFolderUsers(moduleId string) error {
	_, err := InitFolder(moduleId, "-1", "FOLDER.USERS", "Usuarios", "", et.Json{
		"icon":   "users",
		"view":   "users",
		"clase":  "user",
		"help":   "help/module/users",
		"title":  "Usuario",
		"url":    "user/all?state=0&search={search}&page={page}&rows={rows}",
		"state":  "",
		"filter": []et.Json{},
		"states": []et.Json{},
		"detail": et.Json{
			"title":    []string{"$1", "full_name"},
			"subtitle": []string{"$1", "phone"},
			"datetime": []string{"$1", "date_update"},
			"code":     []string{"$1", "full_name"},
			"new_code": "Nuevo",
			"state_color": et.Json{
				"field_name": "_state",
				"warning":    "",
				"alert":      "",
				"info":       "",
			},
			"email":  []string{"$1", "email"},
			"avatar": []string{"$1", "avatar"},
		},
		"showNew":   true,
		"showPrint": false,
		"order":     10,
	})
	if err != nil {
		return err
	}

	return nil
}
