package ws

import (
	"net/http"
	"sync"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/race"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
)

type TypeNode int

const (
	NotNode TypeNode = iota
	NodeMaster
	NodeWorker
)

func (t TypeNode) String() string {
	switch t {
	case NotNode:
		return "notnode"
	case NodeMaster:
		return "master"
	case NodeWorker:
		return "worker"
	}

	return "unknown"
}

type AdapterConfig struct {
	Schema    string
	Host      string
	Path      string
	TypeNode  TypeNode
	Reconcect int
}

func clusterChannel(channel string) string {
	result := strs.Format(`cluster/%s`, channel)
	return utility.ToBase64(result)
}

/**
* Join
* @param config *ClientConfig
**/
func (h *Hub) Join(config AdapterConfig) error {
	if h.master != nil {
		return nil
	}
	client := &Client{
		Channels:  make(map[string]func(Message)),
		Attempts:  race.NewValue(0),
		Connected: race.NewValue(false),
		mutex:     &sync.Mutex{},
		clientId:  h.Id,
		name:      h.Name,
		schema:    config.Schema,
		host:      config.Host,
		path:      config.Path,
		header: http.Header{
			"Authorization": []string{"Bearer " + h.token},
		},
		reconcect: config.Reconcect,
		typeNode:  config.TypeNode,
	}
	err := client.Connect()
	if err != nil {
		return err
	}

	h.master = client
	h.TypeNode = config.TypeNode

	h.SetClusterConnected(func(clientId string) {
		channel := clusterChannel(clientId)
		h.master.Subscribe(channel, func(msg Message) {
			h.SendMessage(msg.Id, msg)
		})
	})

	h.SetClusterSubscribed(func(channel string) {
		channel = clusterChannel(channel)
		h.master.Subscribe(channel, func(msg Message) {
			h.Publish(msg.Channel, msg.Queue, msg, msg.Ignored, msg.From)
		})
	})

	h.SetClusterUnSubscribed(func(channel string) {
		channel = clusterChannel(channel)
		h.master.Unsubscribe(channel)
	})

	h.master.SetReconnectCallback(func(c *Client) {
		logs.Debug("ReconnectCallback:", "Hola")
	})

	h.master.SetDirectMessage(func(msg Message) {
		logs.Debug("DirectMessage:", msg.ToString())
	})

	return nil
}

/**
* SetClusterConnected
* @param fn func(*Subscriber)
**/
func (h *Hub) SetClusterConnected(fn func(string)) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clusterConnected = fn
}

/**
* ClusterConnected
* @param sub *Subscriber
**/
func (h *Hub) ClusterConnected(clienId string) {
	if h.clusterConnected != nil {
		h.clusterConnected(clienId)
	}
}

/**
* SetClusterUnSubscribed
* @param fn func(string)
**/
func (h *Hub) SetClusterUnSubscribed(fn func(string)) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clusterUnSubscribed = fn
}

/**
* ClusterUnSubscribed
* @param sub channel string
**/
func (h *Hub) ClusterUnSubscribed(channel string) {
	if h.clusterUnSubscribed != nil {
		h.clusterUnSubscribed(channel)
	}
}

/**
* ClusterUnSubscribed
* @param sub channel string
**/
func (h *Hub) ClusterPublish(channel string, msg Message) {
	if h.master != nil {
		channel = clusterChannel(channel)
		h.master.Publish(channel, msg)
	}
}

/**
* SetClusterSubscribed
* @param fn func(string)
**/
func (h *Hub) SetClusterSubscribed(fn func(string)) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clusterSubscribed = fn
}

/**
* ClusterSubscribed
* @param channel string
**/
func (h *Hub) ClusterSubscribed(channel string) {
	if h.clusterSubscribed != nil {
		h.clusterSubscribed(channel)
	}
}
