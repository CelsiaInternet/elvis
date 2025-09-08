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
	EVENT_WORKFLOW_SET                  = "workflow:set"
	EVENT_WORKFLOW_DELETE               = "workflow:delete"
	EVENT_WORKFLOW_STATUS               = "workflow:status"
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
	workFlows     *WorkFlows           `json:"-"`
	err           error                `json:"-"`
	resilence     *resilience.Instance `json:"-"`
}

/**
* newFlow
* @param workFlows *WorkFlows, tag, version, name, description string, fn FnContext, totalAttempts int, timeAttempts, retentionTime time.Duration, createdBy string
* @return *Flow
**/
func newFlow(workFlows *WorkFlows, tag, version, name, description string, fn FnContext, createdBy string) *Flow {
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
		workFlows:     workFlows,
	}
	logs.Logf("Workflow", MSG_FLOW_CREATED, tag, version, name)
	flow.Step("Start", MSG_START_WORKFLOW, fn, false)

	return flow
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

	return et.Json{
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

	s.Status = status
	s.UpdatedAt = utility.NowTime()
	if s.Status == FlowStatusDone {
		s.DoneAt = s.UpdatedAt
	}

	switch s.Status {
	case FlowStatusFailed:
		errMsg := ""
		if s.err != nil {
			errMsg = s.err.Error()
		}
		logs.Errorf("Workflow", MSG_INSTANCE_FAILED, s.Id, s.Tag, s.Status, errMsg)
	default:
		logs.Logf("Workflow", MSG_INSTANCE_STATUS, s.Id, s.Tag, s.Status)
	}

	return s.save()
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

	s.Results[s.Current] = Result{
		Step:    s.Current,
		Ctx:     s.Ctx,
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
* setGoto
* @param step int
**/
func (s *Flow) setGoto(step int, message string) {
	s.Current = step
	s.setStatus(FlowStatusRunning)
	logs.Logf("Workflow", MSG_INSTANCE_GOTO, s.Id, s.Tag, step, message)
}

/**
* cloneInstance
* @param id string
* @return *Flow, error
**/
func (s *Flow) cloneInstance(id string) *Flow {
	id = utility.GenId(id)
	result := *s
	result.CreatedAt = utility.NowTime()
	result.WorkerHost = workerHost
	result.Id = id
	result.setStatus(FlowStatusPending)
	return &result
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
		s.Ctxs[s.Current] = ctx
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
		s.Current++

		if step.Stop {
			return result, nil
		}

		if step.Expression != "" {
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

	if !s.startRollback() {
		return et.Json{}, nil
	}

	if s.Status == FlowStatusDone {
		return s.ToJson(), fmt.Errorf("flow already done")
	} else if s.Status == FlowStatusRunning {
		return et.Json{}, fmt.Errorf("flow already running")
	} else if s.Status == FlowStatusPending {
		return et.Json{}, fmt.Errorf("flow is pending")
	}

	var result et.Json
	var err error
	for i := idx - 1; i >= 0; i-- {
		logs.Logf("Workflow", MSG_INSTANCE_ROLLBACK_STEP, i)
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

		ctx := s.Ctxs[i]
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
* startRollback
* @return bool
**/
func (s *Flow) startRollback() bool {
	if s.TotalAttempts == 0 {
		return true
	}

	if s.resilence != nil {
		return s.resilence.IsFailed()
	}

	description := fmt.Sprintf("flow: %s,  %s", s.Name, s.Description)
	s.resilence = resilience.AddCustom(s.Id, s.Tag, description, s.TotalAttempts, s.TimeAttempts, s.RetentionTime, s.run, s.Ctx)
	return false
}

/**
* Step
* @param name, description string, fn FnContext, retries, retryDelay int, stop bool
* @return *Fn
**/
func (s *Flow) Step(name, description string, fn FnContext, stop bool) *Flow {
	result, _ := newStep(name, description, fn, stop)
	s.Steps = append(s.Steps, result)
	event.Publish(EVENT_WORKFLOW_SET, s.ToJson())
	logs.Logf("Workflow", MSG_INSTANCE_STEP_CREATED, len(s.Steps)-1, name, s.Tag)

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
	logs.Logf("Workflow", MSG_INSTANCE_ROLLBACK_CREATED, n-1, step.Name, s.Tag)

	return s
}

/**
* Consistency
* @param consistency TpConsistency
* @return *Flow
**/
func (s *Flow) Consistency(consistency TpConsistency) *Flow {
	s.TpConsistency = consistency
	logs.Logf("Workflow", MSG_INSTANCE_CONSISTENCY, s.Tag, s.TpConsistency)

	return s
}

/**
* Resilence
* @param totalAttempts int, timeAttempts, retentionTime time.Duration
* @return *Flow
**/
func (s *Flow) Resilence(totalAttempts int, timeAttempts, retentionTime time.Duration) *Flow {
	s.TotalAttempts = totalAttempts
	s.TimeAttempts = timeAttempts
	s.RetentionTime = retentionTime
	logs.Logf("Workflow", MSG_INSTANCE_RESILIENCE, s.Tag, totalAttempts, timeAttempts, retentionTime)

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
	logs.Logf("Workflow", MSG_INSTANCE_IFELSE, n-1, step.Name, expression, yesGoTo, noGoTo, s.Tag)

	return s
}

/**
* Run
* @param instanceId string, startId int, ctx et.Json
* @return et.Json, error
**/
func (s *Flow) Run(instanceId string, startId int, ctx et.Json) (et.Json, error) {
	if s.workFlows == nil {
		return et.Json{}, fmt.Errorf(MSG_INSTANCE_WORKFLOWS_IS_NIL)
	}

	return s.workFlows.Run(instanceId, s.Tag, startId, ctx)
}
