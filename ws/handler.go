package ws

import (
	"errors"
	"net/http"

	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/utility"
)

func Connect(w http.ResponseWriter, r *http.Request) (*Client, error) {
	if conn == nil {
		return nil, logs.Errorm(ERR_NOT_WS_SERVICE)
	}

	ctx := r.Context()
	clientId := utility.NewAny(ctx.Value("clientId")).String()
	if clientId == "<nil>" {
		return nil, errors.New(ERR_NOT_DEFINE_CLIENTID)
	}

	idxC := conn.hub.indexClient(clientId)
	if idxC != -1 {
		return conn.hub.clients[idxC], nil
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	userName := utility.NewAny(ctx.Value("username")).String()

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

func Subscribe(clientId string, channel string) bool {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
		return false
	}

	return conn.hub.Subscribe(clientId, channel)
}

func Unsubscribe(clientId string, channel string) bool {
	if conn == nil {
		logs.Errorm(ERR_NOT_WS_SERVICE)
		return false
	}

	return conn.hub.Unsubscribe(clientId, channel)
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
