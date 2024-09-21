package event

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/nats-io/nats.go"
)

var conn *Conn

type Conn struct {
	conn             *nats.Conn
	eventCreatedSub  *nats.Subscription
	eventCreatedChan chan EvenMessage
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

	if conn.eventCreatedChan != nil {
		close(conn.eventCreatedChan)
	}

	console.LogK("Event", `Disconnect...`)
}
