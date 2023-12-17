package core

import "github.com/cgalvisleon/elvis/console"

func InitDefine() error {
	if err := DefineConfig(); err != nil {
		return err
	}
	if err := DefineSync(); err != nil {
		return err
	}
	if err := DefineReference(); err != nil {
		return err
	}
	if err := DefineMode(); err != nil {
		return err
	}
	if err := DefineSeries(); err != nil {
		return err
	}
	if err := DefineRecycling(); err != nil {
		return err
	}
	if err := JoinToMaster(); err != nil {
		return console.Error(err)
	}

	console.LogK("CORE", "Init core")

	return nil
}
