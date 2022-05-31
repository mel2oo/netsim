package app

import (
	"context"
	"errors"
	"netsim/internal/config"
	"netsim/internal/server"
	"netsim/internal/server/tls"
	"netsim/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type App struct {
	ctx    context.Context
	cancel func()

	cnf  *config.Config
	srvs []*server.Server
	sigs []os.Signal
}

func New(cnf *config.Config) *App {
	ctx, cancel := context.WithCancel(context.Background())
	srvs := make([]*server.Server, 0)

	tlscnf, err := tls.Load(cnf.TLS.Cert, cnf.TLS.Key, cnf.TLS.Ca)
	if err != nil {
		panic(err)
	}
	service.DefaultTLSConfig = tlscnf

	for _, sc := range cnf.Listener {
		srv, err := server.New(sc)
		if err != nil {
			logrus.Warnf("[netsim] create proxy server fail, %s", err.Error())
		} else {
			srvs = append(srvs, srv)
		}
	}

	return &App{
		ctx:    ctx,
		cancel: cancel,
		cnf:    cnf,
		srvs:   srvs,
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
}

func (a *App) Run() error {
	g, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.srvs {
		srv := srv

		g.Go(func() error {
			<-ctx.Done()
			return srv.Close()
		})

		g.Go(func() error {
			return srv.ListenAndServe()
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
