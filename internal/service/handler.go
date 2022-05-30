package service

import (
	"net"
)

type Handler interface {
	TCPHandler
	UDPHandler
}

type TCPHandler interface {
	Addr() string
	Dial(network, addr string) (c net.Conn, err error)
}

type UDPHandler interface {
	Addr() string
	DialUDP(network, addr string) (pc net.PacketConn, err error)
}
