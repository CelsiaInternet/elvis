package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
)

var MasterIdx int = 0

func InitDefine() error {
	if err := DefineSync(); err != nil {
		return err
	}
	if err := DefineSeries(); err != nil {
		return err
	}
	if err := DefineRecycling(); err != nil {
		return err
	}

	console.LogK("CORE", "Init core")

	return nil
}

func InitModel(model *linq.Model) error {
	err := model.Init()
	if err != nil {
		return err
	}

	for _, name := range model.Index {
		_, err := jdb.CreateIndex(model.Db, model.Schema, model.Table, name)
		if err != nil {
			return err
		}
	}

	if model.UseSync {
		SetSyncTrigger(model.Schema, model.Table)
	}

	if model.UseRecycle {
		SetRecycligTrigger(model.Schema, model.Table)
	}

	if model.UseIndex {
		SetSerie(model.Name)
	}

	model.Trigger(linq.BeforeInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		if model.UseIndex {
			index := GetSerie(model.Name)
			new.Set("index", index)
		}

		return nil
	})

	return nil
}

func SetMasterIdx(idx int) {
	MasterIdx = idx
}
