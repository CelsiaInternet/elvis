package event

import (
	"net/http"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/response"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// Basic function to publish a message to a channel
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

// Basic function to subscribe to a channel
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

// Basic function to subscrite kind stack to a channel
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

		ok := conn.Lock(key)
		if !ok {
			return
		}

		f(msg)
	})

	return
}

/**
* Event Worker
**/

// Publish event
func Worker(event string, data et.Json) {
	// Publish event
	go Publish("service_event", event, data)

	// Publish log a eente
	go Publish("service_event", "event/publish", et.Json{
		"event": event,
		"data":  data,
	})

	logs.Log("Service event", "event:", event)
}

// Publish event asigned to a worker
func Work(worker, work_id string, data et.Json) {
	go Publish("service_event", "event/work", et.Json{
		"work":    worker,
		"work_id": work_id,
		"data":    data,
	})

	logs.Log("Service event", "worker:", worker)
}

// Publish event begin work
func Working(worker, work_id string) {
	go Publish("service_event", "event/work/begin", et.Json{
		"worker":  worker,
		"work_id": work_id,
	})

	logs.Log("Service event", "worker:", worker, " - worker_id:", work_id)
}

// Done work
func Done(work_id, event string) {
	go Publish("service_event", "event/work/done", et.Json{
		"work_id": work_id,
		"event":   event,
	})

	logs.Log("Service event", "event:", event, " - worker_id:", work_id)
}

// Rejected work
func Rejected(work_id, event string) {
	go Publish("service_event", "event/work/rejected", et.Json{
		"work_id": work_id,
		"event":   event,
	})

	logs.Log("Service event", "event:", event, " - worker_id:", work_id)
}

func Log(event string, data et.Json) {
	// Publish event
	go Publish("service_log", event, data)

	logs.Log("Service log", "event:", event)
}

// http
func Connect(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	event := body.Str("event")
	data := body.Json("data")

	Worker(event, data)

	response.JSON(w, r, http.StatusOK, et.Json{
		"status":  "ok",
		"event":   event,
		"message": "Event published",
	})

}
