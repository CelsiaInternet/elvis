package module

import (
	"time"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
)

var (
	initDefine bool
	TEST_TOKEN string
)

func InitDefine(db *jdb.DB) error {
	if initDefine {
		return nil
	}

	if err := DefineUsers(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineProjects(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineTypes(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineModules(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineFolders(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineProfiles(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineRoles(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineTokens(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineMigration(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineModuleFolders(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineModuleFolders(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineProjectModules(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefineProfileFolders(db); err != nil {
		return logs.Panice(err)
	}
	if err := DefinePermisions(db); err != nil {
		return logs.Panice(err)
	}

	logs.Log(PackageName, "Define models")

	initDefine = true

	return nil
}

func InitData() error {
	// Initial project and module
	project_id := "-1"
	module_id := "-1"
	InitProject(project_id, "System project", et.Json{})
	InitModule(module_id, "System", et.Json{})

	// Initial state types
	InitType(project_id, utility.OF_SYSTEM, utility.OF_SYSTEM, "RECORDS", "Default")
	InitType(project_id, utility.OF_SYSTEM, utility.OF_SYSTEM, "STATE", "System")
	InitType(project_id, utility.FOR_DELETE, utility.OF_SYSTEM, "STATE", "Delete")
	InitType(project_id, utility.ACTIVE, utility.OF_SYSTEM, "STATE", "Activo")
	InitType(project_id, utility.ARCHIVED, utility.OF_SYSTEM, "STATE", "Archivado")
	InitType(project_id, utility.CANCELLED, utility.OF_SYSTEM, "STATE", "Cacnelado")
	InitType(project_id, utility.IN_PROCESS, utility.OF_SYSTEM, "STATE", "En tramite")
	InitType(project_id, utility.PENDING_APPROVAL, utility.OF_SYSTEM, "STATE", "Pendiente de aprobación")
	InitType(project_id, utility.APPROVAL, utility.OF_SYSTEM, "STATE", "Aprobado")
	InitType(project_id, utility.REFUSED, utility.OF_SYSTEM, "STATE", "Rechazado")

	// Initial profile types
	InitType(project_id, "PROFILE.ADMIN", utility.OF_SYSTEM, "PROFILE", "Admin")
	InitType(project_id, "PROFILE.DEV", utility.OF_SYSTEM, "PROFILE", "Develop")
	InitType(project_id, "PROFILE.SUPORT", utility.OF_SYSTEM, "PROFILE", "Suport")

	// Initial permision types
	InitType(project_id, PERMISION_READ, utility.OF_SYSTEM, "PERMISION", "Read")
	InitType(project_id, PERMISION_WRITE, utility.OF_SYSTEM, "PERMISION", "Write")
	InitType(project_id, PERMISION_DELETE, utility.OF_SYSTEM, "PERMISION", "Delete")
	InitType(project_id, PERMISION_EXECUTE, utility.OF_SYSTEM, "PERMISION", "Execute")

	// Initial profile
	InitProfile(module_id, "PROFILE.ADMIN", et.Json{})
	InitProfile(module_id, "PROFILE.DEV", et.Json{})
	InitProfile(module_id, "PROFILE.SUPORT", et.Json{})

	// User Admin
	USER_ADMIN := envar.GetStr("", "USER_ADMIN")
	if len(USER_ADMIN) == 0 {
		return console.NewError("USER_ADMIN is empty")
	}

	PASSWORD_ADMIN := envar.GetStr("", "PASSWORD_ADMIN")
	if len(PASSWORD_ADMIN) == 0 {
		return console.NewError("PAWWOR_ADMIN is empty")
	}

	_, err := InsertUser("USER.ADMIN", USER_ADMIN, "Admin", "", "", "", PASSWORD_ADMIN)
	if err == nil {
		TEST_TOKEN, _ = claim.NewToken(USER_ADMIN, PackageName, USER_ADMIN, USER_ADMIN, "apiREST", 24*time.Hour)
		logs.Logf(PackageName, `Token:%s`, TEST_TOKEN)
	} else {
		key := claim.GetTokenKey(PackageName, "apiREST", USER_ADMIN)
		TEST_TOKEN, _ = claim.GetToken(key)
		if len(TEST_TOKEN) != 0 {
			logs.Logf(PackageName, `Token:%s`, TEST_TOKEN)
		}
	}

	// Initial folder
	InitFolder(module_id, "-1", "-1", "main", et.Json{})
	defaultFolders(module_id)

	CheckProjectModule(project_id, module_id, true)
	CheckRole(project_id, module_id, "PROFILE.ADMIN", "USER.ADMIN", true)
	CheckRole(project_id, module_id, "PROFILE.DEV", "USER.ADMIN", true)
	CheckRole(project_id, module_id, "PROFILE.SUPORT", "USER.ADMIN", true)

	logs.Log(PackageName, "Init data module")

	return nil
}
