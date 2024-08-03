package module

import (
	"database/sql"

	"github.com/cgalvisleon/elvis/linq"
)

var SchemaModule *linq.Schema

func DefineSchemaModule(db *sql.DB) error {
	if SchemaModule != nil {
		return nil
	}

	SchemaModule = linq.NewSchema(db, "module")

	return nil
}
