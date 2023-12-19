package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/core"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
	_ "github.com/joho/godotenv/autoload"
)

var Modules *linq.Model

func DefineModules() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
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
	Modules.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *e.Json, data e.Json) error {
		id := new.Id()
		InitProfile(id, "PROFILE.ADMIN", e.Json{})
		InitProfile(id, "PROFILE.DEV", e.Json{})
		InitProfile(id, "PROFILE.SUPORT", e.Json{})
		CheckProjectModule("-1", id, true)
		CheckRole("-1", id, "PROFILE.ADMIN", "USER.ADMIN", true)

		return nil
	})

	if err := core.InitModel(Modules); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* Module
*	Handler for CRUD data
 */
func GetModuleByName(name string) (e.Item, error) {
	return Modules.Select().
		Where(Modules.Column("name").Eq(name)).
		First()
}

func GetModuleById(id string) (e.Item, error) {
	return Modules.Select().
		Where(Modules.Column("_id").Eq(id)).
		First()
}

func IsInit() (e.Item, error) {
	count := Users.Select().
		Count()

	return e.Item{
		Ok: count > 0,
		Result: e.Json{
			"message": utility.OkOrNot(count > 0, msg.SYSTEM_HAVE_ADMIN, msg.SYSTEM_NOT_HAVE_ADMIN),
		},
	}, nil
}

func InitModule(id, name, description string, data e.Json) (e.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetModuleByName(name)
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
	data.Set("_id", id)
	data.Set("name", name)
	data.Set("description", description)
	item, err := Modules.Upsert(data).
		Where(Modules.Column("_id").Eq(id)).
		Command()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func UpSetModule(id, name, description string, data e.Json) (e.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetModuleByName(name)
	if err != nil {
		return e.Item{}, err
	}

	if current.Ok && current.Id() != id {
		return e.Item{
			Ok: current.Ok,
			Result: e.Json{
				"message": msg.RECORD_FOUND,
				"_id":     current.Id(),
				"index":   current.Index(),
			},
		}, nil
	}

	id = utility.GenId(id)
	data.Set("_id", id)
	data.Set("name", name)
	data.Set("description", description)
	item, err := Modules.Upsert(data).
		Where(Modules.Column("_id").Eq(id)).
		Command()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func StateModule(id, state string) (e.Item, error) {
	if !utility.ValidId(state) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	return Modules.Update(e.Json{
		"_state": state,
	}).
		Where(Modules.Column("_id").Eq(id)).
		And(Modules.Column("_state").Neg(state)).
		Command()
}

func DeleteModule(id string) (e.Item, error) {
	return StateModule(id, utility.FOR_DELETE)
}

func AllModules(state, search string, page, rows int, _select string) (e.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	cols := linq.StrToCols(_select)

	if search != "" {
		return Modules.Select(cols).
			Where(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Modules.Select(cols).
			Where(Modules.Column("_state").Neg(state)).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Modules.Select(cols).
			Where(Modules.Column("_state").In("-1", state)).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	} else {
		return Modules.Select(cols).
			Where(Modules.Column("_state").Eq(state)).
			OrderBy(Modules.Column("name"), true).
			List(page, rows)
	}
}
