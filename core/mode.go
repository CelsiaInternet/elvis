package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
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
	existMode bool
)

func DefineMode() error {
	existMode, _ := ExistTable(0, "core", "MODE")
	if existMode {
		return nil
	}

	if err := DefineCollection(); err != nil {
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
	_, err := jdb.QDDL(sql)
	if err != nil {
		return console.Error(err)
	}

	sql = `
	SELECT A.*
	FROM core.MODE A
	LIMIT 1;`

	item, err := jdb.QueryOne(sql)
	if err != nil {
		return console.Error(err)
	}

	if !item.Ok {
		id := utility.NewId()
		mode := ModeIdle
		data := e.Json{}
		now := utility.Now()
		data["date_make"] = now
		data["date_update"] = now
		data["mode"] = mode
		data["driver"] = jdb.Postgres
		data["host"] = ""
		data["port"] = 5432
		data["dbname"] = ""
		data["user"] = ""

		sql = `
		INSERT INTO core.MODE(DATE_MAKE, DATE_UPDATE, _ID, MODE, _DATA)
		VALUES($1, $1, $2, $3, $4)
		RETURNING *;`

		item, err = jdb.QueryOne(sql, now, id, mode, data)
		if err != nil {
			return console.Error(err)
		}
	}

	ModeId = item.Id()
	ModeTp = item.Int("mode")

	return nil
}

func GetMode() (e.Item, error) {
	sql := `
	SELECT A.*
	FROM core.MODE A
	LIMIT 1;`

	item, err := jdb.QueryOne(sql)
	if err != nil {
		return e.Item{}, err
	}

	delete(item.Result, "password")

	return item, nil
}

func SetMode(mode int, driver, host string, port int, dbname, user, password string) (e.Item, error) {
	if !utility.ValidStr(driver, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "driver")
	}

	if !utility.ValidStr(host, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "host")
	}

	if !utility.ValidStr(dbname, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "dbname")
	}

	if !utility.ValidStr(user, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "user")
	}

	if !utility.ValidStr(password, 0, []string{""}) {
		return e.Item{}, console.AlertF(msg.MSG_ATRIB_REQUIRED, "password")
	}

	id := ModeId
	now := utility.Now()
	data := e.Json{}
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

	item, err := jdb.QueryOne(sql, id, now, mode, password, data)
	if err != nil {
		return e.Item{}, err
	}

	ModeTp = mode

	return e.Item{
		Ok: item.Ok,
		Result: e.Json{
			"message": msg.RECORD_UPDATE,
		},
	}, nil
}
