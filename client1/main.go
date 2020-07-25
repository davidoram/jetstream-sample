package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	nats "github.com/nats-io/nats.go"
)

func main() {
	clientNum := flag.Int("client-num", 1, "an int")
	flag.Parse()
	log.Printf("Client %d started", *clientNum)

	nc, err := nats.Connect("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer nc.Close()
	log.Printf("event listener connected")

	// Listen for 'config.changed'
	go runEventListener(nc, *clientNum)

	// Post 'device.event.%d'	int
	i := 0
	for {
		time.Sleep(time.Second * 10)
		log.Printf("PUB device.event.%d", *clientNum)
		msg := &nats.Msg{
			Subject: fmt.Sprintf("device.event.%d", *clientNum),
			Data:    []byte(fmt.Sprintf("Event on device %d", *clientNum)),
		}
		i = i + 1
		nc.PublishMsg(msg)
	}
}

func runEventListener(nc *nats.Conn, clientNum int) {
	sub, _ := nc.SubscribeSync(fmt.Sprintf("config.changed.%d", clientNum))
	for {
		m, err := sub.NextMsg(5 * time.Second)
		if err == nil {
			log.Printf("MSG %s %s\n", m.Subject, string(m.Data))
		} else {
			//log.Println("NextMsg timed out.")
		}
	}
}
