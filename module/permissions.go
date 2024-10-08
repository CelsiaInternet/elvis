package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Permissions *linq.Model

var PERMISION_READ = "PERMISION.READ"
var PERMISION_WRITE = "PERMISION.WRITE"
var PERMISION_DELETE = "PERMISION.DELETE"
var PERMISION_EXECUTE = "PERMISION.EXECUTE"

func DefinePermisions(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return console.Panic(err)
	}

	if Permissions != nil {
		return nil
	}

	Permissions = linq.NewModel(SchemaModule, "Permissions", "Tabla de permisos", 1)
	Permissions.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Permissions.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	Permissions.DefineColum("model", "", "VARCHAR(80)", "")
	Permissions.DefineColum("profile_tp", "", "VARCHAR(80)", "-1")
	Permissions.DefineColum("permission_tp", "", "VARCHAR(80)", "-1")
	Permissions.DefineColum("index", "", "INTEGER", 0)
	Permissions.DefinePrimaryKey([]string{"project_id", "model", "profile_tp", "permission_tp"})
	Permissions.DefineIndex([]string{
		"date_make",
		"project_id",
		"model",
		"profile_tp",
		"permission_tp",
		"index",
	})

	if err := Permissions.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* GetPermission
* @param projectId string
* @param model string
* @param profileTp string
* @param permissionTp string
* @return et.Item, error
**/
func GetPermission(projectId, model, profileTp, permissionTp string) (et.Item, error) {
	return Permissions.Data().
		Where(Permissions.Column("project_id").Eq(projectId)).
		And(Permissions.Column("model").Eq(model)).
		And(Permissions.Column("profile_tp").Eq(profileTp)).
		And(Permissions.Column("permission_tp").Eq(permissionTp)).
		First()
}

/**
* GetPermissions
* @param projectId string
* @param model string
* @param profileTp string
* @return et.Items, error
**/
func GetPermissions(projectId, model, profileTp string) (map[string]bool, error) {
	result, err := Permissions.Data().
		Where(Permissions.Column("project_id").Eq(projectId)).
		And(Permissions.Column("profile_tp").Eq(profileTp)).
		All()
	if err != nil {
		return nil, err
	}

	if !result.Ok {
		return map[string]bool{
			PERMISION_READ:    true,
			PERMISION_WRITE:   true,
			PERMISION_DELETE:  true,
			PERMISION_EXECUTE: true,
		}, nil
	}

	permissions := map[string]bool{}
	for _, item := range result.Result {
		permission_tp := item.ValStr("", "permission_tp")
		permissions[permission_tp] = true
	}

	return permissions, nil
}

/**
* CheckPermission
* @param projectId string
* @param model string
* @param profileTp string
* @param permissionTp string
* @return et.Item, error
**/
func CheckPermission(projectId, model, profileTp, permissionTp string, chk bool) (et.Item, error) {
	if !utility.ValidId(projectId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "projectId")
	}

	if !utility.ValidStr(model, 0, []string{"", "*"}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "model")
	}

	if !utility.ValidId(profileTp) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "profileTp")
	}

	if !utility.ValidId(permissionTp) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "permissionTp")
	}

	if !chk {
		result, err := Permissions.Delete().
			Where(Permissions.Column("project_id").Eq(projectId)).
			And(Permissions.Column("model").Eq(model)).
			And(Permissions.Column("profile_tp").Eq(profileTp)).
			And(Permissions.Column("permission_tp").Eq(permissionTp)).
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

	current, err := GetPermission(projectId, model, profileTp, permissionTp)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		data := et.Json{}
		data.Set("project_id", projectId)
		data.Set("model", model)
		data.Set("profile_tp", profileTp)
		data.Set("permission_tp", permissionTp)
		result, err := Permissions.Insert(data).
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
