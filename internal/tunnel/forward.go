package tunnel

import (
	"errors"
	"net"
	"netsim/internal/service"
	"time"
)

var (
	errModeNotSupport = errors.New("forward mode not support")
)

type Forwarder struct {
	service.Dialer
	addr string
}

func NewForwarder(mode string, dialTimeout, relayTimeout time.Duration) (*Forwarder, error) {
	if mode == "direct" {
		d, err := NewDirect(dialTimeout, relayTimeout)
		if err != nil {
			return nil, err
		}

		return &Forwarder{Dialer: d, addr: d.Addr()}, nil
	}

	return nil, errModeNotSupport
}

func (f *Forwarder) Dial(network, addr string) (net.Conn, error) {
	c, err := f.Dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return c, nil
}
