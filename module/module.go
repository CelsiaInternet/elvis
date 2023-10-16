package module

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/linq"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
	_ "github.com/joho/godotenv/autoload"
)

var Modules *Model

func DefineModules() error {
	if err := DefineCoreSchema(); err != nil {
		return console.PanicE(err)
	}

	if Modules != nil {
		return nil
	}

	Modules = NewModel(SchemaCore, "MODULES", "Tabla de modulos", 1)
	Modules.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Modules.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Modules.DefineColum("_state", "", "VARCHAR(80)", ACTIVE)
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
	Modules.Trigger(AfterInsert, func(model *Model, old, new *Json, data Json) {
		id := new.Id()
		InitProfile(id, "PROFILE.ADMIN", Json{})
		InitProfile(id, "PROFILE.DEV", Json{})
		InitProfile(id, "PROFILE.SUPORT", Json{})
		CheckProjectModule("-1", id, true)
		CheckRole("-1", id, "PROFILE.ADMIN", "USER.ADMIN", true)
	})

	if err := InitModel(Modules); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* Module
*	Handler for CRUD data
 */
func GetModuleByName(name string) (Item, error) {
	return Modules.Select().
		Where(Modules.Column("name").Eq(name)).
		First()
}

func GetModuleById(id string) (Item, error) {
	return Modules.Select().
		Where(Modules.Column("_id").Eq(id)).
		First()
}

func IsInit() (Item, error) {
	count := Users.Select().
		Count()

	return Item{
		Ok: count > 0,
		Result: Json{
			"message": OkOrNot(count > 0, SYSTEM_HAVE_ADMIN, SYSTEM_NOT_HAVE_ADMIN),
		},
	}, nil
}

func InitModule(id, name, description string, data Json) (Item, error) {
	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetModuleByName(name)
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
	data.Set("_id", id)
	data.Set("name", name)
	data.Set("description", description)
	item, err := Modules.Upsert(data).
		Where(Modules.Column("_id").Eq(id)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func UpSetModule(id, name, description string, data Json) (Item, error) {
	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetModuleByName(name)
	if err != nil {
		return Item{}, err
	}

	if current.Ok && current.Id() != id {
		return Item{
			Ok: current.Ok,
			Result: Json{
				"message": RECORD_FOUND,
				"_id":     current.Id(),
				"index":   current.Index(),
			},
		}, nil
	}

	id = GenId(id)
	data.Set("_id", id)
	data.Set("name", name)
	data.Set("description", description)
	item, err := Modules.Upsert(data).
		Where(Modules.Column("_id").Eq(id)).
		Command()

	return item, nil
}

func StateModule(id, state string) (Item, error) {
	if !ValidId(state) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "state")
	}

	return Modules.Upsert(Json{
		"_state": state,
	}).
		Where(Modules.Column("_id").Eq(id)).
		And(Modules.Column("_state").Neg(state)).
		Command()
}

func DeleteModule(id string) (Item, error) {
	return StateModule(id, FOR_DELETE)
}

func AllModules(state, search string, page, rows int, _select string) (List, error) {
	if state == "" {
		state = ACTIVE
	}

	auxState := state

	cols := StrToColN(_select)

	if auxState == "*" {
		state = FOR_DELETE

		return Modules.Select(cols).
			Where(Modules.Column("_state").Neg(state)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Modules.Select(cols).
			Where(Modules.Column("_state").In("-1", state)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	} else {
		return Modules.Select(cols).
			Where(Modules.Column("_state").Eq(state)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	}
}
