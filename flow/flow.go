package flow

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
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

type FnContext func(ctx et.Json) (et.Item, error)

type Result struct {
	Step    int     `json:"step"`
	Ctx     et.Json `json:"ctx"`
	Attempt int     `json:"attempt"`
	Result  et.Item `json:"result"`
	Error   string  `json:"error"`
}

type Flow struct {
	Tag           string          `json:"tag"`
	Version       string          `json:"version"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Current       int             `json:"current"`
	Retries       int             `json:"retries"`
	RetryDelay    time.Duration   `json:"retry_delay"`
	RetentionTime time.Duration   `json:"retention_time"`
	Ctx           et.Json         `json:"ctx"`
	Steps         []*Step         `json:"steps"`
	Ctxs          map[int]et.Json `json:"ctxs"`
	Results       map[int]Result  `json:"results"`
	Rollbacks     map[int]Result  `json:"rollbacks"`
	LastRollback  int             `json:"last_rollback"`
	Attempt       int             `json:"attempt"`
	TpConsistency TpConsistency   `json:"tp_consistency"`
	Id            string          `json:"id"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DoneAt        time.Time       `json:"done_at"`
	Status        FlowStatus      `json:"status"`
	CreatedBy     string          `json:"created_by"`
	WorkerHost    string          `json:"worker_host"`
	workFlows     *WorkFlows      `json:"-"`
	err           error           `json:"-"`
}

/**
* newFlow
* @param workFlows *WorkFlows, tag, version, name, description string, fn FnContext, retries int, retryDelay, retentionTime time.Duration, createdBy string
* @return *Flow, error
**/
func newFlow(workFlows *WorkFlows, tag, version, name, description string, fn FnContext, retries int, retryDelay, retentionTime time.Duration, createdBy string) (*Flow, error) {
	flow := &Flow{
		Tag:           tag,
		Version:       version,
		Name:          name,
		Description:   description,
		Current:       0,
		Retries:       retries,
		RetryDelay:    retryDelay,
		RetentionTime: retentionTime,
		TpConsistency: TpConsistencyEventual,
		Steps:         make([]*Step, 0),
		Ctx:           et.Json{},
		Ctxs:          make(map[int]et.Json),
		Results:       make(map[int]Result),
		Rollbacks:     make(map[int]Result),
		LastRollback:  -1,
		Attempt:       0,
		CreatedBy:     createdBy,
		workFlows:     workFlows,
	}
	logs.Logf("Workflow", "Flujo creado flowTag:%s version:%s name:%s", tag, version, name)
	flow.Step("Start", "Start the workflow", fn, false)

	return flow, nil
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

	return et.Json{
		"id":             s.Id,
		"tag":            s.Tag,
		"version":        s.Version,
		"name":           s.Name,
		"description":    s.Description,
		"current":        s.Current,
		"retries":        s.Retries,
		"retry_delay":    s.RetryDelay,
		"ctx":            s.Ctx,
		"steps":          steps,
		"ctxs":           s.Ctxs,
		"results":        s.Results,
		"rollbacks":      s.Rollbacks,
		"last_rollback":  s.LastRollback,
		"attempt":        s.Attempt,
		"tp_consistency": s.TpConsistency,
		"created_at":     s.CreatedAt,
		"updated_at":     s.UpdatedAt,
		"done_at":        s.DoneAt,
		"status":         s.Status,
		"worker_host":    s.WorkerHost,
	}
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
		s.doneResilience()
	}

	if s.RetentionTime == 0 {
		s.RetentionTime = 10 * time.Minute
	}

	switch s.Status {
	case FlowStatusPending:
		logs.Logf("Workflow", "Instance creado:%s flowTag:%s status:%s", s.Id, s.Tag, s.Status)
	case FlowStatusRunning:
		logs.Logf("Workflow", "Instance ejecutando:%s flowTag:%s status:%s", s.Id, s.Tag, s.Status)
	case FlowStatusDone:
		logs.Logf("Workflow", "Instance terminado:%s flowTag:%s status:%s", s.Id, s.Tag, s.Status)
	case FlowStatusFailed:
		logs.Errorf("Workflow", "Instance fallido:%s flowTag:%s status:%s, error:%s", s.Id, s.Tag, s.Status, s.err.Error())
	}
	event.Publish(EVENT_WORKFLOW_STATUS, s.ToJson())
	bt, err := json.Marshal(s)
	if err == nil {
		cache.Set(s.Id, string(bt), s.RetentionTime)
	} else {
		console.Debug("statusInstance:", err.Error())
	}

	return nil
}

/**
* setFailed
* @param err error
**/
func (s *Flow) setFailed(err error) {
	s.err = err
	s.setStatus(FlowStatusFailed)
}

