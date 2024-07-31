package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
)

func InitModel(model *linq.Model) error {
	if model == nil {
		return console.PanicM("Model not found")
	}

	if err := model.Init(); err != nil {
		return err
	}

	return nil
}

func NextCode(tag, prefix string) string {
	return jdb.NextCode(tag, prefix)
}
