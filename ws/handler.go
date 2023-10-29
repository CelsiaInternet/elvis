package ws

import (
	"net/http"

	"github.com/cgalvisleon/elvis/logs"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
)

func Connect(w http.ResponseWriter, r *http.Request) (*Client, error) {
	if conn == nil {
		return nil, logs.Errorm(ERR_NOT_WS_SERVICE)
	}

	ctx := r.Context()
	clientId := NewAny(ctx.Value("clientId")).String()

	idxC := conn.hub.indexClient(clientId)
	if idxC != -1 {
		return conn.hub.clients[idxC], nil
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	userName := NewAny(ctx.Value("username")).String()

	return conn.hub.connect(socket, clientId, userName)
}

func Broadcast(message interface{}, ignoreId string) {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
	}

	conn.hub.Broadcast(message, ignoreId)
}

func Publish(channel string, message interface{}, ignoreId string) {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
	}

	conn.hub.Publish(channel, message, ignoreId)
}

func SendMessage(clientId, channel string, message interface{}) bool {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
		return false
	}

	return conn.hub.SendMessage(clientId, channel, message)
}

func Subcribe(clientId string, channel string) bool {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
		return false
	}

	return conn.hub.Subcribe(clientId, channel)
}

func Unsubcribe(clientId string, channel string) bool {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
		return false
	}

	return conn.hub.Unsubcribe(clientId, channel)
}

func GetChannels() []*Channel {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
		return []*Channel{}
	}

	return conn.hub.channels
}

func GetClients() []*Client {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
		return []*Client{}
	}

	return conn.hub.clients
}

func GetSubscribers(channel string) []*Client {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
		return []*Client{}
	}

	return conn.hub.GetSubscribers(channel)
}
