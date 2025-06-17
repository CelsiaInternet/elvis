package resilience

import (
	"time"

	"github.com/celsiainternet/elvis/et"
)

type TransactionStatus string

const (
	StatusPending TransactionStatus = "pending"
	StatusSuccess TransactionStatus = "success"
	StatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	CreatedAt     time.Time         `json:"created_at"`
	LastAttemptAt time.Time         `json:"last_attempt_at"`
	ID            string            `json:"id"`
	Description   string            `json:"description"`
	Status        TransactionStatus `json:"status"`
	Attempts      int               `json:"attempts"`
	fn            interface{}       `json:"-"`
	fnArgs        []interface{}     `json:"-"`
	fnResult      et.Item           `json:"-"`
}

func (s *Transaction) Json() et.Json {
	return et.Json{
		"id":              s.ID,
		"description":     s.Description,
		"status":          s.Status,
		"attempts":        s.Attempts,
		"created_at":      s.CreatedAt,
		"last_attempt_at": s.LastAttemptAt,
	}
}

/**
* Transaction
* @param id, description string, fn interface{}, fnArgs ...interface{}
* @return Transaction
 */
func NewTransaction(id, description string, fn interface{}, fnArgs ...interface{}) *Transaction {
	return &Transaction{
		ID:            id,
		Description:   description,
		Status:        StatusPending,
		fn:            fn,
		fnArgs:        fnArgs,
		fnResult:      et.Item{},
		CreatedAt:     time.Now(),
		LastAttemptAt: time.Now(),
	}
}

/**
* Run
* @return error
**/
/**
func (s *Transaction) Run() (et.Item, error) {
	if s.Status == StatusSuccess {
		return et.Item{
			Ok:     true,
			Result: s.Json(),
		}, nil
	}

	s.LastAttemptAt = time.Now()
	s.Attempts++

	var err error
	s.fnResult, err = s.fn(s.fnArgs...)

	return s.fnResult, err
}
**/
