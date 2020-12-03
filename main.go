package main

import (
	"log"

	"github.com/tarao1006/attendance-slackapp/server"
)

func main() {
	s := server.NewServer()
	s.Init()
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
