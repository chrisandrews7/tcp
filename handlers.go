package main

import (
	"encoding/json"
	"net"

	log "github.com/sirupsen/logrus"
)

type Handlers struct {
	users *UserStore
}

type userRecord struct {
	connection net.Conn
	friends    []int
	userID     int
}

type UserJoinRequest struct {
	Friends []int `json:"friends"`
	UserID  int   `json:"user_id"`
}

type userStatusMessage struct {
	Online bool `json:"online"`
	UserID int  `json:"user_id"`
}

func (h *Handlers) broadcastToFriends(userID int, message []byte) {
	onlineFriends := h.users.GetRelated(userID)

	log.WithFields(log.Fields{
		"userID":     userID,
		"totalUsers": len(onlineFriends),
	}).Info("Broadcasting to other users")

	for _, friend := range onlineFriends {
		go func(connection net.Conn) {
			connection.Write(message)
		}(friend.(userRecord).connection)
	}
}

func (h *Handlers) UserJoinHandler(connection net.Conn, message UserJoinRequest) (userID int) {
	// Add user
	h.users.Add(message.UserID, userRecord{
		connection: connection,
		userID:     message.UserID,
		friends:    message.Friends,
	}, message.Friends)

	log.WithFields(log.Fields{
		"userID": message.UserID,
	}).Info("User connected")

	// Notify online friends
	userJoinedMessage, _ := json.Marshal(userStatusMessage{
		UserID: message.UserID,
		Online: true,
	})
	h.broadcastToFriends(message.UserID, userJoinedMessage)

	return message.UserID
}

func (h *Handlers) UserLeftHandler(connection net.Conn, userID int) {
	// Notify online friends
	userLeftMessage, _ := json.Marshal(userStatusMessage{
		UserID: userID,
		Online: false,
	})
	h.broadcastToFriends(userID, userLeftMessage)

	// Remove user
	h.users.Remove(userID)

	log.WithFields(log.Fields{
		"userID": userID,
	}).Info("User left")
}

func NewHandlers() *Handlers {
	return &Handlers{
		users: NewUserStore(),
	}
}
