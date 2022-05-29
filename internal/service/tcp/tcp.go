package tcp

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type Server struct {
	addr *net.TCPAddr
	ln   *net.TCPListener
}

func New(host, port string) *Server {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
	if err != nil {
		panic(err)
	}

	return &Server{addr: addr}
}

func (s *Server) Start() (err error) {
	s.ln, err = net.ListenTCP("tcp", s.addr)
	if err != nil {
		return err
	}
	defer s.ln.Close()

	logrus.Infof("tcp listen on: %s", s.addr.String())

	for {
		c, err := s.ln.Accept()
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
	return s.ln.Close()
}
