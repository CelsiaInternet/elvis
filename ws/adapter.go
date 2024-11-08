package ws

import (
	"net/http"
	"sync"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/race"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
)

type TypeNode int

const (
	NodeWorker TypeNode = iota
	NodeMaster
)

func (t TypeNode) String() string {
	switch t {
	case NodeMaster:
		return "master"
	case NodeWorker:
		return "worker"
	}

	return "unknown"
}

func (t TypeNode) ToJson() et.Json {
	return et.Json{
		"id":   t,
		"name": t.String(),
	}
}

type Adapter struct {
	Client
	typeNode TypeNode
}

var adapter *Adapter

type AdapterConfig struct {
	Schema    string
	Host      string
	Path      string
	TypeNode  TypeNode
	Reconcect int
	Header    http.Header
}

func clusterChannel(channel string) string {
	result := strs.Format(`cluster/%s`, channel)
	return utility.ToBase64(result)
}

/**
* InitMaster
* @return *Hub
**/
func (h *Hub) InitMaster() {
	if adapter != nil {
		return
	}

	adapter = &Adapter{
		typeNode: NodeMaster,
	}
}

/**
* Join
* @param config *ClientConfig
**/
func (h *Hub) Join(config AdapterConfig) error {
	if adapter != nil {
		return nil
	}

	adapter = &Adapter{
		typeNode: config.TypeNode,
	}
	adapter.Channels = make(map[string]func(Message))
	adapter.Attempts = race.NewValue(0)
	adapter.Connected = race.NewValue(false)
	adapter.mutex = &sync.Mutex{}
	adapter.clientId = h.Id
	adapter.name = h.Name
	adapter.schema = config.Schema
	adapter.host = config.Host
	adapter.path = config.Path
	adapter.header = config.Header
	adapter.reconcect = config.Reconcect
	err := adapter.Connect()
	if err != nil {
		return err
	}

	adapter.SetReconnectCallback(func(c *Client) {
		logs.Debug("ReconnectCallback:", "Hola")
	})

	adapter.SetDirectMessage(func(msg Message) {
		logs.Debug("DirectMessage:", msg.ToString())
	})

	return nil
}

/**
* Live
**/
func (h *Hub) Live() {
	if adapter == nil {
		return
	}

	adapter.Close()
}

/**
* ClusterSubscribed
* @param channel string
**/
func (h *Hub) ClusterSubscribed(channel string) {
	if adapter == nil {
		return
	}

	if !adapter.Connected.Bool() {
		return
	}

	channel = clusterChannel(channel)
	adapter.Subscribe(channel, func(msg Message) {
		if msg.Tp == TpDirect {
			h.SendMessage(msg.Id, msg)
		} else {
			h.Publish(msg.Channel, msg.Queue, msg, msg.Ignored, msg.From)
		}
	})
}

/**
* ClusterUnSubscribed
* @param sub channel string
**/
func (h *Hub) ClusterUnSubscribed(channel string) {
	if adapter == nil {
		return
	}

	if !adapter.Connected.Bool() {
		return
	}

	channel = clusterChannel(channel)
	adapter.Unsubscribe(channel)
}

/**
* ClusterUnSubscribed
* @param sub channel string
**/
func (h *Hub) ClusterPublish(channel string, msg Message) {
	if adapter == nil {
		return
	}

	if !adapter.Connected.Bool() {
		return
	}

	channel = clusterChannel(channel)
	adapter.Publish(channel, msg)
}
