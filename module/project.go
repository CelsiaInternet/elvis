package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Projects *linq.Model

func DefineProjects(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return console.Panic(err)
	}

	if Projects != nil {
		return nil
	}

	Projects = linq.NewModel(SchemaModule, "PROJECTS", "Tabla de projectos", 1)
	Projects.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Projects.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Projects.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Projects.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Projects.DefineColum("name", "", "VARCHAR(250)", "")
	Projects.DefineColum("description", "", "VARCHAR(250)", "")
	Projects.DefineColum("_data", "", "JSONB", "{}")
	Projects.DefineColum("index", "", "INTEGER", 0)
	Projects.DefinePrimaryKey([]string{"_id"})
	Projects.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"name",
		"index",
	})

	if err := Projects.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* GetProjectById
* @param id string
* @return et.Item, error
**/
func GetProjectById(id string) (et.Item, error) {
	return Projects.Data().
		Where(Projects.Column("_id").Eq(id)).
		First()
}

/**
* GetProjectName
* @param name string
* @return et.Item, error
**/
func GetProjectName(name string) (et.Item, error) {
	return Projects.Data().
		Where(Projects.Column("name").Eq(name)).
		First()
}

/**
* InitProject
* @param id string
* @param name string
* @param description string
* @param data et.Json
* @return et.Item, error
**/
func InitProject(id, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetProjectName(name)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		id = utility.GenId(id)
		data.Set("_id", id)
		data.Set("name", name)
		data.Set("description", description)
		item, err := Projects.Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		return item, nil
	}

	return current, nil
}

/**
* UpSetProject
* @param id string
* @param moduleId string
* @param name string
* @param description string
* @param data et.Json
* @return et.Item, error
**/
func UpSetProject(id, moduleId, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidId(moduleId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetProjectName(name)
	if err != nil {
		return et.Item{}, err
	}

	id = utility.GenId(id)
	if !current.Ok {
		data.Set("_id", id)
		data.Set("name", name)
		data.Set("description", description)
		item, err := Projects.Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		if moduleId != "-1" {
			CheckProjectModule(id, moduleId, true)
			CheckRole(id, moduleId, "PROFILE.ADMIN", "USER.ADMIN", true)
		}

		return item, nil
	}

	if current.Id() != id {
		return et.Item{}, console.Alert(msg.RECORD_FOUND)
	}

	if current.State() == utility.OF_SYSTEM {
		return et.Item{}, console.Alert(msg.RECORD_IS_SYSTEM)
	} else if current.State() == utility.FOR_DELETE {
		return et.Item{}, console.Alert(msg.RECORD_DELETE)
	} else if current.State() != utility.ACTIVE {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_ACTIVE, current.State())
	}

	id = utility.GenId(id)
	data.Set("_id", id)
	data.Set("name", name)
	data.Set("description", description)
	data.Set("module_id", moduleId)
	return Projects.Upsert(data).
		Where(Projects.Column("_id").Eq(id)).
		And(Projects.Column("_state").Eq(utility.ACTIVE)).
		CommandOne()
}

/**
* StateProject
* @param id string
* @param state string
* @return et.Item, error
**/
func StateProject(id, state string) (et.Item, error) {
	if !utility.ValidId(id) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "id")
	}

	if !utility.ValidStr(state, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := GetProjectById(id)
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

	return Projects.Update(et.Json{
		"_state": state,
	}).
		Where(Projects.Column("_id").Eq(id)).
		And(Projects.Column("_state").Neg(state)).
		CommandOne()
}

/**
* DeleteProject
* @param id string
* @return et.Item, error
**/
func DeleteProject(id string) (et.Item, error) {
	return StateProject(id, utility.FOR_DELETE)
}

/**
* AllProjects
* @param state string
* @param search string
* @param page int
* @param rows int
* @param _select string
* @return et.List, error
**/
func AllProjects(state, search string, page, rows int, _select string) (et.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Projects.Data(_select).
			Where(Projects.Concat("NAME:", Projects.Column("name"), ":DESCRIPTION:", Projects.Column("description"), ":DATA:", Projects.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Projects.Data(_select).
			Where(Projects.Column("_state").Neg(state)).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Projects.Data(_select).
			Where(Projects.Column("_state").In("-1", state)).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	} else {
		return Projects.Data(_select).
			Where(Projects.Column("_state").Eq(state)).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	}
}
