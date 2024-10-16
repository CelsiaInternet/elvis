package linq

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/utility"
)

func beforeInsert(model *Model, old, new *et.Json, data et.Json) error {
	now := utility.Now()

	if model.UseDateMake {
		new.Set(DateMakeField, now)
	}

	if model.UseDateUpdate {
		new.Set(DateUpdateField, now)
	}

	if model.UseSerie {
		index := jdb.NextSerie(model.db, model.Table)
		new.Set(SerieField, index)
	}

	return nil
}

func afterInsert(model *Model, old, new *et.Json, data et.Json) error {

	return nil
}

func beforeUpdate(model *Model, old, new *et.Json, data et.Json) error {
	now := utility.Now()

	if model.UseDateUpdate {
		new.Set(DateUpdateField, now)
	}

	return nil
}

func afterUpdate(model *Model, old, new *et.Json, data et.Json) error {

	return nil
}

func beforeDelete(model *Model, old, new *et.Json, data et.Json) error {
	return nil
}

func afterDelete(model *Model, old, new *et.Json, data et.Json) error {

	return nil
}
