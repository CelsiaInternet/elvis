package ws

import (
	"net/http"
	"sync"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

const ServiceName = "Websocket"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	Id         string
	Name       string
	Host       string
	clients    []*Subscriber
	channels   []*Channel
	queues     []*Queue
	mutex      *sync.RWMutex
	register   chan *Subscriber
	unregister chan *Subscriber
	main       *Client
	run        bool
}

/**
* NewWs
* @return *Hub
**/
func NewHub() *Hub {
	name := envar.GetStr(ServiceName, "RT_HUB_NAME")

	result := &Hub{
		Id:         utility.UUID(),
		Name:       name,
		clients:    make([]*Subscriber, 0),
		channels:   make([]*Channel, 0),
		register:   make(chan *Subscriber),
		unregister: make(chan *Subscriber),
		mutex:      &sync.RWMutex{},
		run:        false,
	}

	return result
}

func (h *Hub) start() {
	for {
		select {
		case client := <-h.register:
			h.onConnect(client)
		case client := <-h.unregister:
			h.onDisconnect(client)
		}
	}
}

func (h *Hub) getClient(id string) *Subscriber {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	idx := slices.IndexFunc(h.clients, func(c *Subscriber) bool { return c.Id == id })
	if idx == -1 {
		return nil
	}

	return h.clients[idx]
}

func (h *Hub) addClient(value *Subscriber) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients = append(h.clients, value)
}

func (h *Hub) removeClient(value *Subscriber) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	idx := slices.IndexFunc(h.clients, func(c *Subscriber) bool { return c.Id == value.Id })
	if idx == -1 {
		return
	}

	value.close()

	h.clients = append(h.clients[:idx], h.clients[idx+1:]...)
}

func (h *Hub) getChannel(name string) *Channel {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	idx := slices.IndexFunc(h.channels, func(c *Channel) bool { return c.Name == name })
	if idx == -1 {
		return nil
	}

	return h.channels[idx]
}

func (h *Hub) addChannel(value *Channel) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.channels = append(h.channels, value)
}

func (h *Hub) removeChannel(value *Channel) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	idx := slices.IndexFunc(h.channels, func(c *Channel) bool { return c.Name == value.Name })
	if idx == -1 {
		return
	}

	value.close()

	h.channels = append(h.channels[:idx], h.channels[idx+1:]...)
}

func (h *Hub) getQueue(name string) *Queue {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	idx := slices.IndexFunc(h.queues, func(c *Queue) bool { return c.Name == name })
	if idx == -1 {
		return nil
	}

	return h.queues[idx]
}

func (h *Hub) addQueue(value *Queue) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.queues = append(h.queues, value)
}

func (h *Hub) removeQueuel(value *Queue) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	idx := slices.IndexFunc(h.queues, func(c *Queue) bool { return c.Name == value.Name })
	if idx == -1 {
		return
	}

	value.close()

	h.queues = append(h.queues[:idx], h.queues[idx+1:]...)
}

/**
* onConnect
* @param client *Subscriber
**/
func (h *Hub) onConnect(client *Subscriber) {
	h.addClient(client)

	msg := NewMessage(h.From(), et.Json{
		"ok":       true,
		"message":  MSG_CONNECT_SUCCESSFULLY,
		"clientId": client.Id,
		"name":     client.Name,
	}, TpConnect)
	msg.Channel = "ws/connect"

	h.Publish(msg.Channel, "", msg, []string{client.Id}, h.From())
	client.sendMessage(msg)

	logs.Logf(ServiceName, MSG_CLIENT_CONNECT, client.Id, client.Name, h.Id)
}

/**
* onDisconnect
* @param client *Subscriber
**/
func (h *Hub) onDisconnect(client *Subscriber) {
	clientId := client.Id
	name := client.Name
	h.removeClient(client)

	msg := NewMessage(h.From(), et.Json{
		"ok":      true,
		"message": MSG_DISCONNECT_SUCCESSFULLY,
		"client":  client.From(),
	}, TpDisconnect)
	msg.Channel = "ws/disconnect"

	h.Publish(msg.Channel, "", msg, []string{clientId}, h.From())

	logs.Logf(ServiceName, MSG_CLIENT_DISCONNECT, clientId, name, h.Id)
}

/**
* connect
* @param socket *websocket.Conn
* @param clientId string
* @param name string
* @return *Subscriber
* @return error
**/
func (h *Hub) connect(socket *websocket.Conn, clientId, name string) (*Subscriber, error) {
	client := h.getClient(clientId)
	if client != nil {
		return client, nil
	}

	client, isNew := newSubscriber(h, socket, clientId, name)
	if isNew {
		h.register <- client

		go client.write()
		go client.read()
	}

	return client, nil
}

/**
* broadcast
* @param channel string
* @param queue string
* @param msg Message
* @param ignored []string
* @param from et.Json
* @return error
**/
func (h *Hub) broadcast(channel, queue string, msg Message, ignored []string, from et.Json) {
	msg.From = from
	msg.Ignored = ignored

	n := 0
	_channel := h.getChannel(channel)
	if _channel != nil {
		n = _channel.broadcast(msg, ignored)
	}

	_queue := h.getQueue(channel)
	if _queue != nil {
		n = _queue.broadcast(queue, msg, ignored)
	}

	logs.Logf(ServiceName, "Broadcast channel:%s sent:%d", channel, n)
}
