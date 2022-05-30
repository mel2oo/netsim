package forward

import (
	"io"
	"netsim/internal/service"
	"sync"
)

var (
	bpool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 32*1024)
		},
	}
)

type Handler struct {
	Remote  string
	Options *service.HandlerOptions
}

func (h *Handler) Init(options ...service.HandlerOption) {

}

func Transport(rw1, rw2 io.ReadWriter) error {
	errc := make(chan error, 1)
	go func() {
		errc <- CopyBuffer(rw1, rw2)
	}()

	go func() {
		errc <- CopyBuffer(rw2, rw1)
	}()

	err := <-errc
	if err != nil && err == io.EOF {
		err = nil
	}
	return err
}

func CopyBuffer(dst io.Writer, src io.Reader) error {
	buf := bpool.Get().([]byte)
	defer bpool.Put(buf)

	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
