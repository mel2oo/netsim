package tcp

import "netsim/internal/service"

type transporter struct{}

func NewTransporter() service.Transporter {
	return nil
}

func NewListener() service.Listener {
	return nil
}
