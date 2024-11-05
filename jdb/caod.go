package jdb

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
)

/**
* defineCAD create alter or delete
* @param db *linq.DB
* @return error
**/
func defineCAOD(db *DB) error {
	exist, err := ExistTable(db, "core", "CAOD")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return nil
	}

	sql := `
  CREATE TABLE IF NOT EXISTS core.CAOD(
		_ID VARCHAR(80) DEFAULT '-1',
		SQL BYTEA,
		_IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(OPTION, _ID)
	);
	CREATE INDEX IF NOT EXISTS COMMANDS__ID_IDX ON core.CAOD(_ID);
	CREATE INDEX IF NOT EXISTS COMMANDS__IDT_IDX ON core.CAOD(_IDT);
	CREATE INDEX IF NOT EXISTS COMMANDS_INDEX_IDX ON core.CAOD(INDEX);

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_INSERT ON core.CAOD CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_INSERT
	BEFORE INSERT ON core.CAOD
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_INSERT();

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_UPDATE ON core.CAOD CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_UPDATE
	BEFORE UPDATE ON core.CAOD
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_UPDATE();

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_DELETE ON core.CAOD CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_DELETE
	BEFORE DELETE ON core.CAOD
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_DELETE();
	`

	_, err = db.db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* upsertCAOD
* @params query string
**/
func (d *DB) upsertCAOD(id string, query string) error {
	sql := `
	SELECT INDEX
	FROM core.CAOD
	WHERE _ID = $1;`

	item, err := d.QueryOne(sql, id)
	if err != nil {
		return err
	}

	if item.Ok {
		sql = `
		UPDATE core.CAOD SET
		SQL = $2
		WHERE _ID = $1;`

		_, err = d.db.Exec(sql, id, []byte(query))
		if err != nil {
			return err
		}

		return nil
	}

	sql = `
	INSERT INTO core.CAOD (_ID, SQL, INDEX)
	VALUES ($1, $2, $3);`

	id = utility.GenKey(id)
	index := NextSerie(d, "caod")
	_, err = d.db.Exec(sql, id, []byte(query), index)
	if err != nil {
		logs.Alertm(et.Json{
			"_id": id,
			"sql": query,
		}.ToString())
		return err
	}

	return nil
}

/**
* deleteCAOD
* @params query string
**/
func (d *DB) deleteCAOD(id string) error {
	sql := `
	DELETE FROM core.CAOD
	WHERE _ID = $1;`

	_, err := d.db.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}
