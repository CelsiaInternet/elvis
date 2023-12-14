package module

import "github.com/cgalvisleon/elvis/console"

func init() {
	if err := DefineTypes(); err != nil {
		console.PanicE(err)
	}
	if err := DefineProjects(); err != nil {
		console.PanicE(err)
	}
	if err := DefineUsers(); err != nil {
		console.PanicE(err)
	}
	if err := DefineHistorys(); err != nil {
		console.PanicE(err)
	}
	if err := DefineTokens(); err != nil {
		console.PanicE(err)
	}
	if err := DefineModules(); err != nil {
		console.PanicE(err)
	}
	if err := DefineProjectModules(); err != nil {
		console.PanicE(err)
	}
	if err := DefineFolders(); err != nil {
		console.PanicE(err)
	}
	if err := DefineProfiles(); err != nil {
		console.PanicE(err)
	}
	if err := DefineProfileFolders(); err != nil {
		console.PanicE(err)
	}
	if err := DefineRoles(); err != nil {
		console.PanicE(err)
	}
	if err := DefineHistorys(); err != nil {
		console.PanicE(err)
	}

	console.LogK("Module", "Init module")
}
