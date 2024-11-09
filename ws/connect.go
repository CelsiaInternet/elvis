package ws

import (
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
)

/**
* HttpConnect connect to the server using the http
* @param w http.ResponseWriter
* @param r *http.Request
* @return *Subscriber
* @return error
**/
func (h *Hub) HttpConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
	}

	// Identify of the client
	clientId := r.URL.Query().Get("clientId")
	if clientId == "" {
		clientId = utility.UUID()
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Anonimo"
	}

	ctx := r.Context()
	clientId = claim.ClientIdKey.String(ctx, clientId)
	name = claim.NameKey.String(ctx, name)

	_, err = h.connect(conn, clientId, name)
	if err != nil {
		logs.Alert(err)
	}
}
