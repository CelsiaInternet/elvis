package jrpc

import (
	"net/rpc"

	"github.com/cgalvisleon/elvis/console"
	js "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
)

func RpcCall(host string, port int, method string, args []byte) (js.Item, error) {
	var reply *[]byte

	client, err := rpc.DialHTTP("tcp", utility.Format(`%s:%d`, host, port))
	if err != nil {
		return js.Item{}, console.Error(err)
	}
	defer client.Close()

	err = client.Call(method, args, &reply)
	if err != nil {
		return js.Item{}, console.Error(err)
	}

	result := js.Json{}.ToItem(*reply)

	return result, nil
}
