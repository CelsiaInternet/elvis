package module

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

var Migration *linq.Model

func DefineMigration(db *jdb.DB) error {
	if err := DefineSchemaModule(db); err != nil {
		return logs.Panice(err)
	}

	if Migration != nil {
		return nil
	}

	Migration = linq.NewModel(SchemaModule, "MIGRATION", "Tabla de migracion", 1)
	Migration.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Migration.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Migration.DefineColum("old_id", "", "VARCHAR(80)", "")
	Migration.DefineColum("id", "", "VARCHAR(80)", "")
	Migration.DefineColum("tag", "", "VARCHAR(250)", "")
	Migration.DefinePrimaryKey([]string{"old_id", "tag"})
	Migration.DefineIndex([]string{
		"date_make",
		"date_update",
		"tag",
	})

	if err := Migration.Init(); err != nil {
		return logs.Panice(err)
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
		return old_id, nil
	}

	if !utility.ValidNil(tag) {
		return old_id, nil
	}

	item, err := Migration.Select().
		Where(Migration.Col("old_id").Eq(old_id)).
		And(Migration.Col("tag").Eq(tag)).
		First()
	if err != nil {
		return old_id, nil
	}

	if !item.Ok {
		now := utility.Now()
		Migration.Insert(et.Json{
			"data_make":   now,
			"date_update": now,
			"old_id":      old_id,
			"id":          old_id,
			"tag":         tag,
		}).
			CommandOne()

		return old_id, nil
	}

	result := item.ValStr(old_id, "id")

	return result, nil
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
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "old_id")
	}

	if !utility.ValidNil(tag) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "tag")
	}

	if !utility.ValidId(id) {
		return et.Item{}, logs.Alertf(msg.MSG_ATRIB_REQUIRED, "id")
	}

	current, err := Migration.Select().
		Where(Migration.Col("old_id").Eq(old_id)).
		And(Migration.Col("tag").Eq(tag)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		now := utility.Now()
		data := et.Json{
			"data_make":   now,
			"date_update": now,
			"old_id":      old_id,
			"id":          id,
			"tag":         tag,
		}
		result, err := Migration.Insert(data).
			Where(Migration.Col("old_id").Eq(old_id)).
			And(Migration.Col("tag").Eq(tag)).
			CommandOne()
		if err != nil {
			return et.Item{}, err
		}

		return result, nil
	}

	data := et.Json{
		"date_update": utility.Now(),
		"id":          id,
	}
	result, err := Migration.Update(data).
		Where(Migration.Col("old_id").Eq(old_id)).
		And(Migration.Col("tag").Eq(tag)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}
