package ws

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
	"github.com/gorilla/websocket"
)

var conn *ClientWS

type ClientWS struct {
	Host      string
	ClientId  string
	Name      string
	Channels  map[string]func(Message)
	socket    *websocket.Conn
	connected bool
	mutex     *sync.Mutex
}

/**
* LoadFrom
* @params id, name string
* @return erro
**/
func NewClientWS(id, name, schema, host, path string) (*ClientWS, error) {
	path = strs.Format(`%s?clientId=%s&name=%s`, path, id, name)
	u := url.URL{Scheme: schema, Host: host, Path: path}
	header := http.Header{}
	socket, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return nil, err
	}

	client := &ClientWS{
		socket:    socket,
		Host:      host,
		ClientId:  id,
		Name:      name,
		Channels:  make(map[string]func(Message)),
		connected: true,
	}

	go conn.Read()

	logs.Logf("Real time", "Connected host:%s", u.String())

	return client, nil
}

func (c *ClientWS) Close() {
	if c.socket == nil {
		return
	}

	c.socket.Close()
}

/**
* read
**/
func (c *ClientWS) Read() {
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

			f, ok := c.Channels[msg.Channel]
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

	err = c.socket.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	return nil
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

	conn.send(msg)
}

/**
* SetFrom
* @param params et.Json
* @return error
**/
func (c *ClientWS) SetFrom(name string) error {
	if !utility.ValidName(name) {
		return logs.Alertm(ERR_INVALID_NAME)
	}

	conn.Name = name
	msg := NewMessage(c.From(), c.From(), TpSetFrom)
	return conn.send(msg)
}

/**
* Subscribe to a channel
* @param channel string
* @param reciveFn func(message.Message)
**/
func (c *ClientWS) Subscribe(channel string, reciveFn func(Message)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conn.Channels[channel] = reciveFn

	msg := NewMessage(c.From(), et.Json{}, TpSubscribe)
	msg.Channel = channel

	conn.send(msg)
}

/**
* Queue to a channel
* @param channel string
* @param reciveFn func(message.Message)
**/
func (c *ClientWS) Queue(channel, queue string, reciveFn func(Message)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conn.Channels[channel] = reciveFn

	msg := NewMessage(c.From(), et.Json{}, TpStack)
	msg.Channel = channel
	msg.Queue = queue

	conn.send(msg)
}

/**
* Unsubscribe to a channel
* @param channel string
**/
func (c *ClientWS) Unsubscribe(channel string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(conn.Channels, channel)

	msg := NewMessage(c.From(), et.Json{}, TpUnsubscribe)
	msg.Channel = channel

	conn.send(msg)
}

/**
* Publish a message to a channel
* @param channel string
* @param message interface{}
**/
func (c *ClientWS) Publish(channel string, message interface{}) {
	msg := NewMessage(c.From(), message, TpPublish)
	msg.Ignored = []string{conn.ClientId}
	msg.Channel = channel

	conn.send(msg)
}

/**
* SendMessage
* @param clientId string
* @param message interface{}
* @return error
**/
func (c *ClientWS) SendMessage(clientId string, message interface{}) error {
	msg := NewMessage(c.From(), message, TpDirect)
	msg.Ignored = []string{conn.ClientId}
	msg.To = clientId

	return conn.send(msg)
}
