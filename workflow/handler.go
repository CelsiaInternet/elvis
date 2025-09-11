package workflow

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
)

/**
* Load
* @return error
 */
func Load() error {
	if workFlows != nil {
		return nil
	}

	_, err := cache.Load()
	if err != nil {
		return err
	}

	_, err = event.Load()
	if err != nil {
		return err
	}

	workFlows = NewWorkFlows()
	return nil
}

/**
* HealthCheck
* @return bool
**/
func HealthCheck() bool {
	if err := Load(); err != nil {
		return false
	}

	return workFlows.HealthCheck()
}

/**
* New
* @param tag, version, name, description string, fn FnContext, createdBy string
* @return *Flow
**/
func New(tag, version, name, description string, fn FnContext, stop bool, createdBy string) *Flow {
	if err := Load(); err != nil {
		return nil
	}

	return workFlows.NewFlow(tag, version, name, description, fn, stop, createdBy)
}

/**
* Start
* @param instanceId, tag string, startId int, tags et.Json, ctx et.Json
* @return et.Json, error
**/
func Start(instanceId, tag string, startId int, tags et.Json, ctx et.Json) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.Start(instanceId, tag, startId, tags, ctx)
}

/**
* Run
* @param instanceId, tag string, startId int, tags et.Json, ctx et.Json
* @return et.Json, error
**/
func Run(instanceId, tag string, startId int, tags et.Json, ctx et.Json) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.Run(instanceId, tag, startId, tags, ctx)
}

/**
* Continue
* @param instanceId, tag string, ctx et.Json
* @return et.Json, error
**/
func Continue(instanceId, tag string, ctx et.Json) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	result, err := workFlows.Continue(instanceId, tag, ctx)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* Rollback
* @param instanceId, tag string
* @return et.Json, error
**/
func Rollback(instanceId, tag string) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.Rollback(instanceId, tag)
}

/**
* DeleteFlow
* @param tag string
* @return (bool, error)
**/
func DeleteFlow(tag string) (bool, error) {
	if err := Load(); err != nil {
		return false, err
	}

	return workFlows.DeleteFlow(tag), nil
}
