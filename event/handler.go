package event

import (
	"fmt"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func Publish(clientId, channel string, data map[string]interface{}) error {
	if conn == nil {
		return nil
	}

	now := time.Now().UTC()
	id := uuid.NewString()
	msg := CreatedEvenMessage{
		Created_at: now,
		Id:         id,
		ClientId:   clientId,
		Channel:    channel,
		Data:       data,
	}

	dt, err := conn.encodeMessage(msg)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(`%s:%s`, clientId, id)
	cache.Set(key, msg.ToString(), 15)

	return conn.conn.Publish(msg.Type(), dt)
}

func EventPublish(channel string, data map[string]interface{}) {
	go Publish("event", channel, data)
}

func Subscribe(channel string, f func(CreatedEvenMessage)) (err error) {
	if conn == nil {
		return
	}

	msg := CreatedEvenMessage{
		Channel: channel,
	}
	conn.eventCreatedSub, err = conn.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		conn.decodeMessage(m.Data, &msg)
		f(msg)
	})

	return
}

func Stack(channel, workerId string, f func(CreatedEvenMessage)) (err error) {
	if conn == nil {
		return
	}

	msg := CreatedEvenMessage{
		Channel: channel,
	}
	conn.eventCreatedSub, err = conn.conn.Subscribe(channel, func(m *nats.Msg) {
		conn.decodeMessage(m.Data, &msg)
		key := fmt.Sprintf(`%s:%s`, msg.ClientId, msg.Id)

		ok := conn.LockStack(key, workerId)
		if !ok {
			return
		}

		f(msg)
	})

	return
}
