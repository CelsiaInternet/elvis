package core

import (
	"github.com/cgalvisleon/elvis/console"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/linq"
)

var Configs *linq.Model

func DefineConfig() error {
	if err := DefineSchemaCore(); err != nil {
		return console.PanicE(err)
	}

	if Configs != nil {
		return nil
	}

	Configs = linq.NewModel(SchemaCore, "CONFIG", "Tabla de configuración", 1)
	Configs.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Configs.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Configs.DefineColum("_key", "", "VARCHAR(80)", "")
	Configs.DefineColum("_value", "", "TEXT", "")
	Configs.DefineColum("index", "", "INTEGER", 0)
	Configs.DefinePrimaryKey([]string{"_key"})
	Configs.DefineIndex([]string{
		"date_make",
		"date_update",
		"index",
	})

	return InitModel(Configs)
}

func GetConfig(key string, _default string) (string, error) {
	item, err := Configs.Select("_value").
		Where(Configs.Col("_key").Eq(key)).
		First()
	if err != nil {
		return _default, err
	}

	result := item.Str(_default, "_value")

	return result, nil
}

func UpsertConfig(key, value string) (e.Item, error) {
	item, err := Configs.Upsert(e.Json{
		"_key":  key,
		"value": value,
	}).
		Where(Configs.Col("_key").Eq(key)).
		Command()
	if err != nil {
		return e.Item{}, err
	}

	return item, nil
}

func DeleteConfig(key string) (e.Item, error) {
	return Configs.Delete().
		Where(Configs.Col("_key").Eq(key)).
		Command()
}
