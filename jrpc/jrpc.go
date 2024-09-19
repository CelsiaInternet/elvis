package jrpc

import (
	"encoding/json"
	"net"
	"net/http"
	"net/rpc"
	"reflect"

	"github.com/cgalvisleon/elvis/cache"
	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/event"
	"github.com/cgalvisleon/elvis/strs"
)

type Route struct {
	Host string
	Port int
}

type Router struct {
	PackageName string            `json:"packageName"`
	Routes      map[string]*Route `json:"routes"`
}

func NewRouter() *Router {
	return &Router{
		Routes: make(map[string]*Route),
	}
}

type Server struct {
	rpc    *net.Listener
	Router *Router
}

func NewServer(port int) (*Server, error) {
	_, err := cache.Load()
	if err != nil {
		return nil, err
	}

	server := &Server{
		Router: NewRouter(),
	}

	rpc, err := net.Listen("tcp", strs.Format(":%d", port))
	if err != nil {
		return nil, err
	}

	server.rpc = &rpc

	return server, nil
}

func (s *Server) Mount(host string, port int, service any, packageName string) {
	tipoStruct := reflect.TypeOf(service)
	console.Debug("Mounting service", tipoStruct.Name())
	s.Router.PackageName = packageName
	for i := 0; i < tipoStruct.NumMethod(); i++ {
		name := tipoStruct.Method(i).Name
		path := strs.Format(`%s.%s`, tipoStruct.Name(), name)
		s.Router.Routes[path] = &Route{host, port}
	}

	rpc.Register(service)

	bt, err := json.Marshal(s.Router)
	if err != nil {
		return
	}

	var data et.Json
	err = json.Unmarshal(bt, &data)
	if err != nil {
		return
	}

	event.Publish("apigateway/rpc/resolve", data)
}

func (s *Server) Start() {
	if s.rpc == nil {
		return
	}

	go func() {
		console.LogKF("RPC", "Running on tcp:localhost:%s", (*s.rpc).Addr().String())
		http.Serve((*s.rpc), nil)
	}()
}

func (s *Server) Close() {
	if s.rpc == nil {
		return
	}

	(*s.rpc).Close()
}

func (s *Server) Call(method string, data et.Json) (et.Item, error) {
	var result = et.Item{Result: et.Json{}}
	var args []byte = data.ToByte()
	var reply *[]byte

	router := s.Router.Routes[method]
	if router == nil {
		return result, console.NewError(ERR_METHOD_NOT_FOUND)
	}

	client, err := rpc.DialHTTP("tcp", strs.Format(`%s:%d`, router.Host, router.Port))
	if err != nil {
		return et.Item{}, err
	}
	defer client.Close()

	err = client.Call(method, args, &reply)
	if err != nil {
		return et.Item{}, err
	}

	result = et.Json{}.ToItem(*reply)

	return result, nil
}
