package ws

import (
	"os"
	"time"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
)

func (h *Hub) Describe() et.Json {
	return et.Json{
		"id":       h.Id,
		"name":     h.Name,
		"host":     h.Host,
		"type":     h.TypeNode.String(),
		"channels": len(h.channels),
		"clients":  len(h.clients),
		"queues":   len(h.queues),
	}
}

/**
* Close
**/
func (h *Hub) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	logs.Log(ServiceName, "Shutting down server...")

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

	for _, queue := range h.queues {
		queue.close()
	}
	h.queues = nil
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

	go h.start()

	logs.Logf(ServiceName, "Run server on %s", h.Host)
}

/**
* SetName
* @param name string
**/
func (h *Hub) SetName(name string) {
	h.Name = name
}

/**
* InitMaster
**/
func (h *Hub) InitMaster() {
	h.TypeNode = NodeMaster
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
	return h.getChannel(name)
}

/**
* NewChannel
* @param name string
* @param duration time.Duration
* @return *Channel
**/
func (h *Hub) NewChannel(name string, duration time.Duration) *Channel {
	result := h.getChannel(name)
	if result != nil {
		return result
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	result = newChannel(name)
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
func (h *Hub) NewQueue(name, queue string, duration time.Duration) *Queue {
	result := h.getQueue(name, queue)
	if result != nil {
		return result
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	result = newQueue(name, queue)
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
	client := h.getClient(clientId)
	if client == nil {
		return logs.Alertm(ERR_CLIENT_NOT_FOUND)
	}

	ch := h.NewChannel(channel, 0)
	ch.subscribe(client)

	if h.TypeNode != NotNode {
		h.ClusterSubscribed(channel)
	}

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
	client := h.getClient(clientId)
	if client == nil {
		return logs.Alertm(ERR_CLIENT_NOT_FOUND)
	}

	ch := h.NewQueue(channel, queue, 0)
	ch.subscribe(client)

	if h.TypeNode != NotNode {
		h.ClusterSubscribed(channel)
	}

	return nil
}

/**
* Stack
* @param clientId string
* @param channel string
* @return error
**/
func (h *Hub) Stack(clientId string, channel string) error {
	return h.QueueSubscribe(clientId, channel, utility.QUEUE_STACK)
}

/**
* Unsubscribe a client from hub channels
* @param clientId string
* @param channel string
* @return error
**/
func (h *Hub) Unsubscribe(clientId string, channel, queue string) error {
	client := h.getClient(clientId)
	if client == nil {
		return nil
	}

	ch := h.getChannel(channel)
	if ch != nil {
		ch.unsubscribe(client)
	}

	qu := h.getQueue(channel, queue)
	if qu != nil {
		qu.unsubscribe(client)
	}

	if h.TypeNode != NotNode {
		h.ClusterUnSubscribed(channel)
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
func (h *Hub) Publish(channel, queue string, msg Message, ignored []string, from et.Json) {
	h.broadcast(channel, queue, msg, ignored, from)

	if h.TypeNode != NotNode {
		h.ClusterPublish(channel, msg)
	}
}

/**
* SendMessage
* @param clientId string
* @param msg Message
* @return error
**/
func (h *Hub) SendMessage(clientId string, msg Message) error {
	client := h.getClient(clientId)
	if client == nil && h.TypeNode != NotNode {
		channel := clusterChannel(clientId)
		h.master.Publish(channel, msg)
	}

	if client == nil {
		return logs.Alertm(ERR_CLIENT_NOT_FOUND)
	}

	return client.sendMessage(msg)
}

/**
* GetChannels of the hub
* @param key string
* @return et.Items
**/
func (h *Hub) GetChannels(name, queue string) et.Items {
	result := []et.Json{}
	if name == "" {
		for _, channel := range h.channels {
			result = append(result, channel.describe(0))
		}

		for _, queue := range h.queues {
			result = append(result, queue.describe(0))
		}
	} else {
		ch := h.getChannel(name)
		if ch != nil {
			result = append(result, ch.describe(0))
		}

		qu := h.getQueue(name, queue)
		if qu != nil {
			result = append(result, qu.describe(0))
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
		client := h.getClient(key)
		if client != nil {
			result = append(result, client.describe())
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
func (h *Hub) DrainChannel(channel, queue string) error {
	ch := h.getChannel(channel)
	if ch != nil {
		ch.drain()
	}

	qu := h.getQueue(channel, queue)
	if qu != nil {
		qu.drain()
	}

	return nil
}
