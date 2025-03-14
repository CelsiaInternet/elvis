package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"math/rand"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/strs"
	"github.com/celsiainternet/elvis/ws"
)

var conn *ws.Hub
var clients []*ws.Client

func main() {
	if conn != nil {
		return
	}

	envar.SetInt("port", 3300, "Port server", "PORT")
	envar.SetStr("username", "", "Modo cluster: master or worker", "WS_USERNAME")
	envar.SetStr("password", "", "Master host", "WS_PASSWORD")

	port := envar.GetInt(3300, "PORT")
	username := envar.GetStr("", "WS_USERNAME")
	password := envar.GetStr("", "WS_PASSWORD")
	conn = ws.ServerHttp(port, username, password)
	conn.JoinTo(et.Json{
		"adapter":  "redis",
		"host":     "localhost:6379",
		"dbname":   envar.GetInt(0, "REDIS_DB"),
		"password": envar.GetStr("", "REDIS_PASSWORD"),
	})

	// test1(port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	console.LogK("WebSocket", "Shoutdown server...")
}

func test1(port int) {
	url := strs.Format(`ws://localhost:%d/ws`, port)

	n := 10
	for i := 0; i < n; i++ {
		client, err := ws.NewClient(&ws.ClientConfig{
			ClientId:  strs.Format("client-%d", i),
			Name:      strs.Format("Client%d", i),
			Url:       url,
			Reconnect: 3,
		})
		if err != nil {
			console.AlertE(err)
		}

		client.Subscribe("Hola", func(msg ws.Message) {
			console.Debug("client1", msg.ToString())
		})

		clients = append(clients, client)
	}

	rand.NewSource(time.Now().UnixNano())

	t := time.Duration(100)
	for {
		idx := rand.Intn(n)
		client := clients[idx]
		if client != nil {
			client.Publish("Hola", et.Json{
				"msg": strs.Format("Hola %d", idx),
			})
		}
		time.Sleep(t * time.Millisecond)
	}
}
