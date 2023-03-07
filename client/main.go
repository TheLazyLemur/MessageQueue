package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type ServerMessage struct {
	Type      string
	QueueName string
	Message   string
}

var (
	isConsumer bool
)

func main() {

	flag.BoolVar(&isConsumer, "c", true, "consumer mode")
	flag.Parse()

	c, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	if isConsumer {
		sendJoinMessage(c)

		for {
			read(c)
		}
	}

	if !isConsumer {
		publishMessage(c)
	}
}

func RandomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes)[:length]
}

func publishMessage(conn net.Conn) {
	for {
		time.Sleep(time.Second)

		m := ServerMessage{
			QueueName: "test",
			Type:      "pub",
			Message:   RandomString(10),
		}

		jsonBytes, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
			return
		}

		buf := new(bytes.Buffer)

		keyLen := int32(len(jsonBytes))
		_ = binary.Write(buf, binary.LittleEndian, keyLen)

		_ = binary.Write(buf, binary.LittleEndian, jsonBytes)

		_, _ = conn.Write(buf.Bytes())

		fmt.Println("Sending: ", string(m.Message))
	}

}

func sendJoinMessage(conn net.Conn) {
	m := ServerMessage{
		QueueName: "test",
		Type:      "join",
	}

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)

	keyLen := int32(len(jsonBytes))
	_ = binary.Write(buf, binary.LittleEndian, keyLen)

	_ = binary.Write(buf, binary.LittleEndian, jsonBytes)

	_, _ = conn.Write(buf.Bytes())

}

func read(r io.Reader) {
	var mlen int32
	_ = binary.Read(r, binary.LittleEndian, &mlen)

	buf := make([]byte, mlen)
	_ = binary.Read(r, binary.LittleEndian, &buf)

	fmt.Println("Received: ", string(buf))
}
