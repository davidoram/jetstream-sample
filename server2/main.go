package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	nats "github.com/nats-io/nats.go"
)

func main() {
	maxClient := flag.Int("m", 1, "Clients to send to 1..m, an int")
	username := flag.String("u", "", "Nats username")
	password := flag.String("p", "", "Nats password")
	flag.Parse()

	log.Printf("Server 1 started. 1->%d clients", *maxClient)

	opts := nats.Options{
		Servers:  []string{"localhost"},
		User:     *username,
		Password: *password,
	}
	nc, err := opts.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer nc.Close()
	log.Printf("event listener connected")

	// Listen for 'device.event'
	go runEventListener(nc, *username)

	// Post 'config.changed'	int
	for {
		time.Sleep(time.Second * 10)
		for i := 1; i <= *maxClient; i += 1 {
			subject := fmt.Sprintf("config.changed.%d", i)
			log.Printf("%s PUB %s", *username, subject)
			msg := &nats.Msg{
				Subject: subject,
				Data:    []byte(fmt.Sprintf("New config for client %d", i)),
			}
			nc.PublishMsg(msg)
		}
	}
}

func runEventListener(nc *nats.Conn, user string) {
	sub, _ := nc.SubscribeSync("device.event.>")
	for {
		m, err := sub.NextMsg(5 * time.Second)
		if err == nil {
			log.Printf("%s MSG %s %s\n", user, m.Subject, string(m.Data))
		} else {
			//log.Println("NextMsg timed out.")
		}
	}
}
