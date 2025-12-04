package crontab

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
)

type GetInstanceFn func(tag string) (*Job, error)
type SetInstanceFn func(*Job) error
type DeleteInstanceFn func(tag string) error
type QueryInstanceFn func(query string) (et.Items, error)

var getInstance GetInstanceFn
var setInstance SetInstanceFn
var deleteInstance DeleteInstanceFn
var queryInstance QueryInstanceFn

/**
* SetGetInstanceFn
* @param fn GetInstanceFn
**/
func SetGetInstanceFn(fn GetInstanceFn) {
	if fn == nil {
		return
	}

	logs.Log("workflow", "SetGetInstanceFn")
	getInstance = fn
}

/**
* SetSetInstanceFn
* @param fn SetInstanceFn
**/
func SetSetInstanceFn(fn SetInstanceFn) {
	if fn == nil {
		return
	}

	logs.Log("workflow", "SetSetInstanceFn")
	setInstance = fn
}

/**
* SetDeleteInstanceFn
* @param fn DeleteInstanceFn
**/
func SetDeleteInstanceFn(fn DeleteInstanceFn) {
	if fn == nil {
		return
	}

	logs.Log("workflow", "SetDeleteInstanceFn")
	deleteInstance = fn
}

/**
* SetQueryInstanceFn
* @param fn QueryInstanceFn
**/
func SetQueryInstanceFn(fn QueryInstanceFn) {
	if fn == nil {
		return
	}

	logs.Log("workflow", "SetQueryInstanceFn")
	queryInstance = fn
}
