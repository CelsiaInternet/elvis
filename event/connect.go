package event

import (
	"sync"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
	"github.com/nats-io/nats.go"
)

/**
* ConnectTo
* @param host, user, password string
* @return *Conn, error
**/
func ConnectTo(host, user, password string) (*Conn, error) {
	if !utility.ValidStr(host, 0, []string{}) {
		return nil, utility.NewErrorf(msg.MSG_ATRIB_REQUIRED, "host")
	}

	connect, err := nats.Connect(host, nats.UserInfo(user, password))
	if err != nil {
		return nil, err
	}

	logs.Logf("NATS", `Connected host:%s`, host)

	return &Conn{
		Conn:            connect,
		id:              utility.UUID(),
		eventCreatedSub: map[string]*nats.Subscription{},
		mutex:           &sync.RWMutex{},
	}, nil
}

/**
* connect
* @return *Conn, error
**/
func connect() (*Conn, error) {
	host := envar.GetStr("", "NATS_HOST")
	user := envar.GetStr("", "NATS_USER")
	password := envar.GetStr("", "NATS_PASSWORD")
	result, err := ConnectTo(host, user, password)
	if err != nil {
		return nil, err
	}

	return result, nil
}
