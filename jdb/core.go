package jdb

import (
	"github.com/celsiainternet/elvis/console"
)

var makedCore bool

func InitCore(db *DB) error {
	if makedCore {
		return nil
	}

	if db == nil {
		return console.PanicM("Database not found")
	}

	if err := defineCommand(db); err != nil {
		return err
	}

	if err := defineRecords(db); err != nil {
		return err
	}

	if err := defineSeries(db); err != nil {
		return err
	}

	if err := defineRecycling(db); err != nil {
		return err
	}

	makedCore = true

	console.LogK("CORE", "Init core")

	return nil
}
