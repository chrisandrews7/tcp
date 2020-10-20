package main

import (
	"encoding/json"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	Run() error
	Close() error
}

type TCPServer struct {
	address  string
	server   net.Listener
	handlers *Handlers
}

func (s *TCPServer) handleConnection(connection net.Conn) {
	var userID int

	defer func() {
		s.handlers.UserLeftHandler(connection, userID)
		connection.Close()
	}()

	for {
		decoder := json.NewDecoder(connection)
		// @todo Add switch with custom json unmarshaller to handle different json payloads
		var message UserJoinRequest
		err := decoder.Decode(&message)

		if err != nil {
			if err != io.EOF {
				log.Error(err)
			}
			break
		}

		userID = s.handlers.UserJoinHandler(connection, message)
	}
}

func (s *TCPServer) Run() error {
	server, err := net.Listen("tcp", s.address)
	s.server = server

	if err != nil {
		return err
	}

	for {
		connection, err := server.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(connection)
	}
}

func (s *TCPServer) Close() error {
	return s.server.Close()
}

func NewTCPServer(address string) Server {
	return &TCPServer{
		address:  address,
		handlers: NewHandlers(),
	}
}
