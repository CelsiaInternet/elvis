package linq

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
)

/**
* syncListener handler listened
* @param res js.Json
**/
func syncListener(res et.Json) {
	schema := res.Str("schema")
	table := res.Str("table")
	model := Table(schema, table)
	if model != nil && model.OnListener != nil {
		model.OnListener(res)
	}
}

/**
* recyclingListener handler listened
* @param res js.Json
**/
func recyclingListener(res et.Json) {
	schema := res.Str("schema")
	table := res.Str("table")
	model := Table(schema, table)
	if model != nil && model.OnListener != nil {
		model.OnListener(res)
	}
}

/**
* SetListener
* @param db *jdb.DB
**/
func SetListener(db *jdb.DB) {
	if !db.UseCore {
		return
	}

	db.SetListen([]string{"sync"}, syncListener)
	db.SetListen([]string{"recycling"}, recyclingListener)
}
