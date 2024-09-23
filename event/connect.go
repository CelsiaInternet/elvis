package event

import (
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/nats-io/nats.go"
)

func connect() (*Conn, error) {
	host := envar.GetStr("", "NATS_HOST")
	if host == "" {
		return nil, logs.Alertf(msg.ERR_ENV_REQUIRED, "NATS_HOST")
	}

	connect, err := nats.Connect(host)
	if err != nil {
		return nil, err
	}

	logs.Logf("NATS", `Connected host:%s`, host)

	return &Conn{
		conn: connect,
	}, nil
}
