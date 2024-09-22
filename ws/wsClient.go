package ws

import (
	"net/http"
	"net/url"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
	"github.com/gorilla/websocket"
)

var conn *ClientWS

type ClientWS struct {
	socket    *websocket.Conn
	Host      string
	ClientId  string
	Name      string
	channels  map[string]func(Message)
	connected bool
}

/**
* ConnectWs connect to the server using the websocket
* @param host string
* @param scheme string
* @param clientId string
* @param name string
* @return *websocket.Conn
* @return error
**/
func ConnectWs(host, scheme, clientId, name string) error {
	if conn != nil {
		return nil
	}

	if scheme == "" {
		scheme = "ws"
	}

	path := strs.Format("/%s", scheme)
	u := url.URL{Scheme: scheme, Host: host, Path: path}
	header := http.Header{}
	header.Add("clientId", clientId)
	header.Add("name", name)
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return err
	}

	conn = &ClientWS{
		socket:    ws,
		Host:      host,
		ClientId:  clientId,
		Name:      name,
		channels:  make(map[string]func(Message)),
		connected: true,
	}

	go conn.read()

	conn.SetFrom(clientId, name)

	logs.Logf("Real time", "Connected host:%s", u.String())

	return nil
}

/**
* Close
**/
func Close() {
	if conn != nil {
		conn.socket.Close()
	}
}

/**
* read
**/
func (c *ClientWS) read() {
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, data, err := c.socket.ReadMessage()
			if err != nil {
				logs.Alert(err)
				c.connected = false
				return
			}

			msg, err := DecodeMessage(data)
			if err != nil {
				logs.Alert(err)
				return
			}

			f, ok := c.channels[msg.Channel]
			if ok {
				f(msg)
			}
		}
	}()
}

/**
* send
* @param message Message
* @return error
**/
func (c *ClientWS) send(message Message) error {
	if c.socket == nil {
		return logs.Alertm(ERR_NOT_CONNECT_WS)
	}

	msg, err := message.Encode()
	if err != nil {
		return err
	}

	err = conn.socket.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	return nil
}

/**
* IsConnected
* @return bool
**/
func (c *ClientWS) IsConnected() bool {
	return c.connected
}

/**
* From
* @return et.Json
**/
func (c *ClientWS) From() et.Json {
	return et.Json{
		"id":   c.ClientId,
		"name": c.Name,
	}
}

/**
* Ping
**/
func (c *ClientWS) Ping() {
	msg := NewMessage(c.From(), et.Json{}, TpPing)

	c.send(msg)
}

/**
* SetFrom
* @param params et.Json
* @return error
**/
func (c *ClientWS) SetFrom(id, name string) error {
	if !utility.ValidId(id) {
		return logs.Alertm(ERR_INVALID_ID)
	}

	if !utility.ValidName(name) {
		return logs.Alertm(ERR_INVALID_NAME)
	}

	c.ClientId = id
	c.Name = name
	msg := NewMessage(c.From(), c.From(), TpSetFrom)
	return c.send(msg)
}

/**
* Subscribe to a channel
* @param channel string
* @param reciveFn func(message.Message)
**/
func (c *ClientWS) Subscribe(channel string, reciveFn func(Message)) {
	c.channels[channel] = reciveFn

	msg := NewMessage(c.From(), et.Json{}, TpSubscribe)
	msg.Channel = channel

	c.send(msg)
}

/**
* Queue to a channel
* @param channel string
* @param reciveFn func(message.Message)
**/
func (c *ClientWS) Queue(channel, queue string, reciveFn func(Message)) {
	c.channels[channel] = reciveFn

	msg := NewMessage(c.From(), et.Json{}, TpQueue)
	msg.Channel = channel
	msg.Queue = queue

	c.send(msg)
}

/**
* Unsubscribe to a channel
* @param channel string
**/
func (c *ClientWS) Unsubscribe(channel string) {
	delete(c.channels, channel)

	msg := NewMessage(c.From(), et.Json{}, TpUnsubscribe)
	msg.Channel = channel

	c.send(msg)
}

/**
* Publish a message to a channel
* @param channel string
* @param message interface{}
**/
func (c *ClientWS) Publish(channel string, message interface{}) {
	msg := NewMessage(c.From(), message, TpPublish)
	msg.Ignored = []string{c.ClientId}
	msg.Channel = channel

	c.send(msg)
}

/**
* SendMessage
* @param clientId string
* @param message interface{}
* @return error
**/
func (c *ClientWS) SendMessage(clientId string, message interface{}) error {
	msg := NewMessage(c.From(), message, TpDirect)
	msg.Ignored = []string{c.ClientId}
	msg.To = clientId

	return c.send(msg)
}
