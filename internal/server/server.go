package server

import "netsim/internal/service"

type Server interface {
	service.Listener
	service.Handler
}
