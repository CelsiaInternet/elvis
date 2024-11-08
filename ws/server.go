package ws

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/strs"
)

var conn *Hub

/**
* ServerHttp
* @params port int
* @params mode string
* @params master string
* @params schema string
* @params path string
**/
func ServerHttp(port int, mode, master, schema, path string) {
	if conn != nil {
		return
	}

	conn = NewHub()
	conn.Start()
	switch mode {
	case "master":
		conn.InitMaster()
		if master != "" {
			conn.Join(AdapterConfig{
				Schema:    schema,
				Host:      master,
				Path:      path,
				TypeNode:  NodeMaster,
				Reconcect: 3,
				Header:    http.Header{},
			})
		}
	case "worker":
		if master != "" {
			conn.Join(AdapterConfig{
				Schema:    schema,
				Host:      master,
				Path:      path,
				TypeNode:  NodeWorker,
				Reconcect: 3,
				Header:    http.Header{},
			})
		}
	}

	go startHttp(port)

	time.Sleep(1 * time.Second)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	logs.Log("WebSocket", "Shoutdown server...")
}

func startHttp(port int) {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/ws/describe", conn.HttpDescribe)
	http.HandleFunc("/ws/publications", conn.HttpGetPublications)
	http.HandleFunc("/ws/subscribers", conn.HttpGetSubscribers)

	logs.Logf("WebSocket", "Http server in http://localhost:%d/ws", port)
	addr := strs.Format(`:%d`, port)
	logs.Fatal(http.ListenAndServe(addr, nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := conn.HttpConnect(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
	}
}
