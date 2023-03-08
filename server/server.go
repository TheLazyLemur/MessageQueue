package server

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	queueNameToRecepient map[string][]*net.Conn
	messages             *Queue
	queueNameToQueue     map[string]*Queue
	lock                 sync.Mutex
}

type ServerMessage struct {
	Type      string
	QueueName string
	Message   string
}

func NewServer() *Server {
	return &Server{
		queueNameToRecepient: make(map[string][]*net.Conn),
		messages:             NewQueue(),
		lock:                 sync.Mutex{},
		queueNameToQueue:     make(map[string]*Queue),
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
			conn, err := ln.Accept()
			if err != nil {
				log.Fatal(err)
			}

			s.parseMessage(conn, conn)
		}
	}()

	for {
		s.lock.Lock()

		for queueName := range s.queueNameToRecepient {
			connections := s.queueNameToRecepient[queueName]
			queue := s.queueNameToQueue[queueName]
			message := queue.Dequeue()
			for _, conn := range connections {
				sendMessage(*conn, *conn, message)
			}
		}

		s.lock.Unlock()
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
		fmt.Println("Queue name:", m.QueueName)

		if m.Type == "join" {
			s.lock.Lock()
			s.queueNameToRecepient[m.QueueName] = append(s.queueNameToRecepient[m.QueueName], &conn)

			_, ok := s.queueNameToQueue[m.QueueName]
			if !ok {
				s.queueNameToQueue[m.QueueName] = NewQueue()
			}

			log.Printf("Joined queue %s\n", m.QueueName)
			s.lock.Unlock()
		}

		if m.Type == "pub" {
			_, ok := s.queueNameToQueue[m.QueueName]
			if !ok {
				s.queueNameToQueue[m.QueueName] = NewQueue()
			}

			s.messages.Enqueue(m.Message)
			go s.handleQueue(r, conn)
		}
	}
}

func (s *Server) handleQueue(r io.Reader, conn net.Conn) {
	for {
		var keyLen int32
		_ = binary.Read(r, binary.LittleEndian, &keyLen)

		if keyLen > 0 {
			msgBuf := make([]byte, keyLen)
			_ = binary.Read(r, binary.LittleEndian, &msgBuf)

			m := new(ServerMessage)

			err := json.Unmarshal([]byte(msgBuf), &m)
			if err != nil {
				log.Fatal("Error converting to struct:", err)
			}

			s.messages.Enqueue(m.Message)
			s.queueNameToQueue[m.QueueName].Enqueue(m.Message)
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
