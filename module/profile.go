package module

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

var Profiles *linq.Model

var PROFILE_ADMIN = "PROFILE.ADMIN"
var PROFILE_DEV = "PROFILE.DEV"
var PROFILE_SUPORT = "PROFILE.SUPORT"

var profileDefault = map[string]bool{
	PROFILE_ADMIN:  true,
	PROFILE_DEV:    true,
	PROFILE_SUPORT: true,
}

func DefineProfiles(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return console.Panic(err)
	}

	if Profiles != nil {
		return nil
	}

	Profiles = linq.NewModel(SchemaModule, "PROFILES", "Tabla de perfiles", 1)
	Profiles.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Profiles.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Profiles.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	Profiles.DefineColum("profile_tp", "", "VARCHAR(80)", "-1")
	Profiles.DefineColum("_data", "", "JSONB", "{}")
	Profiles.DefineColum("index", "", "INTEGER", 0)
	Profiles.DefinePrimaryKey([]string{"module_id", "profile_tp"})
	Profiles.DefineIndex([]string{
		"date_make",
		"date_update",
		"index",
	})
	Profiles.DefineForeignKey("module_id", Modules.Column("_id"))

	if err := Profiles.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* Profile
*	Handler for CRUD data
**/
func GetProfileById(moduleId, profileTp string) (et.Item, error) {
	return Profiles.Data().
		Where(Profiles.Column("module_id").Eq(moduleId)).
		And(Profiles.Column("profile_tp").Eq(profileTp)).
		First()
}

/**
* InitProfile
* @param moduleId string
* @param profileTp string
* @param data et.Json
* @return et.Item, error
**/
func InitProfile(moduleId, profileTp string, data et.Json) (et.Item, error) {
	if !utility.ValidId(moduleId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "moduleId")
	}

	if !utility.ValidId(profileTp) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "profile_tp")
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return et.Item{}, err
	}

	if !module.Ok {
		return et.Item{}, console.ErrorM(msg.MODULE_NOT_FOUND)
	}

	current, err := GetProfileById(moduleId, profileTp)
	if err != nil {
		return et.Item{}, err
	}

	if current.Ok {
		return et.Item{
			Ok: current.Ok,
			Result: et.Json{
				"message": msg.RECORD_FOUND,
				"_id":     current.Id(),
				"index":   current.Index(),
			},
		}, nil
	}

	data["module_id"] = moduleId
	data["profile_tp"] = profileTp
	return Profiles.Insert(data).
		CommandOne()
}

func UpSetProfile(moduleId, profileTp string, data et.Json) (et.Item, error) {
	if !utility.ValidId(moduleId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "moduleId")
	}

	if !utility.ValidId(profileTp) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "profile_tp")
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return et.Item{}, err
	}

	if !module.Ok {
		return et.Item{}, console.ErrorM(msg.MODULE_NOT_FOUND)
	}

	data["module_id"] = moduleId
	data["profile_tp"] = profileTp
	return Profiles.Upsert(data).
		Where(Profiles.Column("module_id").Eq(moduleId)).
		And(Profiles.Column("profile_tp").Eq(profileTp)).
		CommandOne()
}

func UpSetProfileTp(projectId, moduleId, id, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	profile, err := UpSetType(projectId, id, "PROFILE", name, description)
	if err != nil {
		return et.Item{}, err
	}

	profileTp := profile.Id()
	_, err = UpSetProfile(moduleId, profileTp, data)
	if err != nil {
		return et.Item{}, err
	}

	profile.Set("project_id", projectId)
	profile.Set("module_id", moduleId)
	return profile, nil
}

func DeleteProfile(moduleId, profileTp string) (et.Item, error) {
	current, err := GetProfileById(moduleId, profileTp)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, nil
	}

	return Profiles.Delete().
		Where(Profiles.Column("module_id").Eq(moduleId)).
		And(Profiles.Column("profile_tp").Eq(profileTp)).
		CommandOne()
}

