package ws

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/logs"
)

type AdapterRedis struct {
	conn *cache.Conn
}

/**
* Subscribed
* @param channel string
**/
func (s *AdapterRedis) Subscribed(channel string) {
	channel = clusterChannel(channel)
	s.conn.Sub(channel, func(payload string) {
		msg, err := DecodeMessage([]byte(payload))
		if logs.Alert(err) != nil {
			return
		}

		if msg.tp == TpDirect {
			s.conn.Pub(msg.Id, msg)
		} else {
			s.conn.Pub(msg.Channel, msg)
		}
	})
}

/**
* UnSubscribed
* @param sub channel string
**/
func (s *AdapterRedis) UnSubscribed(channel string) {
	channel = clusterChannel(channel)
	s.conn.Unsub(channel)
}

/**
* Publish
* @param sub channel string
**/
func (s *AdapterRedis) Publish(channel string, msg Message) {
	channel = clusterChannel(channel)
	s.conn.Pub(channel, msg)
}
