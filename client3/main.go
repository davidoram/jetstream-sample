package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	nats "github.com/nats-io/nats.go"
)

func main() {
	clientNum := flag.Int("n", 1, "Client identifier, an int")
	username := flag.String("u", "", "Nats username")
	password := flag.String("p", "", "Nats password")
	server := flag.String("s", "", "Nats server")
	flag.Parse()
	log.Printf("Client v3")
	log.Printf("---------")
	log.Printf("Id %d started, connecting to NATs: %s", *clientNum, *server)

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

	// Listen for 'config.changed'
	go runEventListener(nc, *clientNum)

	// Post 'iot.event'
	i := 0
	subject := "iot.event"
	for {
		time.Sleep(time.Second * 10)
		e := fmt.Sprintf("Event on device %d, %d", *clientNum, rand.Int())
		log.Printf("PUB subject: '%s', msg: '%s'", subject, e)
		msg := &nats.Msg{
			Subject: subject,
			Data:    []byte(e),
		}
		i = i + 1
		nc.PublishMsg(msg)
	}
}

func runEventListener(nc *nats.Conn, clientNum int) {
	subject := "config.changed"
	sub, _ := nc.SubscribeSync(subject)
	for {
		m, err := sub.NextMsg(5 * time.Second)
		if err == nil {
			log.Printf("MSG subject: '%s', data: '%s'\n", m.Subject, string(m.Data))
		} else {
			//log.Println("NextMsg timed out.")
		}
	}
}
