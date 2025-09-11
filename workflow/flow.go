package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/resilience"
	"github.com/celsiainternet/elvis/utility"
)

type FlowStatus string
type TpConsistency string

const (
	FlowStatusPending     FlowStatus    = "pending"
	FlowStatusRunning     FlowStatus    = "running"
	FlowStatusDone        FlowStatus    = "done"
	FlowStatusFailed      FlowStatus    = "failed"
	TpConsistencyStrong   TpConsistency = "strong"
	TpConsistencyEventual TpConsistency = "eventual"
)

var workerHost string

func init() {
	workerHost, _ = os.Hostname()
}

type FnContext func(flow *Flow, ctx et.Json) (et.Json, error)

type Result struct {
	Step    int     `json:"step"`
	Ctx     et.Json `json:"ctx"`
	Attempt int     `json:"attempt"`
	Result  et.Json `json:"result"`
	Error   string  `json:"error"`
}

type Flow struct {
	Tag           string               `json:"tag"`
	Version       string               `json:"version"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Current       int                  `json:"current"`
	TotalAttempts int                  `json:"total_attempts"`
	TimeAttempts  time.Duration        `json:"time_attempts"`
	RetentionTime time.Duration        `json:"retention_time"`
	Ctx           et.Json              `json:"ctx"`
	Steps         []*Step              `json:"steps"`
	Ctxs          map[int]et.Json      `json:"ctxs"`
	Results       map[int]Result       `json:"results"`
	Rollbacks     map[int]Result       `json:"rollbacks"`
	LastRollback  int                  `json:"last_rollback"`
	TpConsistency TpConsistency        `json:"tp_consistency"`
	Id            string               `json:"id"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	DoneAt        time.Time            `json:"done_at"`
	Status        FlowStatus           `json:"status"`
	CreatedBy     string               `json:"created_by"`
	WorkerHost    string               `json:"worker_host"`
	Tags          et.Json              `json:"tags"`
	workFlows     *WorkFlows           `json:"-"`
	done          bool                 `json:"-"`
	goTo          int                  `json:"-"`
	err           error                `json:"-"`
	resilence     *resilience.Instance `json:"-"`
	isDebug       bool                 `json:"-"`
	team          string               `json:"-"`
	level         string               `json:"-"`
}

/**
* newFlow
* @param workFlows *WorkFlows, tag, version, name, description string, fn FnContext, totalAttempts int, timeAttempts, retentionTime time.Duration, createdBy string
* @return *Flow
**/
func newFlow(workFlows *WorkFlows, tag, version, name, description string, fn FnContext, stop bool, createdBy string) *Flow {
	flow := &Flow{
		Tag:           tag,
		Version:       version,
		Name:          name,
		Description:   description,
		Current:       0,
		TpConsistency: TpConsistencyEventual,
		Steps:         make([]*Step, 0),
		Ctx:           et.Json{},
		Ctxs:          make(map[int]et.Json),
		Results:       make(map[int]Result),
		Rollbacks:     make(map[int]Result),
		LastRollback:  -1,
		CreatedBy:     createdBy,
		Tags:          et.Json{},
		workFlows:     workFlows,
		goTo:          -1,
	}
	logs.Logf(packageName, MSG_FLOW_CREATED, tag, version, name)
	flow.Step("Start", MSG_START_WORKFLOW, fn, stop)

	return flow
}

/**
* FlowToJson
* @param flow *Flow
* @return et.Json
**/
func FlowToJson(flow *Flow) et.Json {
	steps := make([]et.Json, len(flow.Steps))
	for i, step := range flow.Steps {
		j := step.ToJson()
		j.Set("_id", i)
		steps[i] = j
	}

	result := et.Json{
		"tag":            flow.Tag,
		"version":        flow.Version,
		"name":           flow.Name,
		"description":    flow.Description,
		"total_attempts": flow.TotalAttempts,
		"time_attempts":  flow.TimeAttempts,
		"retention_time": flow.RetentionTime,
		"steps":          steps,
		"tp_consistency": flow.TpConsistency,
		"worker_host":    flow.WorkerHost,
	}

	return result
}

