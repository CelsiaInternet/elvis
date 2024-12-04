package ws

import (
	"github.com/celsiainternet/elvis/et"
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
