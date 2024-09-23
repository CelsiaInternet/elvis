package rt

import (
	"net/http"
	"net/url"

	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
	"github.com/cgalvisleon/elvis/ws"
	"github.com/gorilla/websocket"
)

var conn *ClientWS

type ClientWS struct {
	socket    *websocket.Conn
	Host      string
	ClientId  string
	Name      string
	channels  map[string]func(ws.Message)
	connected bool
}

/**
* ConnectWs connect to the server using the websocket
* @param host string
* @param scheme string
* @param clientId string
* @param name string
* @return *websocket.Conn
* @return error
**/
func Load() error {
	if conn != nil {
		return nil
	}

	name := envar.GetStr("Real Time", "RT_NAME")
	host := envar.GetStr("localhost", "RT_HOST")
	scheme := envar.GetStr("ws", "RT_SCHEME")

	path := strs.Format("/%s", scheme)
	u := url.URL{Scheme: scheme, Host: host, Path: path}
	header := http.Header{}
	wsocket, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return err
	}

	conn = &ClientWS{
		socket:    wsocket,
		Host:      host,
		ClientId:  utility.UUID(),
		Name:      name,
		channels:  make(map[string]func(ws.Message)),
		connected: true,
	}

	go conn.read()

	SetFrom(name)

	logs.Logf("Real time", "Connected host:%s", u.String())

	return nil
}

/**
* Close
**/
func Close() {
	if conn != nil {
		conn.socket.Close()
	}
}
