package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/et"
)

/**
* defineVars define the vars table
* @param db *sql.DB
* @return error
**/
func defineVars(db *sql.DB) error {
	sql := `
	CREATE SCHEMA IF NOT EXISTS core;

  CREATE TABLE IF NOT EXISTS core.VARS(		
		VAR VARCHAR(80) DEFAULT '',
		VALUE VARCHAR(250) DEFAULT '',
		PRIMARY KEY(VAR)
	);`

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	err = initVar(db, "REPLICA", "1")
	if err != nil {
		return err
	}

	return nil
}

/**
* initVar init a var
* @param db *sql.DB
* @param name string
* @param value string
* @return error
**/
func initVar(db *sql.DB, name string, value string) error {
	sql := `
	SELECT VALUE
	FROM core.VARS
	WHERE VAR = $1;`

	item, err := QueryOne(db, sql, name)
	if err != nil {
		return err
	}

	if !item.Ok {
		sql = `
		INSERT INTO core.VARS (VAR, VALUE)
		VALUES ($1, $2);`

		_, err := QueryOne(db, sql, name, value)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
* setVar set a var
* @param db *sql.DB
* @param name string
* @param value string
* @return error
**/
func SetVar(db *sql.DB, name string, value string) error {
	sql := `
	INSERT INTO core.VARS (VAR, VALUE)
	VALUES ($1, $2)
	ON CONFLICT (VAR) DO UPDATE SET
	VALUE = $2;`

	_, err := db.Exec(sql, name, value)
	if err != nil {
		return err
	}

	return nil
}

/**
* getVar set a var
* @param db *sql.DB
* @param name string
* @param def string
* @return string
* @return error
**/
func GetVar(db *sql.DB, name, def any) (*et.Any, error) {
	result := et.NewAny(def)

	sql := `
	SELECT VALUE
	FROM core.VARS
	WHERE VAR = $1;`

	item, err := QueryOne(db, sql, name)
	if err != nil {
		return result, err
	}

	if !item.Ok {
		return result, nil
	}

	val := item.ValStr(result.Str(), "value")
	result.Set(val)

	return result, nil
}

/**
* delVar delete a var
* @param db *sql.DB
* @param name string
* @return error
**/
func DelVar(db *sql.DB, name string) error {
	sql := `
	DELETE FROM core.VARS
	WHERE VAR = $1;`

	_, err := db.Exec(sql, name)
	if err != nil {
		return err
	}

	return nil
}
