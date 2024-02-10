package core

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/jdb"
)

var makedCore bool

func defineSchemaCore() error {
	var err error
	if makedCore {
		return nil
	}

	makedCore, err = jdb.CreateSchema(0, "core")
	if err != nil {
		return err
	}

	console.LogK("CORE", "Init core")

	return nil
}
