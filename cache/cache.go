package cache

import (
	"context"
	"sync"

	"github.com/celsiainternet/elvis/logs"
	"github.com/redis/go-redis/v9"
)

var conn *Conn

type Conn struct {
	*redis.Client
	ctx     context.Context
	host    string
	dbname  int
	chanels map[string]*redis.PubSub
	mutex   *sync.RWMutex
}

func Load() (*Conn, error) {
	if conn != nil {
		return conn, nil
	}

	var err error
	conn, err = Connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Close() {
	if conn == nil {
		return
	}

	conn.Close()

	logs.Log("Cache", `Disconnect...`)
}
