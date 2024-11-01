package ws

import (
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/utility"
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

	clientId := r.URL.Query().Get("clientId")
	name := r.URL.Query().Get("name")
	if clientId == "" {
		clientId = utility.UUID()
	}
	if name == "" {
		name = "Anonimo"
	}

	ctx := r.Context()
	clientId = claim.ClientIdKey.String(ctx, clientId)
	name = claim.NameKey.String(ctx, name)

	return h.connect(socket, clientId, name)
}
