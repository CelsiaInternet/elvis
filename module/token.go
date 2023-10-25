package module

import (
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/claim"
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/linq"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
)

type Token struct {
	Date_make   time.Time `json:"date_make"`
	Date_update time.Time `json:"date_update"`
	Id          string    `json:"_id"`
	Name        string    `json:"name"`
	App         string    `json:"app"`
	Device      string    `json:"device"`
	Token       string    `json:"token"`
	Index       int       `json:"index"`
}

func (n *Token) Scan(js *Json) error {
	n.Date_make = js.Time("date_make")
	n.Date_update = js.Time("date_update")
	n.Id = js.Str("_id")
	n.Name = js.Str("name")
	n.App = js.Str("app")
	n.Device = js.Str("device")
	n.Token = js.Str("token")
	n.Index = js.Int("index")

	return nil
}

var Tokens *Model

func DefineTokens() error {
	if err := defineSchema(); err != nil {
		return console.PanicE(err)
	}

	if Tokens != nil {
		return nil
	}

	Tokens = NewModel(SchemaModule, "TOKENS", "Tabla de tokens", 1)
	Tokens.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Tokens.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Tokens.DefineColum("_id", "", "VARCHAR(80)", "-1")
	Tokens.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	Tokens.DefineColum("user_id", "", "VARCHAR(80)", "-1")
	Tokens.DefineColum("name", "", "VARCHAR(80)", "")
	Tokens.DefineColum("app", "", "VARCHAR(80)", "")
	Tokens.DefineColum("device", "", "VARCHAR(80)", "")
	Tokens.DefineColum("token", "", "TEXT", "")
	Tokens.DefineColum("index", "", "INTEGER", 0)
	Tokens.DefinePrimaryKey([]string{"_id"})
	Tokens.DefineIndex([]string{
		"date_make",
		"date_update",
		"name",
		"app",
		"device",
		"index",
	})

	if err := InitModel(Tokens); err != nil {
		return console.PanicE(err)
	}

	go LoadTokens()

	return nil
}

func loadToken(token *Token) error {
	key := claim.TokenKey(token.App, token.Device, token.Id)
	err := cache.Set(key, token.Token, 0)
	if err != nil {
		return err
	}

	return nil
}

func unLoadTokenById(app, device, id string) error {
	key := claim.TokenKey(app, device, id)
	_, err := cache.Del(key)
	if err != nil {
		return err
	}

	return nil
}

func getTokenByApp(app, userId string) (Item, error) {
	return Tokens.Select().
		Where(Tokens.Col("app").Eq(app)).
		And(Tokens.Col("user_id").Eq(userId)).
		First()
}

func GetTokenById(id string) (Item, error) {
	item, err := Tokens.Select().
		Where(Tokens.Col("_id").Eq(id)).
		First()
	if err != nil {
		return Item{}, err
	}

	if item.Ok {
		collection, err := GetCollectionById("telemetry.token.last_use", id)
		if err != nil {
			return Item{}, err
		}

		token := item.Result["token"].(string)
		item.Result["token"] = token[0:6] + "..." + token[len(token)-6:]
		item.Result["last_use"] = collection.Str("last_use")
	}

	return item, nil
}

func UpSetToken(projeectId, id, app, device, name, userId string) (Item, error) {
	user, err := GetUserById(userId)
	if err != nil {
		return Item{}, err
	}

	if !user.Ok {
		return Item{}, console.ErrorM(USER_NOT_FONUND)
	}

	id = GenId(id)
	current, err := GetTokenById(id)
	if err != nil {
		return Item{}, err
	}

	if current.Ok {
		id := current.Id()
		now := Now()

		sql := `
		UPDATE core.TOKENS SET
		DATE_UPDATE=$2,
		NAME=$3
		WHERE _ID=$1
		RETURNING *;`

		item, err := QueryOne(sql, id, now, name)
		if err != nil {
			return Item{}, err
		}

		if item.Ok {
			collection, err := GetCollectionById("telemetry.token.last_use", id)
			if err != nil {
				return Item{}, err
			}

			token := item.Result["token"].(string)
			item.Result["token"] = token[0:6] + "..." + token[len(token)-6:]
			item.Result["last_use"] = collection.Str("last_use")
		}

		return Item{
			Ok: item.Ok,
			Result: OkOrNotJson(item.Ok, item.Result, Json{
				"message": RECORD_NOT_UPDATE,
				"_id":     id,
			}),
		}, nil
	} else {
		id := NewId()

		token, err := claim.GenToken(id, app, name, "token", app, device, 0)
		if err != nil {
			return Item{}, console.Error(err)
		}

		sql := `
		INSERT INTO core.TOKENS(PROJECT_ID, _ID, USER_ID, APP, DEVICE, NAME, TOKEN)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING *;`

		item, err := QueryOne(sql, projeectId, id, userId, app, device, name, token)
		if err != nil {
			return Item{}, console.Error(err)
		}

		var result Token
		err = result.Scan(&item.Result)
		if err != nil {
			return Item{}, console.Error(err)
		}

		err = loadToken(&result)
		if err != nil {
			return Item{}, console.Error(err)
		}

		if item.Ok {
			token := item.Result["token"].(string)
			item.Result["long_token"] = token
			item.Result["token"] = token[0:6] + "..." + token[len(token)-6:]
		}

		return Item{
			Ok: item.Ok,
			Result: OkOrNotJson(item.Ok, item.Result, Json{
				"message": RECORD_NOT_CREATE,
				"_id":     id,
			}),
		}, nil
	}
}

