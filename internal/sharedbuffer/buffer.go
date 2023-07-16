package sharedbuffer

import (
	"go.uber.org/zap/buffer"
)

type SharedBuffer struct {
	buf       *buffer.Buffer
	writeChan chan<- struct{}
}

func (sb *SharedBuffer) Write(p []byte) (int, error) {
	n, err := sb.buf.Write(p)
	if err == nil {
		sb.writeChan <- struct{}{}
	}

	return n, err
}

func (sb *SharedBuffer) RawBuffer() *buffer.Buffer {
	return sb.buf
}