/**
* setGoto
* @param step int
**/
func (s *Flow) setGoto(step int, message string) {
	s.Current = step
	s.setStatus(FlowStatusRunning)
	logs.Logf("Workflow", "Instance %s flowTag:%s ir al step:%d, message:%s", s.Id, s.Tag, step, message)
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
* @param ctx et.Json, attempt int
* @return et.Item, error
**/
func (s *Flow) run(ctx et.Json) (et.Item, error) {
	if s.Status == FlowStatusDone {
		return et.Item{
			Ok:     true,
			Result: s.ToJson(),
		}, fmt.Errorf("flow already done")
	} else if s.Status == FlowStatusRunning {
		return et.Item{
			Ok:     true,
			Result: s.ToJson(),
		}, fmt.Errorf("flow already running")
	}

	var result et.Item
	var err error
	ctx = s.setCtx(ctx)
	s.Attempt++
	for i := s.Current; i < len(s.Steps); i++ {
		step := s.Steps[i]
		s.Ctxs[i] = ctx
		s.setStatus(FlowStatusRunning)
		result, err = step.run(ctx)
		if err != nil {
			s.setFailed(err)
			s.Results[i] = Result{
				Step:    i,
				Ctx:     ctx,
				Attempt: s.Attempt,
				Result:  result,
				Error:   err.Error(),
			}

			resultRb, errRb := s.rollback(i)
			if errRb != nil {
				return resultRb, errRb
			}

			return result, err
		}

		s.Results[i] = Result{
			Ctx:     ctx,
			Attempt: s.Attempt,
			Result:  result,
			Error:   "",
		}

		ctx = s.setCtx(result.Result)
		s.Current = i + 1

		if step.Stop {
			return result, nil
		}

		if step.Expression != "" {
			ok, err := step.Evaluate(ctx, s)
			if err != nil {
				resultRb, errRb := s.rollback(s.Current)
				if errRb != nil {
					return resultRb, errRb
				}

				return result, err
			}

			if ok {
				s.setGoto(step.YesGoTo, "Resultado de la expresion es true")
			} else {
				s.setGoto(step.NoGoTo, "Resultado de la expresion es false")
			}
		}
	}

	s.setStatus(FlowStatusDone)

	return result, err
}

/**
* rollback
* @param idx int
* @return et.Item, error
**/
func (s *Flow) rollback(idx int) (et.Item, error) {
	if idx < 0 {
		return et.Item{}, fmt.Errorf("Esta intentando hacer rollback de un step que no existe")
	}

	if s.Status == FlowStatusDone {
		return et.Item{
			Ok:     true,
			Result: s.ToJson(),
		}, fmt.Errorf("flow already done")
	} else if s.Status == FlowStatusRunning {
		return et.Item{
			Ok:     true,
			Result: s.ToJson(),
		}, fmt.Errorf("flow already running")
	} else if s.Status == FlowStatusPending {
		return et.Item{
			Ok:     true,
			Result: s.ToJson(),
		}, fmt.Errorf("flow is pending")
	}

	var result et.Item
	var err error
	for i := idx - 1; i >= 0; i-- {
		logs.Log("Workflow", "haciendo rollback del step:", i)
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
		result, err = step.rollbacks(ctx)
		if err != nil {
			s.Rollbacks[i] = Result{
				Step:    i,
				Ctx:     ctx,
				Attempt: s.Attempt,
				Result:  result,
				Error:   err.Error(),
			}

			if s.TpConsistency == TpConsistencyStrong {
				return et.Item{
					Ok:     false,
					Result: ctx,
				}, err
			}
		}
	}

	return result, err
}

/**
* AddResilience
* @param ctx et.Json
**/
func (s *Flow) addResilience(ctx et.Json) {
	if s.workFlows == nil {
		return
	}

	s.workFlows.addResilience(s, ctx)
}

/**
* DoneResilience
* @param ctx et.Json
**/
func (s *Flow) doneResilience() {
	if s.workFlows == nil {
		return
	}

	s.workFlows.doneResilience(s)
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
	logs.Logf("Workflow", "Step creado:%d name:%s flowTag:%s", len(s.Steps)-1, name, s.Tag)

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
	logs.Logf("Workflow", "Rollback creado:%d name:%s flowTag:%s", n-1, step.Name, s.Tag)

	return s
}

/**
* Consistency
* @param consistency TpConsistency
* @return *Flow
**/
func (s *Flow) Consistency(consistency TpConsistency) *Flow {
	s.TpConsistency = consistency
	logs.Logf("Workflow", "Consistencia definida flowTag:%s consistency:%s", s.Tag, s.TpConsistency)

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
	logs.Logf("Workflow", "IfElse definido step:%d name:%s expresion:%s ? %d : %d flowTag:%s", n-1, step.Name, expression, yesGoTo, noGoTo, s.Tag)

	return s
}

/**
* Run
* @param instanceId string, startId int, ctx et.Json
* @return et.Item, error
**/
func (s *Flow) Run(instanceId string, startId int, ctx et.Json) (et.Item, error) {
	if s.workFlows == nil {
		return et.Item{}, fmt.Errorf("workFlows is nil")
	}

	return s.workFlows.Run(instanceId, s.Tag, startId, ctx)
}
