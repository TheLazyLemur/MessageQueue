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

		wg := sync.WaitGroup{}
		for queueName, cons := range s.queueNameToRecepient {
			queue := s.GetQueueOrCreateIfNotExists(queueName)
			message := queue.Dequeue()

			for _, conn := range cons {
				wg.Add(1)
				go func(conn net.Conn) {
					defer wg.Done()
					sendMessage(conn, conn, message)
				}(*conn)
			}
		}
		wg.Wait()
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

			_ = s.GetQueueOrCreateIfNotExists(m.QueueName)
			log.Printf("Joined queue %s\n", m.QueueName)
			s.lock.Unlock()
		}

		if m.Type == "pub" {
			q := s.GetQueueOrCreateIfNotExists(m.QueueName)

			q.Enqueue(m.Message)
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

			q := s.GetQueueOrCreateIfNotExists(m.QueueName)
			q.Enqueue(m.Message)
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

func (s *Server) GetQueueOrCreateIfNotExists(queueName string) *Queue {
	q, ok := s.queueNameToQueue[queueName]
	if !ok {
		s.queueNameToQueue[queueName] = NewQueue()
		return s.queueNameToQueue[queueName]
	}

	return q
}
