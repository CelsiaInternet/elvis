package ws

import (
	"os"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
)

/**
* Close
**/
func (h *Hub) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	logs.Log("Websocket", "Shutting down server...")

	if !h.run {
		return
	}

	h.run = false

	close(h.register)
	close(h.unregister)

	for _, client := range h.clients {
		client.close()
	}
	h.clients = nil

	for _, channel := range h.channels {
		channel.close()
	}
	h.channels = nil
}

/**
* Start
**/
func (h *Hub) Start() {
	if h.run {
		return
	}

	host, _ := os.Hostname()
	h.Host = envar.GetStr(host, "WS_HOST")
	h.run = true
	logs.Logf("Websocket", "Run server on %s", h.Host)

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
* SetName
* @param name string
**/
func (h *Hub) SetName(name string) {
	h.Name = name
}

/**
* Identify the hub
* @return et.Json
**/
func (h *Hub) From() et.Json {
	return et.Json{
		"id":   h.Id,
		"name": h.Name,
	}
}

/**
* GetChanel
* @param name string
* @return *Channel
**/
func (h *Hub) GetChanel(name string) *Channel {
	idx := h.indexChannel(name)
	if idx == -1 {
		return nil
	}

	return h.channels[idx]
}

/**
* NewChannel
* @param name string
* @param duration time.Duration
* @return *Channel
**/
func (h *Hub) NewChannel(name string, duration time.Duration) *Channel {
	idx := h.indexChannel(name)
	if idx != -1 {
		return h.channels[idx]
	}

	h.mutex.Lock() // Bloquear el mutex para evitar condiciones de carrera
	defer h.mutex.Unlock()

	result := newChannel(name)
	h.channels = append(h.channels, result)

	clean := func() {
		result.close()
	}

	if duration > 0 {
		go time.AfterFunc(duration, clean)
	}

	return result
}

/**
* NewQueue
* @param name string
* @param duration time.Duration
* @return *Queue
**/
func (h *Hub) NewQueue(name string, duration time.Duration) *Queue {
	idx := h.indexChannel(name)
	if idx != -1 {
		return h.queues[idx]
	}

	h.mutex.Lock() // Bloquear el mutex para evitar condiciones de carrera
	defer h.mutex.Unlock()

	result := newQueue(name)
	h.queues = append(h.queues, result)

	clean := func() {
		result.close()
	}

	if duration > 0 {
		go time.AfterFunc(duration, clean)
	}

	return result
}

/**
* Subscribe
* @param clientId string
* @param channel string
* @return error
**/
func (h *Hub) Subscribe(clientId string, channel string) error {
	idxCl := h.indexClient(clientId)
	if idxCl == -1 {
		return logs.Alertm(ERR_CLIENT_NOT_FOUND)
	}

	client := h.clients[idxCl]
	ch := h.NewChannel(channel, 0)
	ch.subscribe(client)

	return nil
}

/**
* QueueSubscribe
* @param clientId string
* @param channel string
* @param queue string
* @return error
**/
func (h *Hub) QueueSubscribe(clientId string, channel, queue string) error {
	idxCl := h.indexClient(clientId)
	if idxCl == -1 {
		return logs.Alertm(ERR_CLIENT_NOT_FOUND)
	}

	client := h.clients[idxCl]
	ch := h.NewQueue(channel, 0)
	ch.subscribe(client, queue)

	return nil
}

/**
* Stack
* @param clientId string
* @param channel string
* @return error
**/
func (h *Hub) Stack(clientId string, channel string) error {
	return h.QueueSubscribe(clientId, channel, QUEUE_STACK)
}

/**
* Unsubscribe a client from hub channels
* @param clientId string
* @param channel string
* @return error
**/
func (h *Hub) Unsubscribe(clientId string, channel string) error {
	idxCl := h.indexClient(clientId)
	if idxCl == -1 {
		return nil
	}

	client := h.clients[idxCl]

	idxCh := h.indexChannel(channel)
	if idxCh != -1 {
		ch := h.channels[idxCh]
		ch.unsubscribe(client)
	}

	idxCh = h.indexQueue(channel)
	if idxCh != -1 {
		ch := h.queues[idxCh]
		ch.unsubscribe(client)
	}

	return nil
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
	err := h.broadcast(channel, msg, ignored, from)
	if err != nil {
		return err
	}

	return nil
}

/**
* SendMessage
* @param clientId string
* @param msg Message
* @return error
**/
func (h *Hub) SendMessage(clientId string, msg Message) error {
	idx := h.indexClient(clientId)
	if idx == -1 {
		return console.NewErrorF(ERR_CLIENT_NOT_FOUND)
	}

	if idx != -1 {
		client := h.clients[idx]
		return client.sendMessage(msg)
	}

	return nil
}

/**
* GetChannels of the hub
* @param key string
* @return et.Items
**/
func (h *Hub) GetChannels(key string) et.Items {
	result := []et.Json{}
	if key == "" {
		for _, channel := range h.channels {
			result = append(result, channel.describe())
		}
	} else {
		idx := h.indexChannel(key)
		if idx != -1 {
			result = append(result, h.channels[idx].describe())
		}
	}

	return et.Items{
		Count:  len(result),
		Ok:     len(result) > 0,
		Result: result,
	}
}

/**
* GetClients of the hub
* @param key string
* @return et.Items
**/
func (h *Hub) GetClients(key string) et.Items {
	result := []et.Json{}
	if key == "" {
		for _, client := range h.clients {
			result = append(result, client.describe())
		}
	} else {
		idx := h.indexClient(key)
		if idx != -1 {
			result = append(result, h.clients[idx].describe())
		}
	}

	return et.Items{
		Count:  len(result),
		Ok:     len(result) > 0,
		Result: result,
	}
}

/**
* DrainChannel
* @param channel *Channel
**/
func (h *Hub) DrainChannel(channel string) error {
	idx := h.indexChannel(channel)
	if idx != -1 {
		ch := h.channels[idx]
		ch.drain()
	}

	idx = h.indexQueue(channel)
	if idx != -1 {
		ch := h.queues[idx]
		ch.drain()
	}

	return nil
}
