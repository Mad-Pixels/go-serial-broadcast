package go_serial_broadcast

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sync"
	"unsafe"

	"github.com/MadPixeles/go-serial-broadcast/port"
	"github.com/pkg/errors"
)

// MessageHandler defines a function signature for handling messages received from the serial port.
// It takes a message as a string and returns an error if the handling fails.
type MessageHandler func(msg string) error

// Broadcast manages serial communication, parsing and routing messages to appropriate handlers.
type Broadcast struct {
	bufferMutex    sync.Mutex
	serial         port.Port
	messages       chan []byte
	semaphore      chan struct{}
	customHandlers map[string]MessageHandler
	defaultHandler MessageHandler
	buffer         bytes.Buffer
}

// NewBroadcast creates a new Broadcast instance with the specified serial port path and rate.
func NewBroadcast(path string, rate, flows int) (*Broadcast, error) {
	serial, err := port.NewPort(path, rate)
	if err != nil {
		return nil, errors.Wrapf(err, "boradcast initialization failed, got")
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
//	    fmt.Println("my custom handler for messages which started with 'CMD1'")
//	    return nil
//	})
func (b *Broadcast) AddHandler(key string, handler MessageHandler) {
	b.customHandlers[key] = handler
}

// SetDefaultHandler sets a default handler that is called when no specific handler is found for a message.
//
// Parameters:
//   - handler: A MessageHandler function to be used as the default handler. It should
//     accept a message as a string and return an error if processing fails.
//
// Usage example:
//
//	bcast, _ := go_serial_broadcast.NewBroadcast("/dev/tty.usbserial-110", 9600, 4)
//	bcast.SetDefaultHandler(func(string) error {
//	    fmt.Println("Default handler received message:", msg)
//	    return nil
//	})
func (b *Broadcast) SetDefaultHandler(handler MessageHandler) {
	b.defaultHandler = handler
}

// HandleMessages listens for incoming messages and dispatches them to the appropriate handlers
// based on their content. If no specific handler is found for a key, the default handler is
// invoked (if it has been set).
//
// Usage example:
//
//	bcast, _ := go_serial_broadcast.NewBroadcast("/dev/tty.usbserial-110", 9600, 4)
//	bcast.SetDefaultHandler(func(msg string) error {
//		fmt.Println(msg)
//		return nil
//	})
//
//	bcast.HandleMessages(nil)
//	// or
//	go bcast.HandleMessages(nil)
//	// or
//	errCh := make(chan err)
//	go bcast.HandleMessages(errCh)
//	for {
//		select:
//		case e := <-errCh:
//			fmt.Println(e)
//		default:
//	}
func (b *Broadcast) HandleMessages(errFlow chan error) {
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
				if errFlow == nil {
					_ = handler(*(*string)(unsafe.Pointer(&msg)))
					return
				}
				errFlow <- handler(*(*string)(unsafe.Pointer(&msg)))
				return
			}
		}(msg)
	}
}

// Read continuously reads from the serial port into an internal buffer.
// Triggers the asynchronous processing of messages and execute MessageHandler functions.
//
// Usage example:
//
//	bcast, _ := go_serial_broadcast.NewBroadcast("/dev/tty.usbserial-110", 9600, 4)
//	bcast.SetDefaultHandler(func(msg string) error {
//		fmt.Println(msg)
//		return nil
//	})
//
//	go bcast.Read(1024)
//	_ = b.HandleMessages()
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
	}
	return nil
}

// Write data to the serial port.
func (b *Broadcast) Write(msg string) (int, error) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&msg))
	slice := (*[1 << 30]byte)(unsafe.Pointer(sh.Data))[:sh.Len:sh.Len]

	return b.serial.Write(append(slice, '\n'))
}
