package ws

import (
	"sync"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/race"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
)

type Adapter interface {
	ConnectTo(params et.Json) error
	Close()
	Subscribed(channel string)
	UnSubscribed(channel string)
	Publish(channel string, msg Message)
}

func clusterChannel(channel string) string {
	result := strs.Format(`cluster/%s`, channel)
	return utility.ToBase64(result)
}

/**
* NewClient
* @config config ConectPatams
* @return erro
**/
func NewNode(config *ClientConfig) (*Client, error) {
	result := &Client{
		Channels:  make(map[string]func(Message)),
		Attempts:  race.NewValue(0),
		Connected: race.NewValue(false),
		mutex:     &sync.Mutex{},
		url:       config.Url,
		header:    config.Header,
		reconnect: config.Reconnect,
	}

	username := envar.GetStr("", "WS_USERNAME")
	password := envar.GetStr("", "WS_PASSWORD")
	us := utility.ToBase64("username")
	ps := utility.ToBase64("username")
	path := strs.Format(`%s?%s=%s&%s=%s`, us, result.url, ps, username, password)
	err := result.connectTo(path)
	if err != nil {
		return nil, err
	}

	return result, nil
}
