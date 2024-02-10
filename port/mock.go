package port

import (
	"github.com/stretchr/testify/mock"
)

type serialMock struct {
	mock.Mock
}

func newSerialMock() Port {
	return &serialMock{}
}

// Read simulates reading data from the serial port mock.
//
// Usage:
//
//	func TestReadSerialPort(t *testing.T) {
//	 	serialMock := port.NewPortMock()
//	 	// Prepare a buffer to read data into and set up the expected behavior for Read.
//	 	// This includes the data to be read into the buffer and the return values
//	 	// (number of bytes read and any error).
//	 	buffer := make([]byte, 10)
//	 	expectedData := []byte("test")
//	 	serialMock.On("Read", mock.Anything).Return(len(expectedData), nil)
//
//	 	// Attempt to read from the mock serial port.
//	 	n, err := serialMock.Read(buffer)
//
//	 	// Verify that the correct number of bytes was read and no error was returned.
//	 	assert.NoError(t, err)
//	 	assert.Equal(t, len(expectedData), n)
//	 	// Optionally, check that the data read matches the expected data.
//	 	assert.Equal(t, expectedData, buffer[:n])
//
//	 	// Check that the Read method was called with the correct parameters.
//	 	serialMock.AssertCalled(t, "Read", mock.Anything)
//	 }
func (m *serialMock) Read(buf []byte) (int, error) {
	args := m.Called(buf)
	return args.Int(0), args.Error(1)
}

// Write simulates writing data to the serial port mock.
//
// Usage:
//
//	func TestWriteSerialPort(t *testing.T) {
//	 	serialMock := port.NewPortMock()
//	 	// Set up the expected behavior for Write, including the data to be written
//	 	// and the return values (number of bytes written and any error).
//	 	serialMock.On("Write", []byte("test data")).Return(len("test data"), nil)
//
//	 	// Attempt to write to the mock serial port.
//	 	n, err := serialMock.Write([]byte("test data"))
//
//	 	// Verify that the correct number of bytes was written and no error was returned.
//	 	assert.NoError(t, err)
//	 	assert.Equal(t, len("test data"), n)
//
//	 	// Check that the Write method was called with the expected parameters.
//	 	serialMock.AssertCalled(t, "Write", []byte("test data"))
//	 }
func (m *serialMock) Write(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

// Close simulates closing the serial port mock.
//
// Usage:
//
//	func TestCloseSerialPort(t *testing.T) {
//	 	serialMock := port.NewPortMock()
//	 	serialMock.On("Close").Return(nil) // Set up the expected behavior.
//
//	 	err := serialMock.Close() // Attempt to close the mock serial port.
//	 	assert.NoError(t, err) // Verify that no error was returned.
//	 }
func (m *serialMock) Close() error {
	args := m.Called()
	return args.Error(0)
}
