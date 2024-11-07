package ws

import (
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
)

/**
* ConnectHttp connect to the server using the http
* @param w http.ResponseWriter
* @param r *http.Request
* @return *Subscriber
* @return error
**/
func (h *Hub) ConnectHttp(w http.ResponseWriter, r *http.Request) (*Subscriber, error) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

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

	return h.connect(socket, clientId, name)
}

/**
* HttpGetPublications
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (h *Hub) HttpGetPublications(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	items := h.GetChannels(key)

	response.ITEMS(w, r, http.StatusOK, items)
}

/**
* HttpGetSubscribers
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (h *Hub) HttpGetSubscribers(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	items := h.GetClients(key)

	response.ITEMS(w, r, http.StatusOK, items)
}
