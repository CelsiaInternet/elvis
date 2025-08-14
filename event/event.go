package event

import (
	"encoding/json"
	"slices"
	"sync"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/logs"
	"github.com/nats-io/nats.go"
)

var conn *Conn

type Conn struct {
	*nats.Conn
	id              string
	eventCreatedSub map[string]*nats.Subscription
	mutex           *sync.RWMutex
	storage         []string
}

/**
* Save
* @return error
**/
func (s *Conn) save() error {
	bt, err := json.Marshal(s.storage)
	if err != nil {
		return err
	}

	cache.Set("event:storage", string(bt), 0)

	return nil
}

/**
* Storage
* @return []string, error
**/
func (s *Conn) load() error {
	bt, err := json.Marshal(s.storage)
	if err != nil {
		return err
	}

	scr, err := cache.Get("event:storage", string(bt))
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(scr), &s.storage)
	if err != nil {
		return err
	}

	return nil
}

/**
* Save
* @return error
**/
func (s *Conn) Add(event string) (bool, error) {
	err := s.load()
	if err != nil {
		return false, err
	}

	idx := slices.IndexFunc(s.storage, func(e string) bool { return e == event })
	if idx == -1 {
		s.storage = append(s.storage, event)
	}

	return idx == -1, s.save()
}

/**
* Remove
* @return error
**/
func (s *Conn) Remove(event string) (bool, error) {
	err := s.load()
	if err != nil {
		return false, err
	}

	idx := slices.IndexFunc(s.storage, func(e string) bool { return e == event })
	if idx == -1 {
		return false, nil
	}

	s.storage = slices.Delete(s.storage, idx, 1)

	return true, s.save()
}

/**
* Load
* @return *Conn, error
**/
func Load() (*Conn, error) {
	if conn != nil {
		return conn, nil
	}

	var err error
	conn, err = connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Close() {
	if conn == nil {
		return
	}

	for _, sub := range conn.eventCreatedSub {
		sub.Unsubscribe()
	}

	conn.Close()

	logs.Log("Event", `Disconnect...`)
}

/**
* Id
* @return string
**/
func Id() string {
	if conn == nil {
		return ""
	}

	return conn.id
}
