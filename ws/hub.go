package ws

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/cgalvisleon/elvis/event"
	. "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/logs"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
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
	clients    []*Client
	channels   []*Channel
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
	run        bool
}

func NewHub() *Hub {
	return &Hub{
		Id:         NewId(),
		clients:    make([]*Client, 0),
		channels:   make([]*Channel, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
		run:        false,
	}
}

func (hub *Hub) Run() {
	if hub.run {
		return
	}

	hub.run = true
	host, _ := os.Hostname()
	logs.Logf("Websocket", "Run server host:%s", host)

	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ignore {
			client.SendMessage(data)
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.clients = append(hub.clients, client)
	client.Addr = client.socket.RemoteAddr().String()
	client.isClose = false

	event.EventPublish("websocket/connect", Json{"hub": hub.Id, "client": client})

	logs.Logf("Websocket", MSG_CLIENT_CONNECT, client.Id, hub.Id)	
}

func (hub *Hub) onDisconnect(client *Client) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	client.Close()
	client.Clear()
	idx := slices.IndexFunc(hub.clients, func(e *Client) bool { return e.Id == client.Id })

	copy(hub.clients[idx:], hub.clients[idx+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]

	event.EventPublish("websocket/disconnect", Json{"hub": hub.Id, "client_id": client.Id})

	logs.Logf("Websocket", MSG_CLIENT_DISCONNECT, client.Id, hub.Id)
}

func (hub *Hub) indexClient(clientId string) int {
	return slices.IndexFunc(hub.clients, func(e *Client) bool { return e.Id == clientId })
}

func (hub *Hub) connect(socket *websocket.Conn, id, name string) (*Client, error) {
	client, isNew := NewClient(hub, socket, id, name)
	if isNew {
		hub.register <- client

		go client.Write()
		go client.Read()
	}

	return client, nil
}

func (hub *Hub) listen(client *Client, messageType int, message []byte) {
	data, err := ToJson(message)
	if err != nil {
		data = Json{
			"type":    messageType,
			"message": bytes.NewBuffer(message).String(),
		}
	}

	client.SendMessage([]byte(data.ToString()))
}

func (hub *Hub) Broadcast(message interface{}, ignoreId string) {
	var client *Client = nil
	idx := slices.IndexFunc(hub.clients, func(e *Client) bool { return e.Id == ignoreId })
	if idx != -1 {
		client = hub.clients[idx]
	}

	hub.broadcast(message, client)
}

func (hub *Hub) Publish(channel string, message interface{}, ignoreId string) {
	data, _ := json.Marshal(message)
	idx := slices.IndexFunc(hub.channels, func(e *Channel) bool { return e.Name == channel })
	if idx != -1 {
		_channel := hub.channels[idx]

		for _, client := range _channel.Subscribers {
			if client.Id != ignoreId {
				client.SendMessage(data)
			}
		}
	}
}

func (hub *Hub) SendMessage(clientId, channel string, message interface{}) bool {
	data, _ := json.Marshal(message)
	idx := slices.IndexFunc(hub.clients, func(e *Client) bool { return e.Id == clientId })
	if idx != -1 {
		client := hub.clients[idx]

		idx = slices.IndexFunc(client.channels, func(e *Channel) bool { return e.Name == channel })
		if idx != -1 {
			client.SendMessage(data)
			return true
		}
	}

	return false
}

func (hub *Hub) Subscribe(clientId string, channel string) bool {
	idx := slices.IndexFunc(hub.clients, func(e *Client) bool { return e.Id == clientId })

	if idx != -1 {
		client := hub.clients[idx]
		client.Subscribe(channel)

		event.EventPublish("websocket/subscribe", Json{"hub": hub.Id, "client": client})
		
		return true
	}

	return false
}

func (hub *Hub) Unsubscribe(clientId string, channel string) bool {
	idx := slices.IndexFunc(hub.clients, func(e *Client) bool { return e.Id == clientId })

	if idx != -1 {
		client := hub.clients[idx]
		client.Unsubscribe(channel)

		event.EventPublish("websocket/unsubscribe", Json{"hub": hub.Id, "client": client})

		return true
	}

	return false
}

func (hub *Hub) GetSubscribers(channel string) []*Client {
	idx := slices.IndexFunc(hub.channels, func(e *Channel) bool { return e.Name == channel })
	if idx != -1 {
		_channel := hub.channels[idx]
		return _channel.Subscribers
	}

	return []*Client{}
}
