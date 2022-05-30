package service

import (
	"net"
)

type Handler interface {
	Init(options ...HandlerOption)
	Handle(net.Conn)
}

type HandlerOptions struct {
	Addr string
}

type HandlerOption func(opts *HandlerOptions)

func WithHandlerAddr(addr string) HandlerOption {
	return func(opts *HandlerOptions) {
		opts.Addr = addr
	}
}
