package jdb

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
)

/**
* defineCore create the core schema
* @param db *linq.DB
* @return error
**/
func defineCommand(db *DB) error {
	exist, err := ExistTable(db, "core", "COMMANDS")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return defineCoreFunction(db)
	}

	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE SCHEMA IF NOT EXISTS core;

  CREATE TABLE IF NOT EXISTS core.COMMANDS(
		OPTION VARCHAR(80) DEFAULT '',
		_ID VARCHAR(80) DEFAULT '-1',
		SQL BYTEA,
		MUTEX INT DEFAULT 0,
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(OPTION, _ID)
	);
	CREATE INDEX IF NOT EXISTS COMMANDS_OPTION_IDX ON core.COMMANDS(OPTION);
	CREATE INDEX IF NOT EXISTS COMMANDS_INDEX_IDX ON core.COMMANDS(INDEX);`

	_, err = db.db.Exec(sql)
	if err != nil {
		return err
	}

	return defineCoreFunction(db)
}

func defineCoreFunction(db *DB) error {
	sql := `
	CREATE OR REPLACE FUNCTION core.COMMANDS_INSERT()
  RETURNS
    TRIGGER AS $$  
  BEGIN
	 	IF NEW.MUTEX = 0 THEN
			PERFORM pg_notify(
			'command',
			json_build_object(        
				'_id', NEW._ID
			)::text
			);
		END IF;
  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS COMMANDS_INSERT ON core.COMMANDS CASCADE;
	CREATE TRIGGER COMMANDS_INSERT
	BEFORE INSERT ON core.COMMANDS
	FOR EACH ROW
	EXECUTE PROCEDURE core.COMMANDS_INSERT();`

	_, err := db.db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* insertCommand
* @params query string
**/
func (d *DB) insertCommand(id string, mutex int, query string) error {
	sql := `
	SELECT INDEX
	FROM core.COMMANDS
	WHERE _ID = $1;`

	item, err := d.QueryOne(sql, id)
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	sql = `
	INSERT INTO core.COMMANDS (_ID, OPTION, SQL, MUTEX, INDEX)
	VALUES ($1, 'INSERT', $2, $3, $4);`

	id = utility.GenKey(id)
	index := NextSerie(d, "commnad")
	_, err = d.db.Exec(sql, id, []byte(query), mutex, index)
	if err != nil {
		logs.Alertm(et.Json{
			"_id":   id,
			"sql":   query,
			"index": index,
		}.ToString())
		return err
	}

	return nil
}

/**
* upsertCommand
* @params query string
**/
func (d *DB) upsertCommand(old, new, id string, mutex int, query string) error {
	sql := `
	SELECT INDEX
	FROM core.COMMANDS
	WHERE OPTION = $1
	AND _ID = $2;`

	item, err := d.QueryOne(sql, old, id)
	if err != nil {
		return err
	}

	if item.Ok {
		index := item.Int64("index")
		sql = `
		UPDATE core.COMMANDS SET
		OPTION = $2,
		SQL = $3
		MUTEX = $4
		WHERE INDEX = $1;`

		_, err = d.db.Exec(sql, index, new, []byte(query), mutex)
		if err != nil {
			return err
		}

		return nil
	}

	sql = `
	INSERT INTO core.COMMANDS (_ID, OPTION, SQL, MUTEX, INDEX)
	VALUES ($1, $2, $3, $4, $5);`

	id = utility.GenKey(id)
	index := NextSerie(d, "commnad")
	_, err = d.db.Exec(sql, id, new, []byte(query), mutex, index)
	if err != nil {
		logs.Alertm(et.Json{
			"_id":   id,
			"sql":   query,
			"index": index,
		}.ToString())
		return err
	}

	return nil
}

/**
* deleteCommand
* @params query string
**/
func (d *DB) deleteCommand(opt, id string) error {
	sql := `
	DELETEFROM core.COMMANDS
	WHERE OPTION = $1
	AND _ID = $2;`

	_, err := d.db.Exec(sql, opt, id)
	if err != nil {
		return err
	}

	return nil
}

/**
* SetCommand
* @params query string
**/
func (d *DB) SetCommand(opt, id, query string) error {
	opt = strs.Uppcase(opt)
	switch opt {
	case CommandInsert:
		return d.insertCommand(id, 0, query)
	case CommandUpdate:
		return d.upsertCommand(opt, opt, id, 0, query)
	case CommandDelete:
		if d.dm != nil {
			return d.upsertCommand(CommandInsert, opt, id, 0, query)
		}
		return d.deleteCommand(CommandInsert, id)
	default:
		return d.upsertCommand(opt, opt, id, 0, query)
	}
}

/**
* SetMutex
* @params id string
* @params query string
* @params index int64
* @return error
**/
func (d *DB) SetMutex(opt, id, query string, index int64) error {
	opt = strs.Uppcase(opt)
	switch opt {
	case CommandInsert:
		return d.insertCommand(id, 1, query)
	case CommandUpdate:
		return d.upsertCommand(opt, opt, id, 1, query)
	case CommandDelete:
		d.deleteCommand(CommandUpdate, id)
		return d.upsertCommand(CommandInsert, opt, id, 1, query)
	default:
		return d.upsertCommand(opt, opt, id, 1, query)
	}
}

/**
* GetCommand
* @params id string
* @return js.Item
* @return error
**/
func (d *DB) GetCommand(id string) (et.Item, error) {
	var result et.Item = et.Item{}

	query := `
	SELECT _ID, SQL, INDEX
	FROM core.COMMANDS
	WHERE _ID = $1 LIMIT 1;`

	rows, err := d.db.Query(query, id)
	if err != nil {
		return result, err
	}

	var _id string
	var sql []byte
	var index int64
	for rows.Next() {
		rows.Scan(&_id, &sql, &index)
		result = et.Item{
			Ok: true,
			Result: et.Json{
				"_id":   _id,
				"sql":   string(sql),
				"index": index,
			},
		}
	}

	return result, nil
}

/**
* getLastCommand
* @return int64
* @return error
**/
func (d *DB) getLastCommand() (int64, error) {
	var result int64 = 0

	sql := `
	SELECT MAX(INDEX) AS result
	FROM core.COMMANDS;`

	rows, err := d.db.Query(sql)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&result)
	}

	return result, nil
}

/**
* SyncCommand
* @return error
**/
func (d *DB) SyncCommand() error {
	if d.dm == nil {
		return logs.Alertm("Database master not found")
	}

	var ok bool = true
	var page int = 1
	var rows int = 1000
	var total int = 0
	lastIndex, err := d.getLastCommand()
	if err != nil {
		return err
	}
	logs.Info(`Sync commands`)

	for ok {
		ok = false

		logs.Debug("Sync commands page:", page)

		offset := (page - 1) * rows
		sql := `
		SELECT A.OPTION, A._ID, A.SQL, A.INDEX
		FROM core.COMMANDS A
		WHERE A.INDEX>=$3
		ORDER BY A.index
		LIMIT $1 OFFSET $2;`

		rows, err := d.dm.Query(sql, rows, offset, lastIndex)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item et.Item
			item.ScanRows(rows)
			opt := item.Str("option")
			id := item.Str("_id")
			sql := item.Str("sql")
			lastIndex = item.Int64("index")

			err = d.SetMutex(opt, id, sql, lastIndex)
			if err != nil {
				return err
			} else {
				total++
			}

			ok = true
		}

		page++
	}

	logs.Infof(`Sync commands total: %d`, total)

	return nil
}
