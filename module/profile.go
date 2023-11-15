package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/core"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Profiles *linq.Model
var ProfileFolders *linq.Model

func DefineProfiles() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
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

	if err := core.InitModel(Profiles); err != nil {
		return console.PanicE(err)
	}

	return nil
}

func DefineProfileFolders() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
	}

	if ProfileFolders != nil {
		return nil
	}

	ProfileFolders = linq.NewModel(SchemaModule, "PROFILE_FOLDERS", "Tabla de carpetas por perfil", 1)
	ProfileFolders.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	ProfileFolders.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	ProfileFolders.DefineColum("profile_tp", "", "VARCHAR(80)", "-1")
	ProfileFolders.DefineColum("folder_id", "", "VARCHAR(80)", "-1")
	ProfileFolders.DefineColum("index", "", "INTEGER", 0)
	ProfileFolders.DefinePrimaryKey([]string{"module_id", "profile_tp", "folder_id"})
	ProfileFolders.DefineIndex([]string{
		"date_make",
		"index",
	})
	ProfileFolders.DefineForeignKey("module_id", Modules.Column("_id"))

	if err := core.InitModel(ProfileFolders); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* Profile
*	Handler for CRUD data
**/
func GetProfileById(moduleId, profileTp string) (e.Item, error) {
	return Profiles.Select().
		Where(Profiles.Column("module_id").Eq(moduleId)).
		And(Profiles.Column("profile_tp").Eq(profileTp)).
		First()
}

func InitProfile(moduleId, profileTp string, data e.Json) (e.Item, error) {
	if !utility.ValidId(moduleId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "moduleId")
	}

	if !utility.ValidId(profileTp) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "profile_tp")
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return e.Item{}, err
	}

	if !module.Ok {
		return e.Item{}, console.ErrorM(msg.MODULE_NOT_FOUND)
	}

	current, err := GetProfileById(moduleId, profileTp)
	if err != nil {
		return e.Item{}, err
	}

	if current.Ok {
		return e.Item{
			Ok: current.Ok,
			Result: e.Json{
				"message": msg.RECORD_FOUND,
				"_id":     current.Id(),
				"index":   current.Index(),
			},
		}, nil
	}

	data["module_id"] = moduleId
	data["profile_tp"] = profileTp
	return Profiles.Insert(data).
		Command()
}

func UpSetProfile(moduleId, profileTp string, data e.Json) (e.Item, error) {
	if !utility.ValidId(moduleId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "moduleId")
	}

	if !utility.ValidId(profileTp) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "profile_tp")
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return e.Item{}, err
	}

	if !module.Ok {
		return e.Item{}, console.ErrorM(msg.MODULE_NOT_FOUND)
	}

	data["module_id"] = moduleId
	data["profile_tp"] = profileTp
	return Profiles.Upsert(data).
		Where(Profiles.Column("module_id").Eq(moduleId)).
		And(Profiles.Column("profile_tp").Eq(profileTp)).
		Command()
}

func UpSetProfileTp(projectId, moduleId, id, name, description string, data e.Json) (e.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	profile, err := UpSetType(projectId, id, "PROFILE", name, description)
	if err != nil {
		return e.Item{}, err
	}

	profileTp := profile.Id()
	_, err = UpSetProfile(moduleId, profileTp, data)
	if err != nil {
		return e.Item{}, err
	}

	profile.Set("project_id", projectId)
	profile.Set("module_id", moduleId)
	return profile, nil
}

func DeleteProfile(moduleId, profileTp string) (e.Item, error) {
	current, err := GetProfileById(moduleId, profileTp)
	if err != nil {
		return e.Item{}, err
	}

	if !current.Ok {
		return e.Item{}, console.ErrorM(msg.RECORD_NOT_FOUND)
	}

	return Profiles.Delete().
		Where(Profiles.Column("module_id").Eq(moduleId)).
		And(Profiles.Column("profile_tp").Eq(profileTp)).
		Command()
}

func GetProfileFolderById(moduleId, profileTp, folderId string) (e.Item, error) {
	return ProfileFolders.Select().
		Where(ProfileFolders.Column("module_id").Eq(moduleId)).
		And(ProfileFolders.Column("profile_tp").Eq(profileTp)).
		And(ProfileFolders.Column("folder_id").Eq(folderId)).
		First()
}

