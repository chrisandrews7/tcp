package main

import (
	"encoding/json"
	"net"
	"os"
	"reflect"
	"testing"
)

const TestAddress = ":3030"

func TestMain(m *testing.M) {
	server := NewTCPServer(TestAddress)
	defer server.Close()

	go func() {
		server.Run()
	}()

	os.Exit(m.Run())
}

func TestConnection(t *testing.T) {
	connection, err := net.Dial("tcp", TestAddress)

	if err != nil {
		t.Error("Unable to connect", err)
	}

	connection.Close()
}

func TestFriendOnlineNotification(t *testing.T) {
	userA, _ := net.Dial("tcp", TestAddress)
	userB, _ := net.Dial("tcp", TestAddress)
	defer userA.Close()
	defer userB.Close()

	userAJoinMessage, _ := json.Marshal(&UserJoinRequest{
		UserID:  1,
		Friends: []int{2, 3},
	})
	if _, err := userA.Write(userAJoinMessage); err != nil {
		t.Error("Could not send message", err)
	}

	userBJoinMessage, _ := json.Marshal(&UserJoinRequest{
		UserID:  2,
		Friends: []int{1, 3},
	})
	if _, err := userB.Write(userBJoinMessage); err != nil {
		t.Error("Could not send message", err)
	}

	decoder := json.NewDecoder(userA)
	var message userStatusMessage
	if err := decoder.Decode(&message); err != nil {
		t.Error(err)
	}

	expected := userStatusMessage{
		Online: true,
		UserID: 2,
	}
	if !reflect.DeepEqual(message, expected) {
		t.Errorf("Expected user online message, got '%v', expected '%v'", message, expected)
	}
}

func TestFriendOfflineNotification(t *testing.T) {
	userA, _ := net.Dial("tcp", TestAddress)
	userB, _ := net.Dial("tcp", TestAddress)
	defer userA.Close()

	userAJoinMessage, _ := json.Marshal(&UserJoinRequest{
		UserID:  3,
		Friends: []int{4},
	})
	if _, err := userA.Write(userAJoinMessage); err != nil {
		t.Error("Could not send message", err)
	}

	userBJoinMessage, _ := json.Marshal(&UserJoinRequest{
		UserID:  4,
		Friends: []int{3},
	})
	if _, err := userB.Write(userBJoinMessage); err != nil {
		t.Error("Could not send message", err)
	}

	// Exit UserB
	userB.Close()

	decoder := json.NewDecoder(userA)
	var message userStatusMessage
	if err := decoder.Decode(&message); err != nil {
		t.Error(err)
	}
	if err := decoder.Decode(&message); err != nil {
		t.Error(err)
	}

	expected := userStatusMessage{
		Online: false,
		UserID: 4,
	}
	if !reflect.DeepEqual(message, expected) {
		t.Errorf("Expected user offline message, got '%v', expected '%v'", message, expected)
	}
}
