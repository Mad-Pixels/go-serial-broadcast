package port

import (
	"fmt"
	"sync"

	bugst "go.bug.st/serial"
)

type serial struct {
	mu   sync.Mutex
	port bugst.Port
	path string
}

func NewSerial() (*serial, error) {
	c := &bugst.Mode{BaudRate: 9600}
	port, err := bugst.Open("/dev/tty.usbserial-110", c)
	if err != nil {
		return nil, err
	}
	return &serial{
		mu:   sync.Mutex{},
		port: port,
		path: "/dev/tty.usbserial-110",
	}, nil
}

// Read ...
func (s *serial) Read(buf []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//if err := s.port.ResetInputBuffer(); err != nil {
	//	return 0, fmt.Errorf("error resetting input buffer: %s", err)
	//}

	n, err := s.port.Read(buf)
	if err != nil {
		return n, fmt.Errorf("error reading from serial port: %s", err)
	}
	return n, nil
}

func (s *serial) Write(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//if err := s.port.ResetOutputBuffer(); err != nil {
	//	return 0, fmt.Errorf("error resetting output buffer: %s", err)
	//}

	n, err := s.port.Write(data)
	if err != nil {
		return n, fmt.Errorf("error writing to serial port: %s", err)
	}

	if err := s.port.Drain(); err != nil {
		return n, fmt.Errorf("error draining serial port: %s", err)
	}
	return n, nil
}

// Close ...
func (s *serial) Close() error {
	return s.port.Close()
}
