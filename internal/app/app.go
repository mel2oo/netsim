package app

import (
	"netsim/internal/config"
	"netsim/internal/server/tcp"
	"netsim/internal/server/udp"
)

type App struct {
	c config.Listener

	tcp *tcp.Server
	udp *udp.Server
}

func New(c config.Listener) *App {
	return &App{
		c:   c,
		tcp: tcp.New(),
		udp: udp.New(),
	}
}

func (a *App) Run() error {
	if err := a.tcp.Run(a.c.Tcp.Host, a.c.Tcp.Port); err != nil {
		return err
	}

	if err := a.udp.Run(a.c.Udp.Host, a.c.Udp.Port); err != nil {
		return err
	}

	select {}
}
