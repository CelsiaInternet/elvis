package ws

import (
	"sync"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"golang.org/x/exp/slices"
)

const QUEUE_STACK = "stack"

type Queue struct {
	Name        string         `json:"name"`
	Queue       map[string]int `json:"queue"`
	Subscribers []*Client      `json:"subscribers"`
	mutex       *sync.RWMutex
}

/**
* newQueue
* @param name string
* @return *Queue
**/
func newQueue(name string) *Queue {
	result := &Queue{
		Name:        strs.Lowcase(name),
		Queue:       map[string]int{},
		Subscribers: []*Client{},
		mutex:       &sync.RWMutex{},
	}

	return result
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
	c.Subscribers = []*Client{}
}

/**
* close
**/
func (c *Queue) close() {
	c.mutex.Lock()         // Bloquea la escritura en Subscribers y Queue
	defer c.mutex.Unlock() // Asegura el desbloqueo al final de la función

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
func (c *Queue) describe() et.Json {
	result, err := et.Object(c)
	if err != nil {
		logs.Error(err)
	}

	return result
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
* @return *Client
**/
func (c *Queue) nextTurn(queue string) *Client {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	n := c.Count()
	if n == 0 {
		return nil
	}

	_, exist := c.Queue[queue]
	if !exist {
		c.Queue[queue] = 0
	}

	turn := c.Queue[queue]
	if turn >= n {
		turn = 0
		c.Queue[queue] = turn
	}

	result := c.Subscribers[turn]
	c.Queue[queue]++

	return result
}

/**
* queueSubscribe a client to channel
* @param client *Client
**/
func (c *Queue) subscribe(client *Client, queue string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if queue == "" {
		return
	}

	_, exist := c.Queue[queue]
	if !exist {
		c.Queue[queue] = 0
	}

	idx := slices.IndexFunc(c.Subscribers, func(e *Client) bool { return e.Id == client.Id })
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
func (c *Queue) unsubscribe(client *Client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	idx := slices.IndexFunc(c.Subscribers, func(e *Client) bool { return e.Id == client.Id })
	if idx == -1 {
		return
	}

	c.Subscribers = append(c.Subscribers[:idx], c.Subscribers[idx+1:]...)
	delete(client.Channels, c.Name)
}
