package ws

import (
	"errors"
	"net/http"

	"github.com/celsiainternet/elvis/claim"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
)

/**
* HttpCluster connect to the server using the http
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (h *Hub) HttpLogin(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	password := r.Header.Get("password")

	logs.Debug(et.Json{
		"username": username,
		"password": password,
	}.ToString())

	ws_username := envar.GetStr("", "WS_USERNAME")
	if !utility.ValidStr(ws_username, 0, []string{}) {
		response.HTTPError(w, r, http.StatusInternalServerError, errors.New(ERR_NOT_SIGNATURE).Error())
	}

	ws_password := envar.GetStr("", "WS_PASSWORD")
	if !utility.ValidStr(ws_password, 0, []string{}) {
		response.HTTPError(w, r, http.StatusInternalServerError, errors.New(ERR_NOT_SIGNATURE).Error())
	}

	if username != ws_username || password != ws_password {
		response.HTTPError(w, r, http.StatusInternalServerError, errors.New(ERR_NOT_SIGNATURE).Error())
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
	}

	clientId := utility.UUID()
	name := "Anonimo"
	_, err = h.connect(conn, clientId, name)
	if err != nil {
		logs.Alert(err)
	}
}

/**
* HttpConnect connect to the server using the http
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (h *Hub) HttpConnect(w http.ResponseWriter, r *http.Request) {
	query := response.GetQuery(r)
	clientId := query.ValStr("", "clientId")
	if clientId == "" {
		clientId = utility.UUID()
	}

	name := query.ValStr("", "name")
	if name == "" {
		name = "Anonimo"
	}

	ctx := r.Context()
	clientId = claim.ClientIdKey.String(ctx, clientId)
	name = claim.NameKey.String(ctx, name)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
	}

	_, err = h.connect(conn, clientId, name)
	if err != nil {
		logs.Alert(err)
	}
}

/**
* HttpStream connect to the server using the http
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (h *Hub) HttpStream(w http.ResponseWriter, r *http.Request) {
	query := response.GetQuery(r)
	clientId := query.ValStr("", "clientId")
	if clientId == "" {
		clientId = utility.UUID()
	}

	name := query.ValStr("", "name")
	if name == "" {
		name = "Anonimo"
	}

	ctx := r.Context()
	clientId = claim.ClientIdKey.String(ctx, clientId)
	name = claim.NameKey.String(ctx, name)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
	}

	_, err = h.streaming(conn, clientId, name)
	if err != nil {
		logs.Alert(err)
	}
}
