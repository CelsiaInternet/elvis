package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/ws"
)

var conn *ws.Hub
var client1 *ws.Client
var client2 *ws.Client
var client3 *ws.Client

func main() {
	if conn != nil {
		return
	}

	conn = ws.NewHub()
	conn.Start()

	go startHttp()
	go startHttp2()

	time.Sleep(3 * time.Second)
	test1()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	console.LogK("WebSocket", "Shoutdown server...")
}

func startHttp() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/ws/publications", conn.HttpGetPublications)
	http.HandleFunc("/ws/subscribers", conn.HttpGetSubscribers)
	console.LogK("WebSocket", "Http server in http://localhost:3500/ws")
	console.Fatal(http.ListenAndServe(":3500", nil))
}

func startHttp2() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/ws/publications", conn.HttpGetPublications)
	http.HandleFunc("/ws/subscribers", conn.HttpGetSubscribers)
	console.LogK("WebSocket", "Http server in http://localhost:3600/ws")
	console.Fatal(http.ListenAndServe(":3600", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := conn.ConnectHttp(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
	}
}

func test1() {
	var err error
	client1, err = ws.NewClient(&ws.ClientConfig{
		ClientId:  "client1",
		Name:      "client1",
		Schema:    "ws",
		Host:      "localhost:3500",
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
		Host:      "localhost:3500",
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
		Host:      "localhost:3500",
		Path:      "/ws",
		Reconcect: 3,
	})
	if err != nil {
		console.Fatal(err)
	}

	client1.DirectMessage = func(msg ws.Message) {
		console.Debug("DirectMessage:", msg.ToString())
	}

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
	n := 100
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

func test2() {
	var err error
	client1, err = ws.NewClient(&ws.ClientConfig{
		ClientId:  "client1",
		Name:      "client1",
		Schema:    "ws",
		Host:      "localhost:3500",
		Path:      "/ws",
		Reconcect: 3,
	})
	if err != nil {
		console.Fatal(err)
	}

	client1.DirectMessage = func(msg ws.Message) {
		console.Debug("DirectMessage:", msg.ToString())
	}

	client1.Subscribe("Hola", func(msg ws.Message) {
		console.DebugF("Channel:%s :: %s", msg.Channel, msg.ToString())
	})

}
