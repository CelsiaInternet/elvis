package crontab

import (
	"fmt"

	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
)

var (
	EVENT_CRONTAB_SET    = "event:crontab:set"
	EVENT_CRONTAB_DELETE = "event:crontab:delete"
	EVENT_CRONTAB_STOP   = "event:crontab:stop"
	EVENT_CRONTAB_START  = "event:crontab:start"
)

/**
* eventInit
* @return error
**/
func (s *Jobs) eventInit() error {
	EVENT_CRONTAB_SET = fmt.Sprintf("event:crontab:set:%s", s.Tag)
	EVENT_CRONTAB_DELETE = fmt.Sprintf("event:crontab:delete:%s", s.Tag)
	EVENT_CRONTAB_STOP = fmt.Sprintf("event:crontab:stop:%s", s.Tag)
	EVENT_CRONTAB_START = fmt.Sprintf("event:crontab:start:%s", s.Tag)

	err := event.Stack(EVENT_CRONTAB_SET, s.eventSet)
	if err != nil {
		return err
	}

	err = event.Subscribe(EVENT_CRONTAB_DELETE, s.eventDelete)
	if err != nil {
		return err
	}

	err = event.Subscribe(EVENT_CRONTAB_STOP, s.eventStop)
	if err != nil {
		return err
	}

	err = event.Subscribe(EVENT_CRONTAB_START, s.eventStart)
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
func (s *Jobs) eventSet(msg event.EvenMessage) {
	data := msg.Data
	tpStr := data.Str("type")
	tag := data.Str("tag")
	spec := data.Str("spec")
	channel := data.Str("channel")
	started := data.Bool("started")
	params := data.Json("params")
	repetitions := data.Int("repetitions")
	tp := TypeJob(tpStr)
	_, err := s.addJob(tp, tag, spec, channel, started, params, repetitions)
	if err != nil {
		logs.Logf(packageName, fmt.Sprintf("%s: %s; Error adding job %s", tpStr, tag, err))
		return
	}

	logs.Logf(packageName, fmt.Sprintf("%s: %s added spec %s", tpStr, tag, spec))
}

/**
* eventDelete
* @param msg event.EvenMessage
* @return error
**/
func (s *Jobs) eventDelete(msg event.EvenMessage) {
	data := msg.Data
	tag := data.Str("tag")
	err := s.removeJob(tag)
	if err != nil {
		logs.Logf(packageName, fmt.Sprintf("Crontab %s; Error removing job %s", tag, err))
		return
	}

	logs.Logf(packageName, fmt.Sprintf("Crontab %s removed", tag))
}

/**
* eventStop
* @param msg event.EvenMessage
* @return error
**/
func (s *Jobs) eventStop(msg event.EvenMessage) {
	data := msg.Data
	tag := data.Str("tag")
	err := s.stopJob(tag)
	if err != nil {
		logs.Logf(packageName, fmt.Sprintf("Crontab %s; Error stopping job %s", tag, err))
		return
	}

	logs.Logf(packageName, fmt.Sprintf("Crontab %s stopped", tag))
}

/**
* eventStart
* @param msg event.EvenMessage
* @return error
**/
func (s *Jobs) eventStart(msg event.EvenMessage) {
	data := msg.Data
	tag := data.Str("tag")
	err := s.startJob(tag)
	if err != nil {
		logs.Logf(packageName, fmt.Sprintf("Crontab %s; Error starting job %s", tag, err))
		return
	}

	logs.Logf(packageName, fmt.Sprintf("Crontab %s started", tag))
}
