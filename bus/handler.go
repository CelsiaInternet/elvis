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
		return console.Panic(err)
	}

	if Apibus != nil {
		return nil
	}

	Apibus = linq.NewModel(SchemaBus, "API_MANAGER", "Tabla", 1)
	Apibus.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Apibus.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Apibus.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Apibus.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Apibus.DefineColum("package_name", "", "VARCHAR(250)", "")
	Apibus.DefineColum("package_path", "", "VARCHAR(250)", "")
	Apibus.DefineColum("kind", "", "VARCHAR(80)", "")
	Apibus.DefineColum("method", "", "VARCHAR(80)", "")
	Apibus.DefineColum("host", "", "VARCHAR(250)", "")
	Apibus.DefineColum("path", "", "TEXT", "")
	Apibus.DefineColum("_data", "", "JSONB", "{}")
	Apibus.DefineColum("index", "", "SERIAL", 0)
	Apibus.DefinePrimaryKey([]string{"_id"})
	Apibus.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"package_name",
		"package_path",
		"kind",
		"method",
		"host",
		"path",
		"index",
	})

	if err := core.InitModel(Apibus); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
*	Handler for CRUD data
 */
func GetApimanagerById(id string) (e.Item, error) {
	return Apibus.Data().
		Where(Apibus.Column("_id").Eq(id)).
		First()
}

func GetApimanagerByPath(method, path string) (e.Item, error) {
	return Apibus.Data().
		Where(Apibus.Column("path").Eq(path)).
		And(Apibus.Column("method").Eq(method)).
		First()
}

func UpSertApimanager(package_name, kind, method, path string, data e.Json) (e.Item, error) {
	if !utility.ValidStr(package_name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "package_name")
	}

	if !utility.ValidStr(kind, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "kind")
	}

	if !utility.ValidStr(method, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "method")
	}

	current, err := GetApimanagerByPath(method, path)
	if err != nil {
		return e.Item{}, err
	}

	var id string
	if current.Ok {
		id = current.Key("_id")
	} else {
		id = utility.NewId()
	}

	data["package_name"] = package_name
	data["kind"] = kind
	data["method"] = method
	data["path"] = path
	data["_id"] = id
	return Apibus.Upsert(data).
		Where(Apibus.Column("_id").Eq(id)).
		CommandOne()
}

func StateApimanager(id, state string) (e.Item, error) {
	if !utility.ValidId(state) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	return Apibus.Update(e.Json{
		"_state": state,
	}).
		Where(Apibus.Column("_id").Eq(id)).
		And(Apibus.Column("_state").Neg(state)).
		CommandOne()
}

func DeleteApimanager(id string) (e.Item, error) {
	return StateApimanager(id, utility.FOR_DELETE)
}

func AllApimanager(package_name, state, search string, page, rows int, _select string) (e.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Apibus.Data(_select).
			Where(Apibus.Column("package_name").Like("%"+package_name+"%")).
			And(Apibus.Concat("kind:", Apibus.Column("kind"), "METHOD:", Apibus.Column("method"), "PATH:", Apibus.Column("path"), "DATA:", Apibus.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Apibus.Column("path"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Apibus.Data(_select).
			Where(Apibus.Column("_state").Neg(state)).
			And(Apibus.Column("package_name").Like("%"+package_name+"%")).
			OrderBy(Apibus.Column("path"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Apibus.Data(_select).
			Where(Apibus.Column("_state").In("-1", state)).
			And(Apibus.Column("package_name").Like("%"+package_name+"%")).
			OrderBy(Apibus.Column("path"), true).
			List(page, rows)
	} else {
		return Apibus.Data(_select).
			Where(Apibus.Column("_state").Eq(state)).
			And(Apibus.Column("package_name").Like("%"+package_name+"%")).
			OrderBy(Apibus.Column("path"), true).
			List(page, rows)
	}
}
