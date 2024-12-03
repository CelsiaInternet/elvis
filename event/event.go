package event

import (
	"sync"

	"github.com/celsiainternet/elvis/logs"
	"github.com/nats-io/nats.go"
)

var conn *Conn
var FromId string

type Conn struct {
	*nats.Conn
	_id             string
	eventCreatedSub map[string]*nats.Subscription
	mutex           *sync.RWMutex
}

func Load() (*Conn, error) {
	if conn != nil {
		return conn, nil
	}

	var err error
	conn, err = connect()
	if err != nil {
		return nil, err
	}

	FromId = conn._id

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
