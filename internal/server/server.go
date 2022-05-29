package server

import (
	"netsim/internal/config"
	"netsim/internal/service"
)

type Server struct {
	service.Listener
	service.Handler
}

func New(c config.Listener) *Server {
	return nil
}
