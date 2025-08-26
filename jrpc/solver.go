package jrpc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/celsiainternet/elvis/logs"
)

type Solver struct {
	PackageName string   `json:"packageName"`
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Method      string   `json:"method"`
	Inputs      []string `json:"inputs"`
	Output      []string `json:"outputs"`
}

/**
* Mount
* @param host string
* @param port int
* @param service any
**/
func Mount(services any) error {
	if pkg == nil {
		return logs.Alertm(ERR_PACKAGE_NOT_FOUND)
	}

	return pkg.Mount(services)
}

/**
* UnMount
* @return error
**/
func UnMount(host, name string) error {
	routers, err := getRouters()
	if err != nil {
		return logs.Alert(err)
	}

	idx := slices.IndexFunc(routers, func(e *Package) bool { return e.Name == name && e.Host == host })
	if idx != -1 {
		routers = append(routers[:idx], routers[idx+1:]...)
	}

	err = setRoutes(routers)
	if err != nil {
		return logs.Alert(err)
	}

	return nil
}

/**
* GetSolver
* @param method string
* @return *Solver
* @return error
**/
func GetSolver(method string) (*Solver, error) {
	method = strings.TrimSpace(method)
	routers, err := getRouters()
	if err != nil {
		return nil, err
	}

	lst := strings.Split(method, ".")
	if len(lst) != 3 {
		return nil, fmt.Errorf(ERR_METHOD_NAME_INVALID, method)
	}

	packageName := lst[0]
	idx := slices.IndexFunc(routers, func(e *Package) bool { return e.Name == packageName })
	if idx == -1 {
		return nil, fmt.Errorf(ERR_METHOD_NOT_FOUND, method)
	}

	router := routers[idx]
	solver := router.Solvers[method]

	if solver == nil {
		return nil, fmt.Errorf(ERR_METHOD_NOT_FOUND, method)
	}

	return solver, nil
}
