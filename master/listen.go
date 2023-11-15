package master

import (
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
)

func listenSync(res e.Json) {
	idT := res.Str("_idt")
	nodeId := res.Str("nodo")

	node := master.GetNodeByID(nodeId)
	if node == nil {
		return
	}

	go node.SyncIdT(idT)
}

func listenNode(res e.Json) {
	action := res.Str("action")
	nodeId := res.Str("nodo")

	switch utility.Uppcase(action) {
	case "INSERT":
		go master.LoadNodeById(nodeId)
	case "UPDATE":
		go master.LoadNodeById(nodeId)
	case "DELETE":
		go master.UnloadNodeById(nodeId)
	}
}
