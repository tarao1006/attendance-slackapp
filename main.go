package main

import (
	"log"

	"github.com/tarao1006/attendance-slackapp/server"
)

func main() {
	s := server.NewServer()
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
