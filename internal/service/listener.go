package service

type Listener interface {
	Start() error
	Stop() error
}
