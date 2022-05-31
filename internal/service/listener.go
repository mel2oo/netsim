package service

import "net"

type Listener interface {
	ListenAndServe() error
	Serve(c net.Conn)
	Close() error
}

var listeners = make(map[string]Listener)

func RegisterListener(name string, ln Listener) {
	listeners[name] = ln
}
