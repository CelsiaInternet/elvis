package ws

import (
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/et"
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
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
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

	_, err = h.connect(socket, clientId, name)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
	}

	response.JSON(w, r, http.StatusOK, et.Json{"message": "Connected"})
}

/**
* HttpGetPublications
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (h *Hub) HttpDescribe(w http.ResponseWriter, r *http.Request) {
	result := h.Describe()

	response.JSON(w, r, http.StatusOK, result)
}

/**
* HttpGetPublications
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (h *Hub) HttpGetPublications(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	queue := r.URL.Query().Get("queue")
	items := h.GetChannels(name, queue)

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
