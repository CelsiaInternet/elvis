package cache

import (
	"context"

	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/logs"
	"github.com/cgalvisleon/elvis/msg"
	"github.com/redis/go-redis/v9"
)

func connect() (*Conn, error) {
	host := envar.GetStr("", "REDIS_HOST")
	password := envar.GetStr("", "REDIS_PASSWORD")
	dbname := envar.GetInt(0, "REDIS_DB")

	if host == "" {
		return nil, logs.Alertf(msg.ERR_ENV_REQUIRED, "REDIS_HOST")
	}

	if password == "" {
		return nil, logs.Alertf(msg.ERR_ENV_REQUIRED, "REDIS_PASSWORD")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       dbname,
	})

	logs.Logf("Redis", "Connected host:%s", host)

	return &Conn{
		ctx:    context.Background(),
		host:   host,
		dbname: dbname,
		db:     client,
	}, nil
}
