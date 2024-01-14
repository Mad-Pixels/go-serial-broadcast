package identifier

type Config interface {
	// GetIncomeDeviceMsgPattern variable which used in case when connected device start write to serial port some data.
	// This variable is a regexp pattern for detect current device by data which it will send to serial port.
	// Supposed that device will send "serial number" or "some code" which should be match with pattern for determinate it.
	GetIncomeDeviceMsgPattern() string
}
