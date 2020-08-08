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
	server := flag.String("s", "", "Nats server")

	flag.Parse()

	log.Printf("Server v3 started")
	log.Printf("-----------------")
	log.Printf("Connecting to NATs: %s, 1->%d clients", *server, *maxClient)

	opts := nats.Options{
		Servers:  []string{*server},
		User:     *username,
		Password: *password,
	}
	nc, err := opts.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer nc.Close()
	log.Printf("event listener connected")

	// Listen for 'iot.event'
	go runEventListener(nc)

	// Post 'config.changed'	int
	for {
		time.Sleep(time.Second * 10)
		for i := 1; i <= *maxClient; i += 1 {
			subject := fmt.Sprintf("config.changed.%d", i)
			e := fmt.Sprintf("Config change for client %d", i)
			log.Printf("PUB subject: %s, data: '%s'", subject, e)
			msg := &nats.Msg{
				Subject: subject,
				Data:    []byte(e),
			}
			nc.PublishMsg(msg)
		}
	}
}

func runEventListener(nc *nats.Conn) {
	sub, _ := nc.SubscribeSync("remote.events.>")
	for {
		m, err := sub.NextMsg(5 * time.Second)
		if err == nil {
			log.Printf("MSG subject: '%s', data: '%s'\n", m.Subject, string(m.Data))
		} else {
			//log.Println("NextMsg timed out.")
		}
	}
}
