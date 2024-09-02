package gateway

import (
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
)

func initEvents() {
	err := event.Stack("gateway/upsert", eventAction)
	if err != nil {
		console.Error(err)
	}

}

func eventAction(m event.EvenMessage) {
	data, err := et.ToJson(m.Data)
	if err != nil {
		console.Error(err)
	}

	kind := data.ValStr("HTTP", "kind")
	method := data.Str("method")
	path := data.Str("path")
	resolve := data.Str("resolve")
	packageName := data.Str("package")

	conn.http.AddRoute(method, path, resolve, kind, packageName)

	console.LogKF("Api gateway", `[%s] %s -> %s - %s`, method, path, resolve, packageName)
}
