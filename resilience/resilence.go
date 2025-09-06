package resilience

import (
	"slices"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
)

const (
	EVENT_RESILIENCE_NOTIFY = "resilience:notify"
)

type Resilence struct {
	CreatedAt       time.Time
	Id              string
	Attempts        []*Attempt
	TotalAttempts   int
	TimeAttempts    time.Duration
	SynchronousMode bool // Nuevo campo para deshabilitar agentes en segundo plano
}

func (s *Resilence) Json() et.Json {
	attempts := make([]et.Json, len(s.Attempts))
	for i, attempt := range s.Attempts {
		attempts[i] = attempt.Json()
	}

	return et.Json{
		"id":               s.Id,
		"created_at":       s.CreatedAt,
		"attempts":         attempts,
		"total_attempts":   s.TotalAttempts,
		"time_attempts":    s.TimeAttempts,
		"synchronous_mode": s.SynchronousMode,
	}
}

var resilience *Resilence

/**
* NewResilence
* @return *Resilience
 */
func NewResilence() *Resilence {
	totalAttempts := envar.EnvarInt(3, "RESILIENCE_TOTAL_ATTEMPTS")
	timeAttempts := envar.EnvarNumber(30, "RESILIENCE_TIME_ATTEMPTS")
	interval := time.Duration(timeAttempts) * time.Second
	synchronousMode := envar.EnvarBool(false, "RESILIENCE_SYNCHRONOUS_MODE")

	return &Resilence{
		CreatedAt:       time.Now(),
		Id:              utility.UUID(),
		Attempts:        make([]*Attempt, 0),
		TotalAttempts:   totalAttempts,
		TimeAttempts:    interval,
		SynchronousMode: synchronousMode,
	}
}

/**
* HealthCheck
* @return bool
**/
func (s *Resilence) HealthCheck() bool {
	ok := event.HealthCheck()
	if !ok {
		return false
	}

	ok = cache.HealthCheck()
	if !ok {
		return false
	}

	return true
}

/**
* Notify
* @param attempt *Attempt
 */
func (s *Resilence) Notify(attempt *Attempt) {
	event.Publish(EVENT_RESILIENCE_NOTIFY, attempt.Json())
}

/**
* Done
* @param attempt *Attempt
 */
func (s *Resilence) Done(attempt *Attempt) {
	idx := slices.IndexFunc(s.Attempts, func(t *Attempt) bool { return t.Id == attempt.Id })
	if idx != -1 {
		s.Attempts = append(s.Attempts[:idx], s.Attempts[idx+1:]...)
	}

	logs.Log("resilience", "done:", attempt.Json().ToString())
}

/**
* Run
* @param attempt *Attempt
 */
func (s *Resilence) Run(attempt *Attempt) {
	if attempt.TimeAttempts == 0 {
		return
	}

	// Si está en modo síncrono, no ejecutar en segundo plano
	if s.SynchronousMode {
		logs.Log("resilience", "synchronous mode enabled, background execution disabled for:", attempt.Tag)
		return
	}

	time.AfterFunc(attempt.TimeAttempts, func() {
		if attempt.Status != StatusSuccess && attempt.Attempt < attempt.TotalAttempts {
			_, err := attempt.Run()
			if err == nil {
				s.Done(attempt)
			} else {
				if attempt.Attempt == attempt.TotalAttempts {
					s.Notify(attempt)
				} else {
					s.Run(attempt)
				}
			}
		}
	})
}

/**
* GetById
* @param id string
* @return *Attempt
 */
func (s *Resilence) GetById(id string) *Attempt {
	idx := slices.IndexFunc(s.Attempts, func(t *Attempt) bool { return t.Id == id })
	if idx != -1 {
		return s.Attempts[idx]
	}

	return nil
}

/**
* GetByTag
* @param tag string
* @return *Attempt
 */
func (s *Resilence) GetByTag(tag string) *Attempt {
	idx := slices.IndexFunc(s.Attempts, func(t *Attempt) bool { return t.Tag == tag })
	if idx != -1 {
		return s.Attempts[idx]
	}

	return nil
}

/**
* SetSynchronousMode
* @param enabled bool
* Permite habilitar/deshabilitar el modo síncrono dinámicamente
 */
func (s *Resilence) SetSynchronousMode(enabled bool) {
	s.SynchronousMode = enabled
	if enabled {
		logs.Log("resilience", "synchronous mode enabled - background agents disabled")
	} else {
		logs.Log("resilience", "synchronous mode disabled - background agents enabled")
	}
}

/**
* IsSynchronousMode
* @return bool
* Devuelve si el modo síncrono está habilitado
 */
func (s *Resilence) IsSynchronousMode() bool {
	return s.SynchronousMode
}
