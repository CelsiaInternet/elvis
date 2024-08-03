package module

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
)

var initDefine bool

func InitDefine(db *sql.DB) error {
	if initDefine {
		return nil
	}

	if err := DefineUsers(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineProjects(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineTypes(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineModules(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineFolders(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineProfiles(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineRoles(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineModuleFolders(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineProjectModules(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineProfileFolders(db); err != nil {
		return console.Panic(err)
	}
	if err := DefineTokens(db); err != nil {
		return console.Panic(err)
	}

	console.LogK("Module", "Define models")

	initDefine = true

	return nil
}

func InitData() error {
	if _, err := Projects.Upsert(et.Json{
		"_id":  "-1",
		"name": "My project",
	}).
		Where(Projects.Column("_id").Eq("-1")).
		Debug().
		CommandOne(); err != nil {
		return err
	}

	/*
		if _, err := Modules.Upsert(et.Json{
			"_id":  "-1",
			"name": "Admin",
		}).
			Where(Modules.Column("_id").Eq("-1")).
			CommandOne(); err != nil {
			return err
		}

		if _, err := Types.Upsert(et.Json{
			"_id": "-1",
		}).
			Where(Types.Column("_id").Eq("-1")).
			CommandOne(); err != nil {
			return err
		}

		// Initial state types
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

		InitProfile("-1", "PROFILE.ADMIN", et.Json{})
		InitProfile("-1", "PROFILE.DEV", et.Json{})
		InitProfile("-1", "PROFILE.SUPORT", et.Json{})

		// User Admin
		ADMIN_COUNTRY := envar.EnvarStr("", "ADMIN_COUNTRY")
		ADMIN_PHONE := envar.EnvarStr("", "ADMIN_PHONE")
		ADMIN_NAME := envar.EnvarStr("", "ADMIN_NAME")
		ADMIN_EMAIL := envar.EnvarStr("", "ADMIN_EMAIL")
		_, err := InitAdmin(ADMIN_NAME, ADMIN_COUNTRY, ADMIN_PHONE, ADMIN_EMAIL)
		if err != nil {
			return err
		}

		// Default token
		defaultToken()

		if _, err := Folders.Upsert(et.Json{
			"_id": "-1",
		}).
			Where(Folders.Column("_id").Eq("-1")).
			CommandOne(); err != nil {
			return err
		}

		CheckProjectModule("-1", "-1", true)
		CheckRole("-1", "-1", "PROFILE.ADMIN", "USER.ADMIN", true)

	*/

	console.LogK("Module", "Init data module")

	return nil
}
