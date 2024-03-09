package apigateway

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
)

func initEvents() {
	console.LogK("Events", "Running svents stack")

	err := event.Stack("apigateway/upsert", eventAction)
	if err != nil {
		console.Error(err)
	}

}

func eventAction(m event.CreatedEvenMessage) {
	data, err := et.ToJson(m.Data)
	if err != nil {
		console.Error(err)
	}

	method := data.Str("method")
	path := data.Str("path")
	resolve := data.Str("resolve")

	AddRoute(method, path, resolve)

	console.LogK("Event", m.Channel)
}
