package jrpc

import (
	"encoding/json"
	"net"
	"net/rpc"
	"reflect"
	"slices"
	"strings"

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
* AddSolver
* @param method string, solver *Solver
* @return error
**/
func (s *Package) Mount(services any) error {
	tipoStruct := reflect.TypeOf(services)
	structName := tipoStruct.String()
	list := strings.Split(structName, ".")
	structName = list[len(list)-1]
	for i := 0; i < tipoStruct.NumMethod(); i++ {
		metodo := tipoStruct.Method(i)
		numInputs := metodo.Type.NumIn()
		numOutputs := metodo.Type.NumOut()

		inputs := []string{}
		for i := 0; i < numInputs; i++ {
			inputs = append(inputs, metodo.Type.In(i).String())
		}

		outputs := []string{}
		for o := 0; o < numOutputs; o++ {
			outputs = append(outputs, metodo.Type.Out(o).String())
		}

		structName = strs.DaskSpace(structName)
		name := strs.DaskSpace(metodo.Name)
		method := strs.Format(`%s.%s`, structName, name)
		key := strs.Format(`%s.%s.%s`, s.Name, structName, name)
		solver := &Solver{
			PackageName: s.Name,
			Host:        s.Host,
			Port:        s.Port,
			Method:      method,
			Inputs:      inputs,
			Output:      outputs,
		}
		s.Solvers[key] = solver
	}

	rpc.Register(services)

	return s.Save()
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
