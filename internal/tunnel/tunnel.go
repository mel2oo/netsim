package tunnel

import (
	"netsim/internal/config"
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
