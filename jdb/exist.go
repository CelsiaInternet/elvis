package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/strs"
)

/**
* Exist database component
**/

// Exist database
func ExistDatabase(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_database
		WHERE UPPER(datname) = UPPER($1));`

	item, err := QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// Exist schema
func ExistSchema(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_namespace
		WHERE UPPER(nspname) = UPPER($1));`

	item, err := QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// Exist table
func ExistTable(db *sql.DB, schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2));`

	item, err := QueryOne(db, sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// Exist column
func ExistColum(db *sql.DB, schema, table, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.columns
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2)
		AND UPPER(column_name) = UPPER($3));`

	item, err := QueryOne(db, sql, schema, table, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// Exist index
func ExistIndex(db *sql.DB, schema, table, field string) (bool, error) {
	indexName := strs.Format(`%s_%s_IDX`, strs.Uppcase(table), strs.Uppcase(field))
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_indexes
		WHERE UPPER(schemaname) = UPPER($1)
		AND UPPER(tablename) = UPPER($2)
		AND UPPER(indexname) = UPPER($3));`

	item, err := QueryOne(db, sql, schema, table, indexName)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// Exist trigger
func ExistTrigger(db *sql.DB, schema, table, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.triggers
		WHERE UPPER(event_object_schema) = UPPER($1)
		AND UPPER(event_object_table) = UPPER($2)
		AND UPPER(trigger_name) = UPPER($3));`

	item, err := QueryOne(db, sql, schema, table, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// Exist serie
func ExistSerie(db *sql.DB, schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_sequences
		WHERE UPPER(schemaname) = UPPER($1)
		AND UPPER(sequencename) = UPPER($2));`

	item, err := QueryOne(db, sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// Exist user
func ExistUser(db *sql.DB, name string) (bool, error) {
	name = strs.Uppcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_roles
		WHERE UPPER(rolname) = UPPER($1));`

	item, err := QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}
