package sharedbuffer

import (
	"sync"

	"go.uber.org/zap/buffer"
)

type Factory struct {
	pool      buffer.Pool
	bufs      sync.Map
	writeChan chan<- struct{}
}

func NewFactory(writeChan chan<- struct{}) *Factory {
	return &Factory{
		pool:      buffer.NewPool(),
		writeChan: writeChan,
	}
}

func (f *Factory) Get(idx uint) *SharedBuffer {
	buf, _ := f.bufs.LoadOrStore(idx, &SharedBuffer{
		buf:       f.pool.Get(),
		writeChan: f.writeChan,
	})
	return buf.(*SharedBuffer)
}
