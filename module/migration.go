package module

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/utility"
)

var Migration *linq.Model

func DefineMigration(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return console.Panic(err)
	}

	if Migration != nil {
		return nil
	}

	Migration = linq.NewModel(SchemaModule, "MIGRATION", "Tabla de migracion", 1)
	Migration.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Migration.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Migration.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Migration.DefineColum("old_id", "", "VARCHAR(80)", "")
	Migration.DefineColum("id", "", "VARCHAR(80)", "")
	Migration.DefineColum("tag", "", "VARCHAR(250)", "")
	Migration.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
	})

	if err := Migration.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* IdMigration
* @param old_id string
* @param tag string
* @return string, error
**/
func IdMigration(old_id string, tag string) (string, error) {
	if !utility.ValidId(old_id) {
		return old_id,
			console.AlertF("Id invalido: %s", old_id)
	}

	if !utility.ValidNil(tag) {
		return old_id,
			console.AlertF("Tag invalido: %s", tag)
	}

	item, err := Migration.Select().
		Where(Migration.Col("old_id").Eq(old_id)).
		And(Migration.Col("tag").Eq(tag)).
		First()
	if err != nil {
		return old_id, err
	}

	if item.Ok {
		now := utility.Now()
		_, err := Migration.Insert(et.Json{
			"data_make":   now,
			"date_update": now,
			"_state":      utility.ACTIVE,
			"old_id":      old_id,
			"id":          old_id,
			"tag":         tag,
		}).
			CommandOne()
		if err != nil {
			return old_id, err
		}

		result := item.Key("id")
		return result, nil
	}

	return old_id, nil
}

/**
* SetMigration
* @param old_id string
* @param id string
* @param tag string
* @return et.Item, error
**/
func SetMigration(old_id string, id string, tag string) (et.Item, error) {
	if !utility.ValidId(old_id) {
		return et.Item{},
			console.AlertF("Id invalido: %s", old_id)
	}

	if !utility.ValidNil(tag) {
		return et.Item{},
			console.AlertF("Tag invalido: %s", tag)
	}

	if !utility.ValidId(id) {
		return et.Item{},
			console.AlertF("Id invalido: %s", id)
	}

	now := utility.Now()
	updateData := et.Json{
		"data_make":   now,
		"date_update": now,
		"_state":      utility.ACTIVE,
		"old_id":      old_id,
		"id":          id,
		"tag":         tag,
	}

	item, err := Migration.Update(updateData).
		Where(Migration.Col("old_id").Eq(old_id)).
		And(Migration.Col("tag").Eq(tag)).
		CommandOne()

	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}
