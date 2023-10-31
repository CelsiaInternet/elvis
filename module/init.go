package module

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/envar"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
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
	if err := DefineStacks(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineHistorys(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineTokens(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineCollection(); err != nil {
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
	// Initial types
	InitType("-1", OF_SYSTEM, OF_SYSTEM, "STATE", "System", "Record system")
	InitType("-1", FOR_DELETE, OF_SYSTEM, "STATE", "Delete", "To delete record")
	InitType("-1", ACTIVE, OF_SYSTEM, "STATE", "Activo", "")
	InitType("-1", ARCHIVED, OF_SYSTEM, "STATE", "Archivado", "")
	InitType("-1", CANCELLED, OF_SYSTEM, "STATE", "Cacnelado", "")
	InitType("-1", IN_PROCESS, OF_SYSTEM, "STATE", "En tramite", "")
	InitType("-1", PENDING_APPROVAL, OF_SYSTEM, "STATE", "Pendiente de aprobación", "")
	InitType("-1", APPROVAL, OF_SYSTEM, "STATE", "Aprobado", "")
	InitType("-1", REFUSED, OF_SYSTEM, "STATE", "Rechazado", "")

	InitType("-1", "PROFILE.ADMIN", OF_SYSTEM, "PROFILE", "Admin", "")
	InitType("-1", "PROFILE.DEV", OF_SYSTEM, "PROFILE", "Develop", "")
	InitType("-1", "PROFILE.SUPORT", OF_SYSTEM, "PROFILE", "Suport", "")

	// Initial project
	app := EnvarStr("", "APP")
	_, err := InitProject("-1", app, "Initial project", Json{})
	if err != nil {
		return err
	}

	// User Admin
	ADMIN_COUNTRY := EnvarStr("", "ADMIN_COUNTRY")
	ADMIN_PHONE := EnvarStr("", "ADMIN_PHONE")
	ADMIN_NAME := EnvarStr("", "ADMIN_NAME")
	ADMIN_EMAIL := EnvarStr("", "ADMIN_EMAIL")

	_, err = InitAdmin(ADMIN_NAME, ADMIN_COUNTRY, ADMIN_PHONE, ADMIN_EMAIL)
	if err != nil {
		return err
	}

	// Admin module
	const moduleId = "MODULE.ADMIN"
	_, err = InitModule(moduleId, "Admin modules", "Admistración de modulos", Json{})
	if err != nil {
		return err
	}

	// Folders
	_, err = InitFolder(moduleId, "-1", "FOLDER.MODULES", "Modulos", "", Json{
		"icon":  "folder",
		"view":  "list",
		"clase": "module",
		"title": "Modulo",
		"url":   "module/all?state=0&search={search}&page={page}&rows={rows}",
		"state": "0",
		"filter": []Json{
			{"name": "Nombre"},
			{"description": "Descripción"},
		},
		"states": []Json{
			{"_id": "0", "name": "Active"},
		},
		"detail": Json{
			"title":    []string{"$1", "name"},
			"subtitle": []string{"$1", "description"},
			"datetime": []string{"$1", "date_update"},
			"code":     []string{"$1", "index"},
			"new_code": "",
			"state_color": Json{
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

	_, err = InitFolder(moduleId, "-1", "FOLDER.USERS", "Usuarios", "", Json{
		"icon":   "users",
		"view":   "users",
		"clase":  "user",
		"title":  "Usuario",
		"url":    "user/all?state=0&search={search}&page={page}&rows={rows}",
		"state":  "",
		"filter": []Json{},
		"states": []Json{},
		"detail": Json{
			"title":    []string{"$1", "full_name"},
			"subtitle": []string{"$1", "phone"},
			"datetime": []string{"$1", "date_update"},
			"code":     []string{"$1", "full_name"},
			"new_code": "Nuevo",
			"state_color": Json{
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

	_, err = InitFolder(moduleId, "-1", "FOLDER.APIKEYS", "API keys", "", Json{
		"icon":   "key",
		"view":   "apiKeys",
		"clase":  "apiKey",
		"title":  "API keys",
		"url":    "",
		"state":  "",
		"filter": []Json{},
		"states": []Json{},
		"detail": Json{
			"title":    []string{},
			"subtitle": []string{},
			"datetime": []string{},
			"code":     []string{},
			"new_code": "",
			"state_color": Json{
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
