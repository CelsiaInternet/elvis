package linq

import (
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/jdb"
)

/**
* beforeListener handler listened
* @param res js.Json
**/
func commandListener(res et.Json) {

}

/**
* beforeListener handler listened
* @param res js.Json
**/
func beforeListener(res et.Json) {
	schema := res.Str("schema")
	table := res.Str("table")
	model := Table(schema, table)
	if model != nil && model.OnListener != nil {
		model.OnListener(res)
	}
}

/**
* afterListener handler listened
* @param res js.Json
**/
func afterListener(res et.Json) {
	schema := res.Str("schema")
	table := res.Str("table")
	model := Table(schema, table)
	if model != nil && model.OnListener != nil {
		model.OnListener(res)
	}
}

/**
* setListener
* @param db *jdb.DB
**/
func setListener(db *jdb.DB) {
	if !db.UseCore {
		return
	}

	db.SetListen([]string{"before"}, beforeListener)
	db.SetListen([]string{"after"}, afterListener)
	db.SetListen([]string{"command"}, commandListener)
}
