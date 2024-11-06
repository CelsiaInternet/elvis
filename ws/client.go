package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
	"github.com/gorilla/websocket"
)

type ClientConfig struct {
	ClientId  string
	Name      string
	Schema    string
	Host      string
	Path      string
	Header    http.Handler
	Reconcect int
}

/**
* From
* @return et.Json
**/
func (s *ClientConfig) From() et.Json {
	return et.Json{
		"id":   s.ClientId,
		"name": s.Name,
	}
}

type Client struct {
	config        *ClientConfig
	Channels      map[string]func(Message)
	DirectMessage func(Message)
	socket        *websocket.Conn
	connected     bool
	mutex         *sync.Mutex
}

/**
* LoadFrom
* @config config ConectPatams
* @return erro
**/
func NewClient(config *ClientConfig) (*Client, error) {
	result := &Client{
		config:    config,
		Channels:  make(map[string]func(Message)),
		connected: false,
		mutex:     &sync.Mutex{},
	}

	err := result.Connect()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) setChannel(channel string, reciveFn func(Message)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Channels[channel] = reciveFn
}

func (c *Client) getChannel(channel string) (func(Message), bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for c := range c.Channels {
		console.Debug("Channel:", c, " :: ", channel)
	}

	resul, ok := c.Channels[channel]
	return resul, ok
}

func (c *Client) deleteChannel(channel string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.Channels, channel)
}

/**
* Connect
* @return error
**/
func (c *Client) Connect() error {
	if c.connected {
		return nil
	}

	path := strs.Format(`%s://%s%s?clientId=%s&name=%s`, c.config.Schema, c.config.Host, c.config.Path, c.config.ClientId, c.config.Name)
	socket, _, err := websocket.DefaultDialer.Dial(path, nil)
	if err != nil {
		return err
	}

	c.socket = socket
	c.connected = true

	go c.Listener()

	logs.Logf("Real time", "Connected host:%s", path)

	return nil
}

/**
* Close
**/
func (c *Client) Close() {
	if c.socket == nil {
		return
	}

	c.socket.Close()
}

/**
* read
**/
func (c *Client) Listener() {
	done := make(chan struct{})

	reconnect := func() {
		if c.config.Reconcect == 0 {
			return
		}

		ticker := time.NewTicker(time.Duration(c.config.Reconcect) * time.Second)
		for range ticker.C {
			c.mutex.Lock()
			if !c.connected {
				c.Connect()
			}
			c.mutex.Unlock()
		}
	}

	go func() {
		defer close(done)

		for {
			_, data, err := c.socket.ReadMessage()
			if err != nil {
				c.connected = false
				reconnect()
				return
			}

			msg, err := DecodeMessage(data)
			if err != nil {
				logs.Alert(err)
				return
			}

			f, ok := c.getChannel(msg.Channel)
			if ok {
				f(msg)
			} else if c.DirectMessage != nil {
				c.DirectMessage(msg)
			}
		}
	}()
}

/**
* send
* @param message Message
* @return error
**/
func (c *Client) send(message Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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
func (c *Client) From() et.Json {
	return c.config.From()
}

/**
* Ping
**/
func (c *Client) Ping() {
	msg := NewMessage(c.From(), et.Json{}, TpPing)

	c.send(msg)
}

/**
* SetFrom
* @param config et.Json
* @return error
**/
func (c *Client) SetFrom(name string) error {
	if !utility.ValidName(name) {
		return logs.Alertm(ERR_INVALID_NAME)
	}

	c.config.Name = name
	msg := NewMessage(c.From(), c.From(), TpSetFrom)
	return c.send(msg)
}

/**
* Subscribe to a channel
* @param channel string
* @param reciveFn func(message.Message)
**/
func (c *Client) Subscribe(channel string, reciveFn func(Message)) {
	c.setChannel(channel, reciveFn)

	msg := NewMessage(c.From(), et.Json{}, TpSubscribe)
	msg.Channel = channel

	c.send(msg)
}

/**
* Queue to a channel
* @param channel, queue string
* @param reciveFn func(message.Message)
**/
func (c *Client) Queue(channel, queue string, reciveFn func(Message)) {
	c.setChannel(channel, reciveFn)

	msg := NewMessage(c.From(), et.Json{}, TpStack)
	msg.Channel = channel
	msg.Data = queue

	c.send(msg)
}

/**
* Stack to a channel
* @param channel string
* @param reciveFn func(message.Message)
**/
func (c *Client) Stack(channel string, reciveFn func(Message)) {
	c.Queue(channel, utility.QUEUE_STACK, reciveFn)
}

/**
* Unsubscribe to a channel
* @param channel string
**/
func (c *Client) Unsubscribe(channel string) {
	c.deleteChannel(channel)

	msg := NewMessage(c.From(), et.Json{}, TpUnsubscribe)
	msg.Channel = channel

	c.send(msg)
}

/**
* Publish a message to a channel
* @param channel string
* @param message interface{}
**/
func (c *Client) Publish(channel string, message interface{}) {
	msg := NewMessage(c.From(), message, TpPublish)
	msg.Ignored = []string{c.config.ClientId}
	msg.Channel = channel

	c.send(msg)
}

/**
* SendMessage
* @param clientId string
* @param message interface{}
* @return error
**/
func (c *Client) SendMessage(clientId string, message interface{}) error {
	msg := NewMessage(c.From(), message, TpDirect)
	msg.Ignored = []string{c.config.ClientId}
	msg.To = clientId

	return c.send(msg)
}
