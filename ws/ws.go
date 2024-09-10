package ws

var conn *Hub

/**
* Server creates a new Websocket Hub
* @return *Hub
**/
func Server() (*Hub, error) {
	if conn != nil {
		return conn, nil
	}

	conn = NewHub()
	go conn.Run()

	return conn, nil
}

/**
* Close the Websocket Hub
* @return error
**/
func Close() error {
	return nil
}