/**
* ToJson
* @return et.Json
**/
func (s *Flow) ToJson() et.Json {
	steps := make([]et.Json, len(s.Steps))
	for i, step := range s.Steps {
		j := step.ToJson()
		j.Set("_id", i)
		steps[i] = j
	}

	resilence := et.Json{}
	if s.resilence != nil {
		resilence = s.resilence.ToJson()
	}

	result := et.Json{
		"id":             s.Id,
		"tag":            s.Tag,
		"version":        s.Version,
		"name":           s.Name,
		"description":    s.Description,
		"current":        s.Current,
		"total_attempts": s.TotalAttempts,
		"time_attempts":  s.TimeAttempts,
		"retention_time": s.RetentionTime,
		"ctx":            s.Ctx,
		"steps":          steps,
		"ctxs":           s.Ctxs,
		"results":        s.Results,
		"rollbacks":      s.Rollbacks,
		"last_rollback":  s.LastRollback,
		"resilence":      resilence,
		"tp_consistency": s.TpConsistency,
		"created_at":     s.CreatedAt,
		"updated_at":     s.UpdatedAt,
		"done_at":        s.DoneAt,
		"status":         s.Status,
		"worker_host":    s.WorkerHost,
	}

	for k, v := range s.Tags {
		result.Set(k, v)
	}

	return result
}

/**
* save
* @return error
**/
func (s *Flow) save() error {
	event.Publish(EVENT_WORKFLOW_STATUS, s.ToJson())
	bt, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if s.RetentionTime == 0 {
		s.RetentionTime = 10 * time.Minute
	}

	cache.Set(s.Id, string(bt), s.RetentionTime)

	return nil
}

/**
* setStatus
* @param status FlowStatus
* @return error
**/
func (s *Flow) setStatus(status FlowStatus) error {
	if s.Status == status {
		return nil
	}

	done := func() {
		if s.workFlows != nil {
			s.workFlows.Done(s.Id)
		}
	}

	s.Status = status
	s.UpdatedAt = utility.NowTime()

	if s.isDebug {
		logs.Logf(packageName, MSG_INSTANCE_DEBUG, s.Id, s.ToJson().ToString())
	}

	if s.Status == FlowStatusDone {
		s.DoneAt = s.UpdatedAt
		s.done = true
		done()
	}

	switch s.Status {
	case FlowStatusFailed:
		errMsg := ""
		if s.err != nil {
			errMsg = s.err.Error()
		}
		logs.Errorf(packageName, MSG_INSTANCE_FAILED, s.Id, s.Tag, s.Status, s.Current, errMsg)
		if s.resilence != nil && s.resilence.IsEnd() {
			done()
		}
	default:
		logs.Logf(packageName, MSG_INSTANCE_STATUS, s.Id, s.Tag, s.Status, s.Current)
	}

	return s.save()
}

/**
* newInstance
* @param id string, tags et.Json
* @return *Flow, error
**/
func (s *Flow) newInstance(id string, tags et.Json) *Flow {
	id = utility.GenId(id)
	result := *s
	result.CreatedAt = utility.NowTime()
	result.WorkerHost = workerHost
	result.Id = id
	result.Tags = tags
	result.setStatus(FlowStatusPending)
	return &result
}

/**
* existInstance
* @param id string
* @return *Flow, error
**/
func (s *Flow) existInstance(id string) bool {
	return cache.Exists(id)
}

