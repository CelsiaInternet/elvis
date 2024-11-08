package module

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

var ProjectModules *linq.Model

func DefineProjectModules(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return logs.Panice(err)
	}

	if ProjectModules != nil {
		return nil
	}

	ProjectModules = linq.NewModel(SchemaModule, "PROJECT_MODULES", "Tabla de modulos por projecto", 1)
	ProjectModules.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	ProjectModules.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	ProjectModules.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	ProjectModules.DefineColum("index", "", "INTEGER", 0)
	ProjectModules.DefinePrimaryKey([]string{"project_id", "module_id"})
	ProjectModules.DefineForeignKey("project_id", Projects.Col("_id"))
	ProjectModules.DefineForeignKey("module_id", Modules.Col("_id"))
	ProjectModules.DefineIndex([]string{
		"date_make",
		"index",
	})

	if err := ProjectModules.Init(); err != nil {
		return logs.Panice(err)
	}

	return nil
}

/**
* GetProjectByModule
* @param projectId string
* @param moduleId string
* @return et.Item, error
**/
func GetProjectByModule(projectId, moduleId string) (et.Item, error) {
	return ProjectModules.Data(ProjectModules.Column("index")).
		Where(ProjectModules.Column("project_id").Eq(projectId)).
		And(ProjectModules.Column("module_id").Eq(moduleId)).
		First()
}

/**
* GetProjectModules
* @param projectId string
* @param state string
* @param search string
* @param page int
* @param rows int
* @return et.List, error
**/
func GetProjectModules(projectId, state, search string, page, rows int) (et.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if auxState != "" {
		state = utility.FOR_DELETE

		return linq.From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").Neg(state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Data().
			List(page, rows)
	} else if auxState == "0" {
		return linq.From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").In("-1", state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Data().
			List(page, rows)
	} else {
		return linq.From(Modules, "A").
			Join(Modules.As("A"), ProjectModules.As("B"), ProjectModules.Col("module_id").Eq(Modules.Col("_id"))).
			Where(Modules.Column("_state").Eq(state)).
			And(ProjectModules.Column("project_id").Eq(projectId)).
			And(Modules.Concat("NAME:", Modules.Column("name"), ":DESCRIPTION", Modules.Column("description"), ":DATA:", Modules.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Modules.Column("name"), true).
			Data().
			List(page, rows)
	}
}

/**
* CheckProjectModule
* @param projectId string
* @param moduleId string
* @param chk bool
* @return et.Item, error
**/
func CheckProjectModule(project_id, module_id string, chk bool) (et.Item, error) {
	if !utility.ValidId(project_id) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	if !utility.ValidId(module_id) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !chk {
		result, err := ProjectModules.Delete().
			Where(ProjectModules.Column("project_id").Eq(project_id)).
			And(ProjectModules.Column("module_id").Eq(module_id)).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		return et.Item{
			Ok: result.Ok,
			Result: et.Json{
				"message": utility.OkOrNot(result.Ok, msg.RECORD_DELETE, msg.RECORD_NOT_DELETE),
			},
		}, nil
	}

	current, err := GetProjectByModule(project_id, module_id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		data := et.Json{}
		data.Set("project_id", project_id)
		data.Set("module_id", module_id)
		result, err := ProjectModules.Insert(data).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		return et.Item{
			Ok: result.Ok,
			Result: et.Json{
				"message": utility.OkOrNot(result.Ok, msg.RECORD_CREATE, msg.RECORD_NOT_CREATE),
			},
		}, nil
	}

	return et.Item{
		Ok: true,
		Result: et.Json{
			"message": msg.RECORD_FOUND,
		},
	}, nil
}
