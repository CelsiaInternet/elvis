package ws

import (
	"errors"
	"net/http"

	"github.com/cgalvisleon/elvis/generic"
	"github.com/cgalvisleon/elvis/logs"
)

func Connect(w http.ResponseWriter, r *http.Request) (*Client, error) {
	if conn == nil {
		return nil, logs.Errorm(ERR_NOT_WS_SERVICE)
	}

	ctx := r.Context()
	clientId := generic.New(ctx.Value("clientId"))
	if clientId.IsNil() {
		return nil, errors.New(ERR_NOT_DEFINE_CLIENTID)
	}

	idxC := conn.hub.indexClient(clientId.Str())
	if idxC != -1 {
		return conn.hub.clients[idxC], nil
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	userName := generic.New(ctx.Value("username"))
	if userName.IsNil() {
		return nil, errors.New(ERR_NOT_DEFINE_USERNAME)
	}

	return conn.hub.connect(socket, clientId.Str(), userName.Str())
}

func Broadcast(message interface{}, ignoreId string) error {
	if conn == nil {
		return errors.New(ERR_NOT_WS_SERVICE)
	}

	conn.hub.Broadcast(message, ignoreId)

	return nil
}

func Publish(channel string, message interface{}, ignoreId string) error {
	if conn == nil {
		return errors.New(ERR_NOT_WS_SERVICE)
	}

	conn.hub.Publish(channel, message, ignoreId)

	return nil
}

func SendMessage(clientId, channel string, message interface{}) (bool, error) {
	if conn == nil {
		return false, errors.New(ERR_NOT_WS_SERVICE)
	}

	result := conn.hub.SendMessage(clientId, channel, message)

	return result, nil
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
