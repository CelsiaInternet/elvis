package event

import (
	"net/http"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
	"github.com/nats-io/nats.go"
)

/**
* Publish
* @param channel string
* @param data et.Json
* @return error
**/
func Publish(channel string, data et.Json) error {
	if conn == nil {
		return nil
	}

	msg := NewEvenMessage(channel, data)
	dt, err := msg.Encode()
	if err != nil {
		return err
	}

	return conn.conn.Publish(msg.Channel, dt)
}

/**
* Subscribe
* @param channel string
* @param f func(EvenMessage)
* @return error
**/
func Subscribe(channel string, f func(EvenMessage)) (err error) {
	if conn == nil {
		return
	}

	if len(channel) == 0 {
		return
	}

	conn.eventCreatedSub, err = conn.conn.Subscribe(channel,
		func(m *nats.Msg) {
			msg, err := DecodeMessage(m.Data)
			if err != nil {
				return
			}

			f(msg)
		},
	)

	return
}

/**
* Queue
* @param string channel
* @param func(EvenMessage) f
* @return error
**/
func Queue(channel, queue string, f func(EvenMessage)) (err error) {
	if conn == nil {
		return logs.NewError(ERR_NOT_CONNECT)
	}

	if len(channel) == 0 {
		return nil
	}

	conn.eventCreatedSub, err = conn.conn.QueueSubscribe(
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
func WorkState(work_id string, status utility.WorkStatus, data et.Json) {
	work := et.Json{
		"update_at": timezone.Now(),
		"_id":       work_id,
		"status":    status.String(),
		"data":      data,
	}
	switch status {
	case utility.WorkStatusPending:
		work["pending_at"] = utility.Now()
	case utility.WorkStatusAccepted:
		work["accepted_at"] = utility.Now()
	case utility.WorkStatusProcessing:
		work["processing_at"] = utility.Now()
	case utility.WorkStatusCompleted:
		work["completed_at"] = utility.Now()
	case utility.WorkStatusFailed:
		work["failed_at"] = utility.Now()
	}

	go Publish("event/worker/state", work)
}

/**
* Data
* @param string channel
* @param func(Message) reciveFn
* @return error
**/
func Data(channel string, data et.Json) error {
	return Publish(channel, data)
}

/**
* Source
* @param string channel
* @param func(Message) reciveFn
* @return error
**/
func Source(channel string, f func(EvenMessage)) error {
	return Subscribe(channel, f)
}

/**
* Log
* @param event string
* @param data et.Json
**/
func Log(event string, data et.Json) {
	go Publish("log", data)
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
