package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	nats "github.com/nats-io/nats.go"
)

func main() {
	clientNum := flag.Int("n", 1, "Client identifier, an int")
	username := flag.String("u", "", "Nats username")
	password := flag.String("p", "", "Nats password")
	server := flag.String("s", "", "Nats server")
	flag.Parse()
	log.Printf("Client %d started, connecting to NATs: %s", *clientNum, *server)

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
	go runEventListener(nc, *clientNum, *username)

	// Post 'iot.event'
	i := 0
	subject := fmt.Sprintf("iot.event.%d", *clientNum)
	for {
		time.Sleep(time.Second * 10)
		log.Printf("%s PUB %s", *username, subject)
		msg := &nats.Msg{
			Subject: subject,
			Data:    []byte(fmt.Sprintf("Event on device %d", *clientNum)),
		}
		i = i + 1
		nc.PublishMsg(msg)
	}
}

func runEventListener(nc *nats.Conn, clientNum int, user string) {
	subject := "config.changed.>"
	sub, _ := nc.SubscribeSync(subject)
	for {
		m, err := sub.NextMsg(5 * time.Second)
		if err == nil {
			log.Printf("%s MSG %s %s\n", user, m.Subject, string(m.Data))
		} else {
			//log.Println("NextMsg timed out.")
		}
	}
}