/**
* loadInstance
* @param id string
* @return *Flow, error
**/
func (s *Flow) loadInstance(id string) (*Flow, error) {
	if !s.existInstance(id) {
		return nil, fmt.Errorf(MSG_INSTANCE_NOT_FOUND)
	}

	source := &Flow{}
	bt, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	src, err := cache.Get(id, string(bt))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(src), &source)
	if err != nil {
		return nil, err
	}

	result := s.newInstance(id, source.Tags)
	result.Current = source.Current
	result.TotalAttempts = source.TotalAttempts
	result.TimeAttempts = source.TimeAttempts
	result.RetentionTime = source.RetentionTime
	result.Ctxs = source.Ctxs
	result.Results = source.Results
	result.Rollbacks = source.Rollbacks
	result.LastRollback = source.LastRollback
	result.TpConsistency = source.TpConsistency
	result.CreatedAt = source.CreatedAt
	result.UpdatedAt = source.UpdatedAt
	result.DoneAt = source.DoneAt
	result.setCtx(source.Ctx)
	result.setStatus(source.Status)

	return result, nil
}

/**
* regResult
* @param result et.Json, err error
**/
func (s *Flow) regResult(result et.Json, err error) {
	s.err = err
	errMessage := ""
	if s.err != nil {
		errMessage = s.err.Error()
	}

	attempt := 0
	if s.resilence != nil {
		attempt = s.resilence.Attempt
	}

	ctx := s.Ctxs[s.Current].Clone()
	s.Results[s.Current] = Result{
		Step:    s.Current,
		Ctx:     ctx,
		Attempt: attempt,
		Result:  result,
		Error:   errMessage,
	}
}

/**
* setFailed
* @param result et.Json, err error
**/
func (s *Flow) setFailed(result et.Json, err error) {
	s.regResult(result, err)
	s.setStatus(FlowStatusFailed)
}

/**
* setResult
* @param result et.Json, err error
**/
func (s *Flow) setResult(result et.Json, err error) {
	s.regResult(result, err)
	s.setCtx(result)
	s.save()
}

/**
* setCtx
* @param ctx et.Json
**/
func (s *Flow) setCtx(ctx et.Json) et.Json {
	for k, v := range ctx {
		s.Ctx[k] = v
	}

	return s.Ctx
}

/**
* setCurrent
* @param step int
**/
func (s *Flow) setCurrent(step int) {
	s.Current = step
	s.save()
}

/**
* Next
* @return error
**/
func (s *Flow) next() {
	s.setCurrent(s.Current + 1)
}

/**
* setGoto
* @param step int
**/
func (s *Flow) setGoto(step int, message string) {
	s.setCurrent(step)
	s.setStatus(FlowStatusRunning)
	logs.Logf(packageName, MSG_INSTANCE_GOTO, s.Id, s.Tag, step, message)
}

/**
* run
* @param ctx et.Json
* @return et.Json, error
**/
func (s *Flow) run(ctx et.Json) (et.Json, error) {
	if s.Status == FlowStatusDone {
		return s.ToJson(), fmt.Errorf(MSG_INSTANCE_ALREADY_DONE)
	} else if s.Status == FlowStatusRunning {
		return s.ToJson(), fmt.Errorf(MSG_INSTANCE_ALREADY_RUNNING)
	}

	var result et.Json
	var err error
	ctx = s.setCtx(ctx)
	for s.Current < len(s.Steps) {
		step := s.Steps[s.Current]
		s.Ctxs[s.Current] = ctx.Clone()
		s.setStatus(FlowStatusRunning)
		result, err = step.run(s, ctx)
		if err != nil {
			s.setFailed(result, err)
			resultRb, errRb := s.rollback(s.Current)
			if errRb != nil {
				return resultRb, errRb
			}

			return result, err
		}

		s.setResult(result, err)
		if s.done {
			s.next()
			return result, nil
		}

		if s.goTo != -1 {
			s.setCurrent(s.goTo)
			s.goTo = -1
			return result, nil
		}

		if step.Stop {
			s.next()
			return result, nil
		}

		if step.Expression == "" {
			s.next()
			continue
		}

		ok, err := step.Evaluate(ctx, s)
		if err != nil {
			s.setFailed(result, err)
			resultRb, errRb := s.rollback(s.Current)
			if errRb != nil {
				return resultRb, errRb
			}

			return result, err
		}

		if ok {
			s.setGoto(step.YesGoTo, MSG_INSTANCE_EXPRESSION_TRUE)
		} else {
			s.setGoto(step.NoGoTo, MSG_INSTANCE_EXPRESSION_FALSE)
		}
	}

	s.setStatus(FlowStatusDone)

	return result, err
}

