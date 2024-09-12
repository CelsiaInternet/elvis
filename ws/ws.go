package ws

import (
	"sync"

	"github.com/cgalvisleon/elvis/console"
)

var once sync.Once

/**
* Server creates a new Websocket Hub
* @return *Hub
**/
func Server() (*Hub, error) {
	result := NewHub()
	if result == nil {
		return nil, console.Alert("Error creating new Websocket Hub")
	}

	go result.Run()

	return result, nil
}
