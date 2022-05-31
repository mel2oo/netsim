package server

import (
	"netsim/internal/config"
	"netsim/internal/service"
)

type Server struct {
	service.Listener
	service.Dialer

	c config.Listener
}

func New(c config.Listener) (*Server, error) {
	var (
		ln service.Listener
		da service.Dialer
	)

	return &Server{
		Listener: ln,
		Dialer:   da,
		c:        c,
	}, nil
}
