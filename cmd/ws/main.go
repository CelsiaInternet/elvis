package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/ws"
)

var conn *ws.Hub
var client1 *ws.Client
var client2 *ws.Client
var client3 *ws.Client

func main() {
	envar.SetInt("port", 3000, "Port server", "PORT")
	envar.SetStr("mode", "", "Modo cluster master, worker", "MODE")
	envar.SetStr("master", "", "Master host", "MASTER_HOST")
	envar.SetStr("schema", "", "Master host", "MASTER_SCHEMA")
	envar.SetStr("path", "", "Master host", "MASTER_PATH")

	if conn != nil {
		return
	}

	port := envar.GetInt(3600, "PORT")
	mode := envar.GetStr("master", "MODE")
	master := envar.GetStr("", "MASTER_HOST")
	schema := envar.GetStr("ws", "MASTER_HOST")
	path := envar.GetStr("/ws", "MASTER_PATH")

	conn = ws.NewHub()
	conn.Start()
	switch mode {
	case "master":
		conn.InitMaster()
		if master != "" {
			conn.Join(ws.AdapterConfig{
				Schema:    schema,
				Host:      master,
				Path:      path,
				TypeNode:  ws.NodeMaster,
				Reconcect: 3,
				Header:    http.Header{},
			})
		}
	case "worker":
		if master != "" {
			conn.Join(ws.AdapterConfig{
				Schema:    schema,
				Host:      master,
				Path:      path,
				TypeNode:  ws.NodeWorker,
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

	console.LogK("WebSocket", "Shoutdown server...")
}

func startHttp(port int) {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/ws/describe", conn.HttpDescribe)
	http.HandleFunc("/ws/publications", conn.HttpGetPublications)
	http.HandleFunc("/ws/subscribers", conn.HttpGetSubscribers)

	console.LogKF("WebSocket", "Http server in http://localhost:%d/ws", port)
	addr := strs.Format(`:%d`, port)
	console.Fatal(http.ListenAndServe(addr, nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := conn.HttpConnect(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
	}
}

func test1(port int) {
	host := strs.Format(`localhost:%d`, port)

	var err error
	client1, err = ws.NewClient(&ws.ClientConfig{
		ClientId:  "client1",
		Name:      "client1",
		Schema:    "ws",
		Host:      host,
		Path:      "/ws",
		Reconcect: 3,
	})
	if err != nil {
		console.Fatal(err)
	}

	client2, err = ws.NewClient(&ws.ClientConfig{
		ClientId:  "client2",
		Name:      "client2",
		Schema:    "ws",
		Host:      host,
		Path:      "/ws",
		Reconcect: 3,
	})
	if err != nil {
		console.Fatal(err)
	}

	client3, err = ws.NewClient(&ws.ClientConfig{
		ClientId:  "client3",
		Name:      "client3",
		Schema:    "ws",
		Host:      host,
		Path:      "/ws",
		Reconcect: 3,
	})
	if err != nil {
		console.Fatal(err)
	}

	client1.SetDirectMessage(func(msg ws.Message) {
		console.Debug("DirectMessage:", msg.ToString())
	})

	client1.Subscribe("Hola", func(msg ws.Message) {
		console.Debug("client1", msg.ToString())
	})

	client2.Subscribe("Hola", func(msg ws.Message) {
		console.Debug("client2", msg.ToString())
	})

	client3.Subscribe("Hola", func(msg ws.Message) {
		console.Debug("client3:", msg.ToString())
	})

	client1.Stack("cola", func(msg ws.Message) {
		console.Debug("client1", msg.ToString())
	})

	client2.Stack("cola", func(msg ws.Message) {
		console.Debug("client2", msg.ToString())
	})

	client3.Stack("cola", func(msg ws.Message) {
		console.Debug("client3:", msg.ToString())
	})

	t := time.Duration(100)
	n := 1
	sendTest1 := func() {
		if n%2 == 0 {
			go client1.SendMessage("client2", "Hello")
		} else {
			go client1.SendMessage("client3", "Hello")
		}
	}

	sendTest2 := func() {
		if n%2 == 0 {
			go client2.SendMessage("client1", "Hello")
		} else {
			go client2.SendMessage("client3", "Hello")
		}
	}

	sendTest3 := func() {
		if n%2 == 0 {
			go client3.SendMessage("client2", "Hello")
		} else {
			go client3.SendMessage("client1", "Hello")
		}
	}

	sendTest4 := func() {
		if n%2 == 0 {
			go client1.Publish("Hola", "Ping")
			go client2.Publish("Hola", "Pong")
			go client3.Publish("Hola", "Ping")
		} else {
			go client1.Publish("Hola", "Pong")
			go client2.Publish("Hola", "Ping")
			go client3.Publish("Hola", "Pong")
		}
	}

	for i := 0; i < n; i++ {
		go sendTest1()
		time.Sleep(t * time.Millisecond)
		go sendTest2()
		time.Sleep(t * time.Millisecond)
		go sendTest3()
		time.Sleep(t * time.Millisecond)
		go sendTest4()
	}
}

func test2(port int) {
	host := strs.Format(`localhost:%d`, port)

	var err error
	client1, err = ws.NewClient(&ws.ClientConfig{
		ClientId:  "client1",
		Name:      "client1",
		Schema:    "ws",
		Host:      host,
		Path:      "/ws",
		Reconcect: 5,
	})
	if err != nil {
		console.Fatal(err)
	}

	client1.SetDirectMessage(func(msg ws.Message) {
		console.Debug("DirectMessage:", msg.ToString())
	})

	client1.Subscribe("Hola", func(msg ws.Message) {
		console.DebugF("Channel:%s :: %s", msg.Channel, msg.ToString())
	})

}
