package jdb

import "github.com/celsiainternet/elvis/logs"

var makedCore bool

func InitCore(db *DB) error {
	if makedCore {
		return nil
	}

	if db == nil {
		return logs.Panicm("Database not found")
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

	if err := defineDDL(db); err != nil {
		return err
	}

	makedCore = true

	logs.Log("CORE", "Init core")

	return nil
}
