package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/thelazylemur/messagequeue/client"
	"github.com/thelazylemur/messagequeue/server"
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
	publishedMessages := 0
	for {
		time.Sleep(time.Millisecond * 20)
		publisher.PublishMessage(randomString(10))
		publishedMessages++
		fmt.Println("Published:", publishedMessages)
	}
}

func runConsumer() {
	consumer := client.NewSubscriber(":3000", "test")
	consumerChan := make(chan string)
	recievedMessages := make([]string, 0)

	go func() {
		for {
			consumedMessage := <-consumerChan
			recievedMessages = append(recievedMessages, consumedMessage)
			fmt.Println("Recieved:", len(recievedMessages))
		}
	}()

	consumer.ReadFromQueue(consumerChan)
}
