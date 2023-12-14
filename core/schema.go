package core

import (
	"github.com/cgalvisleon/elvis/linq"
)

var SchemaCore *linq.Schema

func defineSchema() error {
	if SchemaCore != nil {
		return nil
	}

	SchemaCore = linq.NewSchema(0, "core")

	return nil
}
