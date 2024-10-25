package router

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/utility"
)

func EventLoad(m event.EvenMessage) {
	for _, item := range router {
		id := item.Key("_id")
		method := item.Str("method")
		path := item.Str("path")
		resolve := item.Str("resolve")
		header := item.Json("header")
		tpHeader := ToTpHeader(item.Int("tpHeader"))
		private := item.Bool("private")
		packageName := item.Str("package_name")

		if !utility.ValidStr(method, 0, []string{""}) {
			console.AlertF(msg.MSG_ATRIB_REQUIRED, "method")
			continue
		}

		if !utility.ValidStr(path, 0, []string{""}) {
			console.AlertF(msg.MSG_ATRIB_REQUIRED, "path")
			continue
		}

		if !utility.ValidStr(resolve, 0, []string{""}) {
			console.AlertF(msg.MSG_ATRIB_REQUIRED, "resolve")
			continue
		}

		if !utility.ValidStr(packageName, 0, []string{""}) {
			console.AlertF(msg.MSG_ATRIB_REQUIRED, "package_name")
			continue
		}

		PushApiGateway(id, method, path, resolve, header, tpHeader, private, packageName)
	}
}
