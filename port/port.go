package port

// Port defines a contract for serial port operations.
// It abstracts the underlying implementation details of serial port communication,
// providing a simplified view for transmitting and receiving data, as well as managing the port's lifecycle.
type Port interface {
	// Write sends data to the serial port.
	// It takes a slice of bytes as input, which represents the data to be transmitted.
	// The method returns the number of bytes successfully written and any error encountered during the operation.
	// Implementations of this method should ensure that data is correctly transmitted over the serial port,
	// handling any necessary protocol or configuration details.
	Write([]byte) (int, error)

	// Read attempts to read data from the serial port into the provided buffer.
	// It takes a slice of bytes as an argument, where the read data will be stored.
	// The method returns the number of bytes successfully read into the buffer and any error encountered.
	// Implementations should monitor the port for incoming data, transferring it into the buffer
	// up to the buffer's capacity, and return the actual amount of data read.
	Read([]byte) (int, error)

	// Close terminates the connection to the serial port and releases any resources associated with it.
	// It should ensure that the port is properly closed and the underlying hardware is left in a clean state.
	// The method returns an error if any issues occur during the closure process.
	Close() error
}

// NewPort initialize serial port object.
func NewPort(path string, rate int) (Port, error) {
	return newSerial(path, rate, 8, 0, 0)
}
