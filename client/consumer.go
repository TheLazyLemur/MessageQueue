package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type ServerMessage struct {
	Type      string
	QueueName string
	Message   string
}

type ConsumerClient struct {
	ServerAddrr string
	conn        net.Conn
	QueueName   string
}

func NewConsumer(serverAddrr string, queueName string) *ConsumerClient {
	c := &ConsumerClient{
		ServerAddrr: serverAddrr,
		QueueName:   queueName,
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn

	c.sendJoinMessage()

	return c
}

func (c *ConsumerClient) sendJoinMessage() {
	m := ServerMessage{
		Type:      "join",
		QueueName: c.QueueName,
		Message:   "",
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

	_, _ = c.conn.Write(buf.Bytes())

}

func (c *ConsumerClient) ReadFromQueue(consumeChan chan string) {
	for {
		var mlen int32
		_ = binary.Read(c.conn, binary.LittleEndian, &mlen)
		if mlen != 0 {
			buf := make([]byte, mlen)
			_ = binary.Read(c.conn, binary.LittleEndian, &buf)

			consumeChan <- string(buf)
			//TODO: The consumer should acknowledge the message was received
		}
	}
}
