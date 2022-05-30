package service

import "net"

type Listener interface {
	ListenAndServe()
	Serve(c net.Conn)
}
