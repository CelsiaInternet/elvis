package module

import (
	. "github.com/cgalvisleon/elvis/linq"
)

var SchemaModule *Schema

func defineSchema() error {
	if SchemaModule != nil {
		return nil
	}

	SchemaModule = NewSchema(0, "module")

	return nil
}
