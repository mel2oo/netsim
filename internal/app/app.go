package app

import (
	"context"
	"errors"
	"netsim/internal/config"
	"netsim/internal/server"
	"netsim/internal/server/tcp"
	"netsim/internal/server/udp"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type App struct {
	ctx    context.Context
	cancel func()

	cnf  *config.Listener
	srvs []server.Server
	sigs []os.Signal
}

func New(cnf *config.Listener) *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		ctx:    ctx,
		cancel: cancel,
		cnf:    cnf,
		srvs: []server.Server{
			tcp.New(cnf.Tcp.Host, cnf.Tcp.Port),
			udp.New(cnf.Udp.Host, cnf.Udp.Port),
		},
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
}

func (a *App) Run() error {
	g, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.srvs {
		srv := srv

		g.Go(func() error {
			<-ctx.Done()
			return srv.Stop()
		})

		g.Go(func() error {
			return srv.Start()
		})
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				a.Stop()
			}
		}
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
