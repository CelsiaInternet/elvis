package event

import (
	"sync"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/nats-io/nats.go"
)

func connect() (*Conn, error) {
	host := envar.GetStr("", "NATS_HOST")
	if host == "" {
		return nil, logs.Alertf(msg.ERR_ENV_REQUIRED, "NATS_HOST")
	}

	user := envar.GetStr("", "NATS_USER")
	password := envar.GetStr("", "NATS_PASSWORD")

	connect, err := nats.Connect(host, nats.UserInfo(user, password))
	if err != nil {
		return nil, err
	}

	logs.Logf("NATS", `Connected host:%s`, host)

	return &Conn{
		Conn:            connect,
		eventCreatedSub: map[string]*nats.Subscription{},
		mutex:           &sync.RWMutex{},
	}, nil
}
