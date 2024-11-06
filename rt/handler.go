package rt

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/ws"
)

/**
* From
* @return et.Json
**/
func From() et.Json {
	if conn == nil {
		return et.Json{}
	}

	return conn.From()
}

/**
* Ping
**/
func Ping() {
	if conn == nil {
		return
	}

	conn.Ping()
}

/**
* SetFrom
* @param params et.Json
* @return error
**/
func SetFrom(name string) error {
	if conn == nil {
		return console.NewError(ERR_NOT_CONNECT_WS)
	}

	return conn.SetFrom(name)
}

/**
* Subscribe to a channel
* @param channel string
* @param reciveFn func(message.ws.Message)
**/
func Subscribe(channel string, reciveFn func(ws.Message)) {
	if conn == nil {
		return
	}

	conn.Subscribe(channel, reciveFn)
}

/**
* Queue to a channel
* @param channel string
* @param reciveFn func(message.ws.Message)
**/
func Queue(channel, queue string, reciveFn func(ws.Message)) {
	if conn == nil {
		return
	}

	conn.Queue(channel, queue, reciveFn)
}

/**
* Unsubscribe to a channel
* @param channel string
**/
func Unsubscribe(channel string) {
	if conn == nil {
		return
	}

	conn.Unsubscribe(channel)
}

/**
* Publish a message to a channel
* @param channel string
* @param message interface{}
**/
func Publish(channel string, message interface{}) {
	if conn == nil {
		return
	}

	conn.Publish(channel, message)
}

/**
* SendMessage
* @param clientId string
* @param message interface{}
* @return error
**/
func SendMessage(clientId string, message interface{}) error {
	if conn == nil {
		return console.NewError(ERR_NOT_CONNECT_WS)
	}

	return conn.SendMessage(clientId, message)
}
