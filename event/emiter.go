package event

import (
	"fmt"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/timezone"
	"github.com/celsiainternet/elvis/utility"
)

type Handler func(message EvenMessage)

type EventEmiter struct {
	channel chan EvenMessage   `json:"-"`
	events  map[string]Handler `json:"-"`
}

var emiter *EventEmiter

func init() {
	emiter = NewEventEmiter()
	emiter.Start()
}

/**
* NewEventEmiter
* @return *EventEmiter
**/
func NewEventEmiter() *EventEmiter {
	return &EventEmiter{
		channel: make(chan EvenMessage),
		events:  make(map[string]Handler),
	}
}

/**
* EventEmiter
* @param message EvenMessage
**/
func (s *EventEmiter) eventEmiter(message EvenMessage) {
	if s.events == nil {
		s.events = make(map[string]Handler)
	}

	eventEmiter, ok := s.events[message.Channel]
	if !ok {
		logs.Alert(fmt.Errorf("event not found (%s)", message.Channel))
		return
	}

	eventEmiter(message)
}

/**
* Start
**/
func (s *EventEmiter) Start() {
	go func() {
		for message := range s.channel {
			s.eventEmiter(message)
		}
	}()
}

/**
* On
* @param channel string, handler Handler
**/
func (s *EventEmiter) On(channel string, handler Handler) {
	if s.events == nil {
		s.events = make(map[string]Handler)
	}

	s.events[channel] = handler
}

/**
* Emit
* @param channel string, data et.Json
**/
func (s *EventEmiter) Emit(channel string, data et.Json) {
	if s.channel == nil {
		return
	}

	message := EvenMessage{
		Created_at: timezone.NowTime(),
		Id:         utility.UUID(),
		Channel:    channel,
		Data:       data,
	}

	s.channel <- message
}

/**
* On
* @param channel string, handler Handler
**/
func On(channel string, handler Handler) {
	if emiter == nil {
		emiter = NewEventEmiter()
	}

	emiter.On(channel, handler)
}

/**
* Emit
* @param channel string, data et.Json
**/
func Emit(channel string, data et.Json) {
	emiter.Emit(channel, data)
}
