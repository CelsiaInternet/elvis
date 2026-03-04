package crontab

import (
	"fmt"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
)

var (
	crontab *Jobs
)

/**
* Load
* @params tag string
* @return error
**/
func Load(tag string) error {
	if crontab != nil {
		return nil
	}

	_, err := event.Load()
	if err != nil {
		return err
	}

	crontab = New()
	err = crontab.start()
	if err != nil {
		return err
	}

	tag = strs.Name(tag)
	err = eventInit(tag)
	if err != nil {
		return err
	}

	return nil
}

/**
* Close
* @return void
**/
func Close() {
	if crontab == nil {
		return
	}

	cache.Delete("crontab:nodes")

	logs.Log(packageName, `Disconnect...`)
}

/**
* addJob
* @param jobType TypeJob, tag, spec, channel string, started bool, params et.Json, repetitions int, fn func(event.EvenMessage)
* @return error
**/
func addJob(jobType TypeJob, tag, spec, channel string, started bool, params et.Json, repetitions int, fn func(event.EvenMessage)) error {
	if crontab == nil {
		return fmt.Errorf(MSG_CRONTAB_UNLOAD)
	}

	tag = strs.Name(tag)
	data := et.Json{
		"type":        jobType,
		"tag":         tag,
		"spec":        spec,
		"channel":     channel,
		"started":     started,
		"params":      params,
		"repetitions": repetitions,
	}

	event.Publish(EVENT_CRONTAB_SET, data)
	err := event.Stack(channel, fn)
	if err != nil {
		return err
	}

	logs.Logf(packageName, "Add EventJob: %s", data.ToString())

	return nil
}

/**
* AddEventJob
* @param tag, spec string, repetitions int, started bool, params et.Json, fn func(event.EvenMessage)
* @return error
**/
func AddEventJob(tag, spec string, repetitions int, started bool, params et.Json, fn func(event.EvenMessage)) error {
	tag = strs.Name(tag)
	channel := fmt.Sprintf("cronjob:%s", tag)
	return addJob(CronJob, tag, spec, channel, started, params, repetitions, fn)
}

/**
* AddCronJob
* @param tag, spec string, repetitions int, started bool, params et.Json, fn func(event.EvenMessage)
* @return error
**/
func AddCronJob(tag, spec string, repetitions int, started bool, params et.Json, fn func(event.EvenMessage)) error {
	return AddEventJob(tag, spec, repetitions, started, params, fn)
}

/**
* AddScheduleJob
* Add job to crontab in execute local
* @param tag, schedule string, params et.Json, repetitions int, started bool, fn func(event.EvenMessage)
* @return error
**/
func AddScheduleJob(tag, schedule string, params et.Json, repetitions int, started bool, fn func(event.EvenMessage)) error {
	tag = strs.Name(tag)
	channel := fmt.Sprintf("schedule:%s", tag)
	return addJob(ScheduleJob, tag, schedule, channel, started, params, repetitions, fn)
}

/**
* DeleteJob
* @param tag string
* @return error
**/
func DeleteJob(tag string) error {
	if crontab == nil {
		return fmt.Errorf(MSG_CRONTAB_UNLOAD)
	}

	err := event.Publish(EVENT_CRONTAB_DELETE, et.Json{"tag": tag})
	if err != nil {
		return err
	}

	return nil
}

/**
* StartJob
* @param tag string
* @return int, error
**/
func StartJob(tag string) (int, error) {
	if crontab == nil {
		return 0, fmt.Errorf(MSG_CRONTAB_UNLOAD)
	}

	err := event.Publish(EVENT_CRONTAB_START, et.Json{"tag": tag})
	if err != nil {
		return 0, err
	}

	return 1, nil
}

/**
* StopJob
* @param tag string
* @return error
**/
func StopJob(tag string) error {
	if crontab == nil {
		return fmt.Errorf(MSG_CRONTAB_UNLOAD)
	}

	err := event.Publish(EVENT_CRONTAB_STOP, et.Json{"tag": tag})
	if err != nil {
		return err
	}

	return nil
}

/**
* Stop
* @return error
**/
func Stop() error {
	if crontab == nil {
		return fmt.Errorf(MSG_CRONTAB_UNLOAD)
	}

	return crontab.stop()
}
