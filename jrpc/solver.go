package jrpc

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/celsiainternet/elvis/envar"
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

var (
	pipeHost string
)

func init() {
	pipeHost = envar.GetStr("", "PIPE_HOST")
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
	packages, err := getPackages()
	if err != nil {
		return logs.Alert(err)
	}

	idx := slices.IndexFunc(packages, func(e *Package) bool { return e.Name == name && e.Host == host })
	if idx != -1 {
		packages = append(packages[:idx], packages[idx+1:]...)
	}

	err = setPackages(packages)
	if err != nil {
		return logs.Alert(err)
	}

	return nil
}

/**
* SetPipeHost
* @param host string
**/
func SetPipeHost(host string) {
	pipeHost = host
}

/**
* GetSolver
* @param method string
* @return *Solver
* @return error
**/
func GetSolver(method string) (*Solver, error) {
	method = strings.TrimSpace(method)
	methodList := strings.Split(method, ".")
	if len(methodList) != 3 {
		return nil, fmt.Errorf(ERR_METHOD_NAME_INVALID, method)
	}

	if pipeHost != "" {
		pipeParam := strings.Split(pipeHost, ":")
		if len(pipeParam) != 2 {
			return nil, fmt.Errorf("PIPE_HOST format is invalid <host>:<port>")
		}

		pHost := pipeParam[0]
		pPort, err := strconv.Atoi(pipeParam[1])
		if err != nil {
			return nil, fmt.Errorf("PIPE_HOST port is invalid")
		}
		result := &Solver{
			PackageName: methodList[0],
			Host:        pHost,
			Port:        pPort,
			Method:      method,
			Inputs:      []string{},
			Output:      []string{},
		}
		return result, nil
	}

	packages, err := getPackages()
	if err != nil {
		return nil, err
	}

	packageName := methodList[0]
	idx := slices.IndexFunc(packages, func(e *Package) bool { return e.Name == packageName })
	if idx == -1 {
		return nil, fmt.Errorf(ERR_METHOD_NOT_FOUND, method)
	}

	router := packages[idx]
	solver := router.Solvers[method]

	if solver == nil {
		return nil, fmt.Errorf(ERR_METHOD_NOT_FOUND, method)
	}

	return solver, nil
}
