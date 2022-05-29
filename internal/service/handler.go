package service

import "netsim/internal/tunnel/forward"

type Handler interface {
	forward.Forwarder
}
