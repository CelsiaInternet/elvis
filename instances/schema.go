package instances

import (
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
)

/**
* defineSchema
* @param db *jdb.DB, name string
* @return (*linq.Schema, error)
**/
func defineSchema(db *jdb.DB, name string) (*linq.Schema, error) {
	schema := linq.NewSchema(db, name)
	return schema, nil
}
