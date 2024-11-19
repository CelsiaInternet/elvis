package jrpc

import (
	"net"
	"net/rpc"
	"net/url"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
)

const RPC_KEY = "apigateway-rpc"

var pkg *Package

/**
* load
**/
func Load(name string) (*Package, error) {
	if pkg != nil {
		return pkg, nil
	}

	_, err := cache.Load()
	if err != nil {
		return nil, err
	}

	host := envar.GetStr("localhost", "RPC_HOST")
	parsedURL, err := url.Parse(host)
	if err != nil {
		host = "localhost"
	} else {
		host = parsedURL.Hostname()
	}
	port := envar.GetInt(4200, "RPC_PORT")
	name = strs.DaskSpace(name)

	pkg = &Package{
		Name:    name,
		Host:    host,
		Port:    port,
		Solvers: make(map[string]*Solver),
	}

	return pkg, nil
}

/**
* Start
**/
func Start() error {
	if pkg == nil {
		return logs.NewError(ERR_SERVER_NOT_FOUND)
	}

	address := strs.Format(`:%d`, pkg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logs.Fatal(err)
	}

	logs.Logf("Rpc", `Running on %s%s`, pkg.Host, listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			logs.Panic(err.Error())
			continue
		}

		go rpc.ServeConn(conn)
	}
}

/**
* Close
**/
func Close() {
	if pkg != nil {
		UnMount()
	}

	logs.Log("Rpc", `Shutting down server...`)
}
