package core

import (
	. "github.com/cgalvisleon/elvis/linq"
)

var SchemaCore *Schema

func DefineCoreSchema() error {
	if SchemaCore != nil {
		return nil
	}

	SchemaCore = NewSchema(0, "core")

	return nil
}
