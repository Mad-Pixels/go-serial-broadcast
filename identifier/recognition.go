package identifier

import (
	"regexp"

	"go.bug.st/serial"
)

type recognition struct {
	config Config
}

func newRecognition(config Config) (*recognition, error) {
	return &recognition{
		config: config,
	}, nil
}

// Ports retrieve the list of available serial ports.
func (r *recognition) Ports() ([]string, error) {
	return serial.GetPortsList()
}

// Check income row data with patterns for determinate current device.
func (r *recognition) Check(data []byte) bool {
	res, err := regexp.Match(r.config.GetIncomeDeviceMsgPattern(), data)
	if err != nil {
		return false
	}
	return res
}
