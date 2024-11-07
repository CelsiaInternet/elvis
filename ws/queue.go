package ws

import (
	"sync"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"golang.org/x/exp/slices"
)

type Queue struct {
	Name        string         `json:"name"`
	Queue       map[string]int `json:"queue"`
	Subscribers []*Subscriber  `json:"subscribers"`
	mutex       *sync.RWMutex
}

/**
* newQueue
* @param name string
* @return *Queue
**/
func newQueue(name string) *Queue {
	result := &Queue{
		Name:        name,
		Queue:       map[string]int{},
		Subscribers: []*Subscriber{},
		mutex:       &sync.RWMutex{},
	}

	return result
}

/**
* nextTurn return the next subscriber
* @return *Subscriber
**/
func (c *Queue) nextTurn(queue string) *Subscriber {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	n := len(c.Subscribers)
	if n == 0 {
		return nil
	}

	turn, exist := c.Queue[queue]
	if !exist {
		turn = 0
	}

	if turn >= n {
		turn = 0
	}

	result := c.Subscribers[turn]
	turn++
	c.Queue[queue] = turn

	return result
}

func (c *Queue) drain() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, client := range c.Subscribers {
		if client == nil {
			continue
		}

		delete(client.Channels, c.Name)
	}

	c.Subscribers = []*Subscriber{}
}

/**
* close
**/
func (c *Queue) close() {
	c.drain()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Subscribers = nil
	c.Queue = nil
}

/**
* describe return the channel name
* @return et.Json
**/
func (c *Queue) describe(mode int) et.Json {
	if mode == 0 {
		subscribers := []et.Json{}
		for _, subscriber := range c.Subscribers {
			subscribers = append(subscribers, subscriber.From())
		}

		return et.Json{
			"name":        c.Name,
			"type":        "queue",
			"subscribers": subscribers,
		}
	}

	return et.Json{
		"name": c.Name,
		"type": "queue",
	}
}

/**
* queueSubscribe a client to channel
* @param client *Subscriber
**/
func (c *Queue) subscribe(client *Subscriber, queue string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if queue == "" {
		return
	}

	_, exist := c.Queue[queue]
	if !exist {
		c.Queue[queue] = 0
	}

	idx := slices.IndexFunc(c.Subscribers, func(e *Subscriber) bool { return e.Id == client.Id })
	if idx != -1 {
		return
	}

	c.Subscribers = append(c.Subscribers, client)
	client.Queue[c.Name] = c
}

/**
* unsubscribe
* @param clientId string
**/
func (c *Queue) unsubscribe(client *Subscriber) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	idx := slices.IndexFunc(c.Subscribers, func(e *Subscriber) bool { return e.Id == client.Id })
	if idx == -1 {
		return
	}

	c.Subscribers = append(c.Subscribers[:idx], c.Subscribers[idx+1:]...)
	delete(client.Channels, c.Name)
}

/**
* broadcast
* @param msg Message
* @param ignored []string
* @return int
**/
func (c *Queue) broadcast(queue string, msg Message, ignored []string) int {
	result := 0
	msg.Channel = c.Name
	client := c.nextTurn(queue)
	if client != nil && !slices.Contains(ignored, client.Id) {
		err := client.sendMessage(msg)
		if err != nil {
			logs.Alert(err)
		} else {
			result++
		}
	}

	return result
}
