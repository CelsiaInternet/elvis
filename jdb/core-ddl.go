package jdb

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
)

/**
* defineDDL create alter or delete
* @param db *linq.DB
* @return error
**/
func defineDDL(db *DB) error {
	exist, err := ExistTable(db, "core", "DDL")
	if err != nil {
		return logs.Panice(err)
	}

	if exist {
		return nil
	}

	sql := `
  CREATE TABLE IF NOT EXISTS core.DDL(
		_ID VARCHAR(80) DEFAULT '-1',
		SQL BYTEA,
		_IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(_ID)
	);
	CREATE INDEX IF NOT EXISTS COMMANDS__ID_IDX ON core.DDL(_ID);
	CREATE INDEX IF NOT EXISTS COMMANDS__IDT_IDX ON core.DDL(_IDT);
	CREATE INDEX IF NOT EXISTS COMMANDS_INDEX_IDX ON core.DDL(INDEX);

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_INSERT ON core.DDL CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_INSERT
	BEFORE INSERT ON core.DDL
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_INSERT();

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_UPDATE ON core.DDL CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_UPDATE
	BEFORE UPDATE ON core.DDL
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_UPDATE();

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_DELETE ON core.DDL CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_DELETE
	BEFORE DELETE ON core.DDL
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
* upsertDDL
* @params query string
**/
func (d *DB) upsertDDL(id string, query string) error {
	sql := `
	UPDATE core.DDL SET
	SQL = $2
	WHERE _ID = $1
	RETURNING _ID;`

	item, err := d.QueryOne(sql, id, []byte(query))
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	sql = `
	INSERT INTO core.DDL(_ID, SQL, INDEX)
	VALUES ($1, $2, $3);`

	id = utility.GenKey(id)
	index := NextSerie(d, "ddl")
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
* deleteDDL
* @params query string
**/
func (d *DB) deleteDDL(id string) error {
	sql := `
	DELETE FROM core.DDL
	WHERE _ID = $1;`

	_, err := d.db.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}
