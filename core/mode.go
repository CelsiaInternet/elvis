package core

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
)

/**
* Mode
*	Handler for CRUD data
 */
const ModeNone = 0   /** Node */
const ModeIdle = 1   /** Libre */
const ModeNode = 2   /** Node */
const ModeBridge = 3 /** Bridge */

var (
	ModeId    string
	ModeTp    int
	MasterIdx int
)

func DefineMode() error {
	if err := DefineCoreSchema(); err != nil {
		return console.PanicE(err)
	}

	sql := `
  -- DROP TABLE IF EXISTS core.MODE CASCADE;

  CREATE TABLE IF NOT EXISTS core.MODE(
		DATE_MAKE TIMESTAMP DEFAULT NOW(),
		DATE_UPDATE TIMESTAMP DEFAULT NOW(),
    _ID VARCHAR(80) DEFAULT '',
    MODE INTEGER DEFAULT 0,
		PASSWORD VARCHAR(250) DEFAULT '',		
    _DATA JSONB DEFAULT '{}',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(_ID)
  );
	`
	_, err := QDDL(sql)
	if err != nil {
		return console.Error(err)
	}

	sql = `
	SELECT A.*
	FROM core.MODE A
	LIMIT 1;`

	item, err := QueryOne(sql)
	if err != nil {
		return console.Error(err)
	}

	if !item.Ok {
		id := NewId()
		mode := ModeIdle
		data := Json{}
		now := Now()
		data["date_make"] = now
		data["date_update"] = now
		data["mode"] = mode
		data["driver"] = Postgres
		data["host"] = ""
		data["port"] = 5432
		data["dbname"] = ""
		data["user"] = ""

		sql = `
		INSERT INTO core.MODE(DATE_MAKE, DATE_UPDATE, _ID, MODE, _DATA)
		VALUES($1, $1, $2, $3, $4)
		RETURNING *;`

		item, err = QueryOne(sql, now, id, mode, data)
		if err != nil {
			return console.Error(err)
		}
	}

	ModeId = item.Id()
	ModeTp = item.Int("mode")

	return nil
}

func GetMode() (Item, error) {
	sql := `
	SELECT A.*
	FROM core.MODE A
	LIMIT 1;`

	item, err := QueryOne(sql)
	if err != nil {
		return Item{}, err
	}

	delete(item.Result, "password")

	return item, nil
}

func SetMode(mode int, driver, host string, port int, dbname, user, password string) (Item, error) {
	if !ValidStr(driver, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "driver")
	}

	if !ValidStr(host, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "host")
	}

	if !ValidStr(dbname, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "dbname")
	}

	if !ValidStr(user, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "user")
	}

	if !ValidStr(password, 0, []string{""}) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "password")
	}

	id := ModeId
	now := Now()
	data := Json{}
	data["date_update"] = now
	data["mode"] = mode
	data["driver"] = driver
	data["host"] = host
	data["port"] = port
	data["dbname"] = dbname
	data["user"] = user

	sql := `
		UPDATE core.MODE SET
		DATE_UPDATE=$2,
		MODE=$3,
		PASSWORD=$4,
		_DATA=$5
		WHERE _ID=$1
		RETURNING _ID;`

	item, err := QueryOne(sql, id, now, mode, password, data)
	if err != nil {
		return Item{}, err
	}

	ModeTp = mode

	return Item{
		Ok: item.Ok,
		Result: Json{
			"message": RECORD_UPDATE,
		},
	}, nil
}
