package workflow

import (
	"fmt"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/resilience"
	"github.com/celsiainternet/elvis/timezone"
)

var (
	packageName           = "workflow"
	ErrorInstanceNotFound = fmt.Errorf(MSG_INSTANCE_NOT_FOUND)
)

type WorkFlows struct {
	Flows     map[string]*Flow     `json:"flows"`
	Instances map[string]*Instance `json:"instances"`
	Results   map[string]et.Json   `json:"results"`
	mu        sync.Mutex           `json:"-"`
	isDebug   bool                 `json:"-"`
}

/**
* newWorkFlows
* @return *WorkFlows
**/
func newWorkFlows() *WorkFlows {
	result := &WorkFlows{
		Flows:     make(map[string]*Flow),
		Instances: make(map[string]*Instance),
		Results:   make(map[string]et.Json),
		mu:        sync.Mutex{},
		isDebug:   envar.GetBool(false, "DEBUG"),
	}

	return result
}

/**
* healthCheck
* @return bool
**/
func (s *WorkFlows) healthCheck() bool {
	ok := resilience.HealthCheck()
	if !ok {
		return false
	}

	return true
}

/**
* Add
* @param instance *Instance
**/
func (s *WorkFlows) Add(instance *Instance) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Instances[instance.Id] = instance
	timeNow := timezone.NowTime()
	second := timeNow.Format("2006-01-02-15:04:05")
	cache.Incr(cache.GenKey("workflow:instance", second), 2*time.Second)
}

/**
* Remove
* @param instance *Instance
**/
func (s *WorkFlows) Remove(instance *Instance) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Instances, instance.Id)
}

/**
* Count
* @return int
**/
func (s *WorkFlows) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.Instances)
}

/**
* newInstance
* @param tag, id string, tags et.Json, step int, createdBy string
* @return *Instance, error
**/
func (s *WorkFlows) newInstance(tag, id string, tags et.Json, step int, createdBy string) (*Instance, error) {
	if id == "" {
		return nil, fmt.Errorf(MSG_INSTANCE_ID_REQUIRED)
	}

	flow := s.Flows[tag]
	if flow == nil {
		return nil, fmt.Errorf(MSG_FLOW_NOT_FOUND)
	}

	if s.isDebug {
		logs.Debug("newInstance:1")
	}

	if step == -1 {
		step = 0
	}

	now := timezone.NowTime()
	result := &Instance{
		Flow:       flow,
		workFlows:  s,
		Tag:        tag,
		CreatedAt:  now,
		UpdatedAt:  now,
		Id:         id,
		CreatedBy:  createdBy,
		UpdatedBy:  createdBy,
		Current:    step,
		Ctx:        et.Json{},
		Ctxs:       make(map[int]et.Json),
		Results:    make(map[int]*Result),
		Rollbacks:  make(map[int]*Result),
		Params:     et.Json{},
		Traces:     []et.Json{},
		Tags:       tags,
		WorkerHost: workerHost,
		goTo:       -1,
		isNew:      true,
	}

	if s.isDebug {
		logs.Debugf("newInstance:2 tag:%s id:%s", tag, id)
	}

	return result, result.setStatus(FlowStatusPending)
}

/**
* Debug
* @return *WorkFlows
**/
func (s *WorkFlows) Debug() *WorkFlows {
	s.isDebug = true
	return s
}

/**
* loadInstance
* @param id, tag string
* @return bool, *Instance, error
**/
func (s *WorkFlows) loadInstance(id, tag string) (*Instance, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	dest := s.Instances[id]
	if dest != nil {
		return dest, nil
	}

	if getInstance != nil {
		var dest *Instance
		exists, err := getInstance(id, &dest)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, nil
		}

		if dest == nil {
			return nil, fmt.Errorf("instance not loaded")
		}

		if tag != "" && tag != dest.Tag {
			dest.Tag = tag
		}

		flow := s.Flows[dest.Tag]
		if flow == nil {
			return nil, fmt.Errorf("flow not found")
		}

		dest.Flow = flow
		dest.goTo = -1
		s.Add(dest)

		if s.isDebug {
			logs.Debugf("loadInstance:4 instance:%s", dest.ToString())
		}

		return dest, nil
	}

	return nil, nil
}

/**
* getOrCreateInstance
* @param id, tag string, step int, tags et.Json, createdBy string
* @return *Instance, error
**/
func (s *WorkFlows) getOrCreateInstance(id, tag string, step int, tags et.Json, createdBy string) (*Instance, error) {
	id = reg.GetUUID(id)
	result, err := s.loadInstance(id, tag)
	if err != nil {
		return nil, err
	} else if result == nil {
		return s.newInstance(tag, id, tags, step, createdBy)
	}

	return result, nil
}

/**
* runInstance
* Si el step es -1 se ejecuta el siguiente paso, si no se ejecuta el paso indicado
* @param instanceId, tag string, step int, tags, ctx et.Json, createdBy string
* @return et.Json, error
**/
func (s *WorkFlows) runInstance(instanceId, tag string, step int, tags, ctx et.Json, createdBy string) (et.Json, error) {
	if s.isDebug {
		logs.Debug("runInstance:1")
	}

	instance, err := s.getOrCreateInstance(instanceId, tag, step, tags, createdBy)
	if err != nil {
		return et.Json{}, err
	}

	if s.isDebug {
		logs.Debug("runInstance:2")
	}

	instance.isDebug = s.isDebug
	instance.UpdatedBy = createdBy
	instance.PutTag(tags)
	if step != instance.Current {
		instance.Current = step
	}
	result, err := instance.run(ctx)
	if err != nil {
		return et.Json{}, err
	}

	s.Remove(instance)
	logs.Logf(packageName, "runInstance:%s tag:%s", instanceId, tag)
	if s.isDebug {
		logs.Debugf("runInstance:3 instance:%s", instance.ToString())
	}

	return result, err
}

/**
* resetInstance
* @param instanceId string
* @return error
**/
func (s *WorkFlows) resetInstance(instanceId, updatedBy string) error {
	inst, err := s.loadInstance(instanceId, "")
	if err != nil {
		return err
	}

	inst.UpdatedBy = updatedBy
	inst.setStatus(FlowStatusPending)
	return nil
}

/**
* Rollback
* @param instanceId string, updatedBy string
* @return et.Json, error
**/
func (s *WorkFlows) rollback(instanceId, updatedBy string) (et.Json, error) {
	inst, err := s.loadInstance(instanceId, "")
	if err != nil {
		return et.Json{}, err
	}

	inst.UpdatedBy = updatedBy
	result, err := inst.rollback(et.Json{}, nil)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* stop
* @param instanceId string, updatedBy string
* @return error
**/
func (s *WorkFlows) stop(instanceId, updatedBy string) error {
	inst, err := s.loadInstance(instanceId, "")
	if err != nil {
		return err
	}

	inst.UpdatedBy = updatedBy
	return inst.Stop()
}

/**
* newFlow
* @param tag, version, name, description string, fn FnContext, stop bool, createdBy string
* @return *Flow
**/
func (s *WorkFlows) newFlow(tag, version, name, description string, fn FnContext, stop bool, createdBy string) *Flow {
	flow := newFlow(tag, version, name, description, fn, stop, createdBy)
	s.Flows[tag] = flow

	return flow
}

/**
* deleteFlow
* @param tag string
* @return bool
**/
func (s *WorkFlows) deleteFlow(tag string) bool {
	if s.Flows[tag] == nil {
		return false
	}

	flow := s.Flows[tag]
	event.Publish(EVENT_WORKFLOW_DELETE, flow.ToJson())
	delete(s.Flows, tag)

	return true
}
