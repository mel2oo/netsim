package udp

import (
	"fmt"
	"net"
)

type Server struct {
	conn *net.UDPConn
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(host, port string) error {
	addr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	s.conn = conn

	for {
		handleConnection(conn)
	}
}

func handleConnection(c *net.UDPConn) {
	fmt.Println(c.RemoteAddr().String())
}

func (s *Server) Stop() error {
	return s.conn.Close()
}
