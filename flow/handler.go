package flow

import (
	"fmt"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/resilience"
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

	err = resilience.Load()
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
	if workFlows == nil {
		return false
	}

	return workFlows.HealthCheck()
}

/**
* SetInstanceAtrib
* @param instanceAtrib string
**/
func SetInstanceAtrib(instanceAtrib string) {
	if workFlows == nil {
		return
	}

	workFlows.SetInstanceAtrib(instanceAtrib)
}

/**
* NewFlow
* @param tag, version, name, description string, fn FnContext, retries int, retryDelay, retentionTime time.Duration, createdBy string
* @return *Flow, error
**/
func NewFlow(tag, version, name, description string, fn FnContext, retries int, retryDelay, retentionTime time.Duration, createdBy string) (*Flow, error) {
	if workFlows == nil {
		return nil, fmt.Errorf("workFlows is nil")
	}

	return workFlows.NewFlow(tag, version, name, description, fn, retries, retryDelay, retentionTime, createdBy)
}

/**
* Run
* @param instanceId, tag string, ctx et.Json
* @return et.Item, error
**/
func Run(instanceId, tag string, startId int, ctx et.Json) (et.Item, error) {
	if workFlows == nil {
		return et.Item{}, fmt.Errorf("workFlows is nil")
	}

	return workFlows.Run(instanceId, tag, startId, ctx)
}

/**
* Rollback
* @param instanceId string
* @return et.Item, error
**/
func Rollback(instanceId string) (et.Item, error) {
	if workFlows == nil {
		return et.Item{}, fmt.Errorf("workFlows is nil")
	}

	return workFlows.Rollback(instanceId)
}
