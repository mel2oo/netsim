package server

type Server interface {
	Run(string, string) error
	Stop() error
}
