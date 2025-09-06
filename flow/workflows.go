package flow

import (
	"fmt"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/resilience"
)

type WorkFlows struct {
	Flows         map[string]*Flow               `json:"flows"`
	Instance      map[string]*Flow               `json:"instance"`
	Resilience    map[string]*resilience.Attempt `json:"resilience"`
	InstanceAtrib string                         `json:"instance_atrib"`
}

var workFlows *WorkFlows

/**
* NewWorkFlows
* @return *WorkFlows
**/
func NewWorkFlows() *WorkFlows {
	return &WorkFlows{
		Flows:         make(map[string]*Flow),
		Instance:      make(map[string]*Flow),
		Resilience:    make(map[string]*resilience.Attempt),
		InstanceAtrib: "instance_id",
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

func (s *WorkFlows) SetInstanceAtrib(instanceAtrib string) {
	s.InstanceAtrib = instanceAtrib
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
* @param instanceId, tag string, ctx et.Json
* @return et.Item, error
**/
func (s *WorkFlows) Run(instanceId, tag string, startId int, ctx et.Json) (et.Item, error) {
	instanceId = reg.GetUUID(instanceId)
	instance, err := s.getInstance(instanceId)
	if err != nil {
		return et.Item{}, err
	}

	if instance == nil {
		instance, err = s.newInstance(instanceId, tag, startId)
		if err != nil {
			return et.Item{}, err
		}
	}

	ctx.Set(s.InstanceAtrib, instanceId)
	result, err := instance.run(ctx)
	if err != nil {
		instance.addResilience(ctx)
		return et.Item{}, err
	}

	return result, err
}

/**
* Rollback
* @param instanceId string
* @return et.Item, error
**/

func (s *WorkFlows) Rollback(instanceId string) (et.Item, error) {
	instanceId = reg.GetUUID(instanceId)
	instance, err := s.getInstance(instanceId)
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
