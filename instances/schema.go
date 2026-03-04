package instances

import (
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
)

var schema *linq.Schema

func defineSchema(db *jdb.DB, name string) error {
	if schema == nil {
		schema = linq.NewSchema(db, name)
	}

	return nil
}
