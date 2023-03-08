package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func NewPublisher(serverAddrr string, queueName string) *Client {
	c := &Client{
		ServerAddrr: serverAddrr,
		QueueName:   queueName,
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn

	return c
}

func (c *Client) PublishMessage(msg string) {
	m := ServerMessage{
		QueueName: c.QueueName,
		Type:      "pub",
		Message:   msg,
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

	fmt.Println("Sending: ", string(m.Message))
}
