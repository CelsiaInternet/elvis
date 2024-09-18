package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/logs"
	m "github.com/cgalvisleon/elvis/message"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type HubParams struct {
	Id   string
	Name string
}

type Hub struct {
	Id         string
	Name       string
	Params     *HubParams
	clients    []*Client
	channels   []*Channel
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
	adapter    *RedisAdapter
	run        bool
}

/**
* NewWs
* @return *Hub
**/
func NewWs() *Hub {
	id := utility.UUID()
	name := "Websocket Hub"

	result := &Hub{
		Id:   utility.UUID(),
		Name: name,
		Params: &HubParams{
			Id:   id,
			Name: name,
		},
		clients:    make([]*Client, 0),
		channels:   make([]*Channel, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
		run:        false,
	}

	return result
}

/**
* Start the Websocket Hub
**/
func (h *Hub) Start() {
	go h.start()
}

/**
* Run
**/
func (h *Hub) start() {
	if h.run {
		return
	}

	h.run = true
	logs.Log("Websocket", "Run websocket hub")

	for {
		select {
		case client := <-h.register:
			h.onConnect(client)
		case client := <-h.unregister:
			h.onDisconnect(client)
		}
	}
}

/**
* Close
**/
func (h *Hub) Close() {
	h.run = false
}

/**
* Identify the hub
* @return et.Json
**/
func (h *Hub) from() et.Json {
	return et.Json{
		"id":   h.Id,
		"name": h.Name,
	}
}

/**
* onConnect
* @param client *Client
**/
func (h *Hub) onConnect(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients = append(h.clients, client)
	client.Addr = client.socket.RemoteAddr().String()

	logs.Debug("Client clientId:", client.Id)
	logs.Debug("Client clientName:", client.Name)

	msg := NewMessage(h.from(), et.Json{
		"ok":      true,
		"message": "Connected successfully",
		"client":  client.From(),
	}, m.TpConnect)
	msg.Channel = "ws/connect"

	h.Mute(msg.Channel, msg, []string{client.Id}, h.from())
	client.sendMessage(msg)

	logs.Logf("Websocket", MSG_CLIENT_CONNECT, client.Id, h.Id)
}

/**
* onDisconnect
* @param client *Client
**/
func (h *Hub) onDisconnect(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	client.close()
	client.clear()
	idx := slices.IndexFunc(h.clients, func(c *Client) bool { return c.Id == client.Id })

	copy(h.clients[idx:], h.clients[idx+1:])
	h.clients[len(h.clients)-1] = nil
	h.clients = h.clients[:len(h.clients)-1]

	msg := NewMessage(h.from(), et.Json{
		"ok":      true,
		"message": "Client disconnected",
		"client":  client.From(),
	}, m.TpDisconnect)
	msg.Channel = "ws/disconnect"

	h.Mute(msg.Channel, msg, []string{client.Id}, h.from())

	logs.Logf("Websocket", MSG_CLIENT_DISCONNECT, client.Id, h.Id)
}

/**
* connect
* @param socket *websocket.Conn
* @param clientId string
* @param name string
* @return *Client
* @return error
**/
func (h *Hub) connect(socket *websocket.Conn, clientId, name string) (*Client, error) {
	idxC := slices.IndexFunc(h.clients, func(c *Client) bool { return c.Id == clientId })
	if idxC != -1 {
		return h.clients[idxC], nil
	}

	client, isNew := newClient(h, socket, clientId, name)
	if isNew {
		h.register <- client

		go client.write()
		go client.read()
	}

	return client, nil
}

/**
* broadcast
* @param channel *Channel
* @param msg Message
* @param ignored []string
* @param from et.Json
* @return error
**/
func (h *Hub) broadcast(channel *Channel, msg Message, ignored []string, from et.Json) error {
	msg.Channel = channel.Low()
	msg.From = from
	msg.Ignored = ignored
	if len(channel.Group) > 0 {
		for queue := range channel.Group {
			client := channel.NextTurn(queue)
			if client != nil {
				return client.sendMessage(msg)
			}
		}
	} else {
		for _, client := range channel.Subscribers {
			if !slices.Contains(ignored, client.Id) {
				client.sendMessage(msg)
			}
		}
	}

	if h.adapter != nil {
		h.adapter.Broadcast(channel.Name, msg, ignored, from)
	}

	return nil
}

/**
* listend
* @param msg interface{}
**/
func (h *Hub) listend(msg interface{}) {
	logs.Log("Broadcast", msg)

	m, err := decodeMessageBroadcat([]byte(msg.(string)))
	if err != nil {
		logs.Alert(err)
		return
	}

	switch m.Kind {
	case TpAll:
		h.Publish(m.To, m.Msg, m.Ignored, m.From)
	case TpDirect:
		idx := slices.IndexFunc(h.clients, func(c *Client) bool { return c.Id == m.To })
		if idx != -1 {
			client := h.clients[idx]
			client.sendMessage(m.Msg)
		}
	}
}

/**
* pruneChanner
* @param channel *Channel
**/
func (h *Hub) pruneChanner(channel *Channel) {
	if channel == nil {
		return
	}

	if channel.Count() == 0 {
		logs.Log("Channel prune", channel.Name)
		idx := slices.IndexFunc(h.channels, func(c *Channel) bool { return c.Low() == channel.Low() })
		if idx != -1 {
			h.channels = append(h.channels[:idx], h.channels[idx+1:]...)
		}
	}
}

/**
* getChanel
* @param name string
* @return *Channel
**/
func (h *Hub) getChanel(name string) *Channel {
	var result *Channel

	clean := func() {
		h.pruneChanner(result)
	}

	idx := slices.IndexFunc(h.channels, func(c *Channel) bool { return c.Low() == strs.Lowcase(name) })
	if idx == -1 {
		logs.Log("New channel", name)
		result = newChannel(name)
		h.channels = append(h.channels, result)
	} else {
		result = h.channels[idx]
	}

	duration := 5 * time.Minute
	go time.AfterFunc(duration, clean)

	return result
}

/**
* subscribe
* @param clientId string
* @param channel string
* @return *Channel
* @return error
**/
func (h *Hub) subscribe(clientId string, channel string) (*Channel, error) {
	idx := slices.IndexFunc(h.clients, func(c *Client) bool { return c.Id == clientId })
	if idx == -1 {
		return nil, logs.Alertm(ERR_CLIENT_NOT_FOUND)
	}

	client := h.clients[idx]
	result := h.getChanel(channel)
	result.Subscribe(client)
	client.subscribe([]string{channel})

	return result, nil
}

/**
* stack
* @param clientId string
* @param channel string
* @param queue string
* @return *Channel
* @return error
**/
func (h *Hub) stack(clientId string, channel, queue string) (*Channel, error) {
	idx := slices.IndexFunc(h.clients, func(c *Client) bool { return c.Id == clientId })
	if idx == -1 {
		return nil, logs.Alertm(ERR_CLIENT_NOT_FOUND)
	}

	client := h.clients[idx]
	result := h.getChanel(channel)
	result.QueueSubscribe(client, queue)
	client.subscribe([]string{channel})

	return result, nil
}

/**
* SetName
* @param name string
**/
func (h *Hub) SetName(name string) {
	h.Params.Name = name
}

/**
* Publish a message to a channel
* @param channel string
* @param msg Message
* @param ignored []string
* @param from et.Json
* @return error
**/
func (h *Hub) Publish(channel string, msg Message, ignored []string, from et.Json) error {
	ch := h.getChanel(channel)
	if len(ch.Subscribers) == 0 {
		return logs.Alertf(ERR_CHANNEL_NOT_SUBSCRIBERS, channel)
	}

	return h.broadcast(ch, msg, ignored, from)
}

/**
* Mute a message to a channel
* @param channel string
* @param msg Message
* @param ignored []string
* @param from et.Json
* @return error
**/
func (h *Hub) Mute(channel string, msg Message, ignored []string, from et.Json) error {
	ch := h.getChanel(channel)

	return h.broadcast(ch, msg, ignored, from)
}

/**
* SendMessage
* @param clientId string
* @param msg Message
* @return error
**/
func (h *Hub) SendMessage(clientId string, msg Message) error {
	idx := slices.IndexFunc(h.clients, func(c *Client) bool { return c.Id == clientId })
	if idx == -1 {
		if h.adapter != nil {
			h.adapter.Direct(clientId, msg)
		}
	}

	client := h.clients[idx]
	return client.sendMessage(msg)
}

/**
* Subscribe a client to hub channels
* @param clientId string
* @param channel string
* @return error
**/
func (h *Hub) Subscribe(clientId string, channel string) error {
	_, err := h.subscribe(clientId, channel)
	return err
}

/**
* Stack a client to hub channels
* @param clientId string
* @param channel string
* @return error
**/
func (h *Hub) Stack(clientId string, channel, queue string) error {
	_, err := h.stack(clientId, channel, queue)
	if err != nil {
		return err
	}

	return nil
}

/**
* Unsubscribe a client from hub channels
* @param clientId string
* @param channel string
* @return error
**/
func (h *Hub) Unsubscribe(clientId string, channel string) error {
	idx := slices.IndexFunc(h.clients, func(c *Client) bool { return c.Id == clientId })
	if idx == -1 {
		return logs.Alertm(ERR_CLIENT_NOT_FOUND)
	}

	client := h.clients[idx]
	client.unsubscribe([]string{channel})

	ch := h.getChanel(channel)
	ch.Unsubcribe(clientId)
	h.pruneChanner(ch)

	return nil
}

/**
* GetSubscribers
* @param channel string
* @return []*Client
**/
func (h *Hub) GetSubscribers(channel string) []*Client {
	ch := h.getChanel(channel)
	return ch.Subscribers
}
