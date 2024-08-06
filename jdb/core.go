package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/console"
)

var makedCore bool

func InitCore(db *sql.DB) error {
	if makedCore {
		return nil
	}

	if db == nil {
		return console.PanicM("Database not found")
	}

	err := CreateSchema(db, "core")
	if err != nil {
		return err
	}

	if err := defineSync(db); err != nil {
		return err
	}

	if err := defineSeries(db); err != nil {
		return err
	}

	if err := defineRecycling(db); err != nil {
		return err
	}

	if err := defineVars(db); err != nil {
		return err
	}

	makedCore = true

	console.LogK("CORE", "Init core")

	return nil
}
