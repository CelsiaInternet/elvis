package controlplane

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/response"
)

// Control Plane Structs
type NodeInfo struct {
	ID         int
	InstanceID string
	LastSeen   time.Time
	Service    string
	Host       string
	Port       int
}

type ControlPlane struct {
	Nodes    map[int]*NodeInfo
	MaxNodes int
	mu       sync.RWMutex
}

/**
* getNodeID
* @param name string, maxNodes int, host string, port int, service, instanceID string
* @return int, error
**/
func getNodeID(name string, maxNodes int, host string, port int, service, instanceID string) (int, error) {
	cp, err := load(name)
	if err != nil {
		return 0, err
	}

	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.MaxNodes = maxNodes
	id := len(cp.Nodes) + 1
	if id > maxNodes {
		id, err = lifeProof(cp)
		if err != nil {
			return 0, err
		}
	}

	cp.Nodes[id] = &NodeInfo{
		ID:         id,
		InstanceID: instanceID,
		LastSeen:   time.Now(),
		Service:    service,
		Host:       host,
		Port:       port,
	}

	if err := save(name, cp); err != nil {
		return 0, err
	}

	console.LogF("Node %s:%d registered with ID %d", host, port, id)

	return id, nil
}

/**
* lifeProof
* @param cp *ControlPlane, name string, maxNodes int, host string
* @return int, error
**/
func lifeProof(cp *ControlPlane) (int, error) {
	for _, node := range cp.Nodes {
		id, err := ping(node.Host, node.Service, node.Port)
		if err != nil {
			return node.ID, err
		}
		if id != node.InstanceID {
			return node.ID, nil
		}
	}

	cp.MaxNodes = len(cp.Nodes) + 1

	return cp.MaxNodes, errors.New("no node found")
}

/**
* ping
* @param host, service string, port int
* @return string, error
**/
func ping(host, service string, port int) (string, error) {
	address := fmt.Sprintf("%s.%s:%d", host, service, port)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return "", err
	}
	defer client.Close()

	args := et.Json{}
	var result string

	err = client.Call("Client.Ping", args, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

/**
* startRPC
* @param listener net.Listener
* @return error
**/
func startRPC(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}

type Server int

/**
* GetNodeID
* @param args et.Json, reply int
* @return error
**/
func (s *Server) GetNodeID(args et.Json, reply *int) error {
	name := args.Str("name")
	maxNodes := args.Int("max_nodes")
	host := args.Str("host")
	port := args.Int("port")
	instanceID := args.Str("instance_id")
	service := args.Str("service")
	id, err := getNodeID(name, maxNodes, host, port, service, instanceID)
	if err != nil {
		return err
	}

	*reply = id
	return nil
}

/**
* LoadServer
* @return error
**/
func LoadServer() error {
	server := new(Server)
	err := rpc.Register(server)
	if err != nil {
		return err
	}

	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	port := envar.GetInt(4800, "CP_PORT")
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	console.LogKF("Control Plane", "Server initialized: %s:%d", host, port)

	go startRPC(listener)

	return nil
}

/**
* HTTPReset
* @param w http.ResponseWriter
* @param r *http.Request
**/
func HTTPReset(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	name := body.ValStr("", "name")
	err := Reset(name)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok: true,
		Result: et.Json{
			"message": "Control plane reset",
		},
	})
}
