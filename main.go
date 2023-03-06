package main

import (
	"lemur/messagequeue/server"
)

func main() {
	s := server.NewServer()

	s.Start()
}
