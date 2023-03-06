package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type JoinMessage struct {
	Type      string
	QueueName string
	Message   string
}

func main() {
	c, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	sendJoinMessage(c)

	for {
		read(c)
	}
}

func sendJoinMessage(conn net.Conn) {
	m := JoinMessage{
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

	if len(buf) > 0 {
		fmt.Printf("%s\n", string(buf))
	}
}
