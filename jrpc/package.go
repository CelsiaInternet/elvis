package jrpc

import (
	"encoding/json"
	"slices"

	"github.com/celsiainternet/elvis/cache"
)

type Package struct {
	Name    string             `json:"name"`
	Host    string             `json:"host"`
	Port    int                `json:"port"`
	Solvers map[string]*Solver `json:"routes"`
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

	idx := slices.IndexFunc(routers, func(e *Package) bool { return e.Host == s.Host && e.Port == s.Port })
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
