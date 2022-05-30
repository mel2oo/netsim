package tunnel

import (
	"context"
	"net"
	"time"
)

type Direct struct {
	ip           net.IP
	dialTimeout  time.Duration
	relayTimeout time.Duration
}

func NewDirect(dialTimeout, relayTimeout time.Duration) (*Direct, error) {
	return &Direct{
		dialTimeout:  dialTimeout,
		relayTimeout: relayTimeout,
	}, nil
}

func (d *Direct) Addr() string {
	return "DIRECT"
}

func (d *Direct) Dial(network, addr string) (c net.Conn, err error) {
	return nil, nil
}

func (d *Direct) dial(network, addr string, localIP net.IP) (net.Conn, error) {
	var la net.Addr
	switch network {
	case "tcp":
		la = &net.TCPAddr{IP: localIP}
	case "udp":
		la = &net.UDPAddr{IP: localIP}
	}

	dialer := &net.Dialer{LocalAddr: la, Timeout: d.dialTimeout}
	// if d.iface != nil {
	// 	dialer.Control = sockopt.Control(sockopt.Bind(d.iface))
	// }

	c, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	if c, ok := c.(*net.TCPConn); ok {
		c.SetKeepAlive(true)
	}

	if d.relayTimeout > 0 {
		c.SetDeadline(time.Now().Add(d.relayTimeout))
	}

	return c, err
}

func (d *Direct) DialUDP(network, addr string) (pc net.PacketConn, err error) {
	var la string
	if d.ip != nil {
		la = net.JoinHostPort(d.ip.String(), "0")
	}

	lc := &net.ListenConfig{}

	return lc.ListenPacket(context.Background(), network, la)
}
