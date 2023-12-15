package core

import "github.com/cgalvisleon/elvis/console"

func InitDefine() {
	if err := DefineConfig(); err != nil {
		console.Panic(err)
	}
	if err := DefineSync(); err != nil {
		console.PanicE(err)
	}
	if err := DefineReference(); err != nil {
		console.PanicE(err)
	}
	if err := DefineMode(); err != nil {
		console.PanicE(err)
	}
	if err := DefineSeries(); err != nil {
		console.PanicE(err)
	}
	if err := DefineRecycling(); err != nil {
		console.PanicE(err)
	}
	if err := DefineCollection(); err != nil {
		console.PanicE(err)
	}
	if err := JoinToMaster(); err != nil {
		console.Error(err)
	}

	console.LogK("CORE", "Init core")
}
