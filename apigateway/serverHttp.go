package apigateway

import (
	"net/http"

	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/strs"
	"github.com/rs/cors"
)

type HttpServer struct {
	addr            string
	handler         http.Handler
	mux             *http.ServeMux
	notFoundHandler http.HandlerFunc
	handlerFn       http.HandlerFunc
	handlerWS       http.HandlerFunc
	// middlewares     []func(http.Handler) http.Handler
}

func NewHttpServer() *HttpServer {
	// Create a new server
	mux := http.NewServeMux()

	port := envar.EnvarInt(3300, "PORT")
	result := &HttpServer{
		addr:    strs.Format(":%d", port),
		handler: cors.AllowAll().Handler(mux),
		mux:     mux,
	}
	result.notFoundHandler = notFounder
	result.handlerFn = handlerFn

	// Handler router
	mux.HandleFunc("/version", version)
	mux.HandleFunc("/", result.handlerFn)

	return result
}

func (s *HttpServer) NotFound(handlerFn http.HandlerFunc) {
	s.notFoundHandler = handlerFn
}

func (s *HttpServer) Handler(handlerFn http.HandlerFunc) {
	s.handlerFn = handlerFn
}

func (s *HttpServer) HandlerWebSocket(handlerFn http.HandlerFunc) {
	s.handlerWS = handlerFn
}
