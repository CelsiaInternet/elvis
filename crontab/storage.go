package crontab

import "github.com/celsiainternet/elvis/logs"

type LoadInstanceFn func(id string, dest any) (bool, error)
type SaveInstanceFn func(id, tag string, obj any) error

var loadInstance LoadInstanceFn
var saveInstance SaveInstanceFn

func SetLoadInstance(fn LoadInstanceFn) {
	if fn == nil {
		return
	}

	logs.Log(packageName, "SetLoadInstance")
	loadInstance = fn
}

func SetSaveInstance(fn SaveInstanceFn) {
	if fn == nil {
		return
	}

	logs.Log(packageName, "SetSaveInstance")
	saveInstance = fn
}
