package ws

import (
	"net/http"
	"net/url"

	"github.com/cgalvisleon/elvis/claim"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/middleware"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/utility"
	"github.com/gorilla/websocket"
)

/**
* ConnectHttp connect to the server using the http
* @param w http.ResponseWriter
* @param r *http.Request
* @return *Client
* @return error
**/
func (h *Hub) ConnectHttp(w http.ResponseWriter, r *http.Request) (*Client, error) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	clientId := middleware.ClientIDKey.String(ctx, utility.UUID())
	name := middleware.NameKey.String(ctx, "Anonimo")

	return h.connect(socket, clientId, name)
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
func ConnectWs(host, scheme, clientId, name string) (*websocket.Conn, error) {
	if scheme == "" {
		scheme = "ws"
	}

	path := strs.Format("/%s", scheme)

	u := url.URL{Scheme: scheme, Host: host, Path: path}
	header := make(http.Header)
	tk, err := claim.GenToken(clientId, "ws", name, "ws", name, "microservice", 0)
	if err != nil {
		return nil, err
	}

	tk = strs.Format(`Bearer %s`, tk)
	header.Add("Authorization", tk)
	result, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return nil, err
	}

	logs.Logf("Real time", "Connected host:%s", host)

	return result, nil
}
