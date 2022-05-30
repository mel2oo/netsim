package tunnel

import (
	"net"
	"netsim/internal/service"
)

type Proxy interface {
	Dial(network, addr string) (c net.Conn, handler service.Handler, err error)
	DialUDP(network, addr string) (pc net.PacketConn, dialer service.UDPHandler, err error)
	NextDialer(dstAddr string) service.Handler
	Record(dialer service.Handler, success bool)
}
