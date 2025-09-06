package flow

import (
	"fmt"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/resilience"
)

type WorkFlows struct {
	Flows      map[string]*Flow
	Instance   map[string]*Flow
	Resilience map[string]*resilience.Attempt
}

var workFlows *WorkFlows

/**
* NewWorkFlows
* @return *WorkFlows
**/
func NewWorkFlows() *WorkFlows {
	return &WorkFlows{
		Flows:      make(map[string]*Flow),
		Instance:   make(map[string]*Flow),
		Resilience: make(map[string]*resilience.Attempt),
	}
}

/**
* NewFlow
* @param tag, version, name, description string, fn FnContext, retries int, retryDelay, retentionTime time.Duration, createdBy string
* @return *Flow, error
**/
func (s *WorkFlows) NewFlow(tag, version, name, description string, fn FnContext, retries int, retryDelay, retentionTime time.Duration, createdBy string) (*Flow, error) {
	flow, err := newFlow(s, tag, version, name, description, fn, retries, retryDelay, retentionTime, createdBy)
	if err != nil {
		return nil, err
	}
	s.Flows[tag] = flow

	return flow, nil
}

/**
* Run
* @param serviceId, tag string, ctx et.Json
* @return et.Item, error
**/
func (s *WorkFlows) Run(serviceId, tag string, ctx et.Json) (et.Item, error) {
	serviceId = reg.GetUUID(serviceId)
	instance, err := s.getInstance(serviceId)
	if err != nil {
		return et.Item{}, err
	}

	if instance == nil {
		instance, err = s.newInstance(serviceId, tag)
		if err != nil {
			return et.Item{}, err
		}
	}

	result, err := instance.run(ctx)
	if err != nil {
		instance.addResilience(ctx)
		return et.Item{}, err
	}

	return result, err
}

/**
* Rollback
* @param serviceId string
* @return et.Item, error
**/

func (s *WorkFlows) Rollback(serviceId string) (et.Item, error) {
	serviceId = reg.GetUUID(serviceId)
	instance, err := s.getInstance(serviceId)
	if err != nil {
		return et.Item{}, err
	}

	if instance == nil {
		return et.Item{}, fmt.Errorf("instance not found")
	}

	result, err := instance.rollback(instance.LastRollback)
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* HealthCheck
* @return bool
**/
func (s *WorkFlows) HealthCheck() bool {
	ok := resilience.HealthCheck()
	if !ok {
		return false
	}

	return true
}
