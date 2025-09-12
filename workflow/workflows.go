package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/reg"
	"github.com/celsiainternet/elvis/resilience"
	"github.com/celsiainternet/elvis/timezone"
)

var (
	errorInstanceNotFound = errors.New(MSG_INSTANCE_NOT_FOUND)
)

const packageName = "workflow"

type instanceFn func(instanceId, tag string, startId int, tags, ctx et.Json, createdBy string) (et.Json, error)

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
	Flows         map[string]*Flow     `json:"flows"`
	Instances     map[string]*Instance `json:"instances"`
	LimitRequests int                  `json:"limit_requests"`
	AwaitingList  []*Awaiting          `json:"awaiting_list"`
	Results       map[string]et.Json   `json:"results"`
	count         int                  `json:"-"`
	mu            sync.Mutex           `json:"-"`
}

/**
* newWorkFlows
* @return *WorkFlows
**/
func newWorkFlows() *WorkFlows {
	result := &WorkFlows{
		Flows:         make(map[string]*Flow),
		Instances:     make(map[string]*Instance),
		LimitRequests: envar.GetInt(1, "WORKFLOW_LIMIT_REQUESTS"),
		AwaitingList:  make([]*Awaiting, 0),
		Results:       make(map[string]et.Json),
		count:         0,
		mu:            sync.Mutex{},
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
* instanceInc
**/
func (s *WorkFlows) instanceInc() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count++
	console.DebugF(MSG_INSTANCE_INSTANCE_INC, s.count, s.LimitRequests)
}

/**
* instanceDec
**/
func (s *WorkFlows) instanceDec() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count--
	console.DebugF(MSG_INSTANCE_INSTANCE_DEC, s.count, s.LimitRequests)
}

/**
* instanceCount
* @return int
**/
func (s *WorkFlows) instanceCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.count
}

/**
* newInstance
* @param tag, id string, tags et.Json, startId int, createdBy string
* @return *Instance, error
**/
func (s *WorkFlows) newInstance(tag, id string, tags et.Json, startId int, createdBy string) (*Instance, error) {
	if id == "" {
		return nil, fmt.Errorf(MSG_INSTANCE_ID_REQUIRED)
	}

	flow := s.Flows[tag]
	if flow == nil {
		return nil, fmt.Errorf(MSG_FLOW_NOT_FOUND)
	}

	now := timezone.NowTime()
	result := &Instance{
		Flow:       flow,
		workFlows:  s,
		CreatedAt:  now,
		UpdatedAt:  now,
		Id:         id,
		CreatedBy:  createdBy,
		Current:    startId,
		Ctx:        et.Json{},
		Ctxs:       make(map[int]et.Json),
		Results:    make(map[int]*Result),
		Rollbacks:  make(map[int]*Result),
		Tags:       tags,
		goTo:       -1,
		WorkerHost: workerHost,
	}
	result.setStatus(FlowStatusPending)
	s.Instances[id] = result

	return result, nil
}

/**
* loadInstance
* @param id string
* @return *Flow, error
**/
func (s *WorkFlows) loadInstance(id string) (*Instance, error) {
	if id == "" {
		return nil, fmt.Errorf(MSG_INSTANCE_ID_REQUIRED)
	}

	if s.Instances[id] != nil {
		return s.Instances[id], nil
	}

	if !cache.Exists(id) {
		return nil, errorInstanceNotFound
	}

	result := &Instance{}
	bt, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	src, err := cache.Get(id, string(bt))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(src), &result)
	if err != nil {
		return nil, err
	}

	flow := s.Flows[result.Tag]
	if flow == nil {
		return nil, fmt.Errorf(MSG_FLOW_NOT_FOUND)
	}

	result.Flow = flow
	result.setStatus(result.Status)
	s.Instances[id] = result

	return result, nil
}

