package jrpc

import (
	"net"
	"net/rpc"

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

	host := envar.GetStr("localhost", "HOST")
	host = envar.GetStr(host, "RPC_HOST")
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
func (s *Package) Start() error {
	address := strs.Format(`:%d`, s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logs.Fatal(err)
	}

	logs.Logf("Rpc", `Running on %s%s`, s.Host, listener.Addr())

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
