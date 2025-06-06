package event

import (
	"sync"

	"github.com/celsiainternet/elvis/logs"
	"github.com/nats-io/nats.go"
)

var conn *Conn

type Conn struct {
	*nats.Conn
	id              string
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
