package module

import (
	"github.com/cgalvisleon/elvis/aws"
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/envar"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/linq"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
	_ "github.com/joho/godotenv/autoload"
)

var Users *Model

func DefineUsers() error {
	if err := defineSchema(); err != nil {
		return console.PanicE(err)
	}

	if Users != nil {
		return nil
	}

	Users = NewModel(SchemaModule, "USERS", "Tabla de usuarios", 1)
	Users.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Users.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Users.DefineColum("_state", "", "VARCHAR(80)", ACTIVE)
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
	Users.Details("last_use", "", "", func(col *Column, data *Json) {
		id := data.Id()
		collection, err := GetCollectionById("telemetry.token.last_use", id)
		if err != nil {
			return
		}

		data.Set(col.Low(), collection.Str("last_use"))
	})
	Users.Details("projects", "", []Json{}, func(col *Column, data *Json) {
		id := data.Id()
		projects, err := GetUserProjects(id)
		if err != nil {
			return
		}

		data.Set(col.Low(), projects)
	})
	Users.Details("modules", "", []Json{}, func(col *Column, data *Json) {
		id := data.Id()
		modules, err := GetUserModules(id)
		if err != nil {
			return
		}

		data.Set(col.Low(), modules)
	})
	Users.Trigger(AfterInsert, func(model *Model, old, new *Json, data Json) {
		id := new.Key("_id")
		if id == "USER.ADMIN" {
			fullName := new.Str("full_name")
			country := new.Str("country")
			phone := new.Str("phone")
			APP := EnvarStr("", "APP")
			message := Format(MSG_ADMIN_WELCOME, fullName, APP)
			go aws.SendSMS(country, phone, message)
		}
	})

	if err := InitModel(Users); err != nil {
		return console.PanicE(err)
	}

	return nil
}

/**
* User
*	Handler for CRUD data
 */
func GetUserByName(name string) (Item, error) {
	item, err := Users.Select().
		Where(Users.Column("name").Eq(name)).
		First()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func GetUserById(id string) (Item, error) {
	item, err := Users.Select().
		Where(Users.Column("_id").Eq(id)).
		First()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func InitAdmin(fullName, country, phone, email string) (Item, error) {
	if !ValidStr(country, 0, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "country")
	}

	if !ValidStr(phone, 9, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "phone")
	}

	if !ValidStr(fullName, 0, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "full_name")
	}

	id := "USER.ADMIN"
	name := country + phone
	data := Json{}
	data["_id"] = id
	data["name"] = name
	data["full_name"] = fullName
	data["country"] = country
	data["phone"] = phone
	data["email"] = email
	data["avatar"] = ""
	item, err := Users.Insert(data).
		Where(Users.Column("_id").Eq(id)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func UpSetAdmin(fullName, country, phone, email string) (Item, error) {
	if !ValidStr(country, 0, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "country")
	}

	if !ValidStr(phone, 9, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "phone")
	}

	if !ValidStr(fullName, 0, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "full_name")
	}

	id := "USER.ADMIN"
	name := country + phone
	data := Json{}
	data["_id"] = id
	data["name"] = name
	data["full_name"] = fullName
	data["country"] = country
	data["phone"] = phone
	data["email"] = email
	data["avatar"] = ""
	item, err := Users.Upsert(data).
		Where(Users.Column("_id").Eq(id)).
		Command()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func SetUser(name, password, fullName, phone, email string) (Item, error) {
	if !ValidStr(name, 0, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "name")
	}

	if !ValidStr(phone, 3, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "phone")
	}

	if !ValidStr(fullName, 3, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "full_name")
	}

	current, err := GetUserByName(name)
	if err != nil {
		return Item{}, err
	}

	if current.Ok {
		return Item{}, console.ErrorM(RECORD_FOUND)
	}

	id := NewId()
	data := Json{}
	data["_id"] = id
	data["full_name"] = fullName
	data["phone"] = phone
	data["email"] = email
	data["avatar"] = ""
	item, err := Users.Insert(data).
		Where(Users.Column("name").Eq(name)).
		Command()
	if err != nil {
		return Item{}, err
	}

	item, err = GetProfile(id)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func UpdateUser(id, fullName, phone, email string, data Json) (Item, error) {
	if !ValidStr(fullName, 3, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "full_name")
	}

	if !ValidStr(phone, 3, []string{""}) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "phone")
	}

	current, err := GetUserById(id)
	if err != nil {
		return Item{}, err
	}

	if !current.Ok {
		return Item{}, console.ErrorM(RECORD_NOT_FOUND)
	}

	data["_id"] = id
	data["full_name"] = fullName
	data["phone"] = phone
	data["email"] = email
	data["avatar"] = ""
	item, err := Users.Insert(data).
		Where(Users.Column("_id").Eq(id)).
		And(Users.Column("_state").Eq(ACTIVE)).
		Command()
	if err != nil {
		return Item{}, err
	}

	item, err = GetProfile(id)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func StateUser(id, state string) (Item, error) {
	if !ValidId(state) {
		return Item{}, console.ErrorF(MSG_ATRIB_REQUIRED, "state")
	}

	return Users.Upsert(Json{
		"_state": state,
	}).
		Where(Users.Column("_id").Eq(id)).
		And(Users.Column("_state").Neg(state)).
		Command()
}

func DeleteUser(id string) (Item, error) {
	return StateUser(id, FOR_DELETE)
}

func AllUsers(state, search string, page, rows int, _select string) (List, error) {
	if state == "" {
		state = ACTIVE
	}

	auxState := state

	cols := StrToColN(_select)

	if auxState == "*" {
		state = FOR_DELETE

		return Users.Select(cols).
			Where(Users.Column("_state").Neg(state)).
			And(Users.Concat("NAME:", Users.Column("name"), ":DATA:", Users.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Users.Column("name"), true).
			List(page, rows)
	} else {
		return Users.Select(cols).
			Where(Users.Column("_state").Eq(state)).
			And(Users.Concat("NAME:", Users.Column("name"), ":DATA:", Users.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Users.Column("name"), true).
			List(page, rows)
	}
}

func GetProfile(userId string) (Item, error) {
	item, err := Users.Select().
		Where(Users.Column("_id").Eq(userId)).
		First()
	if err != nil {
		return Item{}, err
	}

	return item, nil
}
