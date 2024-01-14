package go_serial_broadcast

import (
	"context"
	"sync"

	"github.com/MadPixeles/go-serial-broadcast/identifier"
	"github.com/MadPixeles/go-serial-broadcast/transfer"
)

// NewUART initialize serial port object.
func NewUART(path string, options ...Option) (transfer.Transfer, error) {
	cfg := &config{}
	for _, o := range options {
		o(cfg)
	}
	return transfer.New(path, cfg)
}

// NewUARTAutoDetectByDeviceSerialMsg initialize serial port.
func NewUARTAutoDetectByDeviceSerialMsg(ctx context.Context, options ...Option) (transfer.Transfer, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := &config{}
	for _, o := range options {
		o(cfg)
	}

	i, err := identifier.New(cfg)
	if err != nil {
		return nil, err
	}

	var result transfer.Transfer
	for {
		if result != nil {
			return result, nil
		}
		ports, err := i.Ports()
		if err != nil {
			return nil, err
		}

		var wg sync.WaitGroup
		for _, port := range ports {
			wg.Add(1)

			go func(port string) {
				defer wg.Done()

				uart, err := NewUART(port)
				if err != nil {
					return
				}

				outCh := make(chan []byte, 1)
				defer close(outCh)
				errCh := make(chan error, 1)
				defer close(errCh)

				go uart.ReadToCh(ctx, outCh, errCh)
				for {
					select {
					case <-errCh:
						return
					case out := <-outCh:
						if i.Check(out) {
							result = uart
						}
					}
				}
			}(port)
		}
	}
}
