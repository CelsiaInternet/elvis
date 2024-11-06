package ws

import (
	"sync"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"golang.org/x/exp/slices"
)

/**
* Channel
**/
type Channel struct {
	Name        string    `json:"name"`
	Subscribers []*Client `json:"subscribers"`
	mutex       *sync.RWMutex
}

/**
* newChannel
* @param name string
* @return *Channel
**/
func newChannel(name string) *Channel {
	result := &Channel{
		Name:        strs.Lowcase(name),
		Subscribers: []*Client{},
		mutex:       &sync.RWMutex{},
	}

	return result
}

/**
* drain
**/
func (c *Channel) drain() {
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
func (c *Channel) close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, client := range c.Subscribers {
		if client == nil {
			continue
		}

		delete(client.Channels, c.Name)
	}

	c.Subscribers = nil
}

/**
* describe return the channel name
* @return et.Json
**/
func (c *Channel) describe() et.Json {
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
func (c *Channel) Count() int {
	return len(c.Subscribers)
}

/**
* queueSubscribe a client to channel
* @param client *Client
**/
func (c *Channel) subscribe(client *Client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	idx := slices.IndexFunc(c.Subscribers, func(e *Client) bool { return e.Id == client.Id })
	if idx != -1 {
		return
	}

	c.Subscribers = append(c.Subscribers, client)
	client.Channels[c.Name] = c
}

/**
* unsubscribe
* @param clientId string
**/
func (c *Channel) unsubscribe(client *Client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	idx := slices.IndexFunc(c.Subscribers, func(e *Client) bool { return e.Id == client.Id })
	if idx == -1 {
		return
	}

	c.Subscribers = append(c.Subscribers[:idx], c.Subscribers[idx+1:]...)
	delete(client.Channels, c.Name)
}
