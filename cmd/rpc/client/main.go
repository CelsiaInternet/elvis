package main

import (
	"fmt"
	"net/rpc"

	"github.com/cgalvisleon/elvis/et"
)

func main() {
	// Conectar al servidor RPC
	client, err := rpc.Dial("tcp", "Cesars-MacBook-Pro.local:1234")
	if err != nil {
		fmt.Println("Error al conectarse al servidor:", err)
		return
	}
	defer client.Close()

	// Crear la solicitud y la respuesta
	require := et.Json{} // Si tienes algún dato en la solicitud, rellénalo aquí
	var response et.Item

	// Llamar al método Version
	err = client.Call("Services.Version", require, &response)
	if err != nil {
		fmt.Println("Error al invocar el método RPC:", err)
		return
	}

	// Mostrar la respuesta
	fmt.Println("Respuesta:", response)
}
