package crontab

import (
	"fmt"

	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
)

var (
	EVENT_CRONTAB_SET    = "event:crontab:set"
	EVENT_CRONTAB_REMOVE = "event:crontab:remove"
	EVENT_CRONTAB_STOP   = "event:crontab:stop"
	EVENT_CRONTAB_START  = "event:crontab:start"
)

/**
* eventInit
* @return error
**/
func (s *Jobs) eventInit() error {
	EVENT_CRONTAB_SET = fmt.Sprintf("event:crontab:set:%s", s.Tag)
	EVENT_CRONTAB_REMOVE = fmt.Sprintf("event:crontab:remove:%s", s.Tag)
	EVENT_CRONTAB_STOP = fmt.Sprintf("event:crontab:stop:%s", s.Tag)
	EVENT_CRONTAB_START = fmt.Sprintf("event:crontab:start:%s", s.Tag)

	err := event.Stack(EVENT_CRONTAB_SET, s.eventSet)
	if err != nil {
		return err
	}

	err = event.Subscribe(EVENT_CRONTAB_REMOVE, s.eventRemove)
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
		logs.Logf(packageName, fmt.Sprintf("error adding job: %s:%s; %s", tpStr, tag, err))
		return
	}
}

/**
* eventRemove
* @param msg event.EvenMessage
* @return error
**/
func (s *Jobs) eventRemove(msg event.EvenMessage) {
	data := msg.Data
	tag := data.Str("tag")
	s.removeJob(tag)
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
		logs.Logf(packageName, fmt.Sprintf("job:%s; error stopping job %s", tag, err))
		return
	}
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
		logs.Logf(packageName, fmt.Sprintf("job:%s; error starting job %s", tag, err))
		return
	}
}
