package core

import "github.com/cgalvisleon/elvis/console"

func InitDefine() error {
	if err := DefineConfig(); err != nil {
		console.Panic(err)
		return err
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
		return console.Error(err)
	}

	console.LogK("CORE", "Init core")

	return nil
}
