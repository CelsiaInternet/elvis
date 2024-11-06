package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/ws"
)

var conn *ws.Hub

func main() {
	if conn != nil {
		return
	}

	conn = ws.NewHub()
	conn.Start()

	go startHttp()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	console.LogK("WebSocket", "Shoutdown server...")
}

func startHttp() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/ws/channels", conn.HttpGetChannels)
	http.HandleFunc("/ws/clients", conn.HttpGetClients)
	console.LogK("WebSocket", "Http server in http://localhost:3500/ws")
	console.Fatal(http.ListenAndServe(":3500", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := conn.ConnectHttp(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
	}
}
