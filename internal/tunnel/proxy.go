package tunnel

import (
	"net"
	"netsim/internal/service"
)

type Proxy interface {
	Dial(network, addr string) (net.Conn, service.Dialer, error)
	DialUDP(network, addr string) (net.PacketConn, service.UDPDialer, error)
	// NextDialer(dstAddr string) service.Dialer
	Record(service.Dialer, bool)
}
