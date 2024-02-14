package port

import (
	"sync"

	"github.com/pkg/errors"
	bugst "go.bug.st/serial"
)

type serial struct {
	mu   sync.Mutex
	port bugst.Port
	path string
}

// return Port object based on "go.bug.st/serial" pkg.
func newSerial(path string, rate, dataBits, stopBits, parity int) (Port, error) {
	mode := &bugst.Mode{
		StopBits:          bugst.StopBits(stopBits),
		Parity:            bugst.Parity(parity),
		DataBits:          dataBits,
		BaudRate:          rate,
		InitialStatusBits: nil,
	}
	port, err := bugst.Open(path, mode)
	if err != nil {
		return nil, errors.Wrapf(err, "fail initialize serial port %s", path)
	}
	return &serial{
		mu:   sync.Mutex{},
		port: port,
		path: path,
	}, nil
}

// Stores data received from the serial port into the provided byte array buffer.
func (s *serial) Read(buf []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, err := s.port.Read(buf)
	if err != nil {
		return n, errors.Wrapf(err, "fail reading from serial port %s, got", s.path)
	}
	return n, nil
}

// Send the content of the data byte array to the serial port.
func (s *serial) Write(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.port.ResetOutputBuffer(); err != nil {
		return 0, errors.Wrapf(err, "fail to reset output buffer on serial port %s, got", s.path)
	}

	n, err := s.port.Write(data)
	if err != nil {
		return n, errors.Wrapf(err, "fail writing to serial port %s, got", s.path)
	}
	if err := s.port.Drain(); err != nil {
		return n, errors.Wrapf(err, "fail to drain serial port %s, got", s.path)
	}
	return n, nil
}

// Close the serial port.
func (s *serial) Close() error {
	if err := s.port.Close(); err != nil {
		return errors.Wrapf(err, "fail to close serial port %s, got", s.path)
	}
	return nil
}
