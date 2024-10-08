package module

import (
	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

var Users *linq.Model

func DefineUsers(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return console.Panic(err)
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
	Users.Details("last_use", "", "", func(col *linq.Column, data *et.Json) {
		id := data.Id()
		last_use, err := cache.HGetAtrib(id, "telemetry.token.last_use")
		if err != nil {
			return
		}

		data.Set(col.Low(), last_use)
	})
	Users.Details("projects", "", []et.Json{}, func(col *linq.Column, data *et.Json) {
		id := data.Id()
		projects, err := GetUserProjects(id)
		if err != nil {
			return
		}

		data.Set(col.Low(), projects)
	})
	Users.Details("modules", "", []et.Json{}, func(col *linq.Column, data *et.Json) {
		id := data.Id()
		modules, err := GetUserModules(id)
		if err != nil {
			return
		}

		data.Set(col.Low(), modules)
	})

	if err := Users.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* GetUserByName
* @param name string
* @return et.Item
* @return error
**/
func GetUserByName(name string) (et.Item, error) {
	item, err := Users.Data().
		Where(Users.Column("name").Eq(name)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
*  GetUserByEmail
* @param email string
* @return et.Item
* @return error
**/
func GetUserByEmail(email string) (et.Item, error) {
	item, err := Users.Data().
		Where(Users.Column("email").Eq(email)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* GetUserById
* @param id string
* @return et.Item
* @return error
**/
func GetUserById(id string) (et.Item, error) {
	item, err := Users.Data().
		Where(Users.Column("_id").Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* InsertUser
* @param id string
* @param fullName string
* @param country string
* @param phone string
* @param email string
* @param password string
* @return et.Item
* @return error
**/
func InsertUser(id, fullName, country, phone, email, password string) (et.Item, error) {
	if !utility.ValidStr(country, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "country")
	}

	if !utility.ValidStr(phone, 9, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "phone")
	}

	if !utility.ValidStr(fullName, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "full_name")
	}

	name := country + phone
	current, err := GetUserByName(name)
	if err != nil {
		return et.Item{}, err
	}

	if current.Ok {
		return current, console.AlertF(msg.RECORD_FOUND)
	}

	password, err = utility.PasswordHash(password)
	if err != nil {
		return et.Item{}, err
	}

	id = utility.GenKey(id)
	data := et.Json{}
	data["_id"] = id
	data["_state"] = utility.ACTIVE
	data["name"] = name
	data["full_name"] = fullName
	data["country"] = country
	data["phone"] = phone
	data["email"] = email
	data["password"] = password
	data["avatar"] = ""
	_, err = Users.Insert(data).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	item, err := GetProfile(id)
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* UpdateUser
* @param id string
* @param data et.Json
* @return et.Item
* @return error
**/
func UpdateUser(id string, data et.Json) (et.Item, error) {
	fullName := data.ValStr("", "full_name")
	if !utility.ValidStr(fullName, 3, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "full_name")
	}

	current, err := GetUserById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, console.Alert(msg.RECORD_NOT_FOUND)
	}

	if current.State() == utility.OF_SYSTEM {
		return et.Item{}, console.Alert(msg.RECORD_IS_SYSTEM)
	} else if current.State() == utility.FOR_DELETE {
		return et.Item{}, console.Alert(msg.RECORD_DELETE)
	} else if current.State() != utility.ACTIVE {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_ACTIVE, current.State())
	}

	now := utility.Now()
	delete(data, "country")
	delete(data, "phone")
	delete(data, "name")
	delete(data, "email")
	delete(data, "password")
	data["date_update"] = now
	data["_id"] = id
	data["_state"] = utility.ACTIVE
	_, err = Users.Insert(data).
		Where(Users.Column("_id").Eq(id)).
		And(Users.Column("_state").Eq(utility.ACTIVE)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	item, err := GetProfile(id)
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* StateUser
* @param id string
* @param state string
* @return et.Item
* @return error
**/
func StateUser(id, state string) (et.Item, error) {
	if !utility.ValidId(id) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "id")
	}

	if !utility.ValidStr(state, 0, []string{""}) {
		return et.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := GetUserById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, console.Alert(msg.RECORD_NOT_FOUND)
	}

	if current.State() == utility.OF_SYSTEM {
		return et.Item{}, console.Alert(msg.RECORD_IS_SYSTEM)
	} else if current.State() == utility.FOR_DELETE {
		return et.Item{}, console.Alert(msg.RECORD_DELETE)
	} else if current.State() == state {
		return et.Item{}, console.Alert(msg.RECORD_NOT_CHANGE)
	}

	result, err := Users.Update(et.Json{
		"_state": state,
	}).
		Where(Users.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return et.Item{
		Ok: result.Ok,
		Result: et.Json{
			"message": msg.RECORD_UPDATE,
		},
	}, nil
}

/**
* DeleteUser
* @param id string
* @return et.Item
* @return error
**/
func DeleteUser(id string) (et.Item, error) {
	result, err := StateUser(id, utility.FOR_DELETE)
	if err != nil {
		return et.Item{}, err
	}

	return et.Item{
		Ok: result.Ok,
		Result: et.Json{
			"message": msg.RECORD_DELETE,
		},
	}, nil
}

/**
* AllUsers
* @param state string
* @param search string
* @param page int
* @param rows int
* @param _select string
* @return et.List
* @return error
**/
func AllUsers(state, search string, page, rows int, _select string) (et.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Users.Data(_select).
			Where(Users.Concat("NAME:", Users.Column("name"), ":DATA:", Users.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Users.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Users.Data(_select).
			Where(Users.Column("_state").Neg(state)).
			OrderBy(Users.Column("name"), true).
			List(page, rows)
	} else {
		return Users.Data(_select).
			Where(Users.Column("_state").Eq(state)).
			OrderBy(Users.Column("name"), true).
			List(page, rows)
	}
}

/**
* GetProfile
* @param userId string
* @return et.Item
* @return error
**/
func GetProfile(userId string) (et.Item, error) {
	item, err := Users.Data().
		Where(Users.Column("_id").Eq(userId)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}
