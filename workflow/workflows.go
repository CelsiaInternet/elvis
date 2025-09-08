package workflow

import (
	"fmt"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/resilience"
)

const packageName = "workflow"

type WorkFlows struct {
	Flows    map[string]*Flow `json:"flows"`
	Instance map[string]*Flow `json:"instance"`
}

var workFlows *WorkFlows

/**
* NewWorkFlows
* @return *WorkFlows
**/
func NewWorkFlows() *WorkFlows {
	return &WorkFlows{
		Flows:    make(map[string]*Flow),
		Instance: make(map[string]*Flow),
	}
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

/**
* NewFlow
* @param tag, version, name, description string, fn FnContext, createdBy string
* @return *Flow
**/
func (s *WorkFlows) NewFlow(tag, version, name, description string, fn FnContext, createdBy string) *Flow {
	flow := newFlow(s, tag, version, name, description, fn, createdBy)
	s.Flows[tag] = flow

	return flow
}

/**
* Run
* @param instanceId, tag string, ctx et.Json
* @return et.Json, error
**/
func (s *WorkFlows) Run(instanceId, tag string, startId int, ctx et.Json) (et.Json, error) {
	instanceId = reg.GetUUID(instanceId)
	instance, err := s.getInstance(instanceId)
	if err != nil {
		return et.Json{}, err
	}

	if instance == nil {
		instance, err = s.newInstance(instanceId, tag, startId)
		if err != nil {
			return et.Json{}, err
		}
	}

	result, err := instance.run(ctx)
	if err != nil {
		return et.Json{}, err
	}

	return result, err
}

/**
* Rollback
* @param instanceId string
* @return et.Json, error
**/

func (s *WorkFlows) Rollback(instanceId string) (et.Json, error) {
	instanceId = reg.GetUUID(instanceId)
	instance, err := s.getInstance(instanceId)
	if err != nil {
		return et.Json{}, err
	}

	if instance == nil {
		return et.Json{}, fmt.Errorf("instance not found")
	}

	result, err := instance.rollback(instance.LastRollback)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* DeleteFlow
* @param tag string
* @return bool
**/
func (s *WorkFlows) DeleteFlow(tag string) bool {
	if s.Flows[tag] == nil {
		return false
	}

	flow := s.Flows[tag]
	delete(s.Flows, tag)
	event.Publish(EVENT_WORKFLOW_DELETE, flow.ToJson())

	return true
}
