package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/cgalvisleon/elvis/utility"
)

func InitModel(model *linq.Model) error {
	err := model.Init()
	if err != nil {
		return err
	}

	if model.UseSync {
		SetSyncTrigger(model.Schema, model.Table)
	}

	if model.UseRecycle {
		SetRecycligTrigger(model.Schema, model.Table)
	}

	if model.UseIndex {
		model.Trigger(linq.BeforeInsert, func(model *linq.Model, old, new *e.Json, data e.Json) error {
			index := GetSerie(model.Name)
			new.Set("index", index)

			return nil
		})
	}

	model.References(func(references []*linq.ReferenceValue) {
		SetReferences(references)
	})

	return nil
}

/**
*
 */
func ExistDatabase(db int, name string) (bool, error) {
	name = utility.Lowcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_database
		WHERE UPPER(datname) = UPPER($1));`

	item, err := jdb.DBQueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistSchema(db int, name string) (bool, error) {
	name = utility.Lowcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_namespace
		WHERE UPPER(nspname) = UPPER($1));`

	item, err := jdb.DBQueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistUser(db int, name string) (bool, error) {
	name = utility.Uppcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_roles
		WHERE UPPER(rolname) = UPPER($1));`

	item, err := jdb.DBQueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistTable(db int, schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2));`

	item, err := jdb.DBQueryOne(db, sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistColum(db int, schema, table, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.columns
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2)
		AND UPPER(column_name) = UPPER($3));`

	item, err := jdb.DBQueryOne(db, sql, schema, table, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistIndex(db int, schema, table, field string) (bool, error) {
	indexName := utility.Format(`%s_%s_IDX`, utility.Uppcase(table), utility.Uppcase(field))
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_indexes
		WHERE UPPER(schemaname) = UPPER($1)
		AND UPPER(tablename) = UPPER($2)
		AND UPPER(indexname) = UPPER($3));`

	item, err := jdb.QueryOne(sql, schema, table, indexName)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistSerie(db int, schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_sequences
		WHERE UPPER(schemaname) = UPPER($1)
		AND UPPER(sequencename) = UPPER($2));`

	item, err := jdb.DBQueryOne(db, sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

/**
*
**/
func CreateDatabase(db int, name string) (bool, error) {
	name = utility.Lowcase(name)
	exists, err := ExistDatabase(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := utility.Format(`CREATE DATABASE %s;`, name)

		_, err := jdb.DBQuery(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateSchema(db int, name string) (bool, error) {
	name = utility.Lowcase(name)
	exists, err := ExistSchema(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := utility.Format(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE SCHEMA IF NOT EXISTS "%s";`, name)

		_, err := jdb.DBQuery(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateUser(db int, name, password string) (bool, error) {
	name = utility.Uppcase(name)
	exists, err := ExistUser(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		passwordHash, err := utility.PasswordHash(password)
		if err != nil {
			return false, err
		}

		sql := utility.Format(`CREATE USER %s WITH PASSWORD '%s';`, name, passwordHash)

		_, err = jdb.DBQuery(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func ChangePassword(db int, name, password string) (bool, error) {
	exists, err := ExistUser(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, console.ErrorM(msg.SYSTEM_USER_NOT_FOUNT)
	}

	passwordHash, err := utility.PasswordHash(password)
	if err != nil {
		return false, err
	}

	sql := utility.Format(`ALTER USER %s WITH PASSWORD '%s';`, name, passwordHash)

	_, err = jdb.Query(sql)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateColumn(db int, schema, table, name, kind, defaultValue string) (bool, error) {
	exists, err := ExistColum(db, schema, table, name)
	if err != nil {
		return false, err
	}

	if !exists {
		tableName := utility.Format(`%s.%s`, schema, utility.Uppcase(table))
		sql := jdb.SQLDDL(`
		DO $$
		BEGIN
			BEGIN
				ALTER TABLE $1 ADD COLUMN $2 $3 DEFAULT $4;
			EXCEPTION
				WHEN duplicate_column THEN RAISE NOTICE 'column <column_name> already exists in <table_name>.';
			END;
		END;
		$$;`, tableName, utility.Uppcase(name), utility.Uppcase(kind), defaultValue)

		_, err := jdb.QDDL(sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateIndex(db int, schema, table, field string) (bool, error) {
	exists, err := ExistIndex(db, schema, table, field)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_IDX ON $1.$2($3);`, utility.Uppcase(schema), utility.Uppcase(table), utility.Uppcase(field))

		_, err := jdb.QDDL(sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateSerie(db int, schema, tag string) (bool, error) {
	exists, err := ExistSerie(db, schema, tag)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := utility.Format(`CREATE SEQUENCE IF NOT EXISTS %s START 1;`, tag)

		_, err := jdb.Query(sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}
