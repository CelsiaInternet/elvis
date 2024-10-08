package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Roles *linq.Model

func DefineRoles(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return console.Panic(err)
	}

	if Roles != nil {
		return nil
	}

	Roles = linq.NewModel(SchemaModule, "ROLES", "Tabla de roles", 1)
	Roles.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Roles.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Roles.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	Roles.DefineColum("module_id", "", "VARCHAR(80)", "-1")
	Roles.DefineColum("user_id", "", "VARCHAR(80)", "-1")
	Roles.DefineColum("profile_tp", "", "VARCHAR(80)", "-1")
	Roles.DefineColum("index", "", "INTEGER", 0)
	Roles.DefinePrimaryKey([]string{"project_id", "module_id", "user_id"})
	Roles.DefineIndex([]string{
		"date_make",
		"date_update",
		"profile_tp",
		"index",
	})

	if err := Roles.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* GetRoleById
* @param projectId string
* @param moduleId string
* @param userId string
* @return et.Item, error
**/
func GetRoleById(projectId, moduleId, userId string) (et.Item, error) {
	result, err := Roles.Data().
		Where(Roles.Column("project_id").Eq(projectId)).
		And(Roles.Column("module_id").Eq(moduleId)).
		And(Roles.Column("user_id").Eq(userId)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* GetUserRoleByIndex
* @param idx int64
* @return et.Item, error
**/
func GetUserRoleByIndex(idx int64) (et.Item, error) {
	sql := `
	SELECT
	D._ID AS PROJECT_ID,
	D.NAME AS PROJECT,
	B._ID AS MODULE_ID,
	B.NAME AS MODULE,
	A.PROFILE_TP,
	C.NAME PROFILE,
	A.USER_ID,
	A.INDEX
	FROM module.ROLES A
	INNER JOIN module.MODULES B ON B._ID=A.MODULE_ID
	INNER JOIN module.TYPES C ON C._ID=A.PROFILE_TP
	INNER JOIN module.PROJECTS D ON D._ID=A.PROJECT_ID
	WHERE A.INDEX=$1
	LIMIT 1;`

	item, err := Roles.QueryOne(sql, idx)
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* GetUserProjects
* @param userId string
* @return []et.Json, error
**/
func GetUserProjects(userId string) ([]et.Json, error) {
	sql := `
	SELECT
	B._ID,
	B.NAME,
	MIN(A.INDEX) AS INDEX
	FROM module.ROLES A	
	INNER JOIN module.PROJECTS B ON B._ID=A.PROJECT_ID
	WHERE A.USER_ID=$1
	GROUP BY B._ID, B.NAME
	ORDER BY B.NAME;`

	modules, err := Roles.Query(sql, userId)
	if err != nil {
		return []et.Json{}, err
	}

	return modules.Result, nil
}

/**
* GetUserModules
* @param userId string
* @return []et.Json, error
**/
func GetUserModules(userId string) ([]et.Json, error) {
	sql := `
	SELECT
	D._ID AS PROJECT_ID,
	D.NAME AS PROJECT,
	B._ID AS MODULE_ID,
	B.NAME AS MODULE,
	A.PROFILE_TP,
	C.NAME PROFILE,
	A.USER_ID,
	A.INDEX
	FROM module.ROLES A
	INNER JOIN module.MODULES B ON B._ID=A.MODULE_ID
	INNER JOIN module.TYPES C ON C._ID=A.PROFILE_TP
	INNER JOIN module.PROJECTS D ON D._ID=A.PROJECT_ID
	WHERE A.USER_ID=$1
	GROUP BY D._ID, D.NAME, B._ID, B.NAME, A.PROFILE_TP, C.NAME, USER_ID, A.INDEX
	ORDER BY D.NAME, B.NAME, C.NAME;`

	modules, err := Roles.Query(sql, userId)
	if err != nil {
		return []et.Json{}, err
	}

	return modules.Result, nil
}

/**
* CheckRole
* @param projectId string
* @param moduleId string
* @param profileTp string
* @param userId string
* @param chk bool
* @return et.Item, error
**/
func CheckRole(projectId, moduleId, profileTp, userId string, chk bool) (et.Item, error) {
	if !utility.ValidId(projectId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	if !utility.ValidId(moduleId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "module_id")
	}

	if !utility.ValidId(userId) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "user_id")
	}

	if !utility.ValidId(profileTp) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "profile_tp")
	}

	project, err := GetProjectById(projectId)
	if err != nil {
		return et.Item{}, err
	}

	if !project.Ok {
		return et.Item{}, console.AlertF(msg.PROJECT_NOT_FOUND, projectId)
	}

	module, err := GetModuleById(moduleId)
	if err != nil {
		return et.Item{}, err
	}

	if !module.Ok {
		return et.Item{}, console.Alert(msg.MODULE_NOT_FOUND)
	}

	profile, err := GetProfileById(moduleId, profileTp)
	if err != nil {
		return et.Item{}, err
	}

	if !profile.Ok {
		return et.Item{}, console.AlertF(msg.PROFILE_NOT_FOUND, profileTp)
	}

	if chk {
		current, err := GetRoleById(projectId, moduleId, userId)
		if err != nil {
			return et.Item{}, err
		}

		now := utility.Now()
		if current.Ok {
			index := current.Index64()
			item, err := Roles.Update(et.Json{
				"date_update": now,
				"profile_tp":  profileTp,
			}).Where(Roles.Column("index").Eq(index)).
				CommandOne()
			if err != nil {
				return et.Item{}, err
			}

			item, err = GetUserRoleByIndex(index)
			if err != nil {
				return et.Item{}, err
			}

			return et.Item{
				Ok: item.Ok,
				Result: et.OkOrNotJson(item.Ok, item.Result, et.Json{
					"message": msg.RECORD_NOT_UPDATE,
					"index":   index,
				}),
			}, nil
		}

		index := Roles.NextSerie("module.ROLES")
		item, err := Roles.Insert(et.Json{
			"date_make":   now,
			"date_update": now,
			"project_id":  projectId,
			"module_id":   moduleId,
			"user_id":     userId,
			"profile_tp":  profileTp,
			"index":       index,
		}).CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		item, err = GetUserRoleByIndex(index)
		if err != nil {
			return et.Item{}, err
		}

		return et.Item{
			Ok: item.Ok,
			Result: et.OkOrNotJson(item.Ok, item.Result, et.Json{
				"message": msg.RECORD_NOT_UPDATE,
				"index":   index,
			}),
		}, nil
	} else {
		sql := `
		DELETE FROM module.ROLES
		WHERE PROJECT_ID=$1
		AND MODULE_ID=$2
		AND PROFILE_TP=$3
		AND USER_ID=$4
		RETURNING INDEX;`

		item, err := Roles.Command(sql, projectId, moduleId, profileTp, userId)
		if err != nil {
			return et.Item{}, err
		}

		return et.Item{
			Ok: item.Ok,
			Result: et.Json{
				"message": utility.OkOrNot(item.Ok, msg.RECORD_DELETE, msg.RECORD_NOT_DELETE),
				"index":   item.Index(),
			},
		}, nil
	}
}
