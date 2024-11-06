package ws

import (
	"sync"

	"github.com/celsiainternet/elvis/et"
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

func (c *Queue) setQueue(key string, val int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Queue[key] = val
}

func (c *Queue) getQueue(key string) (int, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	result, ok := c.Queue[key]
	return result, ok
}

func (c *Queue) deleteQueue(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.Queue, key)
}

/**
* drain
**/
func (c *Queue) drain() {
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
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, client := range c.Subscribers {
		if client == nil {
			continue
		}

		delete(client.Channels, c.Name)
	}

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
			"subscribers": subscribers,
		}
	}

	return et.Json{
		"name": c.Name,
	}
}

/**
* Count return the number of subscribers
* @return int
**/
func (c *Queue) Count() int {
	return len(c.Subscribers)
}

/**
* nextTurn return the next subscriber
* @return *Subscriber
**/
func (c *Queue) nextTurn(queue string) *Subscriber {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	n := c.Count()
	if n == 0 {
		return nil
	}

	turn, exist := c.getQueue(queue)
	if !exist {
		turn = 0
		c.setQueue(queue, turn)
	}

	if turn >= n {
		turn = 0
		c.setQueue(queue, turn)
	}

	result := c.Subscribers[turn]
	turn++
	c.setQueue(queue, turn)

	return result
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

	_, exist := c.getQueue(queue)
	if !exist {
		c.setQueue(queue, 0)
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
