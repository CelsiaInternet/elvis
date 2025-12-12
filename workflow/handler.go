package workflow

import (
	"fmt"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
)

var workFlows *WorkFlows

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
* @param instanceId, tag string, step int, tags et.Json, ctx et.Json, createdBy string
* @return et.Json, error
**/
func Run(instanceId, tag string, step int, tags et.Json, ctx et.Json, createdBy string) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.runInstance(instanceId, tag, step, tags, ctx, createdBy)
}

/**
* Reset
* @param instanceId, updatedBy string
* @return error
**/
func Reset(instanceId, updatedBy string) error {
	if err := Load(); err != nil {
		return err
	}

	return workFlows.resetInstance(instanceId, updatedBy)
}

/**
* Rollback
* @param instanceId, updatedBy string
* @return et.Json, error
**/
func Rollback(instanceId, updatedBy string) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.rollback(instanceId, updatedBy)
}

/**
* Stop
* @param instanceId, updatedBy string
* @return error
**/
func Stop(instanceId, updatedBy string) error {
	if err := Load(); err != nil {
		return err
	}

	return workFlows.stop(instanceId, updatedBy)
}

/**
* Status
* @param instanceId, status, updatedBy string
* @return FlowStatus, error
**/
func Status(instanceId, status, updatedBy string) (FlowStatus, error) {
	if err := Load(); err != nil {
		return "", err
	}

	if _, ok := FlowStatusList[FlowStatus(status)]; !ok {
		return "", fmt.Errorf("status %s no es valido", status)
	}

	instance, exists := workFlows.loadInstance(instanceId)
	if !exists {
		return "", fmt.Errorf("instance not found")
	}

	instance.setStatus(FlowStatus(status))
	return instance.Status, nil
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

/**
* GetInstance
* @param instanceId string
* @return (*Instance, error)
**/
func GetInstance(instanceId string) (*Instance, error) {
	if err := Load(); err != nil {
		return nil, err
	}

	instance, exists := workFlows.loadInstance(instanceId)
	if !exists {
		return nil, fmt.Errorf("instance not found")
	}

	return instance, nil
}
