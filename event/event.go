package event

import (
	"sync"

	"github.com/celsiainternet/elvis/logs"
	"github.com/nats-io/nats.go"
)

const (
	EVENT            = "event"
	EVENT_LOG        = "event:log"
	EVENT_OVERFLOW   = "event:requests:overflow"
	EVENT_WORK       = "event:worker"
	EVENT_WORK_STATE = "event:worker:state"
	EVENT_SUBSCRIBED = "event:subscribed"
	EVENT_SOURCE     = "event:source"
)

var conn *Conn

type Conn struct {
	*nats.Conn
	id              string
	eventCreatedSub map[string]*nats.Subscription
	mutex           *sync.RWMutex
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

/**
* Close
* @return void
**/
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

/**
* HealthCheck
* @return bool
**/
func HealthCheck() bool {
	if conn == nil {
		return false
	}

	return conn.IsConnected()
}