func getProfileFolders(userId, projectId, mainId string) []et.Json {
	sql := `
	SELECT DISTINCT A._DATA||jsonb_build_object(
	'date_make', A.DATE_MAKE,
	'date_update', A.DATE_UPDATE,
	'project_id', $2,
	'module_id', A.MODULE_ID,
	'_state', A._STATE,
	'_id', A._ID,
	'main_id', A.MAIN_ID,
	'name', A.NAME,
	'description', A.DESCRIPTION,
	'index', A.INDEX) AS _DATA,
	A._DATA#>>'{order}' AS ORDEN,
	A.INDEX
	FROM module.FOLDERS AS A
	INNER JOIN module.PROFILE_FOLDERS AS B ON B.FOLDER_ID = A._ID
	INNER JOIN module.MODULE_FOLDERS AS C ON C.FOLDER_ID = A._ID
	WHERE A.MAIN_ID = $1
	AND C.MODULE_ID IN (SELECT C.MODULE_ID FROM module.PROJECT_MODULES AS C WHERE C.PROJECT_ID = $2)
	AND B.PROFILE_TP IN (SELECT D.PROFILE_TP FROM module.ROLES AS D WHERE D.PROJECT_ID = $2 AND D.USER_ID = $3)
	ORDER BY A._DATA#>>'{order}' ASC;`

	items, err := Profiles.Source(linq.StateField, sql, mainId, projectId, userId)
	if err != nil {
		return []et.Json{}
	}

	return items.Result
}

func GetProfileFolders(userId, projectId string) ([]et.Json, error) {
	if !utility.ValidId(userId) {
		return []et.Json{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "userId")
	}

	if !utility.ValidId(projectId) {
		return []et.Json{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	mainId := "-1"
	result := getProfileFolders(userId, projectId, mainId)
	for _, item := range result {
		mainId = item.Id()
		item["folders"] = getProfileFolders(userId, projectId, mainId)
	}

	return result, nil
}

func AllModuleProfiles(projectId, moduleId, state, search string, page, rows int) (et.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	_select := Profiles.All()
	_select = append(_select, Types.Column("_state"), Types.Column("_id"), Types.Column("name"), Types.Column("description"))

	if search != "" {
		return linq.From(Profiles, "A").
			Join(Profiles.As("A"), Types.As("B"), Types.Col("_id").Eq(Profiles.Col("profile_tp"))).
			Where(Types.Column("project_id").In("-1", projectId)).
			And(Profiles.Column("module_id").Eq(moduleId)).
			And(Profiles.Concat("NAME:", Types.As("B").Col("name"), ":DESCRIPTION:", Types.As("B").Col("description"), ":DATA:", Profiles.As("A").Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Types.Column("name"), true).
			Data(_select).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return linq.From(Profiles, "A").
			Join(Profiles.As("A"), Types.As("B"), Types.Col("_id").Eq(Profiles.Col("profile_tp"))).
			Where(Types.Column("_state").Neg(state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Profiles.Column("module_id").Eq(moduleId)).
			OrderBy(Types.Column("name"), true).
			Data(_select).
			List(page, rows)
	} else if auxState == "0" {
		return linq.From(Profiles, "A").
			Join(Profiles.As("A"), Types.As("B"), Types.Col("_id").Eq(Profiles.Col("profile_tp"))).
			Where(Types.Column("_state").In("-1", state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Profiles.Column("module_id").Eq(moduleId)).
			OrderBy(Types.Column("name"), true).
			Data(_select).
			List(page, rows)
	} else {
		return linq.From(Profiles, "A").
			Join(Profiles.As("A"), Types.As("B"), Types.Col("_id").Eq(Profiles.Col("profile_tp"))).
			Where(Types.Column("_state").Eq(state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Profiles.Column("module_id").Eq(moduleId)).
			OrderBy(Types.Column("name"), true).
			Data(_select).
			List(page, rows)
	}
}
