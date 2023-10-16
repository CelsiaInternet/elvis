package jdb

import (
	"database/sql"

	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/msg"
)

func rowsItems(rows *sql.Rows) Items {
	var result Items = Items{Result: []Json{}}

	for rows.Next() {
		var item Item
		item.Scan(rows)
		result.Result = append(result.Result, item.Result)
		result.Ok = true
		result.Count++
	}

	return result
}

func rowsItem(rows *sql.Rows) Item {
	var result Item = Item{
		Ok: false,
		Result: Json{
			"message": RECORD_NOT_FOUND,
		},
	}

	for rows.Next() {
		result.Scan(rows)
		result.Ok = true
	}

	return result
}

func atribItems(rows *sql.Rows, atrib string) Items {
	var result Items = Items{Result: []Json{}}

	for rows.Next() {
		var item Item
		item.Scan(rows)
		result.Result = append(result.Result, item.Result.Json(atrib))
		result.Ok = true
		result.Count++
	}

	return result
}

func atribItem(rows *sql.Rows, atrib string) Item {
	var result Item = Item{
		Ok: false,
		Result: Json{
			"message": RECORD_NOT_FOUND,
		},
	}

	for rows.Next() {
		result.Scan(rows)
		result.Result = result.Result.Json(atrib)
		result.Ok = true
	}

	return result
}
