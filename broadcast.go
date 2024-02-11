package go_serial_broadcast

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"unsafe"

	"github.com/MadPixeles/go-serial-broadcast/port"
)

// MessageHandler defines a function signature for handling messages received from the serial port.
// It takes a message as a string and returns an error if the handling fails.
type MessageHandler func(msg string) error

// Broadcast manages serial communication, parsing and routing messages to appropriate handlers.
// It uses a mutex to safely access a buffer that stores incoming serial data until it can be processed.
// Messages are then dispatched based on predefined handlers for specific commands or patterns.
type Broadcast struct {
	// Protects access to the internal buffer, ensuring thread-safe operations.
	bufferMutex sync.Mutex

	// Abstraction over a serial port, allowing for reading from and writing to the serial device.
	serial port.Port

	// A channel for dispatching processed messages to be handled by registered handlers.
	messages chan []byte

	// A channel for control numb of goroutines.
	semaphore chan struct{}

	// Maps command strings or message prefixes to their corresponding handlers.
	customHandlers map[string]MessageHandler

	// A fallback handler used when no specific handler is found for a message.
	defaultHandler MessageHandler

	// Temporarily stores incoming data from the serial port until it can be processed.
	buffer bytes.Buffer
}

// NewBroadcast creates a new Broadcast instance with the specified serial port path and rate.
func NewBroadcast(path string, rate, flows int) (*Broadcast, error) {
	serial, err := port.NewPort(path, rate)
	if err != nil {
		return nil, fmt.Errorf("broadcast error: %w", err)
	}
	return &Broadcast{
		customHandlers: make(map[string]MessageHandler),
		semaphore:      make(chan struct{}, flows),
		messages:       make(chan []byte, flows*3),
		buffer:         bytes.Buffer{},
		serial:         serial,
	}, nil
}

// AddHandler registers a custom handler for messages starting with a specific key (without spaces).
// This method enables the Broadcast instance to dynamically handle a variety of message types
// received through the serial port. Associating a handler with a unique prefix allows the application
// to categorize and process different messages based on their initial sequence or command keyword.
// The corresponding handler is invoked to process any message that begins with the matched prefix.
//
// Parameters:
//   - prefix: A string that uniquely identifies the prefix of the messages the handler is designed to process.
//     This could be a specific command keyword or any identifiable starting sequence within the message data.
//   - handler: A MessageHandler function designated to process messages that match the prefix.
//     The handler must accept a message as a slice of bytes and return an error if the processing encounters any issues.
//
// Usage example:
//
//	b.AddHandler("CMD1:", func(msg []byte) error {
//	    // Logic to handle the message
//	    return nil
//	})
func (b *Broadcast) AddHandler(key string, handler MessageHandler) {
	b.customHandlers[key] = handler
}

// SetDefaultHandler sets a default handler that is called when no specific handler
// is found for a message. This allows for a fallback processing mechanism for messages
// that do not match any of the predefined command patterns. The default handler
// can be used to log unhandled messages, perform cleanup, or even as a catch-all
// processor for generic message handling.
//
// Parameters:
//   - handler: A MessageHandler function to be used as the default handler. It should
//     accept a message as a string and return an error if processing fails.
//
// Usage example:
//
//	b.SetDefaultHandler(func(string) error {
//	    fmt.Println("Default handler received message:", msg)
//	    return nil
//	})
func (b *Broadcast) SetDefaultHandler(handler MessageHandler) {
	b.defaultHandler = handler
}

// HandleMessages listens for incoming messages on the messages channel and dispatches
// them to the appropriate handlers based on their content. It uses the first part of the
// message, separated by a space, as a key to identify the correct handler from the
// customHandlers map.
//
// If no specific handler is found for a key, the default handler is invoked (if it has been set).
//
// Usage example:
//
//	err := b.HandleMessages()
//	if err != nil {
//	    log.Fatalf("Error handling messages: %v", err)
//	}
func (b *Broadcast) HandleMessages() error {
	for msg := range b.messages {
		b.semaphore <- struct{}{}

		go func(msg []byte) {
			defer func() { <-b.semaphore }()

			key := bytes.SplitN(msg, []byte(" "), 2)
			handler, ok := b.customHandlers[*(*string)(unsafe.Pointer(&key[0]))]
			if !ok {
				handler = b.defaultHandler
			}
			if handler != nil {
				if err := handler(*(*string)(unsafe.Pointer(&msg))); err != nil {
					//return fmt.Errorf(
					//	"broadcast handler error with key %s, got %w",
					//	*(*string)(unsafe.Pointer(&key[0])),
					//	err,
					//)
				}
			}
		}(msg)
	}
	return nil
}

// Read continuously reads from the serial port into an internal buffer.
// Triggers the asynchronous processing of messages and execute MessageHandler functions.
// The method exits gracefully upon encountering the io.EOF error, signaling the end of the data stream.
//
// Usage example:
//
//	go broadcastInstance.Read(1024) // Initiates reading from the serial port with a 1024-byte buffer.
//
// Note: It's recommended to run Read in its own goroutine to facilitate continuous reading and processing of serial data.
func (b *Broadcast) Read(bufferSize int) error {
	tmp := make([]byte, bufferSize)
	for {
		n, err := b.serial.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("broadcast error: %w", err)
			}
			break
		}
		if n > 0 {
			b.bufferMutex.Lock()
			b.buffer.Write(tmp[:n])
			b.bufferMutex.Unlock()

			b.processMessages()
		}
	}
	return nil
}

func (b *Broadcast) processMessages() {
	for {
		b.bufferMutex.Lock()
		index := bytes.IndexByte(b.buffer.Bytes(), '\n')
		if index != -1 {
			fromBuffer := b.buffer.Next(index + 1)
			msg := make([]byte, index)
			copy(msg, fromBuffer[:index])

			b.bufferMutex.Unlock()
			b.messages <- msg
			continue
		}
		b.bufferMutex.Unlock()
		break
	}
}
