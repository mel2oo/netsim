package service

type Listener interface {
	Start() error
	Stop() error
}

type Handler interface {
}

type Resolver interface {
}

type Transport interface {
}
