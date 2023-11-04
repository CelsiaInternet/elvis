package event

import (
	"time"

	"github.com/cgalvisleon/elvis/cache"
	. "github.com/cgalvisleon/elvis/json"
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

	key := id
	cache.Set(key, msg.ToString(), 15)

	return conn.conn.Publish(msg.Type(), dt)
}

func Event(project_id, event string, data map[string]interface{}) {
	go Publish("event", "event/publish", Json{
		"project_id": project_id,
		"event":      event,
		"data":       data,
	})
}

func Action(action string, data map[string]interface{}) {
	go Publish("action", action, data)
}

func Work(work, work_id string, data map[string]interface{}) {
	go Publish("work", work, Json{
		"work_id": work_id,
		"data":    data,
	})
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

func Stack(channel string, f func(CreatedEvenMessage)) (err error) {
	if conn == nil {
		return
	}

	msg := CreatedEvenMessage{
		Channel: channel,
	}

	conn.eventCreatedSub, err = conn.conn.Subscribe(channel, func(m *nats.Msg) {
		conn.decodeMessage(m.Data, &msg)
		key := msg.Id

		ok := conn.LockStack(key)
		if !ok {
			return
		}

		f(msg)
	})

	return
}
