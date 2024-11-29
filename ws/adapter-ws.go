package ws

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
)

type AdapterWS struct {
	conn *Client
}

func NewWSAdapter() Adapter {
	return &AdapterWS{}
}

/**
* ConnectTo
* @param params et.Json
* @return error
**/
func (s *AdapterWS) ConnectTo(params et.Json) error {
	if s.conn != nil {
		return nil
	}

	result, err := NewClient(&ClientConfig{
		Url:       params.Str("url"),
		Reconnect: 3,
	})
	if err != nil {
		return err
	}

	s.conn = result
	logs.Debug("AdapterRedis:", params.ToString())

	return nil
}

/**
* Close
**/
func (s *AdapterWS) Close() {}

/**
* Subscribed
* @param channel string
**/
func (s *AdapterWS) Subscribed(channel string) {
	channel = clusterChannel(channel)
	s.conn.Subscribe(channel, func(msg Message) {
		if msg.tp == TpDirect {
			s.conn.SendMessage(msg.Id, msg)
		} else {
			s.conn.Publish(msg.Channel, msg)
		}
	})
}

/**
* UnSubscribed
* @param sub channel string
**/
func (s *AdapterWS) UnSubscribed(channel string) {
	channel = clusterChannel(channel)
	s.conn.Unsubscribe(channel)
}

/**
* Publish
* @param sub channel string
**/
func (s *AdapterWS) Publish(channel string, msg Message) {
	channel = clusterChannel(channel)
	s.conn.Publish(channel, msg)
}
