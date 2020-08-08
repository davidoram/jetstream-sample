package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	nats "github.com/nats-io/nats.go"
)

func main() {
	maxClient := flag.Int("max-client", 1, "an int")
	flag.Parse()

	log.Printf("Server v1 started")
	log.Printf("-----------------")
	log.Printf("1->%d clients", *maxClient)

	nc, err := nats.Connect("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer nc.Close()
	log.Printf("event listener connected")

	// Listen for 'device.event'
	go runEventListener(nc)

	// Post 'config.changed'	int
	for {
		time.Sleep(time.Second * 10)
		for i := 1; i <= *maxClient; i += 1 {
			log.Printf("PUB config.changed %d", i)
			msg := &nats.Msg{
				Subject: fmt.Sprintf("config.changed.%d", i),
				Data:    []byte(fmt.Sprintf("New config for client %d", i)),
			}
			nc.PublishMsg(msg)
		}
	}
}

func runEventListener(nc *nats.Conn) {
	sub, _ := nc.SubscribeSync("device.event.>")
	for {
		m, err := sub.NextMsg(5 * time.Second)
		if err == nil {
			log.Printf("MSG %s %s\n", m.Subject, string(m.Data))
		} else {
			//log.Println("NextMsg timed out.")
		}
	}
}
