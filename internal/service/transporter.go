package service

import (
	"net"
	"time"
)

type Transporter interface {
	Dial(addr string, options ...DialOption) (net.Conn, error)
	Handshake(conn net.Conn, options ...HandshakeOption) (net.Conn, error)
}

// DialOptions
type DialOptions struct {
	Timeout time.Duration
	Host    string
}

type DialOption func(opts *DialOptions)

func WithDialTimeout(timeout time.Duration) DialOption {
	return func(opts *DialOptions) {
		opts.Timeout = timeout
	}
}

func WithDialHost(host string) DialOption {
	return func(opts *DialOptions) {
		opts.Host = host
	}
}

// HandshakeOptions
type HandshakeOptions struct {
	Addr     string
	Host     string
	Timeout  time.Duration
	Interval time.Duration
	Retry    int
}

type HandshakeOption func(opts *HandshakeOptions)

func WithHandshakeAddr(addr string) HandshakeOption {
	return func(opts *HandshakeOptions) {
		opts.Addr = addr
	}
}

func WithHandshakeHost(host string) HandshakeOption {
	return func(opts *HandshakeOptions) {
		opts.Host = host
	}
}

func WithHandshakeTimeout(timeout time.Duration) HandshakeOption {
	return func(opts *HandshakeOptions) {
		opts.Timeout = timeout
	}
}

func WithHandshakeInterval(interval time.Duration) HandshakeOption {
	return func(opts *HandshakeOptions) {
		opts.Interval = interval
	}
}

func WithHandshakeRetry(retry int) HandshakeOption {
	return func(opts *HandshakeOptions) {
		opts.Retry = retry
	}
}
