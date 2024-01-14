package transfer

import (
	"context"
	"errors"

	"go.bug.st/serial"
)

type uart struct {
	port serial.Port
	path string
}

func newPort(path string, config Config) (*uart, error) {
	mode := &serial.Mode{
		InitialStatusBits: config.GetInitialStatusBits(),
		BaudRate:          config.GetBaudRate(),
		DataBits:          config.GetDataBits(),
		StopBits:          config.GetStopBits(),
		Parity:            config.GetParity(),
	}
	port, err := serial.Open(path, mode)
	if err != nil {
		return nil, err
	}
	return &uart{
		port: port,
		path: path,
	}, nil
}

// Path current serial port.
func (u *uart) Path() string {
	return u.path
}

// ReadToCh read data from serial port to []byte channel.
func (u *uart) ReadToCh(ctx context.Context, outCh chan<- []byte, errCh chan<- error) {
	buff := make([]byte, 256)
	for {
		select {
		case <-ctx.Done():
			errCh <- errors.Join(ctx.Err(), u.Close())
			return
		default:
		}

		n, err := u.port.Read(buff)
		if err != nil {
			errCh <- errors.Join(err, u.Close())
			return
		}
		if n == 0 {
			return
		}

		outCh <- buff[:n]
		u.ResetInputBuffer()
	}
}

// Write message to serial port.
func (u *uart) Write(msg []byte) (int, error) {
	n, err := u.port.Write(msg)
	defer func() {
		err = errors.Join(err, u.port.Drain(), u.port.ResetOutputBuffer())
	}()
	return n, err
}

// Close the serial port.
func (u *uart) Close() error {
	return u.port.Close()
}

// ResetInputBuffer Purges port read buffer.
func (u *uart) ResetInputBuffer() error {
	return u.port.ResetInputBuffer()
}

// ResetOutputBuffer Purges port write buffer.
func (u *uart) ResetOutputBuffer() error {
	return u.port.ResetOutputBuffer()
}
