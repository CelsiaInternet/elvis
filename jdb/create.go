package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
)

/**
* Created database component
**/

// Crate database
func CreateDatabase(db *sql.DB, name string) error {
	name = strs.Lowcase(name)
	exists, err := ExistDatabase(db, name)
	if err != nil {
		return err
	}

	if !exists {
		sql := strs.Format(`CREATE DATABASE %s;`, name)

		_, err := DBQuery(db, sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create schema
func CreateSchema(db *sql.DB, name string) error {
	name = strs.Lowcase(name)
	exists, err := ExistSchema(db, name)
	if err != nil {
		return err
	}

	if !exists {
		sql := strs.Format(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE SCHEMA IF NOT EXISTS "%s";`, name)

		_, err := DBQuery(db, sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create column
func CreateColumn(db *sql.DB, schema, table, name, kind, defaultValue string) error {
	exists, err := ExistColum(db, schema, table, name)
	if err != nil {
		return err
	}

	if !exists {
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

		_, err := QDDL(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create index
func CreateIndex(db *sql.DB, schema, table, field string) error {
	exists, err := ExistIndex(db, schema, table, field)
	if err != nil {
		return err
	}

	if !exists {
		sql := SQLDDL(`
		CREATE INDEX IF NOT EXISTS $2_$3_IDX ON $1.$2($3);`,
			strs.Uppcase(schema), strs.Uppcase(table), strs.Uppcase(field))

		_, err := QDDL(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create trigger
func CreateTrigger(db *sql.DB, schema, table, name, when, event, function string) error {
	exists, err := ExistTrigger(db, schema, table, name)
	if err != nil {
		return err
	}

	if !exists {
		sql := SQLDDL(`
		DROP TRIGGER IF EXISTS $3 ON $1.$2 CASCADE;
		CREATE TRIGGER $3
		$4 $5 ON $1.$2
		FOR EACH ROW
		EXECUTE PROCEDURE $6;`,
			strs.Uppcase(schema), strs.Uppcase(table), strs.Uppcase(name), when, event, function)

		_, err := QDDL(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create serie
func CreateSerie(db *sql.DB, schema, tag string) error {
	exists, err := ExistSerie(db, schema, tag)
	if err != nil {
		return err
	}

	if !exists {
		sql := strs.Format(`CREATE SEQUENCE IF NOT EXISTS %s START 1;`, tag)

		_, err := Query(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create user
func CreateUser(db *sql.DB, name, password string) error {
	name = strs.Uppcase(name)
	exists, err := ExistUser(db, name)
	if err != nil {
		return err
	}

	if !exists {
		passwordHash, err := utility.PasswordHash(password)
		if err != nil {
			return err
		}

		sql := strs.Format(`CREATE USER %s WITH PASSWORD '%s';`, name, passwordHash)

		_, err = DBQuery(db, sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// Changue password
func ChangePassword(db *sql.DB, name, password string) error {
	exists, err := ExistUser(db, name)
	if err != nil {
		return err
	}

	if !exists {
		return console.ErrorM(msg.SYSTEM_USER_NOT_FOUNT)
	}

	passwordHash, err := utility.PasswordHash(password)
	if err != nil {
		return err
	}

	sql := strs.Format(`ALTER USER %s WITH PASSWORD '%s';`, name, passwordHash)

	_, err = Query(sql)
	if err != nil {
		return err
	}

	return nil
}