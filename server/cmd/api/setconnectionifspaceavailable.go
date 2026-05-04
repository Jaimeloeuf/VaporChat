package main

import (
	"errors"

	"github.com/gorilla/websocket"
)

func (chatStorage *ChatStorage) setConnectionIfSpaceAvailable(chatID string, newConnection *websocket.Conn) error {
	chatStorage.Lock()
	defer chatStorage.Unlock()

	chatConnections := chatStorage.chats[chatID]

	// Assign connection to first empty slot
	for i, chatConnection := range chatConnections {
		if chatConnection == nil {
			chatConnections[i] = newConnection
			chatStorage.chats[chatID] = chatConnections
			return nil
		}
	}

	return errors.New("Chat room is full")

}
