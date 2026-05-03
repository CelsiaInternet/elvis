package resilience

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

type TpStore string

type Status string

const (
	packageName          = "resilience"
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusStop    Status = "stop"
	StatusFailed  Status = "failed"
)

type Instance struct {
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	LastAttemptAt time.Time       `json:"last_attempt_at"`
	DoneAt        time.Time       `json:"done_at"`
	Id            string          `json:"id"`
	Tag           string          `json:"tag"`
	Description   string          `json:"description"`
	Status        Status          `json:"status"`
	TpStore       TpStore         `json:"store"`
	Attempt       int             `json:"attempt"`
	TotalAttempts int             `json:"total_attempts"`
	TimeAttempts  time.Duration   `json:"time_attempts"`
	Tags          et.Json         `json:"tags"`
	Team          string          `json:"team"`
	Level         string          `json:"level"`
	stop        bool            `json:"-"`
	err         error           `json:"-"`
	fn          interface{}     `json:"-"`
	fnArgs      []interface{}   `json:"-"`
	fnArgsRefl  []reflect.Value `json:"-"` // pre-computed to avoid reflect.ValueOf per retry
	fnResult    []reflect.Value `json:"-"`
}

/**
* Serialize
* @return ([]byte, error)
**/
func (s *Instance) Serialize() ([]byte, error) {
	bt, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Instance) ToJson() et.Json {
	bt, err := s.Serialize()
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	for k, v := range s.Tags {
		result.Set(k, v)
	}

	return result
}

/**
* Save
* @return error
**/
func (s *Instance) Save() error {
	data := s.ToJson()
	event.Publish(EVENT_RESILIENCE_STATUS, data)

	if saveInstance != nil {
		return saveInstance(s)
	}

	return fmt.Errorf("Save: saveInstance is nil")
}

/**
* setStatus
* @param status Status
* @return error
**/
func (s *Instance) setStatus(status Status) error {
	if s.Status == status {
		return nil
	}

	s.Status = status
	s.UpdatedAt = utility.NowTime()
	if s.Status == StatusDone {
		s.DoneAt = s.UpdatedAt
	}

	switch s.Status {
	case StatusFailed:
		errMsg := ""
		if s.err != nil {
			errMsg = s.err.Error()
		}
		if s.Attempt == s.TotalAttempts {
			data := s.ToJson().Clone()
			data.Set("team", s.Team)
			data.Set("level", s.Level)
			message := fmt.Sprintf(MSG_RESILIENCE_FINISHED_ERROR, s.Attempt, s.TotalAttempts, s.Id, s.Tag, s.Status, errMsg)
			event.Publish(EVENT_RESILIENCE_FAILED, data)
			logs.Logf(packageName, message)
		} else {
			logs.Logf(packageName, MSG_RESILIENCE_ERROR, s.Attempt, s.TotalAttempts, s.Id, s.Tag, s.Status, errMsg)
		}
	default:
		if s.Attempt == s.TotalAttempts {
			logs.Logf(packageName, MSG_RESILIENCE_FINISHED, s.Attempt, s.TotalAttempts, s.Id, s.Tag, s.Status)
		} else {
			logs.Logf(packageName, MSG_RESILIENCE_STATUS, s.Attempt, s.TotalAttempts, s.Id, s.Tag, s.Status)
		}
	}

	return s.Save()
}

/**
* setError
* @param err error
* @return error
**/
func (s *Instance) setError(err error) {
	s.err = err
	s.setStatus(StatusFailed)
}

/**
* setStop
* @return et.Item
**/
func (s *Instance) setStop() et.Item {
	s.stop = true
	s.setStatus(StatusStop)

	return et.Item{
		Ok:     true,
		Result: s.ToJson(),
	}
}

/**
* setRestart
* @return et.Item
**/
func (s *Instance) setRestart() et.Item {
	s.stop = false
	s.setStatus(StatusPending)
	go s.run()

	return et.Item{
		Ok:     true,
		Result: s.ToJson(),
	}
}

/**
* run
* @return []reflect.Value, error
**/
func (s *Instance) run() ([]reflect.Value, error) {
	if s.Status == StatusDone {
		return []reflect.Value{reflect.ValueOf(et.Item{
			Ok:     true,
			Result: s.ToJson(),
		})}, nil
	}

	if s.stop {
		return []reflect.Value{reflect.ValueOf(et.Item{
			Ok:     false,
			Result: s.ToJson(),
		})}, nil
	}

	s.LastAttemptAt = utility.NowTime()
	s.Attempt++
	s.setStatus(StatusRunning)

	var err error
	var ok bool
	fn := reflect.ValueOf(s.fn)
	s.fnResult = fn.Call(s.fnArgsRefl)
	for _, r := range s.fnResult {
		if r.Type().Implements(errorInterface) {
			err, ok = r.Interface().(error)
			if ok && err != nil {
				s.setError(err)
			}
		}
	}

	if s.Status != StatusFailed {
		s.done()
	}

	return s.fnResult, err
}

/**
* done
* @return error
**/
func (s *Instance) done() {
	s.setStatus(StatusDone)

	time.AfterFunc(3*time.Second, func() {
		delete(resilience, s.Id)
	})
}

/**
* runAttempt
* @return error
**/
// backoff returns the delay before the next retry attempt using exponential
// backoff capped at 16x the base TimeAttempts, plus up to 25% random jitter
// to avoid thundering herd when multiple instances fail simultaneously.
func (s *Instance) backoff() time.Duration {
	if s.TimeAttempts == 0 {
		return 0
	}
	exp := s.Attempt
	if exp > 4 {
		exp = 4
	}
	dur := s.TimeAttempts * (1 << exp)
	jitter := time.Duration(rand.Int63n(int64(dur)/4 + 1))
	return dur + jitter
}

func (s *Instance) runAttempt() {
	delay := s.backoff()
	if delay == 0 {
		return
	}

	time.AfterFunc(delay, func() {
		if s.Status != StatusDone && s.Attempt < s.TotalAttempts {
			_, err := s.run()
			if err != nil {
				s.runAttempt()
			}
		}
	})
}

/**
* IsFailed
* @return bool
**/
func (s *Instance) IsFailed() bool {
	return s.Status == StatusFailed && s.Attempt == s.TotalAttempts
}

/**
* IsEnd
* @return bool
**/
func (s *Instance) IsEnd() bool {
	return s.Attempt == s.TotalAttempts
}