/**
* rollback
* @param idx int
* @return et.Json, error
**/
func (s *Flow) rollback(idx int) (et.Json, error) {
	if idx < 0 {
		return et.Json{}, fmt.Errorf(MSG_INSTANCE_ROLLBACK)
	}

	if s.startResilence() {
		return et.Json{}, nil
	}

	if s.Status == FlowStatusDone {
		return s.ToJson(), fmt.Errorf(MSG_INSTANCE_ALREADY_DONE)
	} else if s.Status == FlowStatusRunning {
		return et.Json{}, fmt.Errorf(MSG_INSTANCE_ALREADY_RUNNING)
	} else if s.Status == FlowStatusPending {
		return et.Json{}, fmt.Errorf(MSG_INSTANCE_PENDING)
	}

	var result et.Json
	var err error
	for i := idx - 1; i >= 0; i-- {
		logs.Logf(packageName, MSG_INSTANCE_ROLLBACK_STEP, i)
		s.LastRollback = i
		step := s.Steps[i]
		if step == nil {
			continue
		}

		if step.rollbacks == nil {
			continue
		}

		if s.Ctxs[i] == nil {
			continue
		}

		ctx := s.Ctxs[i].Clone()
		result, err = step.rollbacks(s, ctx)
		if err != nil {
			attempt := 0
			if s.resilence != nil {
				attempt = s.resilence.Attempt
			}
			s.Rollbacks[i] = Result{
				Step:    i,
				Ctx:     ctx,
				Attempt: attempt,
				Result:  result,
				Error:   err.Error(),
			}

			if s.TpConsistency == TpConsistencyStrong {
				return ctx, err
			}
		}
	}

	return result, err
}

/**
* startResilence
* @return bool
**/
func (s *Flow) startResilence() bool {
	if s.TotalAttempts == 0 {
		return false
	}

	if s.resilence != nil {
		return !s.resilence.IsFailed()
	}

	description := fmt.Sprintf("flow: %s,  %s", s.Name, s.Description)
	s.resilence = resilience.AddCustom(s.Id, s.Tag, description, s.TotalAttempts, s.TimeAttempts, s.RetentionTime, s.Tags, s.team, s.level, s.run, s.Ctx)
	return true
}

/**
* setConfig
* @return error
**/
func (s *Flow) setConfig(format string, args ...any) {
	event.Publish(EVENT_WORKFLOW_SET, FlowToJson(s))
	logs.Logf(packageName, format, args...)
}

/**
* Debug
* @return *Flow
**/
func (s *Flow) Debug() *Flow {
	s.isDebug = true
	return s
}

/**
* Step
* @param name, description string, fn FnContext, retries, retryDelay int, stop bool
* @return *Fn
**/
func (s *Flow) Step(name, description string, fn FnContext, stop bool) *Flow {
	result, _ := newStep(name, description, fn, stop)
	s.Steps = append(s.Steps, result)
	s.setConfig(MSG_INSTANCE_STEP_CREATED, len(s.Steps)-1, name, s.Tag)

	return s
}

/**
* Rollback
* @params fn FnContext
* @return *Flow
**/
func (s *Flow) Rollback(fn FnContext) *Flow {
	n := len(s.Steps)
	step := s.Steps[n-1]
	step.rollbacks = fn
	s.setConfig(MSG_INSTANCE_ROLLBACK_CREATED, n-1, step.Name, s.Tag)

	return s
}

