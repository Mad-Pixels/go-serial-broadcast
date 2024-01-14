package go_serial_broadcast

import (
	"go.bug.st/serial"
)

type config struct {
	baudRate               int
	dataBits               int
	incomeDeviceMsgPattern string
	parity                 serial.Parity
	stopBits               serial.StopBits
	initialStatusBits      *serial.ModemOutputBits
}

// GetBaudRate variable.
func (c *config) GetBaudRate() int {
	return c.baudRate
}

// GetDataBits variable.
func (c *config) GetDataBits() int {
	return c.dataBits
}

// GetIncomeDeviceMsgPattern variable.
func (c *config) GetIncomeDeviceMsgPattern() string {
	return c.incomeDeviceMsgPattern
}

// GetParity variable.
func (c *config) GetParity() serial.Parity {
	return c.parity
}

// GetStopBits variable.
func (c *config) GetStopBits() serial.StopBits {
	return c.stopBits
}

// GetInitialStatusBits variable.
func (c *config) GetInitialStatusBits() *serial.ModemOutputBits {
	return c.initialStatusBits
}
