package jdb

import (
	"database/sql"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/strs"
)

/**
* rowsItems
* @param rows *sql.Rows
* @return et.Items
**/
func rowsItems(rows *sql.Rows) et.Items {
	var result = et.Items{Result: []et.Json{}}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Count++
		result.Result = append(result.Result, item)
	}

	return result
}

/**
* rowsItem
* @param rows *sql.Rows
* @return et.Item
**/
func rowsItem(rows *sql.Rows) et.Item {
	var result = et.Item{Result: et.Json{}}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Result = item
	}

	return result
}

/**
* sourcseItems
* @param rows *sql.Rows, source string
* @return et.Items
**/
func sourceItems(rows *sql.Rows, source string) et.Items {
	var result = et.Items{Result: []et.Json{}}
	source = strs.Lowcase(source)
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Count++
		result.Result = append(result.Result, item.Json(source))
	}

	return result
}
