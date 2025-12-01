package event

import (
	"fmt"
	"sync"
	"time"

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
		return nil, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "host")
	}

	options := []nats.Option{
		nats.UserInfo(user, password),
		nats.ReconnectWait(5 * time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			logs.Logf("NATS", `Disconnected host:%s error:%s`, host, err.Error())
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			logs.Logf("NATS", `Reconnected host:%s`, host)
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			logs.Logf("NATS", `Closed host:%s`, host)
		}),
	}
	connect, err := nats.Connect(host, options...)
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
