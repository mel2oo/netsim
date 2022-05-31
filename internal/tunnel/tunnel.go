package tunnel

import (
	"net"
	"netsim/internal/config"
	"netsim/internal/service"
	"time"
)

type Tunnel struct {
	forwarders *Forwarder
}

func NewTunnel(c *config.Forwarder) (*Tunnel, error) {
	forwarder, err := NewForwarder(
		c.Mode,
		time.Duration(c.DTimeout)*time.Second,
		time.Duration(c.RTimeout)*time.Second)
	if err != nil {
		return nil, err
	}

	return &Tunnel{forwarder}, nil
}

func (t *Tunnel) Dial(network, addr string) (net.Conn, service.Dialer, error) {
	conn, err := t.forwarders.Dial(network, addr)
	return conn, t.forwarders, err
}

func (t *Tunnel) DialUDP(network, addr string) (net.PacketConn, service.UDPDialer, error) {
	conn, err := t.forwarders.DialUDP(network, addr)
	return conn, t.forwarders, err
}

func (t *Tunnel) Record(dialer service.Dialer, success bool) {

}
