package tcp

import (
	"net"
	"netsim/internal/service"
	"netsim/internal/tunnel/forward"

	"github.com/sirupsen/logrus"
)

// type transporter struct{}

// func NewTransporter() (service.Transporter, error) {
// 	return &transporter{}, nil
// }

// func (t *transporter) Dial(addr string, options ...service.DialOption) (net.Conn, error) {
// 	opts := &service.DialOptions{}
// 	for _, option := range options {
// 		option(opts)
// 	}

// 	timeout := opts.Timeout
// 	if timeout <= 0 {
// 		timeout = service.DefaultDialTimeout
// 	}

// 	return net.DialTimeout("tcp", addr, timeout)
// }

// func (t *transporter) Handshake(conn net.Conn, options ...service.HandshakeOption) (net.Conn, error) {
// 	return conn, nil
// }

type listener struct {
	net.Listener
}

func NewListener(addr string) (service.Listener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}

	return &listener{&keeplive{ln}}, nil
}

type keeplive struct {
	*net.TCPListener
}

func (ln *keeplive) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(service.DefaultKeepAliveTime)
	return tc, nil
}

type handler struct {
	*forward.Handler
}

func NewHandler(opts ...service.HandlerOption) (service.Handler, error) {
	h := &handler{
		Handler: &forward.Handler{
			Options: &service.HandlerOptions{},
		},
	}

	for _, o := range opts {
		o(h.Options)
	}

	return h, nil
}

func (h *handler) Init(options ...service.HandlerOption) {
	h.Handler.Init(options...)
}

func (h *handler) Handle(conn net.Conn) {
	logrus.Infof("[tcp] %s -> %s", conn.LocalAddr(), conn.RemoteAddr())

	defer conn.Close()

	client, err := net.Dial("tcp", conn.RemoteAddr().String())
	if err != nil {
		return
	}
	defer client.Close()

	forward.Transport(conn, client)
}
