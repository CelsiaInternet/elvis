package resilience

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/event"
)

/**
* initEvents
**/
func initEvents() {
	err := event.Subscribe(EVENT_RESILIENCE_STOP, eventStop)
	if err != nil {
		console.Error(err)
	}

	err = event.Subscribe(EVENT_RESILIENCE_RESTART, eventRestart)
	if err != nil {
		console.Error(err)
	}

}

/**
* eventStop
* @param m event.EvenMessage
**/
func eventStop(m event.EvenMessage) {
	data := m.Data
	id := data.Str("id")
	if id == "" {
		console.ErrorM(MSG_ID_REQUIRED)
		return
	}

	Stop(id)
	console.Log("eventStop:", data.ToString())
}

/**
* eventRestart
* @param m event.EvenMessage
**/
func eventRestart(m event.EvenMessage) {
	data := m.Data
	id := data.Str("id")
	if id == "" {
		console.ErrorM(MSG_ID_REQUIRED)
		return
	}

	Restart(id)
	console.Log("eventRestart:", data.ToString())
}
