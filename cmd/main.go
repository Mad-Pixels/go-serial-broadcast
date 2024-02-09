package main

import (
	"fmt"
	"log"

	go_serial_broadcast "github.com/MadPixeles/go-serial-broadcast"
	"github.com/MadPixeles/go-serial-broadcast/verification"
)

func main() {
	//p, err := port.NewPort()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//b := bytes.NewBuffer([]byte{})

	//b := make([]byte, 4095)
	//for {
	//	fmt.Println(p.Read(b))
	//	fmt.Println(string(b))
	//}

	//bcast
	v, _ := verification.NewByMask([]byte{}, "*")
	bcast, err := go_serial_broadcast.NewBroadcast(v)
	if err != nil {
		log.Fatal(err)
	}
	bcast.AddHandler("customCommand", func(msg string) error {
		fmt.Println("Custom command received:", msg)
		return nil
	})
	go bcast.Read()
	go bcast.HandleMessages()

	select {}
}
