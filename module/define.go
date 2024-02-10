package module

import (
	"github.com/cgalvisleon/elvis/console"
)

var initDefine bool

func InitDefine() error {
	if initDefine {
		return nil
	}

	if err := DefineTypes(); err != nil {
		return console.Panic(err)
	}
	if err := DefineProjects(); err != nil {
		return console.Panic(err)
	}
	if err := DefineUsers(); err != nil {
		return console.Panic(err)
	}
	if err := DefineTokens(); err != nil {
		return console.Panic(err)
	}
	if err := DefineModules(); err != nil {
		return console.Panic(err)
	}
	if err := DefineProjectModules(); err != nil {
		return console.Panic(err)
	}
	if err := DefineFolders(); err != nil {
		return console.Panic(err)
	}
	if err := DefineProfiles(); err != nil {
		return console.Panic(err)
	}
	if err := DefineProfileFolders(); err != nil {
		return console.Panic(err)
	}
	if err := DefineRoles(); err != nil {
		return console.Panic(err)
	}

	console.LogK("Module", "Init module")

	initDefine = true

	return nil
}
