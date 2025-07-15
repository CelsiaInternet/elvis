package cache

import (
	"context"
	"sync"

	"github.com/celsiainternet/elvis/logs"
	"github.com/redis/go-redis/v9"
)

var conn *Conn
var FromId string

type Conn struct {
	*redis.Client
	_id      string
	ctx      context.Context
	host     string
	dbname   int
	channels map[string]bool
	mutex    *sync.RWMutex
}

func Load() (*Conn, error) {
	if conn != nil {
		return conn, nil
	}

	var err error
	conn, err = connect()
	if err != nil {
		return nil, err
	}

	FromId = conn._id

	return conn, nil
}

func Close() {
	if conn == nil {
		return
	}

	conn.Close()

	logs.Log("Cache", `Disconnect...`)
}

func IsLoad() bool {
	return conn != nil
}
