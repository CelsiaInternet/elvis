package ws

import (
	"sync"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Subscriber struct {
	hub        *Hub
	Created_at time.Time           `json:"created_at"`
	Id         string              `json:"id"`
	Name       string              `json:"name"`
	Addr       string              `json:"addr"`
	Channels   map[string]*Channel `json:"channels"`
	Queue      map[string]*Queue   `json:"queue"`
	socket     *websocket.Conn
	outbound   chan []byte
	mutex      sync.RWMutex
}

/**
* newSubscriber
* @param *Hub
* @param *websocket.Conn
* @param string
* @param string
* @return *Subscriber
* @return bool
**/
func newSubscriber(hub *Hub, socket *websocket.Conn, id, name string) (*Subscriber, bool) {
	id = utility.GenKey(id)
	return &Subscriber{
		hub:        hub,
		Created_at: timezone.NowTime(),
		Id:         id,
		Name:       name,
		Addr:       socket.RemoteAddr().String(),
		Channels:   make(map[string]*Channel),
		Queue:      make(map[string]*Queue),
		socket:     socket,
		outbound:   make(chan []byte),
	}, true
}

/**
* Describe
* @return et.Json
**/
func (c *Subscriber) describe() et.Json {
	channels := []et.Json{}
	for _, ch := range c.Channels {
		channels = append(channels, ch.describe(1))
	}

	queues := []et.Json{}
	for _, q := range c.Queue {
		queues = append(queues, q.describe(1))
	}

	return et.Json{
		"created_at": strs.FormatDateTime("02/01/2006 03:04:05 PM", c.Created_at),
		"id":         c.Id,
		"name":       c.Name,
		"addr":       c.Addr,
		"channels":   channels,
		"queue":      queues,
	}
}

/**
* close
**/
func (c *Subscriber) close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, channel := range c.Channels {
		channel.unsubscribe(c)
	}

	for _, queue := range c.Queue {
		queue.unsubscribe(c)
	}

	c.socket.Close()
	close(c.outbound)
}

/**
* From
* @return et.Json
**/
func (c *Subscriber) From() et.Json {
	return et.Json{
		"id":   c.Id,
		"name": c.Name,
	}
}

/**
* read
**/
func (c *Subscriber) read() {
	defer func() {
		if c.hub != nil {
			c.hub.unregister <- c
			c.socket.Close()
		}
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			break
		}

		c.listener(message)
	}
}

/**
* write
**/
func (c *Subscriber) write() {
	for message := range c.outbound {
		c.socket.WriteMessage(websocket.TextMessage, message)
	}

	c.socket.WriteMessage(websocket.CloseMessage, []byte{})
}

/**
* sendMessage
* @param message Message
* @return error
**/
func (c *Subscriber) sendMessage(message Message) error {
	message.To = c.Id
	msg, err := message.Encode()
	if err != nil {
		return err
	}

	if c.socket == nil {
		return logs.Alertm(ERR_NOT_WS_SERVICE)
	}

	if c.outbound == nil {
		return logs.Alertm(ERR_NOT_WS_SERVICE)
	}

	c.outbound <- msg

	return nil
}

/**
* listener
* @param message []byte
**/
func (c *Subscriber) listener(message []byte) {
	response := func(ok bool, message string) {
		msg := NewMessage(c.hub.From(), et.Json{
			"ok":      ok,
			"message": message,
		}, TpDirect)

		c.sendMessage(msg)
	}

	msg, err := DecodeMessage(message)
	if err != nil {
		response(false, err.Error())
		return
	}

	msg.From = c.From()
	switch msg.Tp {
	case TpPing:
		response(true, "pong")
	case TpSetFrom:
		data, err := et.Object(msg.Data)
		if err != nil {
			response(false, err.Error())
			return
		}

		name := data.ValStr("", "name")
		if name == "" {
			c.Name = utility.GetOTP(6)
		}

		response(true, PARAMS_UPDATED)
	case TpSubscribe:
		console.Ping()
		channel := msg.Channel
		if channel == "" {
			response(false, ERR_CHANNEL_EMPTY)
			return
		}

		err := c.hub.Subscribe(c.Id, channel)
		if err != nil {
			response(false, err.Error())
			return
		}

		response(true, "Subscribed to channel "+channel)
	case TpQueueSubscribe:
		channel := msg.Channel
		if channel == "" {
			response(false, ERR_CHANNEL_EMPTY)
			return
		}

		queue, ok := msg.Data.(string)
		if !ok {
			response(false, ERR_QUEUE_EMPTY)
			return
		}
		if queue == "" {
			response(false, ERR_QUEUE_EMPTY)
			return
		}

		err := c.hub.QueueSubscribe(c.Id, channel, queue)
		if err != nil {
			response(false, err.Error())
			return
		}

		response(true, "Subscribe to channel "+channel)
	case TpStack:
		channel := msg.Channel
		if channel == "" {
			response(false, ERR_CHANNEL_EMPTY)
			return
		}

		err := c.hub.Stack(c.Id, channel)
		if err != nil {
			response(false, err.Error())
			return
		}

		response(true, "Subscribe to channel "+channel)
	case TpUnsubscribe:
		channel := msg.Channel
		if channel == "" {
			response(false, ERR_CHANNEL_EMPTY)
			return
		}

		err := c.hub.Unsubscribe(c.Id, channel)
		if err != nil {
			response(false, err.Error())
			return
		}

		response(true, "Unsubscribed from channel "+channel)
	case TpPublish:
		channel := msg.Channel
		if channel == "" {
			response(false, ERR_CHANNEL_EMPTY)
			return
		}

		go c.hub.Publish(channel, msg, []string{c.Id}, c.From())
	case TpDirect:
		clientId := msg.To

		msg.From = c.From()
		err := c.hub.SendMessage(clientId, msg)
		if err != nil {
			response(false, err.Error())
			return
		}
	default:
		response(false, ERR_MESSAGE_UNFORMATTED)
	}

	logs.Logf(ServiceName, "Subscriber %s message: %s", c.Id, msg.ToString())
}
