package resilience

import (
	"slices"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/mem"
	"github.com/celsiainternet/elvis/utility"
)

type Resilence struct {
	CreatedAt    time.Time
	Id           string
	Name         string
	Transactions []*Transaction
	Attempts     int
	TimeAttempts time.Duration
}

func (s *Resilence) Json() et.Json {
	transactions := make([]et.Json, len(s.Transactions))
	for i, transaction := range s.Transactions {
		transactions[i] = transaction.Json()
	}

	return et.Json{
		"id":            s.Id,
		"name":          s.Name,
		"created_at":    s.CreatedAt,
		"transactions":  transactions,
		"attempts":      s.Attempts,
		"time_attempts": s.TimeAttempts,
	}
}

var resilience *Resilence

/**
* NewResilence
* @param name string
* @return *Resilience
 */
func NewResilence(name string) *Resilence {
	timeAttempts := envar.EnvarNumber(30, "RESILIENCE_TIME_ATTEMPTS")

	return &Resilence{
		CreatedAt:    time.Now(),
		Id:           utility.UUID(),
		Name:         name,
		Transactions: make([]*Transaction, 0),
		Attempts:     3,
		TimeAttempts: time.Duration(timeAttempts) * time.Second,
	}
}

/**
* Load
* @param name string
* @return *Resilience
 */
func Load(name string) *Resilence {
	if resilience != nil {
		return resilience
	}

	_, err := cache.Load()
	if err != nil {
		mem.Load()
	}

	resilience = NewResilence(name)
	return resilience
}

/**
* Run
* @param result *Transaction
 */
func (s *Resilence) Run(result *Transaction) {
	logs.Log("resilience", "run", result.Json().ToString())
	duration := s.TimeAttempts
	if duration != 0 {
		time.AfterFunc(duration, func() {
			if result.Status != StatusSuccess && result.Attempts < s.Attempts {
				result.Run()
				if result.Attempts < s.Attempts {
					s.Run(result)
				}
			}
		})
	}
}

/**
* Add
* @param tag, description string, fn interface{}, fnArgs ...interface{}
* @return *Transaction
 */
func Add(tag, description string, fn interface{}, fnArgs ...interface{}) *Transaction {
	result := NewTransaction(tag, description, fn, fnArgs...)
	resilience.Transactions = append(resilience.Transactions, result)
	resilience.Run(result)
	logs.Log("resilience", "add", result.Json().ToString())

	return result
}

/**
* GetById
* @param id string
* @return *Transaction
 */
func (s *Resilence) GetById(id string) *Transaction {
	idx := slices.IndexFunc(s.Transactions, func(t *Transaction) bool { return t.Id == id })
	if idx != -1 {
		return s.Transactions[idx]
	}

	return nil
}

/**
* GetByTag
* @param tag string
* @return *Transaction
 */
func (s *Resilence) GetByTag(tag string) *Transaction {
	idx := slices.IndexFunc(s.Transactions, func(t *Transaction) bool { return t.Tag == tag })
	if idx != -1 {
		return s.Transactions[idx]
	}

	return nil
}
