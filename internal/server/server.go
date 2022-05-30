package server

import (
	"netsim/internal/config"
	"netsim/internal/service"
	"netsim/internal/service/tcp"

	"github.com/sirupsen/logrus"
)

type Server struct {
	service.Listener
	service.Handler

	c config.Listener
}

func New(c config.Listener) (*Server, error) {
	var (
		err     error
		ln      service.Listener
		handler service.Handler
	)

	switch c.Transport {
	case service.TransportTCP:
		ln, err = tcp.NewListener(c.Address)
	}
	if err != nil {
		return nil, err
	}

	switch c.Protocol {
	case service.ProroclTCP:
		handler, err = tcp.NewHandler()
	}
	if err != nil {
		return nil, err
	}

	return &Server{
		Listener: ln,
		Handler:  handler,
		c:        c,
	}, nil
}

func (s *Server) Serve() error {
	logrus.Infof("[netsim] @%s listen on \"%s\"", s.c.Protocol, s.c.Address)
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		go s.Handle(conn)
	}
}

func (s *Server) Close() error {
	return s.Listener.Close()
}
