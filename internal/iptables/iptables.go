package iptables

type Options struct {
}

type Option func(opts *Options)

func Set(opts ...Option) {

}
