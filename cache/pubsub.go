package cache

import (
	"context"
	"fmt"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
)

func pubCtx(ctx context.Context, channel string, message interface{}) error {
	if conn == nil {
		return logs.Errorm(msg.ERR_NOT_CACHE_SERVICE)
	}

	err := conn.db.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func Pub(channel string, message interface{}) error {
	ctx := context.Background()
	return pubCtx(ctx, channel, message)
}

func subCtx(ctx context.Context, channel string, f func(interface{})) {
	if conn == nil {
		return
	}

	pubsub := conn.db.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		f(msg.Payload)
	}
}

func Sub(channel string, f func(interface{})) {
	ctx := context.Background()
	subCtx(ctx, channel, f)
}
