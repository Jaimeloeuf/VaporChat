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

// @todo This should be a method on chatRoom after chatRoom have its own independent lock
func (chatStorage *ChatStorage) saveConnection(chatID string, newConnection *websocket.Conn) error {
	chatStorage.Lock()
	defer chatStorage.Unlock()

	chatRoom, chatRoomExists := chatStorage.chatRooms[chatID]
	if !chatRoomExists {
		return errors.New("Chat room not available")
	}

	if chatRoom.currentNumberOfParticipants >= chatRoom.chatConfig.MaxNumberOfParticipants {
		return errors.New("Chat room not available")
	}

	chatRoom.currentNumberOfParticipants++
	chatRoom.websocketConnections = append(chatRoom.websocketConnections, newConnection)
	return nil
}
