package ws

import "sync"

var once sync.Once

/**
* Server creates a new Websocket Hub
* @return *Hub
**/
func Server() (*Hub, error) {
	var result *Hub

	initial := func() {
		result = NewHub()
		go result.Run()
	}

	once.Do(initial)

	return result, nil
}
