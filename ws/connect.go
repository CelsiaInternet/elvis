package ws

import (
	"net/http"

	"github.com/cgalvisleon/elvis/middleware"
	"github.com/cgalvisleon/elvis/utility"
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
	name := middleware.NameKey.String(ctx, "Anonymous")

	return h.connect(socket, clientId, name)
}
