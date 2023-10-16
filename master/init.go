package master

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/core"
	. "github.com/cgalvisleon/elvis/jdb"
)

func InitMaster() error {
	if master != nil {
		return nil
	}

	master = &Master{}

	if err := DefineConfig(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineNodes(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineSeries(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineReference(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineCollection(); err != nil {
		return console.PanicE(err)
	}

	go Listen("master", Jdb(0).URL, "node", listenNode)

	console.LogK("CORE", "Init Master")

	return nil
}
