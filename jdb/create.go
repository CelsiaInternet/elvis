package jdb

import (
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
)

/**
* Created database component
**/

// Crate database
func CreateDatabase(db *DB, name string) error {
	name = strs.Lowcase(name)
	exists, err := ExistDatabase(db, name)
	if err != nil {
		return err
	}

	if !exists {
		sql := strs.Format(`CREATE DATABASE %s;`, name)

		_, err := db.Command(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create schema
func CreateSchema(db *DB, name string) error {
	sql := strs.Format(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE SCHEMA IF NOT EXISTS "%s";`, name)

	_, err := db.Command(sql)
	if err != nil {
		return err
	}

	return nil
}

// Create column
func CreateColumn(db *DB, schema, table, name, kind, defaultValue string) error {
	tableName := strs.Format(`%s.%s`, schema, strs.Uppcase(table))
	sql := SQLDDL(`
	DO $$
	BEGIN
		BEGIN
			ALTER TABLE $1 ADD COLUMN $2 $3 DEFAULT $4;
		EXCEPTION
			WHEN duplicate_column THEN RAISE NOTICE 'column <column_name> already exists in <table_name>.';
		END;
	END;
	$$;`, tableName, strs.Uppcase(name), strs.Uppcase(kind), defaultValue)

	_, err := db.Command(sql)
	if err != nil {
		return err
	}

	return nil
}

// Create index
func CreateIndex(db *DB, schema, table, field string) error {
	sql := SQLDDL(`
	CREATE INDEX IF NOT EXISTS $2_$3_IDX ON $1.$2($3);`,
		strs.Uppcase(schema), strs.Uppcase(table), strs.Uppcase(field))

	_, err := db.Command(sql)
	if err != nil {
		return err
	}

	return nil
}

// Create trigger
func CreateTrigger(db *DB, schema, table, name, when, event, function string) error {
	sql := SQLDDL(`
	DROP TRIGGER IF EXISTS $3 ON $1.$2 CASCADE;
	CREATE TRIGGER $3
	$4 $5 ON $1.$2
	FOR EACH ROW
	EXECUTE PROCEDURE $6;`,
		strs.Uppcase(schema), strs.Uppcase(table), strs.Uppcase(name), when, event, function)

	_, err := db.Command(sql)
	if err != nil {
		return err
	}

	return nil
}

// Create serie
func CreateSequence(db *DB, schema, tag string) error {
	sql := strs.Format(`CREATE SEQUENCE IF NOT EXISTS %s START 1;`, tag)

	_, err := db.Command(sql)
	if err != nil {
		return err
	}

	return nil
}

// Create user
func CreateUser(db *DB, name, password string) error {
	passwordHash, err := utility.PasswordHash(password)
	if err != nil {
		return err
	}

	sql := strs.Format(`CREATE USER %s WITH PASSWORD '%s';`, name, passwordHash)

	_, err = db.Command(sql)
	if err != nil {
		return err
	}

	return nil
}

// Changue password
func ChangePassword(db *DB, name, password string) error {
	passwordHash, err := utility.PasswordHash(password)
	if err != nil {
		return err
	}

	sql := strs.Format(`ALTER USER %s WITH PASSWORD '%s';`, name, passwordHash)

	_, err = db.Command(sql)
	if err != nil {
		return err
	}

	return nil
}
