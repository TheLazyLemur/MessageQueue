package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"lemur/messagequeue/client"
	"lemur/messagequeue/server"
	"time"
)

var (
	isClient   bool
	isConsumer bool
	queueName  string
)

func randomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes)[:length]
}

func main() {
	flag.BoolVar(&isClient, "cl", true, "is client")
	flag.BoolVar(&isConsumer, "c", true, "consumer mode")
	flag.StringVar(&queueName, "q", "", "queue name")
	flag.Parse()

	if isClient {
		if isConsumer {
			client := client.NewClient(":3000", client.Subscriber, queueName)

			client.SendJoinMessage()

			for {
				client.ReadFromQueue()
			}
		}

		if !isConsumer {
			client := client.NewClient(":3000", client.Publisher, queueName)

			for {
				time.Sleep(time.Second)
				client.PublishMessage(randomString(10))
			}
		}
	} else {
		s := server.NewServer()
		s.Start()
	}

}
