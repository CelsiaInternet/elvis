package cache

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
	"github.com/redis/go-redis/v9"
)

/**
* ConnectTo
* @param host, password string, dbname int
* @return *Conn, error
**/
func ConnectTo(host, password string, dbname int) (*Conn, error) {
	if !utility.ValidStr(host, 0, []string{}) {
		return nil, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "redist_host")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       dbname,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	logs.Logf("Redis", "Connected host:%s", host)

	return &Conn{
		Client:   client,
		_id:      utility.UUID(),
		ctx:      ctx,
		host:     host,
		dbname:   dbname,
		channels: make(map[string]bool),
		mutex:    &sync.RWMutex{},
	}, nil
}

/**
* connect
* @return *Conn, error
**/
func connect() (*Conn, error) {
	host := envar.GetStr("", "REDIS_HOST")
	password := envar.GetStr("", "REDIS_PASSWORD")
	dbname := envar.GetInt(0, "REDIS_DB")

	if !utility.ValidStr(host, 0, []string{}) {
		return nil, logs.Alertf(msg.ERR_ENV_REQUIRED, "REDIS_HOST")
	}

	if !utility.ValidStr(password, 0, []string{}) {
		return nil, logs.Alertf(msg.ERR_ENV_REQUIRED, "REDIS_PASSWORD")
	}

	return ConnectTo(host, password, dbname)
}