/**
* Consistency
* @param consistency TpConsistency
* @return *Flow
**/
func (s *Flow) Consistency(consistency TpConsistency) *Flow {
	s.TpConsistency = consistency
	s.setConfig(MSG_INSTANCE_CONSISTENCY, s.Tag, s.TpConsistency)

	return s
}

/**
* Resilence
* @param totalAttempts int, timeAttempts time.Duration
* @return *Flow
**/
func (s *Flow) Resilence(totalAttempts int, timeAttempts time.Duration, team string, level string) *Flow {
	s.TotalAttempts = totalAttempts
	s.TimeAttempts = timeAttempts
	retentionTime := time.Duration(s.TotalAttempts * int(timeAttempts))
	if s.RetentionTime < retentionTime {
		s.RetentionTime = retentionTime
	}
	s.team = team
	s.level = level
	s.setConfig(MSG_INSTANCE_RESILIENCE, s.Tag, totalAttempts, timeAttempts, retentionTime)

	return s
}

/**
* Retention
* @param retentionTime time.Duration
* @return *Flow
**/
func (s *Flow) Retention(retentionTime time.Duration) *Flow {
	s.RetentionTime = retentionTime
	s.setConfig(MSG_INSTANCE_RETENTION, s.Tag, retentionTime)

	return s
}

/**
* IfElse
* @param expression string, yesGoTo int, noGoTo int
* @return *Flow, error
**/
func (s *Flow) IfElse(expression string, yesGoTo int, noGoTo int) *Flow {
	n := len(s.Steps)
	step := s.Steps[n-1]
	step.IfElse(expression, yesGoTo, noGoTo)
	s.setConfig(MSG_INSTANCE_IFELSE, n-1, step.Name, expression, yesGoTo, noGoTo, s.Tag)

	return s
}

/**
* Run
* @param startId int, tags et.Json, ctx et.Json
* @return et.Json, error
**/
func (s *Flow) Run(startId int, ctx et.Json) (et.Json, error) {
	if s.workFlows == nil {
		return et.Json{}, fmt.Errorf(MSG_INSTANCE_WORKFLOWS_IS_NIL)
	}

	if s.Id == "" {
		return et.Json{}, fmt.Errorf(MSG_FLOW_NOT_INSTANCE)
	}

	return s.run(ctx)
}

/**
* Continue
* @param ctx et.Json
* @return et.Json, error
**/
func (s *Flow) Continue(ctx et.Json) (et.Json, error) {
	if s.workFlows == nil {
		return et.Json{}, fmt.Errorf(MSG_INSTANCE_WORKFLOWS_IS_NIL)
	}

	if s.Id == "" {
		return et.Json{}, fmt.Errorf(MSG_FLOW_NOT_INSTANCE)
	}

	s.setCurrent(s.Current)
	return s.run(ctx)
}

/**
* Stop
* @return error
**/
func (s *Flow) Stop() error {
	if s.workFlows == nil {
		return fmt.Errorf(MSG_INSTANCE_WORKFLOWS_IS_NIL)
	}

	if s.Id == "" {
		return fmt.Errorf(MSG_FLOW_NOT_INSTANCE)
	}

	s.Steps[s.Current].Stop = true
	return nil
}

/**
* Done
* @return error
**/
func (s *Flow) Done() error {
	if s.workFlows == nil {
		return fmt.Errorf(MSG_INSTANCE_WORKFLOWS_IS_NIL)
	}

	if s.Id == "" {
		return fmt.Errorf(MSG_FLOW_NOT_INSTANCE)
	}

	s.setStatus(FlowStatusDone)
	return nil
}

/**
* Goto
* @param step int
**/
func (s *Flow) Goto(step int) error {
	if s.workFlows == nil {
		return fmt.Errorf(MSG_INSTANCE_WORKFLOWS_IS_NIL)
	}

	if s.Id == "" {
		return fmt.Errorf(MSG_FLOW_NOT_INSTANCE)
	}

	s.goTo = step
	s.setGoto(step, MSG_INSTANCE_GOTO_USER_DECISION)

	return nil
}
