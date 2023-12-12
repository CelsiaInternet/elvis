package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/envar"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
)

var initModules bool

func InitModules() error {
	if initModules {
		return nil
	}

	if err := DefineTypes(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineProjects(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineUsers(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineHistorys(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineTokens(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineModules(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineProjectModules(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineFolders(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineProfiles(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineProfileFolders(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineRoles(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineHistorys(); err != nil {
		return console.PanicE(err)
	}
	if err := defineModule(); err != nil {
		return console.PanicE(err)
	}

	initModules = true

	return nil
}

func defineModule() error {
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

	// Initial project
	app := envar.EnvarStr("", "APP")
	_, err := InitProject("-1", app, "Initial project", e.Json{})
	if err != nil {
		return err
	}

	// User Admin
	ADMIN_COUNTRY := envar.EnvarStr("", "ADMIN_COUNTRY")
	ADMIN_PHONE := envar.EnvarStr("", "ADMIN_PHONE")
	ADMIN_NAME := envar.EnvarStr("", "ADMIN_NAME")
	ADMIN_EMAIL := envar.EnvarStr("", "ADMIN_EMAIL")

	_, err = InitAdmin(ADMIN_NAME, ADMIN_COUNTRY, ADMIN_PHONE, ADMIN_EMAIL)
	if err != nil {
		return err
	}

	// Admin module
	const moduleId = "MODULE.ADMIN"
	_, err = InitModule(moduleId, "Admin modules", "Admistración de modulos", e.Json{})
	if err != nil {
		return err
	}

	// Folders
	_, err = InitFolder(moduleId, "-1", "FOLDER.MODULES", "Modulos", "", e.Json{
		"icon":  "folder",
		"view":  "list",
		"clase": "module",
		"title": "Modulo",
		"url":   "module/all?state=0&search={search}&page={page}&rows={rows}",
		"state": "0",
		"filter": []e.Json{
			{"name": "Nombre"},
			{"description": "Descripción"},
		},
		"states": []e.Json{
			{"_id": "0", "name": "Active"},
		},
		"detail": e.Json{
			"title":    []string{"$1", "name"},
			"subtitle": []string{"$1", "description"},
			"datetime": []string{"$1", "date_update"},
			"code":     []string{"$1", "index"},
			"new_code": "",
			"state_color": e.Json{
				"field_name": "_state",
				"warning":    "",
				"alert":      "",
				"info":       "",
			},
		},
		"showNew":   true,
		"showPrint": false,
		"order":     1,
	})
	if err != nil {
		return err
	}

	_, err = InitFolder(moduleId, "-1", "FOLDER.USERS", "Usuarios", "", e.Json{
		"icon":   "users",
		"view":   "users",
		"clase":  "user",
		"title":  "Usuario",
		"url":    "user/all?state=0&search={search}&page={page}&rows={rows}",
		"state":  "",
		"filter": []e.Json{},
		"states": []e.Json{},
		"detail": e.Json{
			"title":    []string{"$1", "full_name"},
			"subtitle": []string{"$1", "phone"},
			"datetime": []string{"$1", "date_update"},
			"code":     []string{"$1", "full_name"},
			"new_code": "Nuevo",
			"state_color": e.Json{
				"field_name": "_state",
				"warning":    "",
				"alert":      "",
				"info":       "",
			},
			"email":  []string{"$1", "email"},
			"avatar": []string{"$1", "avatar"},
		},
		"showNew":   true,
		"showPrint": false,
		"order":     2,
	})
	if err != nil {
		return err
	}

	_, err = InitFolder(moduleId, "-1", "FOLDER.APIKEYS", "API keys", "", e.Json{
		"icon":   "key",
		"view":   "apiKeys",
		"clase":  "apiKey",
		"title":  "API keys",
		"url":    "",
		"state":  "",
		"filter": []e.Json{},
		"states": []e.Json{},
		"detail": e.Json{
			"title":    []string{},
			"subtitle": []string{},
			"datetime": []string{},
			"code":     []string{},
			"new_code": "",
			"state_color": e.Json{
				"field_name": "_state",
				"warning":    "",
				"alert":      "",
				"info":       "",
			},
			"email":  []string{"$1", "email"},
			"avatar": []string{"$1", "avatar"},
		},
		"showNew":   true,
		"showPrint": false,
		"order":     3,
	})
	if err != nil {
		return err
	}

	return nil
}
