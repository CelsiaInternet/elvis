package event

import (
	"fmt"
	"net/http"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
	"github.com/nats-io/nats.go"
)

const QUEUE_STACK = "stack"

/**
* publish
* @param channel string, data et.Json
* @return error
**/
func publish(channel string, data et.Json) error {
	if conn == nil {
		return nil
	}

	msg := NewEvenMessage(channel, data)
	msg.FromId = conn.id
	dt, err := msg.Encode()
	if err != nil {
		return err
	}

	conn.Publish(EVENT, dt)
	return conn.Publish(msg.Channel, dt)
}

/**
* Publish
* @param channel string, data et.Json
* @return error
**/
func Publish(channel string, data et.Json) error {
	if conn == nil {
		return nil
	}

	stage := envar.GetStr("local", "STAGE")
	publish(strs.Format(`pipe:%s:%s`, stage, channel), data)

	_, err := conn.Add(channel)
	if err != nil {
		return err
	}

	return publish(channel, data)
}

/**
* Subscribe
* @param channel string, f func(EvenMessage)
* @return error
**/
func Subscribe(channel string, f func(EvenMessage)) error {
	if conn == nil {
		return fmt.Errorf(ERR_NOT_CONNECT)
	}

	if len(channel) == 0 {
		return fmt.Errorf(ERR_CHANNEL_REQUIRED)
	}

	ok, err := conn.Add(channel)
	if err != nil {
		return err
	}

	if ok {
		publish(EVENT_SUBSCRIBED, et.Json{"channel": channel})
	}

	subscribe, err := conn.Subscribe(channel,
		func(m *nats.Msg) {
			msg, err := DecodeMessage(m.Data)
			if err != nil {
				return
			}

			msg.MySelf = msg.FromId == conn.id
			f(msg)
		},
	)
	if err != nil {
		return err
	}

	conn.mutex.Lock()
	conn.eventCreatedSub[channel] = subscribe
	conn.mutex.Unlock()

	return err
}

/**
* Unsubscribe
* @param channel string
* @return error
**/
func Unsubscribe(channel string) error {
	if conn == nil {
		return fmt.Errorf(ERR_NOT_CONNECT)
	}

	if len(channel) == 0 {
		return fmt.Errorf(ERR_CHANNEL_REQUIRED)
	}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	subscribe, ok := conn.eventCreatedSub[channel]
	if !ok {
		return fmt.Errorf("channel %s not found", channel)
	}

	subscribe.Unsubscribe()
	delete(conn.eventCreatedSub, channel)

	return nil
}

/**
* Queue
* @param string channel, string queue, func(EvenMessage) f
* @return error
**/
func Queue(channel, queue string, f func(EvenMessage)) error {
	if conn == nil {
		return fmt.Errorf(ERR_NOT_CONNECT)
	}

	if len(channel) == 0 {
		return fmt.Errorf(ERR_CHANNEL_REQUIRED)
	}

	ok, err := conn.Add(channel)
	if err != nil {
		return err
	}

	if ok {
		publish(EVENT_SUBSCRIBED, et.Json{"channel": channel})
	}

	subscribe, err := conn.QueueSubscribe(
		channel,
		queue,
		func(m *nats.Msg) {
			msg, err := DecodeMessage(m.Data)
			if err != nil {
				return
			}

			f(msg)
		},
	)
	if err != nil {
		return err
	}

	conn.mutex.Lock()
	conn.eventCreatedSub[channel] = subscribe
	conn.mutex.Unlock()

	return nil
}

/**
* Stack
* @param channel string, f func(EvenMessage)
* @return error
**/
func Stack(channel string, f func(EvenMessage)) error {
	return Queue(channel, utility.QUEUE_STACK, f)
}

/**
* Work
* @param event string, data et.Json
**/
func Work(event string, data et.Json) et.Json {
	serviceId := data.Str("service_id")
	if len(serviceId) == 0 {
		serviceId = utility.UUID()
	}
	work := et.Json{
		"created_at": timezone.Now(),
		"_id":        serviceId,
		"event":      event,
		"data":       data,
	}

	Publish(EVENT_WORK, work)
	Publish(event, work)

	return work
}

/**
* WorkState
* @param work_id string, status WorkStatus, data et.Json
**/
func WorkState(work_id string, status WorkStatus, data et.Json) {
	work := et.Json{
		"update_at": timezone.Now(),
		"_id":       work_id,
		"status":    status.String(),
		"data":      data,
	}
	switch status {
	case WorkStatusPending:
		work["pending_at"] = utility.Now()
	case WorkStatusAccepted:
		work["accepted_at"] = utility.Now()
	case WorkStatusProcessing:
		work["processing_at"] = utility.Now()
	case WorkStatusCompleted:
		work["completed_at"] = utility.Now()
	case WorkStatusFailed:
		work["failed_at"] = utility.Now()
	}

	go Publish(EVENT_WORK_STATE, work)
}

/**
* Source
* @param string model, string action, string err, data et.Json
* @return error
**/
func Source(model, action string, data et.Json) et.Json {
	source := et.Json{
		"created_at": timezone.Now(),
		"_id":        utility.UUID(),
		"model":      model,
		"action":     action,
		"data":       data,
	}

	go Publish(EVENT_SOURCE, source)

	return source
}

/**
* Log
* @param event string, data et.Json
**/
func Log(event string, data et.Json) {
	go Publish("log", data)
	go Publish(event, data)
}

/**
* Overflow
* @param data et.Json
**/
func Overflow(data et.Json) {
	go Publish(EVENT_OVERFLOW, data)
}

/**
* HttpEventWork
* @param w http.ResponseWriter, r *http.Request
**/
func HttpEventWork(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	event := body.Str("event")
	data := body.Json("data")
	work := Work(event, data)

	if len(event) == 0 {
		response.JSON(w, r, http.StatusBadRequest, et.Json{"error": "event is required"})
		return
	}

	if len(data) == 0 {
		response.JSON(w, r, http.StatusBadRequest, et.Json{"error": "data is required"})
		return
	}

	response.JSON(w, r, http.StatusOK, work)
}
