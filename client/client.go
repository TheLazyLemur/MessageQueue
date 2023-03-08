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

type Client struct {
	ServerAddrr string
	conn        net.Conn
	QueueName   string
}

func NewSubscriber(serverAddrr string, queueName string) *Client {
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

func (c *Client) SendJoinMessage() {
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

func (c *Client) ReadFromQueue() {
	var mlen int32
	_ = binary.Read(c.conn, binary.LittleEndian, &mlen)
	if mlen == 0 {
		return
	}

	buf := make([]byte, mlen)
	_ = binary.Read(c.conn, binary.LittleEndian, &buf)

	fmt.Println("Received: ", string(buf))
}
