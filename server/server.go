package server

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	queues   map[string][]*net.Conn
	messages *Queue
}

type ServerMessage struct {
	Type      string
	QueueName string
	Message   string
}

func NewServer() *Server {
	return &Server{
		queues:   make(map[string][]*net.Conn),
		messages: NewQueue(),
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

func (s *Server) Start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			s.messages.Enqueue(RandomString(10))
		}
	}()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatal(err)
			}

			s.parseMessage(conn, conn)
		}
	}()

	for {
		time.Sleep(5 * time.Second)

		for _, q := range s.queues {
			m := s.messages.Dequeue()
			for _, conn := range q {
				if len(m) > 0 {
					sendMessage(*conn, *conn, m)
				}
			}
		}
	}
}

func (s *Server) parseMessage(r io.Reader, conn net.Conn) {
	var keyLen int32
	_ = binary.Read(r, binary.LittleEndian, &keyLen)

	msgBuf := make([]byte, keyLen)
	_ = binary.Read(r, binary.LittleEndian, &msgBuf)

	log.Println(string(msgBuf))

	if len(msgBuf) > 0 {
		m := new(ServerMessage)

		err := json.Unmarshal([]byte(msgBuf), &m)
		if err != nil {
			log.Fatal("Error converting to struct:", err)
		}

		if m.Type == "join" {
			s.queues[m.QueueName] = append(s.queues[m.QueueName], &conn)
			log.Printf("Joined queue %s\n", m.QueueName)
		}

		if m.Type == "pub" {
			s.messages.Enqueue(m.Message)
			log.Printf("Published message %s\n", m.Message)
		}
	}
}

func sendMessage(r io.Reader, conn net.Conn, message string) {
	buf := new(bytes.Buffer)

	k := message
	keyLen := int32(len([]byte(k)))

	_ = binary.Write(buf, binary.LittleEndian, keyLen)

	_ = binary.Write(buf, binary.LittleEndian, []byte(k))

	_, _ = conn.Write(buf.Bytes())
}
