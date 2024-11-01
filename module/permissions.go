package module

import (
	"encoding/json"
	"net/http"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/strs"
)

type Permission map[string]bool

/**
* ToString
* @return string, error
**/
func (p Permission) ToString() (string, error) {
	jsonString, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}

func (p Permission) Method(r *http.Request) bool {
	method := r.Method
	switch method {
	case "GET":
		return p[PERMISION_READ]
	case "POST":
		return p[PERMISION_WRITE]
	case "PUT":
		return p[PERMISION_UPDATE]
	case "DELETE":
		return p[PERMISION_DELETE]
	case "PATCH":
		return p[PERMISION_EXECUTE]
	default:
		return false
	}
}

/**
* NewPermision
* @param data string
* @return Permission, error
**/
func NewPermision(data string) (Permission, error) {
	var result = make(Permission)
	if data == "" {
		return result, nil
	}

	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

var Permissions *linq.Model

var PERMISION_READ = "PERMISION.READ"
var PERMISION_WRITE = "PERMISION.WRITE"
var PERMISION_DELETE = "PERMISION.DELETE"
var PERMISION_UPDATE = "PERMISION.UPDATE"
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
	Permissions.DefineColum("profile_tp", "", "VARCHAR(80)", "-1")
	Permissions.DefineColum("model", "", "VARCHAR(80)", "")
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
* ResetPermissions
* @param projectId string
* @param profileTp string
* @param model string
* @return Permission, error
**/
func ResetPermissions(projectId, profileTp, model string) (Permission, error) {
	var result = make(Permission)
	items, err := Permissions.Select().
		Where(Permissions.Column("project_id").Eq(projectId)).
		And(Permissions.Column("model").Eq(model)).
		And(Permissions.Column("profile_tp").Eq(profileTp)).
		OrderBy(Permissions.Col("index"), true).
		All()
	if err != nil {
		return result, err
	}

	if items.Ok {
		for _, item := range items.Result {
			permision := item.ValStr("-1", "permission_tp")
			if permision != "-1" {
				result[permision] = true
			}
		}
	}

	var key = strs.Format("%v-%v-%v", projectId, profileTp, model)
	value, _ := result.ToString()
	cache.SetM(key, value)

	return result, nil
}

/**
* GetPermissions
* @param projectId string
* @param profileTp string
* @param model string
* @return et.Item, error
**/
func GetPermissions(projectId, profileTp, model string) (Permission, error) {
	var key = strs.Format("%v-%v-%v", projectId, profileTp, model)
	value, err := cache.Get(key, "")
	if err != nil {
		return Permission{}, err
	}

	if value != "" {
		result, err := NewPermision(value)
		if err == nil {
			return result, nil
		}
	}

	var result = make(Permission)
	items, err := Permissions.Select().
		Where(Permissions.Column("project_id").Eq(projectId)).
		And(Permissions.Column("model").Eq(model)).
		And(Permissions.Column("profile_tp").Eq(profileTp)).
		OrderBy(Permissions.Col("index"), true).
		All()
	if err != nil {
		return result, err
	}

	if items.Ok {
		for _, item := range items.Result {
			permision := item.ValStr("-1", "permission_tp")
			if permision != "-1" {
				result[permision] = true
			}
		}

		value, _ := result.ToString()
		cache.SetM(key, value)

		return result, nil
	}

	result, err = GetPermissions("-1", profileTp, model)
	if err != nil {
		return result, err
	}

	return result, nil
}

/**
* CheckPermissions
* @param projectId string
* @param profileTp string
* @param model string
* @param permissionTp string
* @param chk bool
* @return error
**/
func CheckPermissions(projectId, profileTp, model, permissionTp string, chk bool) error {
	if chk {
		current, err := Permissions.Select().
			Where(Permissions.Column("project_id").Eq(projectId)).
			And(Permissions.Column("model").Eq(model)).
			And(Permissions.Column("profile_tp").Eq(profileTp)).
			And(Permissions.Column("permission_tp").Eq(permissionTp)).
			First()
		if err != nil {
			return err
		}

		if !current.Ok {
			data := et.Json{
				"project_id":    projectId,
				"profile_tp":    profileTp,
				"model":         model,
				"permission_tp": permissionTp,
			}

			_, err := Permissions.Insert(data).
				CommandOne()
			if err != nil {
				return err
			}

			_, err = ResetPermissions(projectId, profileTp, model)
			if err != nil {
				return err
			}
		}

		return nil
	}

	_, err := Permissions.Delete().
		Where(Permissions.Column("project_id").Eq(projectId)).
		And(Permissions.Column("model").Eq(model)).
		And(Permissions.Column("profile_tp").Eq(profileTp)).
		And(Permissions.Column("permission_tp").Eq(permissionTp)).
		CommandOne()
	if err != nil {
		return err
	}

	_, err = ResetPermissions(projectId, profileTp, model)
	if err != nil {
		return err
	}

	return nil
}

/**
* PermissionsMiddleware
* @param next http.Handler
* @return http.Handler
**/
func PermissionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		project_id := claim.ProjectIdKey.String(ctx, "")
		profile_tp := claim.ProfileTpKey.String(ctx, "")
		model := claim.ModelKey.String(ctx, "")
		permisions, err := GetPermissions(project_id, profile_tp, model)
		if err != nil {
			response.InternalServerError(w, r)
			return
		}

		ok := permisions.Method(r)
		if !ok {
			response.Forbidden(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
