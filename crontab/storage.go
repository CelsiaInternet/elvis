package crontab

import "github.com/celsiainternet/elvis/logs"

type LoadInstanceFn func(id string) (*Job, error)
type SaveInstanceFn func(*Job) error

var loadInstance LoadInstanceFn
var saveInstance SaveInstanceFn

func SetLoadInstance(fn LoadInstanceFn) {
	if fn == nil {
		return
	}

	logs.Log("workflow", "SetLoadInstance")
	loadInstance = fn
}

func SetSaveInstance(fn SaveInstanceFn) {
	if fn == nil {
		return
	}

	logs.Log("workflow", "SetSaveInstance")
	saveInstance = fn
}
