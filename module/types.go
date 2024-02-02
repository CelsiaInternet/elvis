package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/core"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Types *linq.Model

func DefineTypes() error {
	if err := DefineSchemaModule(); err != nil {
		return console.Panic(err)
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
	Types.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"project_id",
		"kind",
		"name",
		"index",
	})

	if err := core.InitModel(Types); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* Types
*	Handler for CRUD data
 */
func GetTypeByName(kind, name string) (e.Item, error) {
	return Types.Data().
		Where(Types.Column("kind").Eq(kind)).
		And(Types.Column("name").Eq(name)).
		First()
}

func GetTypeById(id string) (e.Item, error) {
	return Types.Data().
		Where(Types.Column("_id").Eq(id)).
		First()
}

func GetTypeByIndex(idx int) (e.Item, error) {
	return Types.Data().
		Where(Types.Column("index").Eq(idx)).
		First()
}

func InitType(projectId, id, state, kind, name, description string) (e.Item, error) {
	if !utility.ValidId(kind) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "kind")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetTypeByName(kind, name)
	if err != nil {
		return e.Item{}, err
	}

	if current.Ok && current.Id() != id {
		return e.Item{
			Ok: current.Ok,
			Result: e.Json{
				"message": msg.RECORD_FOUND,
			},
		}, nil
	}

	id = utility.GenId(id)
	data := e.Json{}
	data["project_id"] = projectId
	data["_id"] = id
	data["kind"] = kind
	data["name"] = name
	data["description"] = description
	return Types.Upsert(data).
		Where(Types.Column("_id").Eq(id)).
		CommandOne()
}

func UpSetType(projectId, id, kind, name, description string) (e.Item, error) {
	if !utility.ValidId(id) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "_id")
	}

	if !utility.ValidId(kind) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "kind")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetTypeByName(kind, name)
	if err != nil {
		return e.Item{}, err
	}

	if current.Ok && current.Id() != id {
		return e.Item{
			Ok: current.Ok,
			Result: e.Json{
				"message": msg.RECORD_FOUND,
				"_id":     id,
				"index":   current.Index(),
			},
		}, nil
	}

	id = utility.GenId(id)
	data := e.Json{}
	data["project_id"] = projectId
	data["_id"] = id
	data["kind"] = kind
	data["name"] = name
	data["description"] = description
	return Types.Upsert(data).
		Where(Types.Column("_id").Eq(id)).
		CommandOne()
}

func StateType(id, state string) (e.Item, error) {
	if !utility.ValidId(state) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	return Types.Update(e.Json{
		"_state": state,
	}).
		Where(Types.Column("_id").Eq(id)).
		And(Types.Column("_state").Neg(state)).
		CommandOne()
}

func DeleteType(id string) (e.Item, error) {
	return StateType(id, utility.FOR_DELETE)
}

func AllTypes(projectId, kind, state, search string, page, rows int, _select string) (e.List, error) {
	if !utility.ValidId(kind) {
		return e.List{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "kind")
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
