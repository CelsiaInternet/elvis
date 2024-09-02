package event

import (
	"github.com/cgalvisleon/elvis/cache"
	"github.com/nats-io/nats.go"
)

var conn *Conn

type Conn struct {
	conn             *nats.Conn
	eventCreatedSub  *nats.Subscription
	eventCreatedChan chan EvenMessage
}

func (c *Conn) Lock(key string) bool {
	val, err := cache.Del(key)
	if err != nil {
		return false
	}

	return val == 1
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
	if conn.conn != nil {
		conn.conn.Close()
	}

	if conn.eventCreatedSub != nil {
		conn.eventCreatedSub.Unsubscribe()
	}

	if conn.eventCreatedChan == nil {
		return
	}

	close(conn.eventCreatedChan)
}
