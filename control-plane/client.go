package controlplane

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
)

type Client int

var nodeID *NodeInfo

func NewNodeID() *NodeInfo {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	return &NodeInfo{
		ID:       0,
		LastSeen: time.Now(),
		Host:     host,
		Port:     envar.GetInt(4800, "CP_PORT"),
	}
}

func init() {
	if nodeID != nil {
		return
	}

	nodeID = NewNodeID()
}

/**
* Ping
* @param args et.Json, reply *int
* @return error
**/
func (s *Client) Ping(args et.Json, reply *int) error {
	if nodeID == nil {
		nodeID = NewNodeID()
	}

	*reply = nodeID.ID
	return nil
}

/**
* LoadClient
* @return error
**/
func LoadClient() error {
	client := new(Client)
	err := rpc.Register(client)
	if err != nil {
		return err
	}

	port := envar.GetInt(4800, "CP_PORT")
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	console.LogKF("Control Plane", "Client initialized: %s", address)

	go startRPC(listener)

	return nil
}

/**
* GetNodeID
* @param name string, maxNodes int, serverHost string, serverPort int
* @return int, error
**/
func GetNodeID(name string, maxNodes int, serverHost string, serverPort int) (int, error) {
	address := fmt.Sprintf("%s:%d", serverHost, serverPort)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	args := et.Json{
		"name":      name,
		"max_nodes": maxNodes,
		"host":      nodeID.Host,
		"port":      nodeID.Port,
	}
	var result int

	err = client.Call("Server.GetNodeID", args, &result)
	if err != nil {
		return 0, err
	}

	console.LogF("GetNodeID: %d", result)
	nodeID.ID = result
	nodeID.LastSeen = time.Now()

	return result, nil
}