func CheckProfileFolder(moduleId, profileTp, folderId string, chk bool) (e.Item, error) {
	if !utility.ValidId(moduleId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidId(profileTp) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "profile_tp")
	}

	if !utility.ValidId(folderId) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "folder_id")
	}

	profile, err := GetTypeById(profileTp)
	if err != nil {
		return e.Item{}, err
	}

	if !profile.Ok {
		return e.Item{}, console.AlertF(msg.PROFILE_NOT_FOUND, profileTp)
	}

	data := e.Json{}
	data.Set("module_id", moduleId)
	data.Set("profile_tp", profileTp)
	data.Set("folder_id", folderId)
	if chk {
		return ProfileFolders.Insert(data).
			Where(ProfileFolders.Column("module_id").Eq(moduleId)).
			And(ProfileFolders.Column("profile_tp").Eq(profileTp)).
			And(ProfileFolders.Column("folder_id").Eq(folderId)).
			Returns(ProfileFolders.Column("index")).
			Command()
	} else {
		return ProfileFolders.Delete().
			Where(ProfileFolders.Column("module_id").Eq(moduleId)).
			And(ProfileFolders.Column("profile_tp").Eq(profileTp)).
			And(ProfileFolders.Column("folder_id").Eq(folderId)).
			Command()
	}
}

func getProfileFolders(userId, projectId, mainId string) []e.Json {
	items, err := linq.From(Folders, "A").
		Join(Folders.As("A"), ProfileFolders.As("B"), ProfileFolders.Column("folder_id").Eq(Folders.Column("_id"))).
		Where(Folders.Column("main_id").Eq(mainId)).
		And(Folders.Column("module_id").In(
			linq.From(ProjectModules, "C").
				Where(ProjectModules.Column("project_id").Eq(projectId)).
				Select(ProjectModules.Column("module_id")).SQL())).
		And(ProfileFolders.Column("profile_tp").In(
			linq.From(Roles, "D").
				Where(Roles.Column("project_id").Eq(projectId)).
				And(Roles.Column("user_id").Eq(userId)).
				Select(Roles.Column("profile_tp")).SQL())).
		OrderBy(Folders.Column("index"), true).
		Find()
	if err != nil {
		return []e.Json{}
	}

	for _, item := range items.Result {
		item["project_id"] = projectId
	}

	return items.Result
}

func GetProfileFolders(userId, projectId string) ([]e.Json, error) {
	if !utility.ValidId(userId) {
		return []e.Json{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "clientId")
	}

	if !utility.ValidId(projectId) {
		return []e.Json{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	mainId := "-1"
	result := getProfileFolders(userId, projectId, mainId)
	for _, item := range result {
		mainId = item.Id()
		item["folders"] = getProfileFolders(userId, projectId, mainId)
	}

	return result, nil
}

func AllModuleProfiles(projectId, moduleId, state, search string, page, rows int) (e.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	_select := Profiles.All()
	_select = append(_select, Types.Column("_state"), Types.Column("_id"), Types.Column("name"), Types.Column("description"))

	if auxState == "*" {
		state = utility.FOR_DELETE

		return linq.From(Profiles, "A").
			Join(Profiles.As("A"), Types.As("B"), Types.Col("_id").Eq(Profiles.Col("profile_tp"))).
			Where(Types.Column("_state").Neg(state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Profiles.Column("module_id").Eq(moduleId)).
			And(Profiles.Concat("NAME:", Types.As("B").Col("name"), ":DESCRIPTION:", Types.As("B").Col("description"), ":DATA:", Profiles.As("A").Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Types.Column("name"), true).
			Select(_select).
			List(page, rows)
	} else if auxState == "0" {
		return linq.From(Profiles, "A").
			Join(Profiles.As("A"), Types.As("B"), Types.Col("_id").Eq(Profiles.Col("profile_tp"))).
			Where(Types.Column("_state").In("-1", state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Profiles.Column("module_id").Eq(moduleId)).
			And(Profiles.Concat("NAME:", Types.As("B").Col("name"), ":DESCRIPTION:", Types.As("B").Col("description"), ":DATA:", Profiles.As("A").Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Types.Column("name"), true).
			Select(_select).
			List(page, rows)
	} else {
		return linq.From(Profiles, "A").
			Join(Profiles.As("A"), Types.As("B"), Types.Col("_id").Eq(Profiles.Col("profile_tp"))).
			Where(Types.Column("_state").Eq(state)).
			And(Types.Column("project_id").In("-1", projectId)).
			And(Profiles.Column("module_id").Eq(moduleId)).
			And(Profiles.Concat("NAME:", Types.As("B").Col("name"), ":DESCRIPTION:", Types.As("B").Col("description"), ":DATA:", Profiles.As("A").Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Types.Column("name"), true).
			Select(_select).
			List(page, rows)
	}
}