/**
* runNextInstance
**/
func (s *WorkFlows) runNextInstance() (et.Json, error) {
	if len(s.AwaitingList) == 0 {
		return et.Json{}, nil
	}

	awaiting := s.AwaitingList[0]
	s.AwaitingList = s.AwaitingList[1:]
	logs.Logf(packageName, MSG_INSTANCE_RUN, awaiting.Id, awaiting.ToJson().ToString())
	awaiting.ExecutedAt = time.Now()
	args := awaiting.fnArgs
	result, err := awaiting.fn(args[0].(string), args[1].(string), args[2].(int), args[3].(et.Json), args[4].(et.Json), args[5].(string))
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* getOrCreateInstance
* @param id, tag string, startId int, tags et.Json, createdBy string
* @return *Instance, error
**/
func (s *WorkFlows) getOrCreateInstance(id, tag string, startId int, tags et.Json, createdBy string) (*Instance, error) {
	id = reg.GetUUID(id)
	if result, err := s.loadInstance(id); err == nil {
		return result, nil
	} else if errors.Is(err, errorInstanceNotFound) {
		return s.newInstance(tag, id, tags, startId, createdBy)
	}

	return nil, fmt.Errorf(MSG_INSTANCE_NOT_FOUND)
}

/**
* instanceRun
* @param instanceId, tag string, startId int, tags, ctx et.Json, createdBy string
* @return et.Json, error
**/
func (s *WorkFlows) instanceRun(instanceId, tag string, startId int, tags, ctx et.Json, createdBy string) (et.Json, error) {
	s.instanceInc()
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

			if result.Error != "" {
				return et.Json{}, errors.New(result.Error)
			}

			return result.Result, nil
		}
	}

	instance, err := s.getOrCreateInstance(instanceId, tag, startId, tags, createdBy)
	if err != nil {
		return et.Json{}, err
	}

	result, err := instance.run(ctx)
	if err != nil {
		return et.Json{}, err
	}

	if instance.isDebug {
		logs.Debugf("Flow instance:%s", instance.ToJson().ToString())
	}

	return result, err
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
* run
* @param instanceId, tag string, startId int, tags, ctx et.Json, createdBy string
* @return et.Json, error
**/
func (s *WorkFlows) run(instanceId, tag string, startId int, tags, ctx et.Json, createdBy string) (et.Json, error) {
	response := func(result et.Json, err error) (et.Json, error) {
		s.instanceDec()
		delete(s.Instances, instanceId)
		logs.Logf(packageName, MSG_WORKFLOW_DONE_INSTANCE, instanceId)
		go s.runNextInstance()

		return result, err
	}

	if s.LimitRequests == 0 {
		return response(s.instanceRun(instanceId, tag, startId, tags, ctx, createdBy))
	}

	totalInstances := s.instanceCount()
	if totalInstances < s.LimitRequests {
		return response(s.instanceRun(instanceId, tag, startId, tags, ctx, createdBy))
	}

	instanceId = reg.GetUUID(instanceId)
	awaiting := &Awaiting{
		CreatedAt: time.Now(),
		Id:        instanceId,
		Tag:       tag,
		fn:        s.run,
		fnArgs:    []interface{}{instanceId, tag, startId, tags, ctx, createdBy},
	}
	s.AwaitingList = append(s.AwaitingList, awaiting)
	event.Publish(EVENT_WORKFLOW_AWAITING, awaiting.ToJson())

	return et.Json{}, fmt.Errorf(MSG_WORKFLOW_LIMIT_REQUESTS, instanceId)
}

/**
* Rollback
* @param instanceId string
* @return et.Json, error
**/

func (s *WorkFlows) rollback(instanceId string) (et.Json, error) {
	instance, err := s.loadInstance(instanceId)
	if err != nil {
		return et.Json{}, err
	}

	result, err := instance.rollback(et.Json{}, nil)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* stop
* @param instanceId string
* @return error
**/
func (s *WorkFlows) stop(instanceId string) error {
	instance, err := s.loadInstance(instanceId)
	if err != nil {
		return err
	}

	return instance.Stop()
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
