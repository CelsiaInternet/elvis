package resilience

import (
	"time"

	"github.com/celsiainternet/elvis/utility"
)

type Resilence struct {
	CreatedAt    time.Time
	ID           string
	Name         string
	Transactions map[string]*Transaction
}

var resilience *Resilence

/**
* NewResilence
* @param name string
* @return *Resilience
 */
func NewResilence(name string) *Resilence {
	return &Resilence{
		CreatedAt:    time.Now(),
		ID:           utility.UUID(),
		Name:         name,
		Transactions: make(map[string]*Transaction),
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

	resilience = NewResilence(name)
	return resilience
}
