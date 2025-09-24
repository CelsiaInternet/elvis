package crontab

import (
	"fmt"

	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
)

var (
	EVENT_CRONTAB_SET    = "event:crontab:set"
	EVENT_CRONTAB_STATUS = "event:crontab:status"
	EVENT_CRONTAB_STOP   = "event:crontab:stop"
	EVENT_CRONTAB_START  = "event:crontab:start"
	EVENT_CRONTAB_DELETE = "event:crontab:delete"
)

/**
* eventInit
* @return error
**/
func eventInit() error {
	err := event.Subscribe(EVENT_CRONTAB_SET, eventSet)
	if err != nil {
		return err
	}

	err = event.Subscribe(EVENT_CRONTAB_DELETE, eventDelete)
	if err != nil {
		return err
	}

	err = event.Subscribe(EVENT_CRONTAB_STOP, eventStop)
	if err != nil {
		return err
	}

	err = event.Subscribe(EVENT_CRONTAB_START, eventStart)
	if err != nil {
		return err
	}

	return nil
}

/**
* eventSet
* @param msg event.EvenMessage
* @return error
**/
func eventSet(msg event.EvenMessage) {
	if crontab == nil {
		return
	}

	if !crontab.isServer {
		return
	}

	data := msg.Data
	id := data.Str("id")
	name := data.Str("name")
	spec := data.Str("spec")
	channel := data.Str("channel")
	params := data.Json("params")
	_, err := crontab.addEventJob(id, name, spec, channel, params, true)
	if err != nil {
		logs.Logf(packageName, fmt.Sprintf("Error adding job %s", err))
		return
	}

	logs.Logf(packageName, fmt.Sprintf("Crontab %s added", name))
}

/**
* eventDelete
* @param msg event.EvenMessage
* @return error
**/
func eventDelete(msg event.EvenMessage) {
	if crontab == nil {
		return
	}

	if !crontab.isServer {
		return
	}

	data := msg.Data
	id := data.Str("id")
	err := crontab.deleteJobById(id)
	if err != nil {
		logs.Logf(packageName, fmt.Sprintf("Error deleting job %s", err))
		return
	}

	logs.Logf(packageName, fmt.Sprintf("Crontab %s deleted", id))
}

/**
* eventStop
* @param msg event.EvenMessage
* @return error
**/
func eventStop(msg event.EvenMessage) {
	if crontab == nil {
		return
	}

	if !crontab.isServer {
		return
	}

	data := msg.Data
	id := data.Str("id")
	err := crontab.stopJobById(id)
	if err != nil {
		logs.Logf(packageName, fmt.Sprintf("Error stopping job %s", err))
		return
	}

	logs.Logf(packageName, fmt.Sprintf("Crontab %s stopped", id))
}

/**
* eventStart
* @param msg event.EvenMessage
* @return error
**/
func eventStart(msg event.EvenMessage) {
	if crontab == nil {
		return
	}

	if !crontab.isServer {
		return
	}

	data := msg.Data
	id := data.Str("id")
	err := crontab.startJobById(id)
	if err != nil {
		logs.Logf(packageName, fmt.Sprintf("Error starting job %s", err))
		return
	}

	logs.Logf(packageName, fmt.Sprintf("Crontab %s started", id))
}
