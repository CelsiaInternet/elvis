package ws

import (
	"net/http"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/utility"
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

	url := params.Str("url")
	if url == "" {
		return nil
	}

	username := params.Str("username")
	if username == "" {
		return utility.NewError("WS Adapter, username is required")
	}

	password := envar.GetStr("", "WS_PASSWORD")
	if password == "" {
		return utility.NewError("WS Adapter, password is required")
	}

	name := params.ValStr("Anonime", "name")
	result, err := Login(&ClientConfig{
		ClientId:  utility.UUID(),
		Name:      name,
		Url:       url,
		Reconnect: envar.GetInt(3, "RT_RECONCECT"),
		Header: http.Header{
			"username": []string{username},
			"password": []string{password},
		},
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
