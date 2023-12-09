package bus

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/core"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Apibus *linq.Model

func DefineApimanager() error {
	if err := defineSchema(); err != nil {
		return console.PanicE(err)
	}

	if Apibus != nil {
		return nil
	}

	Apibus = linq.NewModel(SchemaBus, "API_MANAGER", "Tabla", 1)
	Apibus.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Apibus.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Apibus.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Apibus.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Apibus.DefineColum("project", "", "VARCHAR(80)", "-1")
	Apibus.DefineColum("stage", "", "VARCHAR(80)", "development")
	Apibus.DefineColum("method", "", "VARCHAR(80)", "")
	Apibus.DefineColum("path", "", "TEXT", "")
	Apibus.DefineColum("_data", "", "JSONB", "{}")
	Apibus.DefineColum("index", "", "SERIAL", 0)
	Apibus.DefinePrimaryKey([]string{"_id"})
	Apibus.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"project",
		"stage",
		"method",
		"path",
		"index",
	})

	if err := core.InitModel(Apibus); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
*	Handler for CRUD data
 */
func GetApi_managerById(id string) (e.Item, error) {
	return Apibus.Select().
		Where(Apibus.Column("_id").Eq(id)).
		First()
}

func UpSertApi_manager(project, stage, method, path, id string, data e.Json) (e.Item, error) {
	if !utility.ValidStr(project, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "project")
	}

	if !utility.ValidStr(stage, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "stage")
	}

	if !utility.ValidStr(method, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "method")
	}

	id = utility.GenId(id)
	data["project"] = project
	data["_id"] = id
	return Apibus.Upsert(data).
		Where(Apibus.Column("_id").Eq(id)).
		Command()
}

func StateApi_manager(id, state string) (e.Item, error) {
	if !utility.ValidId(state) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	return Apibus.Upsert(e.Json{
		"_state": state,
	}).
		Where(Apibus.Column("_id").Eq(id)).
		And(Apibus.Column("_state").Neg(state)).
		Command()
}

func DeleteApi_manager(id string) (e.Item, error) {
	return StateApi_manager(id, utility.FOR_DELETE)
}

func AllApi_manager(project, state, search string, page, rows int, _select string) (e.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	cols := linq.StrToCols(_select)

	if search != "" {
		return Apibus.Select(cols).
			Where(Apibus.Column("project").Like("%"+project+"%")).
			And(Apibus.Concat("STAGE:", Apibus.Column("stage"), "METHOD:", Apibus.Column("method"), "PATH:", Apibus.Column("path"), "DATA:", Apibus.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Apibus.Column("path"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Apibus.Select(cols).
			Where(Apibus.Column("_state").Neg(state)).
			And(Apibus.Column("project").Like("%"+project+"%")).
			OrderBy(Apibus.Column("path"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Apibus.Select(cols).
			Where(Apibus.Column("_state").In("-1", state)).
			And(Apibus.Column("project").Like("%"+project+"%")).
			OrderBy(Apibus.Column("path"), true).
			List(page, rows)
	} else {
		return Apibus.Select(cols).
			Where(Apibus.Column("_state").Eq(state)).
			And(Apibus.Column("project").Like("%"+project+"%")).
			OrderBy(Apibus.Column("path"), true).
			List(page, rows)
	}
}
