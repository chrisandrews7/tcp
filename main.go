package main

import (
	log "github.com/sirupsen/logrus"
)

const address = ":8080"

func main() {
	server := NewTCPServer(address)
	err := server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
