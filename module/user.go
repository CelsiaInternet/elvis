package module

import (
	"github.com/cgalvisleon/elvis/aws"
	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/core"
	"github.com/cgalvisleon/elvis/envar"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
	_ "github.com/joho/godotenv/autoload"
)

var Users *linq.Model

func DefineUsers() error {
	if err := DefineSchemaModule(); err != nil {
		return console.PanicE(err)
	}

	if Users != nil {
		return nil
	}

	Users = linq.NewModel(SchemaModule, "USERS", "Tabla de usuarios", 1)
	Users.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Users.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Users.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Users.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Users.DefineColum("name", "", "VARCHAR(250)", "")
	Users.DefineColum("password", "", "VARCHAR(250)", "")
	Users.DefineColum("_data", "", "JSONB", "{}")
	Users.DefineColum("index", "", "INTEGER", 0)
	Users.DefineAtrib("full_name", "", "text", "")
	Users.DefineAtrib("country", "", "text", "")
	Users.DefineAtrib("phone", "", "text", "")
	Users.DefineAtrib("email", "", "text", "")
	Users.DefineAtrib("avatar", "", "text", "")
	Users.DefinePrimaryKey([]string{"_id"})
	Users.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"name",
		"index",
	})
	Users.DefineHidden([]string{"password"})
	Users.Details("last_use", "", "", func(col *linq.Column, data *e.Json) {
		id := data.Id()
		last_use, err := cache.HGetAtrib(id, "telemetry.token.last_use")
		if err != nil {
			return
		}

		data.Set(col.Low(), last_use)
	})
	Users.Details("projects", "", []e.Json{}, func(col *linq.Column, data *e.Json) {
		id := data.Id()
		projects, err := GetUserProjects(id)
		if err != nil {
			return
		}

		data.Set(col.Low(), projects)
	})
	Users.Details("modules", "", []e.Json{}, func(col *linq.Column, data *e.Json) {
		id := data.Id()
		modules, err := GetUserModules(id)
		if err != nil {
			return
		}

		data.Set(col.Low(), modules)
	})
	Users.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *e.Json, data e.Json) error {
		id := new.Key("_id")
		if id == "USER.ADMIN" {
			fullName := new.Str("full_name")
			country := new.Str("country")
			phone := new.Str("phone")
			APP := envar.EnvarStr("", "APP")
			message := strs.Format(msg.MSG_ADMIN_WELCOME, fullName, APP)
			go aws.SendSMS(country, phone, message)
		}

		return nil
	})

	if err := core.InitModel(Users); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* User
*	Handler for CRUD data
 */
func GetUserByName(name string) (e.Item, error) {
	item, err := Users.Select().
		Where(Users.Column("name").Eq(name)).
		First()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func GetUserById(id string) (e.Item, error) {
	item, err := Users.Select().
		Where(Users.Column("_id").Eq(id)).
		First()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func InitAdmin(fullName, country, phone, email string) (e.Item, error) {
	if !utility.ValidStr(country, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "country")
	}

	if !utility.ValidStr(phone, 9, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "phone")
	}

	if !utility.ValidStr(fullName, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "full_name")
	}

	id := "USER.ADMIN"
	current, err := GetUserById(id)
	if err != nil {
		return e.Item{}, err
	}

	if current.Ok {
		return current, nil
	}

	name := country + phone
	data := e.Json{}
	data["_id"] = id
	data["name"] = name
	data["full_name"] = fullName
	data["country"] = country
	data["phone"] = phone
	data["email"] = email
	data["avatar"] = ""
	item, err := Users.Insert(data).
		Where(Users.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func UpSetAdmin(fullName, country, phone, email string) (e.Item, error) {
	if !utility.ValidStr(country, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "country")
	}

	if !utility.ValidStr(phone, 9, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "phone")
	}

	if !utility.ValidStr(fullName, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "full_name")
	}

	id := "USER.ADMIN"
	name := country + phone
	data := e.Json{}
	data["_id"] = id
	data["name"] = name
	data["full_name"] = fullName
	data["country"] = country
	data["phone"] = phone
	data["email"] = email
	data["avatar"] = ""
	item, err := Users.Upsert(data).
		Where(Users.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func SetUser(name, password, fullName, phone, email string) (e.Item, error) {
	if !utility.ValidStr(name, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "name")
	}

	if !utility.ValidStr(phone, 3, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "phone")
	}

	if !utility.ValidStr(fullName, 3, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "full_name")
	}

	current, err := GetUserByName(name)
	if err != nil {
		return e.Item{}, err
	}

	if current.Ok {
		return e.Item{}, console.Alert(msg.RECORD_FOUND)
	}

	id := utility.NewId()
	data := e.Json{}
	data["_id"] = id
	data["full_name"] = fullName
	data["phone"] = phone
	data["email"] = email
	data["avatar"] = ""
	_, err = Users.Insert(data).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	item, err := GetProfile(id)
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func UpdateUser(id, fullName, phone, email string, data e.Json) (e.Item, error) {
	if !utility.ValidStr(fullName, 3, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "full_name")
	}

	if !utility.ValidStr(phone, 3, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "phone")
	}

	current, err := GetUserById(id)
	if err != nil {
		return e.Item{}, err
	}

	if !current.Ok {
		return e.Item{}, console.ErrorM(msg.RECORD_NOT_FOUND)
	}

	data["_id"] = id
	data["full_name"] = fullName
	data["phone"] = phone
	data["email"] = email
	data["avatar"] = ""
	_, err = Users.Insert(data).
		Where(Users.Column("_id").Eq(id)).
		And(Users.Column("_state").Eq(utility.ACTIVE)).
		CommandOne()
	if err != nil {
		return e.Item{}, err
	}

	item, err := GetProfile(id)
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func StateUser(id, state string) (e.Item, error) {
	if !utility.ValidId(state) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	return Users.Update(e.Json{
		"_state": state,
	}).
		Where(Users.Column("_id").Eq(id)).
		And(Users.Column("_state").Neg(state)).
		CommandOne()
}

func DeleteUser(id string) (e.Item, error) {
	return StateUser(id, utility.FOR_DELETE)
}

func AllUsers(state, search string, page, rows int, _select string) (e.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Users.Select(_select).
			Where(Users.Concat("NAME:", Users.Column("name"), ":DATA:", Users.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Users.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Users.Select(_select).
			Where(Users.Column("_state").Neg(state)).
			OrderBy(Users.Column("name"), true).
			List(page, rows)
	} else {
		return Users.Select(_select).
			Where(Users.Column("_state").Eq(state)).
			OrderBy(Users.Column("name"), true).
			List(page, rows)
	}
}

func GetProfile(userId string) (e.Item, error) {
	item, err := Users.Select().
		Where(Users.Column("_id").Eq(userId)).
		First()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}
