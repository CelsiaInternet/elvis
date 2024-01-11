package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/core"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Projects *linq.Model
var ProjectModules *linq.Model

func DefineProjects() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
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
	Projects.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *e.Json, data e.Json) error {
		moduleId := data.Key("module_id")
		if moduleId != "" {
			id := new.Id()
			CheckProjectModule(id, moduleId, true)
			CheckRole(id, moduleId, "PROFILE.ADMIN", "USER.ADMIN", true)
		}

		return nil
	})

	if err := core.InitModel(Projects); err != nil {
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

	ProjectModules = linq.NewModel(SchemaModule, "PROJECT_MODULES", "Tabla de moduloes por projecto", 1)
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

	if err := core.InitModel(ProjectModules); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* Project
*	Handler for CRUD data
 */
func GetProjectById(id string) (e.Item, error) {
	return Projects.Select().
		Where(Projects.Column("_id").Eq(id)).
		First()
}

func GetProjectName(name string) (e.Item, error) {
	return Projects.Select().
		Where(Projects.Column("name").Eq(name)).
		First()
}

func GetProjectByModule(projectId, moduleId string) (e.Item, error) {
	return ProjectModules.Select(ProjectModules.Column("index")).
		Where(ProjectModules.Column("project_id").Eq(projectId)).
		And(ProjectModules.Column("module_id").Eq(moduleId)).
		First()
}

func InitProject(id, name, description string, data e.Json) (e.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	id = utility.GenId(id)
	data.Set("_id", id)
	data.Set("name", name)
	data.Set("description", description)
	item, err := Projects.Upsert(data).
		Where(Projects.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func UpSetProject(id, moduleId, name, description string, data e.Json) (e.Item, error) {
	if !utility.ValidId(moduleId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	current, err := GetProjectName(name)
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
	data.Set("module_id", moduleId)
	item, err := Projects.Upsert(data).
		Where(Projects.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func StateProject(id, state string) (e.Item, error) {
	if !utility.ValidId(state) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	return Projects.Update(e.Json{
		"_state": state,
	}).
		Where(Projects.Column("_id").Eq(id)).
		And(Projects.Column("_state").Neg(state)).
		CommandOne()
}

func DeleteProject(id string) (e.Item, error) {
	return StateProject(id, utility.FOR_DELETE)
}

func AllProjects(state, search string, page, rows int, _select string) (e.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Projects.Select(_select).
			Where(Projects.Concat("NAME:", Projects.Column("name"), ":DESCRIPTION:", Projects.Column("description"), ":DATA:", Projects.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Projects.Select(_select).
			Where(Projects.Column("_state").Neg(state)).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Projects.Select(_select).
			Where(Projects.Column("_state").In("-1", state)).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	} else {
		return Projects.Select(_select).
			Where(Projects.Column("_state").Eq(state)).
			OrderBy(Projects.Column("name"), true).
			List(page, rows)
	}
}

func GetProjectModules(projectId, state, search string, page, rows int) (e.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if auxState == "*" {
		state = utility.FOR_DELETE

		return linq.From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").Neg(state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Select().
			List(page, rows)
	} else if auxState == "0" {
		return linq.From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").In("-1", state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Select().
			List(page, rows)
	} else {
		return linq.From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").Eq(state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Select().
			List(page, rows)
	}
}

func CheckProjectModule(project_id, module_id string, chk bool) (e.Item, error) {
	if !utility.ValidId(project_id) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	if !utility.ValidId(module_id) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	data := e.Json{}
	data.Set("project_id", project_id)
	data.Set("module_id", module_id)
	if chk {
		current, err := GetProjectByModule(project_id, module_id)
		if err != nil {
			return e.Item{}, err
		}

		if current.Ok {
			return e.Item{
				Ok: current.Ok,
				Result: e.Json{
					"message": msg.RECORD_NOT_UPDATE,
					"index":   current.Index(),
				},
			}, nil
		}

		return ProjectModules.Insert(data).
			CommandOne()
	} else {
		return ProjectModules.Delete().
			Where(ProjectModules.Column("project_id").Eq(project_id)).
			And(ProjectModules.Column("module_id").Eq(module_id)).
			CommandOne()
	}
}
