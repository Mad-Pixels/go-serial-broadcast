package go_serial_broadcast

import (
	"go.bug.st/serial"
)

// Option object.
type Option func(*config)

// WithBaudRate set serial port bitrate.
func WithBaudRate(v int) func(*config) {
	return func(c *config) {
		c.baudRate = v
	}
}

// WithDataBits set size of the character (must be 5, 6, 7 or 8).
func WithDataBits(v int) func(*config) {
	return func(c *config) {
		c.dataBits = v
	}
}

// WithParity set serial.Parity.
func WithParity(v serial.Parity) func(*config) {
	return func(c *config) {
		c.parity = v
	}
}

// WithStopBits set serial.StopBits.
func WithStopBits(v serial.StopBits) func(*config) {
	return func(c *config) {
		c.stopBits = v
	}
}

// WithInitialStatusBits set output modem bits status (if nil defaults to DTR=true and RTS=true).
func WithInitialStatusBits(v *serial.ModemOutputBits) func(*config) {
	return func(c *config) {
		c.initialStatusBits = v
	}
}

// WithDeviceMsgPattern is an option that set the pattern which will use for match outputs from device for determinate it.
func WithDeviceMsgPattern(regexpPattern string) func(*config) {
	return func(c *config) {
		c.incomeDeviceMsgPattern = regexpPattern
	}
}
