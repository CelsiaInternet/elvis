package resilience

import (
	"reflect"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/mem"
	"github.com/celsiainternet/elvis/utility"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

type TpStore int

const (
	TpStoreCache TpStore = iota
	TpStoreMemory
)

func (s TpStore) String() string {
	return []string{"cache", "memory"}[s]
}

type AttemptStatus string

const (
	StatusPending AttemptStatus = "pending"
	StatusSuccess AttemptStatus = "success"
	StatusRunning AttemptStatus = "running"
	StatusFailed  AttemptStatus = "failed"
)

type Attempt struct {
	CreatedAt     time.Time       `json:"created_at"`
	LastAttemptAt time.Time       `json:"last_attempt_at"`
	Id            string          `json:"id"`
	Tag           string          `json:"tag"`
	Description   string          `json:"description"`
	Status        AttemptStatus   `json:"status"`
	TpStore       TpStore         `json:"store"`
	Attempt       int             `json:"attempt"`
	TotalAttempts int             `json:"total_attempts"`
	TimeAttempts  time.Duration   `json:"time_attempts"`
	fn            interface{}     `json:"-"`
	fnArgs        []interface{}   `json:"-"`
	fnResult      []reflect.Value `json:"-"`
}

/**
* Json
* @return et.Json
**/
func (s *Attempt) Json() et.Json {
	return et.Json{
		"created_at":      s.CreatedAt,
		"last_attempt_at": s.LastAttemptAt,
		"id":              s.Id,
		"tag":             s.Tag,
		"description":     s.Description,
		"status":          s.Status,
		"tp_store":        s.TpStore,
		"store":           s.TpStore.String(),
		"attempt":         s.Attempt,
		"total_attempts":  s.TotalAttempts,
		"time_attempts":   s.TimeAttempts,
	}
}

/**
* Attempt
* @param id, tag, description string, fn interface{}, fnArgs ...interface{}
* @return Attempt
 */
func NewAttempt(id, tag, description string, totalAttempts int, timeAttempts time.Duration, fn interface{}, fnArgs ...interface{}) *Attempt {
	id = utility.GenId(id)
	result := &Attempt{
		Id:            id,
		Tag:           tag,
		Description:   description,
		Status:        StatusPending,
		fn:            fn,
		fnArgs:        fnArgs,
		fnResult:      []reflect.Value{},
		CreatedAt:     time.Now(),
		LastAttemptAt: time.Now(),
		TotalAttempts: totalAttempts,
		TimeAttempts:  timeAttempts,
	}

	result.save()
	return result
}

/**
* save
* @return error
**/
func (s *Attempt) save() error {
	err := cache.Set(s.Id, s.Json(), 0)
	if err != nil {
		mem.Set(s.Id, s.Json().ToString(), 0)
		s.TpStore = TpStoreMemory
	} else {
		s.TpStore = TpStoreCache
	}

	return nil
}

/**
* done
* @return error
**/
func (s *Attempt) Done() error {
	if s.TpStore == TpStoreCache {
		_, err := cache.Delete(s.Id)
		if err != nil {
			return err
		}
	} else {
		mem.Delete(s.Id)
	}

	return nil
}

/**
* SetStatus
* @param status AttemptStatus
* @return error
**/
func (s *Attempt) setStatus(status AttemptStatus) error {
	s.Status = status
	return s.save()
}

/**
* Run
* @return error
**/
func (s *Attempt) Run() ([]reflect.Value, error) {
	if s.Status == StatusSuccess {
		return []reflect.Value{reflect.ValueOf(et.Item{
			Ok:     true,
			Result: s.Json(),
		})}, nil
	}

	s.LastAttemptAt = utility.NowTime()
	s.Attempt++
	s.setStatus(StatusRunning)

	argsValues := make([]reflect.Value, len(s.fnArgs))
	for i, arg := range s.fnArgs {
		argsValues[i] = reflect.ValueOf(arg)
	}

	var err error
	var ok bool
	fn := reflect.ValueOf(s.fn)
	s.fnResult = fn.Call(argsValues)
	for _, r := range s.fnResult {
		if r.Type().Implements(errorInterface) {
			err, ok = r.Interface().(error)
			if ok && err != nil {
				s.setStatus(StatusFailed)
			}
		}
	}

	if s.Status != StatusFailed {
		s.setStatus(StatusSuccess)
	}

	logs.Log("resilience", "run:", s.Json().ToString())
	return s.fnResult, err
}
