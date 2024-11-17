package module

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

var Types *linq.Model

func DefineTypes(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return logs.Panice(err)
	}

	if Types != nil {
		return nil
	}

	Types = linq.NewModel(SchemaModule, "TYPES", "Tabla de tipo", 1)
	Types.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Types.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Types.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Types.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Types.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	Types.DefineColum("kind", "", "VARCHAR(80)", "")
	Types.DefineColum("name", "", "VARCHAR(250)", "")
	Types.DefineColum("description", "", "TEXT", "")
	Types.DefineColum("_data", "", "JSONB", "{}")
	Types.DefineColum("index", "", "INTEGER", 0)
	Types.DefinePrimaryKey([]string{"_id"})
	Types.DefineForeignKey("project_id", Projects.Col("_id"))
	Types.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"project_id",
		"kind",
		"name",
		"index",
	})

	if err := Types.Init(); err != nil {
		return logs.Panice(err)
	}

	return nil
}

/**
* Types
*	Handler for CRUD data
**/
func GetTypeByName(kind, name string) (et.Item, error) {
	return Types.Data().
		Where(Types.Column("kind").Eq(kind)).
		And(Types.Column("name").Eq(name)).
		First()
}

/**
* GetTypeById
* @param string id
* @return et.Item, error
**/
func GetTypeById(id string) (et.Item, error) {
	return Types.Data().
		Where(Types.Column("_id").Eq(id)).
		First()
}

/**
* GetTypeByIndex
* @param int idx
* @return et.Item, error
**/
func GetTypeByIndex(idx int) (et.Item, error) {
	return Types.Data().
		Where(Types.Column("index").Eq(idx)).
		First()
}

/**
* InitType
* @param string projectId
* @param string id
* @param string state
* @param string kind
* @param string name
* @return et.Item, error
**/
func InitType(projectId, id, state, kind, name string) (et.Item, error) {
	if !utility.ValidId(kind) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "kind")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetTypeByName(kind, name)
	if err != nil {
		return et.Item{}, err
	}

	id = utility.GenId(id)
	if !current.Ok {
		data := et.Json{}
		data["project_id"] = projectId
		data["_state"] = state
		data["_id"] = id
		data["kind"] = kind
		data["name"] = name
		result, err := Types.Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		return result, nil
	}

	return current, nil
}

/**
* UpSetType
* @param string projectId
* @param string id
* @param string kind
* @param string name
* @param string description
* @return et.Item, error
**/
func UpSetType(projectId, id, kind, name, description string) (et.Item, error) {
	if !utility.ValidId(id) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "_id")
	}

	if !utility.ValidId(kind) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "kind")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetTypeByName(kind, name)
	if err != nil {
		return et.Item{}, err
	}

	id = utility.GenKey(id)
	if !current.Ok {
		data := et.Json{}
		data["project_id"] = projectId
		data["_state"] = utility.ACTIVE
		data["_id"] = id
		data["kind"] = kind
		data["name"] = name
		data["description"] = description
		return Types.Insert(data).
			CommandOne()
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

	data := et.Json{}
	data["project_id"] = projectId
	data["_id"] = id
	data["kind"] = kind
	data["name"] = name
	data["description"] = description
	return Types.Update(data).
		Where(Types.Column("_id").Eq(id)).
		And(Types.Column("_state").Eq(utility.ACTIVE)).
		CommandOne()
}

/**
* StateType
* @param string id
* @param string state
* @return et.Item, error
**/
func StateType(id, state string) (et.Item, error) {
	if !utility.ValidId(id) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "id")
	}

	if !utility.ValidStr(state, 0, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := GetTypeById(id)
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

	return Types.Update(et.Json{
		"_state": state,
	}).
		Where(Types.Column("_id").Eq(id)).
		And(Types.Column("_state").Neg(state)).
		CommandOne()
}

func DeleteType(id string) (et.Item, error) {
	return StateType(id, utility.FOR_DELETE)
}

func AllTypes(projectId, kind, state, search string, page, rows int, _select string) (et.List, error) {
	if !utility.ValidId(kind) {
		return et.List{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "kind")
	}

	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Types.Data(_select).
			Where(Types.Column("kind").Eq(kind)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Types.Concat("NAME:", Types.Column("name"), ":DESCRIPTION", Types.Column("description"), ":DATA:", Types.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Types.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Types.Data(_select).
			Where(Types.Column("kind").Eq(kind)).
			And(Types.Column("_state").Neg(state)).
			And(Types.Column("project_id").In("-1", projectId)).
			OrderBy(Types.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Types.Data(_select).
			Where(Types.Column("kind").Eq(kind)).
			And(Types.Column("_state").In("-1", state)).
			And(Types.Column("project_id").In("-1", projectId)).
			OrderBy(Types.Column("name"), true).
			List(page, rows)
	} else {
		return Types.Data(_select).
			Where(Types.Column("kind").Eq(kind)).
			And(Types.Column("_state").Eq(state)).
			And(Types.Column("project_id").In("-1", projectId)).
			OrderBy(Types.Column("name"), true).
			List(page, rows)
	}
}
