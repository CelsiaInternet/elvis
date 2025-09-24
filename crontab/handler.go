package crontab

import (
	"fmt"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
)

var crontab *Jobs

/**
* Load
**/
func Load() error {
	if crontab != nil {
		return nil
	}

	crontab = New()
	err := crontab.load()
	if err != nil {
		return err
	}

	return crontab.start()
}

/**
* Server
* @return error
**/
func Server() error {
	if crontab != nil {
		return nil
	}

	crontab = New()
	crontab.isServer = true
	err := crontab.load()
	if err != nil {
		return err
	}

	err = crontab.start()
	if err != nil {
		return err
	}

	err = eventInit()
	if err != nil {
		return err
	}

	return nil
}

/**
* AddJob
* @param id, name, spec, channel string, params et.Json, fn func()
* @return *Job, error
**/
func AddJob(id, name, spec, channel string, params et.Json, fn func(job *Job)) (*Job, error) {
	err := Load()
	if err != nil {
		return nil, err
	}

	if crontab.isServer {
		return nil, fmt.Errorf("crontab is server")
	}

	result, err := crontab.addJob(id, name, spec, channel, params, fn)
	if err != nil {
		return nil, err
	}

	err = result.Start()
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* AddEventJob
* @param id, name, spec, channel string, params et.Json
* @return *Job, error
**/
func AddEventJob(id, name, spec, channel string, params et.Json) (*Job, error) {
	err := Server()
	if err != nil {
		return nil, err
	}

	return crontab.addEventJob(id, name, spec, channel, params)
}

/**
* EventJob
* @param id, name, spec, channel string, params et.Json
* @return *Job, error
**/
func EventJob(id, name, spec, channel string, params et.Json, fn func(event.EvenMessage)) error {
	event.Publish(EVENT_CRONTAB_SET, et.Json{
		"id":      id,
		"name":    name,
		"spec":    spec,
		"channel": channel,
		"params":  params,
	})

	err := event.Stack(channel, fn)
	if err != nil {
		return err
	}

	return nil
}

/**
* DeleteJob
* @param name string
* @return error
**/
func DeleteJob(name string) error {
	err := Load()
	if err != nil {
		return err
	}

	return crontab.deleteJobByName(name)
}

/**
* DeleteJobById
* @param id string
* @return error
**/
func DeleteJobById(id string) error {
	err := Load()
	if err != nil {
		return err
	}

	return crontab.deleteJobById(id)
}

/**
* StartJob
* @param name string
* @return int, error
**/
func StartJob(name string) (int, error) {
	err := Load()
	if err != nil {
		return 0, err
	}

	return crontab.startJobByName(name)
}

/**
* StartJobById
* @param id string
* @return error
**/
func StartJobById(id string) error {
	err := Load()
	if err != nil {
		return err
	}

	return crontab.startJobById(id)
}

/**
* StopJob
* @param name string
* @return error
**/
func StopJob(name string) error {
	err := Load()
	if err != nil {
		return err
	}

	return crontab.stopJobByName(name)
}

/**
* StopJobById
* @param id string
* @return error
**/
func StopJobById(id string) error {
	err := Load()
	if err != nil {
		return err
	}

	return crontab.stopJobById(id)
}

/**
* ListJobs
* @return et.Items, error
**/
func ListJobs() (et.Items, error) {
	err := Load()
	if err != nil {
		return et.Items{}, err
	}

	return crontab.list(), nil
}

/**
* Start
* @return error
**/
func Start() error {
	err := Load()
	if err != nil {
		return err
	}

	return crontab.start()
}

/**
* Stop
* @return error
**/
func Stop() error {
	err := Load()
	if err != nil {
		return err
	}

	return crontab.stop()
}

/**
* EventStatusRunning
* @param data et.Json
* @return error
**/
func EventStatusRunning(data et.Json) error {
	data.Set("status", StatusRunning)
	return event.Publish(EVENT_CRONTAB_STATUS, data)
}

/**
* EventStatusPending
* @param data et.Json
* @return error
**/
func EventStatusDone(data et.Json) error {
	data.Set("status", StatusDone)
	return event.Publish(EVENT_CRONTAB_STATUS, data)
}

/**
* EventStatusFailed
* @param data et.Json
* @return error
**/
func EventStatusFailed(data et.Json) error {
	data.Set("status", StatusFailed)
	return event.Publish(EVENT_CRONTAB_STATUS, data)
}
