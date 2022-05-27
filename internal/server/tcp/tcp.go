package tcp

import (
	"fmt"
	"net"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(host, port string) error {
	l, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}
	defer l.Close()

	for {

		c, err := l.Accept()
		if err != nil {
			return err
		}

		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	fmt.Println(c.RemoteAddr().String())
}
