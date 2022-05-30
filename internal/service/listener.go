package service

import "net"

type Listener interface {
	net.Listener
}
