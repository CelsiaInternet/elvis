package workflow

import (
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
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
* @param tag, version, name, description string, fn FnContext, stop bool, createdBy string
* @return *Flow
**/
func (s *WorkFlows) NewFlow(tag, version, name, description string, fn FnContext, stop bool, createdBy string) *Flow {
	flow := newFlow(s, tag, version, name, description, fn, stop, createdBy)
	s.Flows[tag] = flow

	return flow
}

/**
* Start
* @param instanceId, tag string, startId int, tags et.Json, ctx et.Json
* @return et.Json, error
**/
func (s *WorkFlows) Start(instanceId, tag string, startId int, tags, ctx et.Json) (et.Json, error) {
	instance, err := s.createInstance(instanceId, tag, startId, tags)
	if err != nil {
		return et.Json{}, err
	}

	result, err := instance.run(ctx)
	if err != nil {
		return et.Json{}, err
	}

	return result, err
}

/**
* Run
* @param instanceId, tag string, startId int, tags, ctx et.Json
* @return et.Json, error
**/
func (s *WorkFlows) Run(instanceId, tag string, startId int, tags, ctx et.Json) (et.Json, error) {
	instance, err := s.getOrCreateInstance(instanceId, tag, tags)
	if err != nil {
		return et.Json{}, err
	}

	result, err := instance.Run(startId, ctx)
	if err != nil {
		return et.Json{}, err
	}

	if instance.isDebug {
		console.Debug("Flow instance:", instance.ToJson().ToString())
	}

	return result, err
}

/**
* Continue
* @param instanceId, tag string, ctx et.Json
* @return et.Json, error
**/
func (s *WorkFlows) Continue(instanceId, tag string, ctx et.Json) (et.Json, error) {
	instance, err := s.getInstance(instanceId, tag)
	if err != nil {
		return et.Json{}, err
	}

	result, err := instance.Continue(ctx)
	if err != nil {
		return et.Json{}, err
	}

	return result, err
}

/**
* Rollback
* @param instanceId, tag string
* @return et.Json, error
**/

func (s *WorkFlows) Rollback(instanceId, tag string) (et.Json, error) {
	instance, err := s.getInstance(instanceId, tag)
	if err != nil {
		return et.Json{}, err
	}

	result, err := instance.rollback(instance.LastRollback)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* Done
* @param instanceId string
* @return bool
**/
func (s *WorkFlows) Done(instanceId string) bool {
	if s.Instance[instanceId] == nil {
		return false
	}

	time.AfterFunc(300*time.Millisecond, func() {
		delete(s.Instance, instanceId)
		logs.Logf(packageName, MSG_WORKFLOW_DONE_INSTANCE, instanceId)
	})

	return true
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
	event.Publish(EVENT_WORKFLOW_DELETE, flow.ToJson())
	time.AfterFunc(300*time.Millisecond, func() {
		delete(s.Flows, tag)
	})

	return true
}
