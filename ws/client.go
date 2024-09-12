package ws

import (
	"time"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/logs"
	m "github.com/cgalvisleon/elvis/message"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

type WsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Client struct {
	Created_at time.Time
	hub        *Hub
	Id         string
	Name       string
	Addr       string
	socket     *websocket.Conn
	Channels   []string
	outbound   chan []byte
	closed     bool
	allowed    bool
}

/**
* NewClient
* @param *Hub
* @param *websocket.Conn
* @param string
* @param string
* @return *Client
* @return bool
**/
func newClient(hub *Hub, socket *websocket.Conn, id, name string) (*Client, bool) {
	return &Client{
		Created_at: time.Now(),
		hub:        hub,
		Id:         id,
		Name:       name,
		socket:     socket,
		Channels:   make([]string, 0),
		outbound:   make(chan []byte),
		closed:     false,
		allowed:    true,
	}, true
}

/**
* read
**/
func (c *Client) read() {
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

		c.listen(message)
	}
}

/**
* write
**/
func (c *Client) write() {
	for {
		select {
		case message, ok := <-c.outbound:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

/**
* subscribe a client to a channel
**/
func (c *Client) subscribe(channels []string) {
	for _, channel := range channels {
		idx := slices.IndexFunc(c.Channels, func(e string) bool { return e == strs.Lowcase(channel) })
		if idx == -1 {
			c.Channels = append(c.Channels, strs.Lowcase(channel))
		}
	}
}

/**
* unsubscribe a client from a channel
**/
func (c *Client) unsubscribe(channels []string) {
	for _, channel := range channels {
		idx := slices.IndexFunc(c.Channels, func(e string) bool { return e == strs.Lowcase(channel) })
		if idx != -1 {
			c.Channels = append(c.Channels[:idx], c.Channels[idx+1:]...)
		}
	}
}

/**
* sendMessage
* @param Message
* @return error
**/
func (c *Client) sendMessage(message Message) error {
	msg, err := message.Encode()
	if err != nil {
		return err
	}

	if c.closed {
		return logs.Alertm(ERR_CLIENT_IS_CLOSED)
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
* clear
**/
func (c *Client) clear() {
	c.unsubscribe(c.Channels)
}

/**
* listen
* @param []byte
**/
func (c *Client) listen(message []byte) {
	send := func(ok bool, message string) {
		msg := NewMessage(c.hub.from(), et.Json{
			"ok":      ok,
			"message": message,
		}, m.TpDirect)
		c.sendMessage(msg)
	}

	msg, err := DecodeMessage(message)
	if err != nil {
		send(false, err.Error())
		return
	}

	tp := msg.Type()
	switch tp {
	case m.TpPing:
		send(true, "pong")
	case m.TpParams:
		params, err := msg.Json()
		if err != nil {
			send(false, err.Error())
			return
		}

		name := params.ValStr("", "name")
		if name != "" {
			c.Name = name
		}

		send(true, PARAMS_UPDATED)
	case m.TpSubscribe:
		channel := msg.Channel
		if channel == "" {
			send(false, ERR_CHANNEL_EMPTY)
			return
		}

		err := c.hub.Subscribe(c.Id, channel)
		if err != nil {
			send(false, err.Error())
			return
		}

		send(true, "Subscribed to channel "+channel)
	case m.TpStack:
		channel := msg.Channel
		if channel == "" {
			send(false, ERR_CHANNEL_EMPTY)
			return
		}

		queue := msg.Queue
		if queue == "" {
			queue = "worker"
		}

		err := c.hub.Stack(c.Id, channel, queue)
		if err != nil {
			send(false, err.Error())
			return
		}

		send(true, "Stacked to channel "+channel)
	case m.TpUnsubscribe:
		channel := msg.Channel
		if channel == "" {
			send(false, ERR_CHANNEL_EMPTY)
			return
		}

		err := c.hub.Subscribe(c.Id, channel)
		if err != nil {
			send(false, err.Error())
			return
		}

		send(true, "Unsubscribed from channel "+channel)
	case m.TpPublish:
		channel := msg.Channel
		if channel == "" {
			send(false, ERR_CHANNEL_EMPTY)
			return
		}

		go c.hub.Publish(channel, msg, []string{c.Id}, c.From())
		send(true, "Message published to "+channel)
	case m.TpDirect:
		clientId := msg.to

		msg.From = c.From()
		err := c.hub.SendMessage(clientId, msg)
		if err != nil {
			send(false, err.Error())
			return
		}

		send(true, "Message sent to "+clientId)
	default:
		send(false, ERR_MESSAGE_UNFORMATTED)
	}

	logs.Logf("Websocket", "Client %s message: %s", c.Id, msg.ToString())
}

/**
* close
**/
func (c *Client) close() {
	c.closed = true
	c.socket.Close()
	close(c.outbound)
}

/**
* From
* @return et.Json
**/
func (c *Client) From() et.Json {
	return et.Json{
		"id":   c.Id,
		"name": c.Name,
	}
}
