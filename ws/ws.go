package ws

import (
	"net/http"
	"sync"

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
	clients    []*Client
	channels   []*Channel
	mutex      *sync.Mutex
	register   chan *Client
	unregister chan *Client
	adapter    *RedisAdapter
	run        bool
}

/**
* NewWs
* @return *Hub
**/
func NewWs() *Hub {
	name := envar.GetStr("Websocket", "RT_HUB_NAME")

	result := &Hub{
		Id:         utility.UUID(),
		Name:       name,
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
* indexChannel
* @param name string
* @return int
**/
func (h *Hub) indexChannel(name string) int {
	return slices.IndexFunc(h.channels, func(c *Channel) bool { return c.Low() == strs.Lowcase(name) })
}

/**
* indexClient
* @param id string
* @return int
**/
func (h *Hub) indexClient(id string) int {
	return slices.IndexFunc(h.clients, func(c *Client) bool { return c.Id == id })
}

/**
* onConnect
* @param client *Client
**/
func (h *Hub) onConnect(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients = append(h.clients, client)

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
* @param client *Client
**/
func (h *Hub) onDisconnect(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	client.close()
	idx := h.indexClient(client.Id)
	if idx == -1 {
		return
	}

	copy(h.clients[idx:], h.clients[idx+1:])
	h.clients[len(h.clients)-1] = nil
	h.clients = h.clients[:len(h.clients)-1]

	msg := NewMessage(h.From(), et.Json{
		"ok":      true,
		"message": "Client disconnected",
		"client":  client.From(),
	}, TpDisconnect)
	msg.Channel = "ws/disconnect"

	h.Publish(msg.Channel, msg, []string{client.Id}, h.From())

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
	h.mutex.Lock()
	defer h.mutex.Unlock()

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
* @param channel *Channel
* @param msg Message
* @param ignored []string
* @param from et.Json
* @return error
**/
func (h *Hub) broadcast(channel *Channel, msg Message, ignored []string, from et.Json) error {
	logs.Log("Broadcast", msg)
	msg.Channel = channel.Low()
	msg.From = from
	msg.Ignored = ignored
	if len(channel.Queue) > 0 {
		for queue := range channel.Queue {
			client := channel.nextTurn(queue)
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
	m, err := decodeMessageBroadcat([]byte(msg.(string)))
	if err != nil {
		logs.Alert(err)
		return
	}

	switch m.Kind {
	case BroadcastAll:
		h.Publish(m.To, m.Msg, m.Ignored, m.From)
	case BroadcastDirect:
		idx := h.indexClient(m.To)
		if idx != -1 {
			client := h.clients[idx]
			client.sendMessage(m.Msg)
		}
	}
}