func LoadTokens() error {
	var ok bool = true
	var rows int = 30
	var page int = 1
	for ok {
		ok = false

		offset := (page - 1) * rows
		sql := Format(`
		SELECT *
		FROM core.TOKENS
		ORDER BY INDEX
		LIMIT %d OFFSET %d;`, rows, offset)

		items, err := Query(sql)
		if err != nil {
			return console.Error(err)
		}

		for _, item := range items.Result {
			var result Token
			err = result.Scan(&item)
			if err != nil {
				return console.Error(err)
			}

			err = loadToken(&result)
			if err != nil {
				return console.Error(err)
			}

			ok = true
		}

		page++
	}

	return nil
}

func UnLoadTokens() error {
	var ok bool = true
	var rows int = 30
	var page int = 1
	for ok {
		ok = false

		offset := (page - 1) * rows
		sql := Format(`
		SELECT APP, DEVICE, _ID
		FROM core.TOKENS
		ORDER BY INDEX
		LIMIT %d OFFSET %d;`, rows, offset)

		items, err := Query(sql)
		if err != nil {
			return console.Error(err)
		}

		for _, item := range items.Result {
			app := item.Str("app")
			device := item.Str("device")
			id := item.Id()
			err = unLoadTokenById(app, device, id)
			if err != nil {
				return console.Error(err)
			}

			ok = true
		}

		page++
	}

	return nil
}

func GetTokensByUserId(userId, search string, page, rows int) (List, error) {
	sql := `
  SELECT COUNT(*) AS COUNT
  FROM core.TOKENS A
  WHERE A.USER_ID=$1
	AND CONCAT('NAME:', A.NAME, ':APP:', A.APP, ':DEVICE:', A.DEVICE, ':') ILIKE CONCAT('%', $2, '%');`

	all := QueryCount(sql, userId, search)

	offset := (page - 1) * rows
	sql = `
  SELECT A.*
  FROM core.TOKENS A
	WHERE A.USER_ID=$1
  AND CONCAT('NAME:', A.NAME, ':APP:', A.APP, ':DEVICE:', A.DEVICE, ':') ILIKE CONCAT('%', $2, '%')
	ORDER BY A.APP, A.DEVICE, A.NAME
  LIMIT $3 OFFSET $4;`

	items, err := Query(sql, userId, search, rows, offset)
	if err != nil {
		return List{}, err
	}

	for _, item := range items.Result {
		id := item.Id()
		collection, err := GetCollectionById("telemetry.token.last_use", id)
		if err != nil {
			return List{}, err
		}

		token := item["token"].(string)
		item["token"] = token[0:6] + "..." + token[len(token)-6:]
		item["last_use"] = collection.Str("last_use")
	}

	return items.ToList(all, page, rows), nil
}

func DeleteToken(id string) (Item, error) {
	current, err := GetTokenById(id)
	if err != nil {
		return Item{}, err
	}

	if !current.Ok {
		return Item{}, console.ErrorM(RECORD_NOT_FOUND)
	}

	sql := `
  DELETE FROM core.TOKENS
  WHERE _ID=$1
  RETURNING *;`

	item, err := QueryOne(sql, id)
	if err != nil {
		return Item{}, err
	}

	app := item.Str("app")
	device := item.Str("device")
	err = unLoadTokenById(app, device, id)
	if err != nil {
		return Item{}, err
	}

	return Item{
		Ok: item.Ok,
		Result: Json{
			"message": OkOrNot(item.Ok, RECORD_DELETE, RECORD_NOT_DELETE),
			"index":   item.Index(),
		},
	}, nil
}
