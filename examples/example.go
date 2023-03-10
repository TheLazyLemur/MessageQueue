package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"lemur/messagequeue/client"
	"lemur/messagequeue/server"
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
	go func() {
		runServer()
	}()

	go func() {
		time.Sleep(time.Second * 3)
		runConsumer()
	}()

	go func() {
		time.Sleep(time.Second * 2)
		runPublisher()

	}()

	select {}
}

func runServer() {
	s := server.NewServer()
	s.Start()
}

func runPublisher() {
	publisher := client.NewPublisher(":3000", "test")
	for {
		time.Sleep(time.Millisecond * 20)
		publisher.PublishMessage(randomString(10))
	}
}

func runConsumer() {
	consumer := client.NewSubscriber(":3000", "test")
	consumerChan := make(chan string)

	go func() {
		for {
			x := <-consumerChan
			fmt.Println("Consuming", x)
		}
	}()

	consumer.ReadFromQueue(consumerChan)
}
