package workflow

import (
	"errors"
	"fmt"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/resilience"
)

const packageName = "workflow"

type instanceFn func(instanceId, tag string, startId int, tags, ctx et.Json) (et.Json, error)

type Awaiting struct {
	CreatedAt  time.Time     `json:"created_at"`
	ExecutedAt time.Time     `json:"executed_at"`
	Id         string        `json:"id"`
	Tag        string        `json:"tag"`
	fn         instanceFn    `json:"-"`
	fnArgs     []interface{} `json:"-"`
}

func (s *Awaiting) ToJson() et.Json {
	return et.Json{
		"created_at":  s.CreatedAt,
		"id":          s.Id,
		"tag":         s.Tag,
		"args":        s.fnArgs,
		"executed_at": s.ExecutedAt,
	}
}

type WorkFlows struct {
	Flows         map[string]*Flow   `json:"flows"`
	Instances     map[string]*Flow   `json:"instances"`
	LimitRequests int                `json:"limit_requests"`
	AwaitingList  []*Awaiting        `json:"awaiting_list"`
	Results       map[string]et.Json `json:"results"`
	retentionTime time.Duration      `json:"-"`
	count         chan int           `json:"-"`
}

var workFlows *WorkFlows

/**
* newWorkFlows
* @return *WorkFlows
**/
func newWorkFlows() *WorkFlows {
	result := &WorkFlows{
		Flows:         make(map[string]*Flow),
		Instances:     make(map[string]*Flow),
		LimitRequests: envar.GetInt(1, "WORKFLOW_LIMIT_REQUESTS"),
		AwaitingList:  make([]*Awaiting, 0),
		Results:       make(map[string]et.Json),
		retentionTime: 100 * time.Millisecond,
		count:         make(chan int),
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
* instanceRun
* @param instanceId, tag string, startId int, tags, ctx et.Json
* @return et.Json, error
**/
func (s *WorkFlows) instanceRun(instanceId, tag string, startId int, tags, ctx et.Json) (et.Json, error) {
	instance, err := s.getOrCreateInstance(instanceId, tag, startId, tags)
	if err != nil {
		return et.Json{}, err
	}

	result, err := instance.run(ctx)
	if err != nil {
		return et.Json{}, err
	}

	if instance.isDebug {
		console.Debug("Flow instance:", instance.ToJson().ToString())
	}

	return result, err
}

/**
* newFlow
* @param tag, version, name, description string, fn FnContext, stop bool, createdBy string
* @return *Flow
**/
func (s *WorkFlows) newFlow(tag, version, name, description string, fn FnContext, stop bool, createdBy string) *Flow {
	flow := newFlow(s, tag, version, name, description, fn, stop, createdBy)
	s.Flows[tag] = flow

	return flow
}

/**
* run
* @param instanceId, tag string, tags, ctx et.Json
* @return et.Json, error
**/
func (s *WorkFlows) run(instanceId, tag string, startId int, tags, ctx et.Json) (et.Json, error) {
	if instanceId != "" {
		key := fmt.Sprintf("workflow:result:%s", instanceId)
		if cache.Exists(key) {
			scr, err := cache.Get(key, "")
			if err != nil {
				return et.Json{}, err
			}

			result, err := loadResult(scr)
			if err != nil {
				return et.Json{}, err
			}

			if result == nil {
				return et.Json{}, nil
			}

			if result.Error != "" {
				return et.Json{}, errors.New(result.Error)
			}

			return result.Result, nil
		}
	}

	instanceId = reg.GetUUID(instanceId)
	if s.LimitRequests == 0 {
		return s.instanceRun(instanceId, tag, startId, tags, ctx)
	}

	totalInstances := s.instanceCount()
	if totalInstances < s.LimitRequests {
		return s.instanceRun(instanceId, tag, startId, tags, ctx)
	}

	awaiting := &Awaiting{
		CreatedAt: time.Now(),
		Id:        instanceId,
		Tag:       tag,
		fn:        s.run,
		fnArgs:    []interface{}{instanceId, tag, startId, tags, ctx},
	}
	s.AwaitingList = append(s.AwaitingList, awaiting)
	event.Publish(EVENT_WORKFLOW_AWAITING, awaiting.ToJson())

	return et.Json{}, fmt.Errorf(MSG_WORKFLOW_LIMIT_REQUESTS, instanceId)
}

/**
* Rollback
* @param instanceId, tag string
* @return et.Json, error
**/

func (s *WorkFlows) rollback(instanceId, tag string) (et.Json, error) {
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
* stop
* @param instanceId, tag string
* @return error
**/
func (s *WorkFlows) stop(instanceId, tag string) error {
	instance, err := s.getInstance(instanceId, tag)
	if err != nil {
		return err
	}

	return instance.stop()
}

/**
* done
* @param instanceId string
* @return bool
**/
func (s *WorkFlows) done(instanceId string) bool {
	if s.Instances[instanceId] == nil {
		return false
	}

	instance := s.Instances[instanceId]
	n := len(instance.Results)
	if n > 0 {
		result := instance.Results[n-1]
		if result != nil {
			key := fmt.Sprintf("workflow:result:%s", instanceId)
			src, err := result.Serialize()
			if err != nil {
				console.ErrorF("WorkFlows.done, Error serializing result:%s", err.Error())
			}
			cache.Set(key, src, instance.RetentionTime)
			event.Publish(EVENT_WORKFLOW_RESULTS, result.ToJson())
		}
	}

	time.AfterFunc(s.retentionTime, func() {
		s.instanceRemove(instanceId)
		logs.Logf(packageName, MSG_WORKFLOW_DONE_INSTANCE, instanceId)
	})

	if len(s.AwaitingList) == 0 {
		return true
	}

	awaiting := s.AwaitingList[0]
	awaiting.ExecutedAt = time.Now()
	s.AwaitingList = s.AwaitingList[1:]
	args := awaiting.fnArgs
	go awaiting.fn(args[0].(string), args[1].(string), args[2].(int), args[3].(et.Json), args[4].(et.Json))
	logs.Logf(packageName, "Run instance:%s, flow:%s", awaiting.Id, awaiting.ToJson().ToString())

	return true
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
	time.AfterFunc(s.retentionTime, func() {
		delete(s.Flows, tag)
	})

	return true
}
