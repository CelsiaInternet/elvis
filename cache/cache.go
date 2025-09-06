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

/**
* Load
* @return *Conn, error
**/
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

/**
* Close
* @return void
**/
func Close() {
	if conn == nil {
		return
	}

	conn.Close()

	logs.Log("Cache", `Disconnect...`)
}

/**
* IsLoad
* @return bool
**/
func IsLoad() bool {
	return conn != nil
}

/**
* HealthCheck
* @return bool
**/
func HealthCheck() bool {
	if conn == nil {
		return false
	}

	err := conn.Ping(conn.ctx).Err()
	if err != nil {
		return false
	}

	return true
}
