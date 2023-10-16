package core

import (
	"github.com/cgalvisleon/elvis/console"
)

var initCore bool

func InitCore() error {
	if initCore {
		return nil
	}

	if err := DefineConfig(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineSync(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineReference(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineMode(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineSeries(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineRecycling(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineCollection(); err != nil {
		return console.PanicE(err)
	}
	if err := JoinToMaster(); err != nil {
		console.Error(err)
	}

	initCore = true

	console.LogK("CORE", "Init core")

	return nil
}
