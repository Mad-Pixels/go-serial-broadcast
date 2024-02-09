package go_serial_broadcast

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/MadPixeles/go-serial-broadcast/port"
	"github.com/MadPixeles/go-serial-broadcast/verification"
)

type MessageHandler func(msg string) error

// Broadcast ...
type Broadcast struct {
	verification verification.Interface
	serial       port.Interface
	buffer       bytes.Buffer
	messages     chan string
	bufferMutex  sync.Mutex

	customHandlers map[string]MessageHandler
}

// NewBroadcast ...
func NewBroadcast(verifyMethod verification.Interface) (*Broadcast, error) {
	serial, err := port.NewPort()
	if err != nil {
		return nil, err
	}
	return &Broadcast{
		verification: verifyMethod,
		serial:       serial,
		messages:     make(chan string, 100),
		buffer:       bytes.Buffer{},
	}, nil
}

func (b *Broadcast) Read() {
	tmp := make([]byte, 1024)
	for {
		n, err := b.serial.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: %s", err)
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
}

func (b *Broadcast) processMessages() {
	for {
		b.bufferMutex.Lock()
		msg, err := b.buffer.ReadString('\n')
		b.bufferMutex.Unlock()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error processing messages: %s", err)
			}
			if msg != "" {
				b.bufferMutex.Lock()
				b.buffer.WriteString(msg)
				b.bufferMutex.Unlock()
			}
			break
		}
		b.messages <- msg
	}
}

func (b *Broadcast) AddHandler(command string, handler MessageHandler) {
	if b.customHandlers == nil {
		b.customHandlers = make(map[string]MessageHandler)
	}
	b.customHandlers[command] = handler
}

//func (b *Broadcast) processMessages() {
//	for {
//		msg, err := b.buffer.ReadString('\n')
//		if err != nil {
//			if err == io.EOF {
//				b.buffer.WriteString(msg)
//			} else {
//				log.Printf("Error processing messages: %s", err)
//			}
//			break
//		}
//		fmt.Println("Received message:", msg)
//	}
//}

func (b *Broadcast) HandleMessages() {
	for msg := range b.messages {
		trimmedMsg := strings.TrimSuffix(msg, "\n")
		// Обработка сообщения
		switch {
		case msg == "command1\n":
			fmt.Println("Handling command1")
			// Ваш код для обработки command1
		case msg == "command2\n":
			fmt.Println("Handling command2")
			// Ваш код для обработки command2
		default:
			if handler, exists := b.customHandlers[trimmedMsg]; exists {
				err := handler(msg)
				if err != nil {
					fmt.Printf("Error handling custom command '%s': %v\n", trimmedMsg, err)
				}
			} else {
				fmt.Println("Unknown command:", msg)
			}
		}

	}
}
