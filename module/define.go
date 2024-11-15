package module

import (
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
)

var initDefine bool

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

	logs.Log("Module", "Define models")

	initDefine = true

	return nil
}

func InitData() error {
	// Initial project and module
	InitProject("-1", "My project", "", et.Json{})
	InitModule("-1", "Admin", "", et.Json{})

	// Initial state types
	InitType("-1", utility.OF_SYSTEM, utility.OF_SYSTEM, "RECORDS", "Default", "Record default")
	InitType("-1", utility.OF_SYSTEM, utility.OF_SYSTEM, "STATE", "System", "Record system")
	InitType("-1", utility.FOR_DELETE, utility.OF_SYSTEM, "STATE", "Delete", "To delete record")
	InitType("-1", utility.ACTIVE, utility.OF_SYSTEM, "STATE", "Activo", "")
	InitType("-1", utility.ARCHIVED, utility.OF_SYSTEM, "STATE", "Archivado", "")
	InitType("-1", utility.CANCELLED, utility.OF_SYSTEM, "STATE", "Cacnelado", "")
	InitType("-1", utility.IN_PROCESS, utility.OF_SYSTEM, "STATE", "En tramite", "")
	InitType("-1", utility.PENDING_APPROVAL, utility.OF_SYSTEM, "STATE", "Pendiente de aprobación", "")
	InitType("-1", utility.APPROVAL, utility.OF_SYSTEM, "STATE", "Aprobado", "")
	InitType("-1", utility.REFUSED, utility.OF_SYSTEM, "STATE", "Rechazado", "")

	// Initial profile types
	InitType("-1", "PROFILE.ADMIN", utility.OF_SYSTEM, "PROFILE", "Admin", "")
	InitType("-1", "PROFILE.DEV", utility.OF_SYSTEM, "PROFILE", "Develop", "")
	InitType("-1", "PROFILE.SUPORT", utility.OF_SYSTEM, "PROFILE", "Suport", "")

	// Initial permision types
	InitType("-1", PERMISION_READ, utility.OF_SYSTEM, "PERMISION", "Read", "")
	InitType("-1", PERMISION_WRITE, utility.OF_SYSTEM, "PERMISION", "Write", "")
	InitType("-1", PERMISION_DELETE, utility.OF_SYSTEM, "PERMISION", "Delete", "")
	InitType("-1", PERMISION_EXECUTE, utility.OF_SYSTEM, "PERMISION", "Execute", "")

	// Initial profile
	InitProfile("-1", "PROFILE.ADMIN", et.Json{})
	InitProfile("-1", "PROFILE.DEV", et.Json{})
	InitProfile("-1", "PROFILE.SUPORT", et.Json{})

	// User Admin
	USER_ADMIN := envar.GetStr("", "USER_ADMIN")
	PASWWOR_ADMIN := envar.GetStr("", "PAWWOR_ADMIN")
	InsertUser("USER.ADMIN", USER_ADMIN, "", "", "", "", PASWWOR_ADMIN)

	// Initial folder
	InitFolder("-1", "-1", "-1", "", "", et.Json{})

	CheckProjectModule("-1", "-1", true)
	CheckRole("-1", "-1", "PROFILE.ADMIN", "USER.ADMIN", true)
	CheckRole("-1", "-1", "PROFILE.DEV", "USER.ADMIN", true)
	CheckRole("-1", "-1", "PROFILE.SUPORT", "USER.ADMIN", true)

	logs.Log("Module", "Init data module")

	return nil
}
