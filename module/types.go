package module

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/linq"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
)

var Types *Model

func DefineTypes() error {
	if err := defineSchema(); err != nil {
		return console.PanicE(err)
	}

	if Types != nil {
		return nil
	}

	Types = NewModel(SchemaModule, "TYPES", "Tabla de tipo", 1)
	Types.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Types.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Types.DefineColum("_state", "", "VARCHAR(80)", ACTIVE)
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

	if err := InitModel(Types); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* Types
*	Handler for CRUD data
 */
func GetTypeByName(kind, name string) (Item, error) {
	return Types.Select().
		Where(Types.Column("kind").Eq(kind)).
		And(Types.Column("name").Eq(name)).
		First()
}

func GetTypeById(id string) (Item, error) {
	return Types.Select().
		Where(Types.Column("_id").Eq(id)).
		First()
}

func GetTypeByIndex(idx int) (Item, error) {
	return Types.Select().
		Where(Types.Column("index").Eq(idx)).
		First()
}

func InitType(projectId, id, state, kind, name, description string) (Item, error) {
	if !ValidId(kind) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "kind")
	}

	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetTypeByName(kind, name)
	if err != nil {
		return Item{}, err
	}

	if current.Ok && current.Id() != id {
		return Item{
			Ok: current.Ok,
			Result: Json{
				"message": RECORD_FOUND,
			},
		}, nil
	}

	id = GenId(id)
	data := Json{}
	data["project_id"] = projectId
	data["_id"] = id
	data["kind"] = kind
	data["name"] = name
	data["description"] = description
	return Types.Upsert(data).
		Where(Types.Column("_id").Eq(id)).
		Command()
}

func UpSetType(projectId, id, kind, name, description string) (Item, error) {
	if !ValidId(id) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "_id")
	}

	if !ValidId(kind) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "kind")
	}

	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetTypeByName(kind, name)
	if err != nil {
		return Item{}, err
	}

	if current.Ok && current.Id() != id {
		return Item{
			Ok: current.Ok,
			Result: Json{
				"message": RECORD_FOUND,
				"_id":     id,
				"index":   current.Index(),
			},
		}, nil
	}

	id = GenId(id)
	data := Json{}
	data["project_id"] = projectId
	data["_id"] = id
	data["kind"] = kind
	data["name"] = name
	data["description"] = description
	return Types.Upsert(data).
		Where(Types.Column("_id").Eq(id)).
		Command()
}

func StateType(id, state string) (Item, error) {
	if !ValidId(state) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "state")
	}

	return Types.Upsert(Json{
		"_state": state,
	}).
		Where(Types.Column("_id").Eq(id)).
		And(Types.Column("_state").Neg(state)).
		Command()
}

func DeleteType(id string) (Item, error) {
	return StateType(id, FOR_DELETE)
}

func AllTypes(projectId, kind, state, search string, page, rows int, _select string) (List, error) {
	if !ValidId(kind) {
		return List{}, console.AlertF(MSG_ATRIB_REQUIRED, "kind")
	}

	if state == "" {
		state = ACTIVE
	}

	auxState := state

	cols := StrToColN(_select)

	if auxState == "*" {
		state = FOR_DELETE

		return Types.Select(cols).
			Where(Types.Column("kind").Eq(kind)).
			And(Types.Column("_state").Neg(state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Types.Concat("NAME:", Types.Column("name"), ":DESCRIPTION", Types.Column("description"), ":DATA:", Types.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Types.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Types.Select(cols).
			Where(Types.Column("kind").Eq(kind)).
			And(Types.Column("_state").In("-1", state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Types.Concat("NAME:", Types.Column("name"), ":DESCRIPTION", Types.Column("description"), ":DATA:", Types.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Types.Column("name"), true).
			List(page, rows)
	} else {
		return Types.Select(cols).
			Where(Types.Column("kind").Eq(kind)).
			And(Types.Column("_state").Eq(state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Types.Concat("NAME:", Types.Column("name"), ":DESCRIPTION", Types.Column("description"), ":DATA:", Types.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Types.Column("name"), true).
			List(page, rows)
	}
}
