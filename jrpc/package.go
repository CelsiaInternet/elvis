package jrpc

import (
	"encoding/json"
	"net"
	"net/rpc"
	"slices"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
)

type Package struct {
	Name    string             `json:"name"`
	Host    string             `json:"host"`
	Port    int                `json:"port"`
	Solvers map[string]*Solver `json:"routes"`
}

/**
* NewPackage
* @param name string, host string, port int
* @return *Package
**/
func NewPackage(name string, host string, port int) *Package {
	return &Package{
		Name:    name,
		Host:    host,
		Port:    port,
		Solvers: make(map[string]*Solver),
	}
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
* Save
* @return error
**/
func (s *Package) Save() error {
	routers, err := getRouters()
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(routers, func(e *Package) bool { return e.Host == s.Host && e.Name == s.Name })
	if idx == -1 {
		routers = append(routers, s)
	} else {
		routers[idx] = s
	}

	err = setRoutes(routers)
	if err != nil {
		return err
	}

	return nil
}

/**
* getRouters
* @return []*Router
* @return error
**/
func getRouters() ([]*Package, error) {
	routers := make([]*Package, 0)
	bt, err := json.Marshal(routers)
	if err != nil {
		return nil, err
	}

	str, err := cache.Get(RPC_KEY, string(bt))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(str), &routers)
	if err != nil {
		return nil, err
	}

	return routers, nil
}

/**
* setRoutes
* @param routers []*Router
* @return error
**/
func setRoutes(routers []*Package) error {
	bt, err := json.Marshal(routers)
	if err != nil {
		return err
	}

	err = cache.Set(RPC_KEY, string(bt), 0)
	if err != nil {
		return err
	}

	return nil
}
