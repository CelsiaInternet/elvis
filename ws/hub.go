package ws

import (
	"net/http"
	"sync"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

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
	mutex      *sync.Mutex
	register   chan *Subscriber
	unregister chan *Subscriber
	main       *ClientWS
	run        bool
}

/**
* NewWs
* @return *Hub
**/
func NewHub() *Hub {
	name := envar.GetStr("Websocket", "RT_HUB_NAME")

	result := &Hub{
		Id:         utility.UUID(),
		Name:       name,
		clients:    make([]*Subscriber, 0),
		channels:   make([]*Channel, 0),
		register:   make(chan *Subscriber),
		unregister: make(chan *Subscriber),
		mutex:      &sync.Mutex{},
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

/**
* indexChannel
* @param name string
* @return int
**/
func (h *Hub) indexChannel(name string) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return slices.IndexFunc(h.channels, func(c *Channel) bool { return c.Name == strs.Lowcase(name) })
}

/**
* indexQueue
* @param name string
* @return int
**/
func (h *Hub) indexQueue(name string) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return slices.IndexFunc(h.queues, func(c *Queue) bool { return c.Name == strs.Lowcase(name) })
}

/**
* indexClient
* @param id string
* @return int
**/
func (h *Hub) indexClient(id string) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return slices.IndexFunc(h.clients, func(c *Subscriber) bool { return c.Id == id })
}

func (h *Hub) addClient(client *Subscriber) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients = append(h.clients, client)
}

func (h *Hub) removeClient(client *Subscriber) {
	idx := h.indexClient(client.Id)
	if idx == -1 {
		return
	}

	client.close()

	h.mutex.Lock()
	defer h.mutex.Unlock()

	copy(h.clients[idx:], h.clients[idx+1:])
	h.clients[len(h.clients)-1] = nil
	h.clients = h.clients[:len(h.clients)-1]
}

/**
* onConnect
* @param client *Subscriber
**/
func (h *Hub) onConnect(client *Subscriber) {
	h.addClient(client)

	msg := NewMessage(h.From(), et.Json{
		"ok":       true,
		"message":  "Connected successfully",
		"clientId": client.Id,
		"name":     client.Name,
	}, TpConnect)
	msg.Channel = "ws/connect"

	h.Publish(msg.Channel, msg, []string{client.Id}, h.From())
	client.sendMessage(msg)

	logs.Logf("Websocket", MSG_CLIENT_CONNECT, client.Id, client.Name, h.Id)
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
		"message": "Subscriber disconnected",
		"client":  client.From(),
	}, TpDisconnect)
	msg.Channel = "ws/disconnect"

	h.Publish(msg.Channel, msg, []string{clientId}, h.From())

	logs.Logf("Websocket", MSG_CLIENT_DISCONNECT, clientId, name, h.Id)
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
	idxC := h.indexClient(clientId)
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
* @param channel string
* @param msg Message
* @param ignored []string
* @param from et.Json
* @return error
**/
func (h *Hub) broadcast(channel string, msg Message, ignored []string, from et.Json) error {
	msg.From = from
	msg.Ignored = ignored

	n := 0
	idx := h.indexChannel(channel)
	if idx != -1 {
		_channel := h.channels[idx]
		msg.Channel = _channel.Name
		for _, client := range _channel.Subscribers {
			if !slices.Contains(ignored, client.Id) {
				err := client.sendMessage(msg)
				if err != nil {
					console.AlertE(err)
				} else {
					n++
				}
			}
		}
	}

	idx = h.indexQueue(channel)
	if idx != -1 {
		_channel := h.queues[idx]
		msg.Channel = _channel.Name
		client := _channel.nextTurn(QUEUE_STACK)
		if client != nil {
			err := client.sendMessage(msg)
			if err != nil {
				console.AlertE(err)
			} else {
				n++
			}
		}
	}

	console.LogF("Broadcast channel:%s sent:%d", channel, n)

	return nil
}

/**
* LoadMain
**/
func (h *Hub) LoadMain() {

}
