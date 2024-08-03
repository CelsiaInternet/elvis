package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/et"
)

func rowsItems(rows *sql.Rows) et.Items {
	var result et.Items = et.Items{}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Count++
		result.Result = append(result.Result, item)
	}

	return result
}

func rowsItem(rows *sql.Rows) et.Item {
	var result et.Item = et.Item{}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Result = item
	}

	return result
}

func sourceItems(rows *sql.Rows, source string) et.Items {
	var result et.Items = et.Items{}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Count++
		result.Result = append(result.Result, item.Json(source))
	}

	return result
}

func sourceItem(rows *sql.Rows, source string) et.Item {
	var result et.Item = et.Item{}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Result = item.Json(source)
	}

	return result
}
