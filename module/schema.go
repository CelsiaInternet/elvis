package module

import (
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
)

var SchemaModule *linq.Schema

func DefineSchemaModule(db *jdb.DB) error {
	if SchemaModule != nil {
		return nil
	}

	SchemaModule = linq.NewSchema(db, PackageName)

	return nil
}
