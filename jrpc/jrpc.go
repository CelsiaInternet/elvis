package jrpc

import (
	"net/rpc"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/utilities"
)

type Args struct {
	require Json
}

func RpcCall(host string, port int, method string, args []byte) (Item, error) {
	var reply *[]byte

	client, err := rpc.DialHTTP("tcp", Format(`%s:%d`, host, port))
	if err != nil {
		return Item{}, console.Error(err)
	}
	defer client.Close()

	err = client.Call(method, args, &reply)
	if err != nil {
		return Item{}, console.Error(err)
	}

	result := Json{}.ToItem(*reply)

	return result, nil
}
