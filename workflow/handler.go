package workflow

import (
	"fmt"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
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

	console.Debug("Run", et.Json{
		"instanceId": instanceId,
		"tag":        tag,
		"step":       step,
		"tags":       tags,
		"ctx":        ctx,
		"createdBy":  createdBy,
	}.ToString())

	return workFlows.run(instanceId, tag, step, tags, ctx, createdBy)
}

/**
* Reset
* @param instanceId string
* @return error
**/
func Reset(instanceId string) error {
	if err := Load(); err != nil {
		return err
	}

	return workFlows.reset(instanceId)
}

/**
* Rollback
* @param instanceId string
* @return et.Json, error
**/
func Rollback(instanceId, tag string) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.rollback(instanceId, tag)
}

/**
* Status
* @param instanceId, status string
* @return FlowStatus, error
**/
func Status(instanceId, status string) (FlowStatus, error) {
	if err := Load(); err != nil {
		return "", err
	}

	if _, ok := FlowStatusList[FlowStatus(status)]; !ok {
		return "", fmt.Errorf("status %s no es valido", status)
	}

	instance, err := workFlows.getInstance(instanceId)
	if err != nil {
		return "", err
	}

	instance.setStatus(FlowStatus(status))
	return instance.Status, nil
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

/**
* GetInstance
* @param instanceId string
* @return (*Instance, error)
**/
func GetInstance(instanceId string) (*Instance, error) {
	if err := Load(); err != nil {
		return nil, err
	}

	return workFlows.getInstance(instanceId)
}
