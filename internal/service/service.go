package service

import (
	"crypto/tls"
	"time"
)

const (
	DefaultKeepAliveTime = 180 * time.Second
	DefaultDialTimeout   = 5 * time.Second
)

const (
	TransportTCP = "tcp"
)

const (
	ProroclTCP = "tcp"
)

var (
	DefaultTLSConfig *tls.Config
)
