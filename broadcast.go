package go_serial_broadcast

import (
	"context"
	"sync"
	"time"

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

func NewUARTAutoDetectByDeviceSerialMsg(ctx context.Context, options ...Option) (transfer.Transfer, error) {
	cfg := &config{watchDogTimeout: 2}
	for _, o := range options {
		o(cfg)
	}

	deviceIdentity, err := identifier.New(cfg)
	if err != nil {
		return nil, err
	}
	var devicePath string
	for {
		if devicePath != "" {
			break
		}
		ports, err := deviceIdentity.Ports()
		if err != nil {
			return nil, err
		}

		var wg sync.WaitGroup
		for _, port := range ports {
			wg.Add(1)

			go func(ctx context.Context, port string) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
				defer cancel()

				outCh := make(chan []byte, 1)
				defer close(outCh)
				errCh := make(chan error, 1)
				defer close(errCh)

				uart, err := NewUART(port)
				if err != nil {
					return
				}
				defer uart.Close()

				go uart.ReadToCh(ctx, outCh, errCh)
				for {
					select {
					case <-errCh:
					case <-ctx.Done():
						return
					case out := <-outCh:
						if deviceIdentity.Check(out) {
							devicePath = port
							return
						}
					}
				}
			}(ctx, port)
		}
		wg.Wait()
		time.Sleep(time.Second * cfg.watchDogTimeout)
	}

	return NewUART(devicePath)
}
