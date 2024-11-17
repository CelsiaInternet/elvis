package module

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

const PackageName = "module"

var Modules *linq.Model

func DefineModules(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return logs.Panice(err)
	}

	if Modules != nil {
		return nil
	}

	Modules = linq.NewModel(SchemaModule, "MODULES", "Tabla de modulos", 1)
	Modules.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Modules.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Modules.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Modules.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Modules.DefineColum("name", "", "VARCHAR(250)", "")
	Modules.DefineColum("description", "", "VARCHAR(250)", "")
	Modules.DefineColum("_data", "", "JSONB", "{}")
	Modules.DefineColum("index", "", "INTEGER", 0)
	Modules.DefinePrimaryKey([]string{"_id"})
	Modules.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"name",
		"index",
	})

	if err := Modules.Init(); err != nil {
		return logs.Panice(err)
	}

	return nil
}

/**
* GetModuleByName
* @param name string
* @return et.Item, error
**/
func GetModuleByName(name string) (et.Item, error) {
	return Modules.Data().
		Where(Modules.Column("name").Eq(name)).
		First()
}

/**
* GetModuleById
* @param id string
* @return et.Item, error
**/
func GetModuleById(id string) (et.Item, error) {
	return Modules.Data().
		Where(Modules.Column("_id").Eq(id)).
		First()
}

/**
* InitModule
* @param id string
* @param name string
* @param data et.Json
* @return et.Item, error
**/
func InitModule(id, name string, data et.Json) (et.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetModuleByName(name)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		id = utility.GenId(id)
		data.Set("_id", id)
		data.Set("name", name)
		return Modules.Insert(data).
			CommandOne()
	}

	return current, nil
}

/**
* UpSetModule
* @param id string
* @param name string
* @param description string
* @param data et.Json
* @return et.Item, error
**/
func UpSetModule(id, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetModuleByName(name)
	if err != nil {
		return et.Item{}, err
	}

	id = utility.GenId(id)
	if !current.Ok {
		data.Set("_id", id)
		data.Set("name", name)
		data.Set("description", description)
		item, err := Modules.Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		if item.Ok {
			InitProfile(id, "PROFILE.ADMIN", et.Json{})
			InitProfile(id, "PROFILE.DEV", et.Json{})
			InitProfile(id, "PROFILE.SUPORT", et.Json{})
			CheckProjectModule("-1", id, true)
			CheckRole("-1", id, "PROFILE.ADMIN", "USER.ADMIN", true)
		}

		return item, nil
	}

	if current.Id() != id {
		return et.Item{}, logs.Alertm(msg.RECORD_FOUND)
	}

	if current.State() == utility.OF_SYSTEM {
		return et.Item{}, logs.Alertm(msg.RECORD_IS_SYSTEM)
	} else if current.State() == utility.FOR_DELETE {
		return et.Item{}, logs.Alertm(msg.RECORD_DELETE)
	} else if current.State() != utility.ACTIVE {
		return et.Item{}, logs.Alertf(msg.RECORD_NOT_ACTIVE, current.State())
	}

	data.Set("_id", id)
	data.Set("name", name)
	data.Set("description", description)
	return Modules.Update(data).
		Where(Modules.Column("_id").Eq(id)).
		And(Modules.Column("_state").Eq(utility.ACTIVE)).
		CommandOne()
}

/**
* StateModule
* @param id string
* @param state string
* @return et.Item, error
**/
func StateModule(id, state string) (et.Item, error) {
	if !utility.ValidId(id) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "id")
	}

	if !utility.ValidStr(state, 0, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := GetModuleById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, logs.Alertm(msg.RECORD_NOT_FOUND)
	}

	if current.State() == utility.OF_SYSTEM {
		return et.Item{}, logs.Alertm(msg.RECORD_IS_SYSTEM)
	} else if current.State() == utility.FOR_DELETE {
		return et.Item{}, logs.Alertm(msg.RECORD_DELETE)
	} else if current.State() == state {
		return et.Item{}, logs.Alertm(msg.RECORD_NOT_CHANGE)
	}

	return Modules.Update(et.Json{
		"_state": state,
	}).
		Where(Modules.Column("_id").Eq(id)).
		And(Modules.Column("_state").Neg(state)).
		CommandOne()
}

/**
* DeleteModule
* @param id string
* @return et.Item, error
**/
func DeleteModule(id string) (et.Item, error) {
	return StateModule(id, utility.FOR_DELETE)
}

/**
* AllModules
* @param state string
* @param search string
* @param page int
* @param rows int
* @param _select string
* @return et.List, error
**/
func AllModules(state, search string, page, rows int, _select string) (et.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Modules.Data(_select).
			Where(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Modules.Data(_select).
			Where(Modules.Column("_state").Neg(state)).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Modules.Data(_select).
			Where(Modules.Column("_state").In("-1", state)).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	} else {
		return Modules.Data(_select).
			Where(Modules.Column("_state").Eq(state)).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	}
}
