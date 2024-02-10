package core

import (
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
)

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

	if model.UseSync() {
		SetSyncTrigger(model)
	} else {
		SetListenerTrigger(model)
	}

	if model.UseRecycle {
		SetRecycligTrigger(model)
	}

	if model.UseSerie {
		DefineSerie(model)
	}

	return nil
}
