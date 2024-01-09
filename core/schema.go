package core

import (
	"github.com/cgalvisleon/elvis/jdb"
)

var makeCore bool

func DefineSchemaCore() error {
	var err error
	if makeCore {
		return nil
	}

	makeCore, err = jdb.CreateSchema(0, "core")
	if err != nil {
		return err
	}

	return nil
}
