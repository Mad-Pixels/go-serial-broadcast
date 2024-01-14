package transfer

import (
	"go.bug.st/serial"
)

// Config object.
type Config interface {
	// GetBaudRate The serial port bitrate (aka Baudrate).
	GetBaudRate() int

	// GetDataBits Size of the character (must be 5, 6, 7 or 8).
	GetDataBits() int

	// GetParity serial.Parity object (see Parity type for more info).
	GetParity() serial.Parity

	// GetStopBits Stop bits (see StopBits type for more info).
	GetStopBits() serial.StopBits

	// GetInitialStatusBits Initial output modem bits status (if nil defaults to DTR=true and RTS=true).
	GetInitialStatusBits() *serial.ModemOutputBits
}
