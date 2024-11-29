package cache

import (
	"context"
	"fmt"
)

func (s *Conn) pubCtx(ctx context.Context, channel string, message interface{}) error {
	err := s.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Conn) subCtx(ctx context.Context, channel string, f func(string)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	pubsub := s.Subscribe(ctx, channel)
	s.chanels[channel] = pubsub

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			return
		}

		fmt.Println(msg.Channel, msg.Payload)
		f(msg.Payload)
	}
}

/**
* Pub
* @param channel string
* @param message interface{}
* @return error
**/
func (s *Conn) Pub(channel string, message interface{}) error {
	ctx := context.Background()
	return s.pubCtx(ctx, channel, message)
}

/**
* Sub
* @param channel string
* @param f func(interface{})
**/
func (s *Conn) Sub(channel string, f func(string)) {
	ctx := context.Background()
	s.subCtx(ctx, channel, f)
}

/**
* Unsub
* @param channel string
* @return error
**/
func (s *Conn) Unsub(channel string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	pubsub := s.chanels[channel]
	if pubsub == nil {
		return nil
	}

	return pubsub.Close()
}
