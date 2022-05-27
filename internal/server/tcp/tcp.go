package tcp

import (
	"fmt"
	"net"
)

type Server struct {
	conn net.Listener
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(host, port string) error {
	conn, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}
	defer conn.Close()
	s.conn = conn

	for {
		c, err := conn.Accept()
		if err != nil {
			return err
		}

		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	fmt.Println(c.RemoteAddr().String())
}

func (s *Server) Stop() error {
	return s.conn.Close()
}
