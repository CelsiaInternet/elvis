package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/strs"
)

/**
* Drop database component
**/

// Drop database
func DropDatabase(db *sql.DB, name string) error {
	name = strs.Lowcase(name)
	sql := strs.Format(`DROP DATABASE %s;`, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop schema
func DropSchema(db *sql.DB, name string) error {
	name = strs.Lowcase(name)
	sql := strs.Format(`DROP SCHEMA %s CASCADE;`, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop table
func DropTable(db *sql.DB, schema, name string) error {
	sql := strs.Format(`DROP TABLE %s.%s CASCADE;`, schema, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop column
func DropColumn(db *sql.DB, schema, table, name string) error {
	sql := strs.Format(`ALTER TABLE %s.%s DROP COLUMN %s;`, schema, table, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop index
func DropIndex(db *sql.DB, schema, table, field string) error {
	indexName := strs.Format(`%s_%s_IDX`, strs.Uppcase(table), strs.Uppcase(field))
	sql := strs.Format(`DROP INDEX %s.%s CASCADE;`, schema, indexName)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop trigger
func DropTrigger(db *sql.DB, schema, table, name string) error {
	sql := strs.Format(`DROP TRIGGER %s.%s CASCADE;`, schema, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop serie
func DropSerie(db *sql.DB, schema, name string) error {
	sql := strs.Format(`DROP SEQUENCE %s.%s CASCADE;`, schema, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop user
func DropUser(db *sql.DB, name string) error {
	name = strs.Uppcase(name)
	sql := strs.Format(`DROP USER %s;`, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}
