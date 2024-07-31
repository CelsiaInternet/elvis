package jdb

import (
	"github.com/cgalvisleon/elvis/console"
)

var makedCore bool

func InitCore() error {
	if makedCore {
		return nil
	}

	db := DB(0)
	if db == nil {
		return console.PanicM("Database not found")
	}

	err := CreateSchema(db.Db, "core")
	if err != nil {
		return err
	}

	if err := defineSync(); err != nil {
		return err
	}

	if err := defineSeries(); err != nil {
		return err
	}

	if err := defineRecycling(); err != nil {
		return err
	}

	makedCore = true

	console.LogK("CORE", "Init core")

	return nil
}
