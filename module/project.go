package module

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/linq"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utility"
)

var Projects *Model
var ProjectModules *Model

func DefineProjects() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
	}

	if Projects != nil {
		return nil
	}

	Projects = NewModel(SchemaModule, "PROJECTS", "Tabla de projectos", 1)
	Projects.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Projects.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Projects.DefineColum("_state", "", "VARCHAR(80)", ACTIVE)
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
	Projects.Trigger(AfterInsert, func(model *Model, old, new *Json, data Json) error {
		moduleId := data.Key("module_id")
		if moduleId != "" {
			id := new.Id()
			CheckProjectModule(id, moduleId, true)
			CheckRole(id, moduleId, "PROFILE.ADMIN", "USER.ADMIN", true)
		}

		return nil
	})

	if err := InitModel(Projects); err != nil {
		return console.PanicE(err)
	}

	return nil
}

func DefineProjectModules() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
	}

	if ProjectModules != nil {
		return nil
	}

	ProjectModules = NewModel(SchemaModule, "PROJECT_MODULES", "Tabla de moduloes por projecto", 1)
	ProjectModules.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	ProjectModules.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	ProjectModules.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	ProjectModules.DefineColum("index", "", "INTEGER", 0)
	ProjectModules.DefinePrimaryKey([]string{"project_id", "module_id"})
	ProjectModules.DefineIndex([]string{
		"date_make",
		"index",
	})
	ProjectModules.DefineForeignKey("project_id", Projects.Column("_id"))
	ProjectModules.DefineForeignKey("module_id", Modules.Column("_id"))

	if err := InitModel(ProjectModules); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* Project
*	Handler for CRUD data
 */
func GetProjectById(id string) (Item, error) {
	return Projects.Select().
		Where(Projects.Column("_id").Eq(id)).
		First()
}

func GetProjectName(name string) (Item, error) {
	return Projects.Select().
		Where(Projects.Column("name").Eq(name)).
		First()
}

func GetProjectByModule(projectId, moduleId string) (Item, error) {
	return ProjectModules.Select(ProjectModules.Column("index")).
		Where(ProjectModules.Column("project_id").Eq(projectId)).
		And(ProjectModules.Column("module_id").Eq(moduleId)).
		First()
}

func InitProject(id, name, description string, data Json) (Item, error) {
	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "name")
	}

	id = GenId(id)
	data.Set("_id", id)
	data.Set("name", name)
	data.Set("description", description)
	item, err := Projects.Upsert(data).
		Where(Projects.Column("_id").Eq(id)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func UpSetProject(id, moduleId, name, description string, data Json) (Item, error) {
	if !ValidId(moduleId) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "module_id")
	}

	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetProjectName(name)
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
	data.Set("module_id", moduleId)
	item, err := Projects.Upsert(data).
		Where(Projects.Column("_id").Eq(id)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func StateProject(id, state string) (Item, error) {
	if !ValidId(state) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "state")
	}

	return Projects.Upsert(Json{
		"_state": state,
	}).
		Where(Projects.Column("_id").Eq(id)).
		And(Projects.Column("_state").Neg(state)).
		Command()
}

func DeleteProject(id string) (Item, error) {
	return StateProject(id, FOR_DELETE)
}

func AllProjects(state, search string, page, rows int, _select string) (List, error) {
	if state == "" {
		state = ACTIVE
	}

	auxState := state

	cols := StrToCols(_select)

	if auxState == "*" {
		state = FOR_DELETE

		return Projects.Select(cols).
			Where(Projects.Column("_state").Neg(state)).
			And(Projects.Concat("NAME:", Projects.Column("name"), ":DESCRIPTION:", Projects.Column("description"), ":DATA:", Projects.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Projects.Select(cols).
			Where(Projects.Column("_state").In("-1", state)).
			And(Projects.Concat("NAME:", Projects.Column("name"), ":DESCRIPTION:", Projects.Column("description"), ":DATA:", Projects.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	} else {
		return Projects.Select(cols).
			Where(Projects.Column("_state").Eq(state)).
			And(Projects.Concat("NAME:", Projects.Column("name"), ":DESCRIPTION:", Projects.Column("description"), ":DATA:", Projects.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	}
}

func GetProjectModules(projectId, state, search string, page, rows int) (List, error) {
	if state == "" {
		state = ACTIVE
	}

	auxState := state

	if auxState == "*" {
		state = FOR_DELETE

		return From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").Neg(state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Select().
			List(page, rows)
	} else if auxState == "0" {
		return From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").In("-1", state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Select().
			List(page, rows)
	} else {
		return From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").Eq(state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Select().
			List(page, rows)
	}
}

func CheckProjectModule(project_id, module_id string, chk bool) (Item, error) {
	if !ValidId(project_id) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "project_id")
	}

	if !ValidId(module_id) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "module_id")
	}

	data := Json{}
	data.Set("project_id", project_id)
	data.Set("module_id", module_id)
	if chk {
		return ProjectModules.Insert(data).
			Where(ProjectModules.Column("project_id").Eq(project_id)).
			And(ProjectModules.Column("module_id").Eq(module_id)).
			Command()
	} else {
		return ProjectModules.Delete().
			Where(ProjectModules.Column("project_id").Eq(project_id)).
			And(ProjectModules.Column("module_id").Eq(module_id)).
			Command()
	}
}
