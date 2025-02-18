package event

import (
	"errors"
	"net/http"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
	"github.com/nats-io/nats.go"
)

func publish(channel string, data et.Json) error {
	if conn == nil {
		return nil
	}

	msg := NewEvenMessage(channel, data)
	dt, err := msg.Encode()
	if err != nil {
		return err
	}

	return conn.Publish(msg.Channel, dt)
}

/**
* Publish
* @param channel string
* @param data et.Json
* @return error
**/
func Publish(channel string, data et.Json) error {
	stage := envar.GetStr("local", "STAGE")
	publish(strs.Format(`event:chanels:%s`, stage), et.Json{"channel": channel})
	publish(strs.Format(`pipe:%s:%s`, stage, channel), data)

	return publish(channel, data)
}

/**
* Subscribe
* @param channel string
* @param f func(EvenMessage)
* @return error
**/
func Subscribe(channel string, f func(EvenMessage)) error {
	if conn == nil {
		return errors.New(ERR_NOT_CONNECT)
	}

	if len(channel) == 0 {
		return errors.New(ERR_CHANNEL_REQUIRED)
	}

	subscribe, err := conn.Subscribe(channel,
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

	return err
}

/**
* Queue
* @param string channel
* @param func(EvenMessage) f
* @return error
**/
func Queue(channel, queue string, f func(EvenMessage)) error {
	if conn == nil {
		return errors.New(ERR_NOT_CONNECT)
	}

	if len(channel) == 0 {
		return errors.New(ERR_CHANNEL_REQUIRED)
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
* @param channel string
* @param f func(EvenMessage)
* @return error
**/
func Stack(channel string, f func(EvenMessage)) error {
	return Queue(channel, utility.QUEUE_STACK, f)
}

/**
* Work
* @param event string
* @param data et.Json
**/
func Work(event string, data et.Json) et.Json {
	work := et.Json{
		"created_at": timezone.Now(),
		"_id":        utility.UUID(),
		"from_id":    conn.Id,
		"event":      event,
		"data":       data,
	}

	go Publish("event/worker", work)
	go Publish(event, work)

	return work
}

/**
* WorkState
* @param work_id string
* @param status WorkStatus
* @param data et.Json
**/
func WorkState(work_id string, status WorkStatus, data et.Json) {
	work := et.Json{
		"update_at": timezone.Now(),
		"_id":       work_id,
		"from_id":   conn.Id,
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

	go Publish("event/worker/state", work)
}

/**
* Source
* @param string channel
* @param data et.Json
* @return error
**/
func Source(model, action, err string, data et.Json) et.Json {
	source := et.Json{
		"created_at": timezone.Now(),
		"_id":        utility.UUID(),
		"from_id":    conn.Id,
		"model":      model,
		"action":     action,
		"error":      err,
		"data":       data,
	}

	go Publish("event/source", source)
	if len(err) > 0 {
		go Publish("event/source/error", source)
	}

	return source
}

/**
* Log
* @param event string
* @param data et.Json
**/
func Log(event string, data et.Json) {
	go Publish("log", data)
	go Publish(event, data)
}

/**
* Telemetry
* @param data et.Json
**/
func Telemetry(data et.Json) {
	go Publish("telemetry", data)
}

/**
* Overflow
* @param data et.Json
**/
func Overflow(data et.Json) {
	go Publish("requests/overflow", data)
}

/**
* TokenLastUse
* @param data et.Json
**/
func TokenLastUse(data et.Json) {
	go Publish("telemetry.token.last_use", data)
}

/**
* HttpEventWork
* @param w http.ResponseWriter
* @param r *http.Request
**/
func HttpEventWork(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	event := body.Str("event")
	data := body.Json("data")
	work := Work(event, data)

	response.JSON(w, r, http.StatusOK, work)
}
