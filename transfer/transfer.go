package transfer

import (
	"context"
)

// Transfer object.
type Transfer interface {
	ReadToCh(context.Context, chan<- []byte, chan<- error)
	Write([]byte) (int, error)
	ResetOutputBuffer() error
	ResetInputBuffer() error
}

// New initialize Transfer object.
func New(path string, config Config) (Transfer, error) {
	return newPort(path, config)
}
