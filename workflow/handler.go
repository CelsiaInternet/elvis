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

	workFlows = newWorkFlows()
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

	return workFlows.healthCheck()
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

	return workFlows.newFlow(tag, version, name, description, fn, stop, createdBy)
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

	return workFlows.run(instanceId, tag, startId, tags, ctx)
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

	return workFlows.rollback(instanceId, tag)
}

/**
* Stop
* @param instanceId, tag string
* @return error
**/
func Stop(instanceId, tag string) error {
	if err := Load(); err != nil {
		return err
	}

	return workFlows.stop(instanceId, tag)
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

	return workFlows.deleteFlow(tag), nil
}
