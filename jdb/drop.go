package jdb

import (
	"github.com/celsiainternet/elvis/strs"
)

/**
* Drop database component
**/

// Drop database
func DropDatabase(db *DB, name string) error {
	name = strs.Lowcase(name)
	sql := strs.Format(`DROP DATABASE %s;`, name)

	id := strs.Format(`drop-db-%s`, name)
	err := db.Exec(id, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop schema
func DropSchema(db *DB, name string) error {
	name = strs.Lowcase(name)
	sql := strs.Format(`DROP SCHEMA %s CASCADE;`, name)

	id := strs.Format(`drop-schema-%s`, name)
	err := db.Exec(id, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop table
func DropTable(db *DB, schema, name string) error {
	sql := strs.Format(`DROP TABLE %s.%s CASCADE;`, schema, name)

	id := strs.Format(`drop-table-%s-%s`, schema, name)
	err := db.Exec(id, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop column
func DropColumn(db *DB, schema, table, name string) error {
	sql := strs.Format(`ALTER TABLE %s.%s DROP COLUMN %s;`, schema, table, name)

	id := strs.Format(`drop-column-%s-%s-%s`, schema, table, name)
	err := db.Exec(id, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop index
func DropIndex(db *DB, schema, table, field string) error {
	indexName := strs.Format(`%s_%s_IDX`, strs.Uppcase(table), strs.Uppcase(field))
	sql := strs.Format(`DROP INDEX %s.%s CASCADE;`, schema, indexName)

	id := strs.Format(`drop-index-%s-%s-%s`, schema, table, field)
	err := db.Exec(id, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop trigger
func DropTrigger(db *DB, schema, table, name string) error {
	sql := strs.Format(`DROP TRIGGER %s.%s CASCADE;`, schema, name)

	id := strs.Format(`drop-trigger-%s-%s-%s`, schema, table, name)
	err := db.Exec(id, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop serie
func DropSerie(db *DB, schema, name string) error {
	sql := strs.Format(`DROP SEQUENCE %s.%s CASCADE;`, schema, name)

	id := strs.Format(`drop-sequence-%s-%s`, schema, name)
	err := db.Exec(id, sql)
	if err != nil {
		return err
	}

	return nil
}

// Drop user
func DropUser(db *DB, name string) error {
	name = strs.Uppcase(name)
	sql := strs.Format(`DROP USER %s;`, name)

	id := strs.Format(`drop-user-%s`, name)
	err := db.Exec(id, sql)
	if err != nil {
		return err
	}

	return nil
}
