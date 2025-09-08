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
func New(tag, version, name, description string, fn FnContext, createdBy string) *Flow {
	if err := Load(); err != nil {
		return nil
	}

	return workFlows.NewFlow(tag, version, name, description, fn, createdBy)
}

/**
* Run
* @param instanceId, tag string, ctx et.Json
* @return et.Json, error
**/
func Run(instanceId, tag string, startId int, ctx et.Json) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.Run(instanceId, tag, startId, ctx)
}

/**
* Rollback
* @param instanceId string
* @return et.Json, error
**/
func Rollback(instanceId string) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.Rollback(instanceId)
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
