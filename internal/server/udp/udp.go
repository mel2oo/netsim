package udp

import (
	"net"

	"github.com/sirupsen/logrus"
)

type Server struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func New(host, port string) *Server {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))
	if err != nil {
		panic(err)
	}

	return &Server{addr: addr}
}

func (s *Server) Start() (err error) {

	s.conn, err = net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}
	defer s.conn.Close()

	logrus.Infof("udp listen on: %s", s.addr.String())

	for {
		handleConnection(s.conn)
	}
}

func handleConnection(c *net.UDPConn) {
	d := make([]byte, 1024)
	_, _, err := c.ReadFromUDP(d)
	if err != nil {
		return
	}
}

func (s *Server) Stop() error {
	return s.conn.Close()
}
