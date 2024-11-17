package module

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

var Users *linq.Model

func DefineUsers(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return logs.Panice(err)
	}

	if Users != nil {
		return nil
	}

	Users = linq.NewModel(SchemaModule, "USERS", "Tabla de usuarios", 1)
	Users.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Users.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Users.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Users.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Users.DefineColum("username", "", "VARCHAR(250)", "")
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
		"username",
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
		return logs.Panice(err)
	}

	return nil
}

/**
* GetUserByUserName
* @param username string
* @return et.Item
* @return error
**/
func GetUserByUserName(username string) (et.Item, error) {
	item, err := Users.Data().
		Where(Users.Column("username").Eq(username)).
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
* @param username string
* @param fullName string
* @param country string
* @param phone string
* @param email string
* @param password string
* @return et.Item
* @return error
**/
func InsertUser(id, username, fullName, country, phone, email, password string) (et.Item, error) {
	if !utility.ValidStr(username, 0, []string{""}) {
		return et.Item{}, logs.NewErrorf(msg.MSG_ATRIB_REQUIRED, "username")
	}

	if !utility.ValidStr(fullName, 0, []string{""}) {
		return et.Item{}, logs.NewErrorf(msg.MSG_ATRIB_REQUIRED, "fullName")
	}

	current, err := GetUserByUserName(username)
	if err != nil {
		return et.Item{}, err
	}

	if current.Ok {
		return current, logs.NewErrorf(msg.RECORD_FOUND)
	}

	password = utility.PasswordSha256(password)
	id = utility.GenKey(id)
	data := et.Json{}
	data["_id"] = id
	data["_state"] = utility.ACTIVE
	data["username"] = username
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
func UpdateUser(id, fullName, country, phone, email string) (et.Item, error) {
	if !utility.ValidStr(fullName, 3, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "full_name")
	}

	current, err := GetUserById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, logs.Alertm(msg.RECORD_NOT_FOUND)
	}

	if current.State() != utility.ACTIVE {
		return et.Item{}, logs.Alertf(msg.RECORD_NOT_ACTIVE, current.State())
	}

	now := utility.Now()
	data := et.Json{}
	data["date_update"] = now
	data["_id"] = id
	data["_state"] = utility.ACTIVE
	data["full_name"] = fullName
	data["country"] = country
	data["phone"] = phone
	data["email"] = email
	_, err = Users.Insert(data).
		Where(Users.Column("_id").Eq(id)).
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
* UpdatePassword
* @param id string
* @param oldPassword string
* @param newPassword string
* @param confirmPassword string
* @return et.Item
* @return error
**/
func UpdatePassword(id string, oldPassword, newPassword, confirmPassword string) (et.Item, error) {
	if !utility.ValidStr(oldPassword, 6, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "oldPassword")
	}

	if !utility.ValidStr(newPassword, 6, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "newPassword")
	}

	if !utility.ValidStr(confirmPassword, 6, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "confirmPassword")
	}

	if newPassword != confirmPassword {
		return et.Item{}, logs.Alertm(msg.PASSWORD_NOT_MATCH)
	}

	oldPassword = utility.PasswordSha256(oldPassword)
	current, err := Users.Data("_state").
		Where(Users.Column("_id").Eq(id)).
		And(Users.Column("_state").Eq(utility.ACTIVE)).
		And(Users.Column("password").Eq(oldPassword)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, logs.Alertm(msg.RECORD_NOT_FOUND)
	}

	password := utility.PasswordSha256(newPassword)
	now := utility.Now()
	data := et.Json{}
	data["date_update"] = now
	data["_id"] = id
	data["password"] = password
	_, err = Users.Insert(data).
		Where(Users.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return et.Item{
		Ok: true,
		Result: et.Json{
			"message": msg.RECORD_UPDATE,
		},
	}, nil
}

/**
* SetPassword
* @param id string
* @param newPassword string
* @param confirmPassword string
* @return et.Item
* @return error
**/
func SetPassword(id string, newPassword, confirmPassword string) (et.Item, error) {
	if !utility.ValidStr(newPassword, 6, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "newPassword")
	}

	if !utility.ValidStr(confirmPassword, 6, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "confirmPassword")
	}

	if newPassword != confirmPassword {
		return et.Item{}, logs.Alertm(msg.PASSWORD_NOT_MATCH)
	}

	current, err := Users.Data("_state").
		Where(Users.Column("_id").Eq(id)).
		And(Users.Column("_state").Eq(utility.ACTIVE)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, logs.Alertm(msg.RECORD_NOT_FOUND)
	}

	password := utility.PasswordSha256(newPassword)
	now := utility.Now()
	data := et.Json{}
	data["date_update"] = now
	data["_id"] = id
	data["password"] = password
	_, err = Users.Insert(data).
		Where(Users.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return et.Item{
		Ok: true,
		Result: et.Json{
			"message": msg.RECORD_UPDATE,
		},
	}, nil
}

/**
* StateUser
* @param id string
* @param state string
* @return et.Item
* @return error
**/
func StateUser(id, state string) (et.Item, error) {
	if !utility.ValidStr(state, 0, []string{""}) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := GetUserById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, logs.Alertm(msg.RECORD_NOT_FOUND)
	}

	if current.State() == utility.OF_SYSTEM {
		return et.Item{}, logs.Alertm(msg.RECORD_IS_SYSTEM)
	} else if current.State() == state {
		return et.Item{}, logs.Alertm(msg.RECORD_NOT_CHANGE)
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
	current, err := GetUserById(id)
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, logs.Alertm(msg.RECORD_NOT_FOUND)
	}

	_, err = Roles.Delete().
		Where(Roles.Column("user_id").Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	_, err = Tokens.Delete().
		Where(Tokens.Column("user_id").Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	_, err = Users.Delete().
		Where(Users.Column("_id").Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return et.Item{
		Ok: true,
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
			Where(Users.Concat("USERNAME:", Users.Column("username"), ":DATA:", Users.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Users.Column("username"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Users.Data(_select).
			Where(Users.Column("_state").Neg(state)).
			OrderBy(Users.Column("username"), true).
			List(page, rows)
	} else {
		return Users.Data(_select).
			Where(Users.Column("_state").Eq(state)).
			OrderBy(Users.Column("username"), true).
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
