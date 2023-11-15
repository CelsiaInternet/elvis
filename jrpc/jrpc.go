package jrpc

import (
	"net/rpc"

	"github.com/cgalvisleon/elvis/console"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/utility"
)

func RpcCall(host string, port int, method string, data e.Json) (e.Item, error) {
	var args []byte = data.ToByte()
	var reply *[]byte

	client, err := rpc.DialHTTP("tcp", utility.Format(`%s:%d`, host, port))
	if err != nil {
		return e.Item{}, console.Error(err)
	}
	defer client.Close()

	err = client.Call(method, args, &reply)
	if err != nil {
		return e.Item{}, console.Error(err)
	}

	result := e.Json{}.ToItem(*reply)

	return result, nil
}
