package jdb

import (
	"database/sql"

	e "github.com/cgalvisleon/elvis/json"
)

/**
* Data Definition Language
**/

func rowsItems(rows *sql.Rows) e.Items {
	var result e.Items = e.Items{Result: []e.Json{}}

	for rows.Next() {
		var item e.Item
		item.Scan(rows)
		result.Result = append(result.Result, item.Result)
		result.Ok = true
		result.Count++
	}

	return result
}

func atribItems(rows *sql.Rows, atrib string) e.Items {
	var result e.Items = e.Items{Result: []e.Json{}}

	for rows.Next() {
		var item e.Item
		item.Scan(rows)
		result.Result = append(result.Result, item.Result.Json(atrib))
		result.Ok = true
		result.Count++
	}

	return result
}
