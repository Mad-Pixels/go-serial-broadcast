package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	go_serial_broadcast "github.com/MadPixeles/go-serial-broadcast"
	"github.com/MadPixeles/go-serial-broadcast/middleware"
)

func printStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Printf("\tGoroutines = %d\n", runtime.NumGoroutine())
}

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

	go func() {
		for {
			time.Sleep(time.Second)
			printStats()
		}
	}()

	//bcast
	v, _ := middleware.NewVerifyByMask([]byte{}, "*")
	bcast, err := go_serial_broadcast.NewBroadcast("/dev/tty.usbserial-110", 9600, 4, v)
	if err != nil {
		log.Fatal(err)
	}
	bcast.AddHandler("customCommand", func(msg string) error {
		fmt.Println("Custom command received:", msg)
		return nil
	})
	rand.Seed(time.Now().UnixNano())
	bcast.SetDefaultHandler(func(msg string) error {
		randomInt := rand.Intn(20)
		time.Sleep(time.Second * time.Duration(randomInt))
		fmt.Println(msg)
		//
		//
		return nil
	})
	go bcast.Read(1024)
	//_ = bcast.HandleMessages()

	select {}
}
