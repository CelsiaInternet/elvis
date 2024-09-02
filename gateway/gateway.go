package gateway

import (
	"fmt"
	"os"
	"time"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/cgalvisleon/elvis/ws"
	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
)

type Server struct {
	http *HttpServer
	ws   *ws.Hub
}

var PackageName = "gateway"
var PackageTitle = envar.EnvarStr("Apigateway", "PACKAGE_TITLE")
var PackagePath = envar.EnvarStr("/api/gateway", "PATH_URL")
var PackageVersion = envar.EnvarStr("0.0.1", "VERSION")
var Company = envar.EnvarStr("", "COMPANY")
var Web = envar.EnvarStr("", "WEB")
var HostName, _ = os.Hostname()
var Host = strs.Format(`%s:%d`, envar.EnvarStr("http://localhost", "HOST"), envar.EnvarInt(3300, "PORT"))
var conn *Server

func New() (*Server, error) {
	if conn != nil {
		return conn, nil
	}

	// Cache
	_, err := cache.Load()
	if err != nil {
		return nil, err
	}

	// Event
	_, err = event.Load()
	if err != nil {
		return nil, err
	}

	// WS server
	ws, err := ws.Server()
	if err != nil {
		panic(err)
	}

	// HTTP server
	http := newHttpServer()

	// Create a new server
	conn = &Server{
		http: http,
		ws:   ws,
	}

	return conn, nil
}

func (serv *Server) Close() error {
	return nil
}

func (serv *Server) Start() {
	// Start HTTP server
	InitHttp(serv)

	// Init events
	initEvents()

	// Banner
	Banner()

	<-make(chan struct{})
}

func Version() et.Json {
	service := et.Json{
		"version": envar.EnvarStr("", "VERSION"),
		"service": PackageName,
		"host":    HostName,
		"company": Company,
		"web":     Web,
		"help":    "",
	}

	return service
}

func Banner() {
	time.Sleep(3 * time.Second)
	templ := fmt.Sprintf(`{{ .Title "%s V%s" "" 4 }}`, PackageName, PackageVersion)
	banner.InitString(colorable.NewColorableStdout(), true, true, templ)
	fmt.Println()
}
