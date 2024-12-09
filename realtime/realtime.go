package realtime

import (
	"errors"
	"net/http"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/utility"
	"github.com/celsiainternet/elvis/ws"
)

const ServiceName = "Real Time"

var conn *ws.Client
var FromId string

/**
* Load
* @return erro
**/
func Load(id, name string) (*ws.Client, error) {
	if conn != nil {
		return conn, nil
	}

	url := envar.GetStr("", "RT_URL")
	if url == "" {
		return nil, errors.New(MSG_RT_URL_REQUIRED)
	}

	username := envar.GetStr("", "WS_USERNAME")
	if username == "" {
		return nil, utility.NewError(ERR_WS_USERNAME_REQUIRED)
	}

	password := envar.GetStr("", "WS_PASSWORD")
	if password == "" {
		return nil, utility.NewError(ERR_WS_PASSWORD_REQUIRED)
	}

	client, err := ws.Login(&ws.ClientConfig{
		ClientId:  id,
		Name:      name,
		Url:       strs.Format(`%s?clientId=%s&name=%s`, url, id, name),
		Reconnect: envar.GetInt(3, "RT_RECONCECT"),
		Header: http.Header{
			"username": []string{username},
			"password": []string{password},
		},
	})
	if err != nil {
		return nil, err
	}

	conn = client
	FromId = client.ClientId

	logs.Logf(ServiceName, `Connected host:%s`, url)

	return conn, nil
}

/**
* Close
**/
func Close() {
	if conn == nil {
		return
	}

	conn.Close()
}
