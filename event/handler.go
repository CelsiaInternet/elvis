package event

import (
	"net/http"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/et"
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
	msg := EvenMessage{
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
func Subscribe(channel string, f func(EvenMessage)) (err error) {
	if conn == nil {
		return
	}

	if len(channel) == 0 {
		return
	}

	msg := EvenMessage{
		Channel: channel,
	}
	conn.eventCreatedSub, err = conn.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		conn.decodeMessage(m.Data, &msg)
		f(msg)
	})

	return
}

// Basic function to subscrite kind stack to a channel
func Stack(channel string, f func(EvenMessage)) (err error) {
	if conn == nil {
		return
	}

	if len(channel) == 0 {
		return
	}

	msg := EvenMessage{
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
* Worker
* @param event string
* @param data et.Json
**/
func Worker(event string, data et.Json) {
	go Publish("service_event", event, data)

	go Publish("service_event", "event/publish", et.Json{
		"event": event,
		"data":  data,
	})
}

/**
* Work
* @param worker string
* @param work_id string
* @param data et.Json
**/
func Work(worker, work_id string, data et.Json) {
	go Publish("service_event", "event/work", et.Json{
		"work":    worker,
		"work_id": work_id,
		"data":    data,
	})
}

/**
* Working
* @param worker string
* @param work_id string
**/
func Working(worker, work_id string) {
	go Publish("service_event", "event/work/begin", et.Json{
		"worker":  worker,
		"work_id": work_id,
	})
}

/**
* Done
* @param work_id string
* @param event string
**/
func Done(work_id, event string) {
	go Publish("service_event", "event/work/done", et.Json{
		"work_id": work_id,
		"event":   event,
	})
}

/**
* Rejected
* @param work_id string
* @param event string
**/
func Rejected(work_id, event string) {
	go Publish("service_event", "event/work/rejected", et.Json{
		"work_id": work_id,
		"event":   event,
	})
}

/**
* Log
* @param event string
* @param data et.Json
**/
func Log(event string, data et.Json) {
	go Publish("service_log", event, data)
}

/**
* Test event, testing message broker
* @param w http.ResponseWriter
* @param r *http.Request
**/
func Test(w http.ResponseWriter, r *http.Request) {
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
