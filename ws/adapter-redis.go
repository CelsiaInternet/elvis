package ws

import (
	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
)

type AdapterRedis struct {
	conn *cache.Conn
}

func NewRedisAdapter() Adapter {
	return &AdapterRedis{}
}

/**
* ConnectTo
* @param params et.Json
* @return error
**/
func (s *AdapterRedis) ConnectTo(params et.Json) error {
	if s.conn != nil {
		return nil
	}

	host := params.Str("host")
	password := params.Str("password")
	dbname := params.Int("dbname")
	result, err := cache.ConnectTo(host, password, dbname)
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
func (s *AdapterRedis) Close() {}

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
