package tcp

import "net"

type TCP struct {
	addr string
}

func (s *TCP) ListenAndServe() error {
	return nil
}

func (s *TCP) Serve(c net.Conn) {

}

func (s *TCP) Close() error {
	return nil
}
