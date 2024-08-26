package module

import (
	"github.com/cgalvisleon/elvis/jdb"
	"github.com/cgalvisleon/elvis/linq"
)

var SchemaModule *linq.Schema

func DefineSchemaModule(db *jdb.DB) error {
	if SchemaModule != nil {
		return nil
	}

	SchemaModule = linq.NewSchema(db, "module")

	return nil
}
