package router

import (
	"github.com/celsiainternet/elvis/event"
)

const (
	APIMANAGER_LOADED = "apimanager/loaded"
)

func EventLoad(m event.EvenMessage) {
	for _, item := range router {
		id := item.Key("_id")
		method := item.Str("method")
		path := item.Str("path")
		resolve := item.Str("resolve")
		header := item.Json("header")
		tpHeader := ToTpHeader(item.Int("tp_header"))
		excludeHeader := item.ArrayStr("exclude_header")
		private := item.Bool("private")
		packageName := item.Str("package_name")

		PushApiGateway(id, method, path, resolve, header, tpHeader, excludeHeader, private, packageName)
	}
}
