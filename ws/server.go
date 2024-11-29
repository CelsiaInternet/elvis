package ws

import (
	"net/http"
	"time"

	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/strs"
)

/**
* ServerHttp
* @params port int
* @params mode string
* @params master string
* @params schema string
* @params path string
* @return *Hub
**/
func ServerHttp(port int) *Hub {
	result := NewHub()
	result.Start()
	go startHttp(result, port)
	time.Sleep(1 * time.Second)

	return result
}

func startHttp(hub *Hub, port int) {
	http.HandleFunc("/ws", hub.HttpConnect)
	http.HandleFunc("/ws/describe", hub.HttpDescribe)
	http.HandleFunc("/ws/publications", hub.HttpGetPublications)
	http.HandleFunc("/ws/subscribers", hub.HttpGetSubscribers)

	logs.Logf("WebSocket", "Http server in http://localhost:%d/ws", port)
	addr := strs.Format(`:%d`, port)
	logs.Fatal(http.ListenAndServe(addr, nil))
}
