package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/utility"
)

// Definir el tipo Services
type Services struct{}

// Función Version que se expone vía RPC
func (c *Services) Version(require et.Json, response *et.Item) error {
	// Rellenando la respuesta
	response.Ok = true
	response.Result = et.Json{
		"id":      utility.UUID(),
		"service": PackageName(),
		"host":    HostName(),
	}

	return nil
}

// Función para obtener el nombre del paquete (simulado)
func PackageName() string {
	return "ExampleService"
}

// Función para obtener el nombre del host
func HostName() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return host
}

func main() {
	// Crear una nueva instancia del servicio
	services := new(Services)

	// Registrar el servicio con el servidor RPC
	rpc.Register(services)

	// Iniciar el servidor RPC en una goroutine
	go startRPCServer()

	// Iniciar el servidor HTTP en una goroutine
	go startHTTPServer()

	// Bloquear el hilo principal para que el servidor siga corriendo
	select {}
}

// startRPCServer inicia el servidor RPC
func startRPCServer() {
	// Escuchar en un puerto TCP para RPC
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error al iniciar el servidor RPC:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor RPC escuchando en el puerto 1234")

	// Aceptar y manejar conexiones RPC
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar la conexión RPC:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// startHTTPServer inicia el servidor HTTP
func startHTTPServer() {
	// Definir un manejador simple para la ruta principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Servidor HTTP funcionando en el puerto 8080")
	})

	// Iniciar el servidor HTTP en el puerto 8080
	fmt.Println("Servidor HTTP escuchando en el puerto 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor HTTP:", err)
	}
}
