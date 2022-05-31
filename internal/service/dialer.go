package service

import (
	"net"
)

type Dialer interface {
	TCPDialer
	UDPDialer
}

type TCPDialer interface {
	Addr() string
	Dial(network, addr string) (c net.Conn, err error)
}

type UDPDialer interface {
	Addr() string
	DialUDP(network, addr string) (pc net.PacketConn, err error)
}

var dialers = make(map[string]Dialer)

func RegisterDialer(name string, di Dialer) {
	dialers[name] = di
}
