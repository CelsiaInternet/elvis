package module

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/core"
)

var initDefine bool

func InitDefine() error {
	if initDefine {
		return nil
	}

	if err := core.InitDefine(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineTypes(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineProjects(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineUsers(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineHistorys(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineTokens(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineModules(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineProjectModules(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineFolders(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineProfiles(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineProfileFolders(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineRoles(); err != nil {
		return console.PanicE(err)
	}
	if err := DefineHistorys(); err != nil {
		return console.PanicE(err)
	}

	console.LogK("Module", "Init module")

	initDefine = true

	return nil
}
