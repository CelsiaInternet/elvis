package cache

import (
	"context"

	"github.com/celsiainternet/elvis/console"
	"github.com/redis/go-redis/v9"
)

var conn *Conn

type Conn struct {
	ctx    context.Context
	host   string
	dbname int
	db     *redis.Client
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

	return conn, nil
}

func Close() {
	if conn != nil && conn.db != nil {
		conn.db.Close()
	}

	console.LogK("Cache", `Disconnect...`)
}
