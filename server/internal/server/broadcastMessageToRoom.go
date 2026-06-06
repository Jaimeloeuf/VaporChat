package server

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

// @todo This should be a method on the chatRoom itself, and lock the chatRoom itself
func broadcastMessageToRoom(roomID string, message string, selfWebsocketConnection *websocket.Conn) error {
	chatStorage.RLock()
	defer chatStorage.RUnlock()

	chatRoom, chatRoomExists := chatStorage.chatRooms[roomID]
	if !chatRoomExists {
		return errors.New("Chat room not available")
	}

	chatMessageAsByteSlice := []byte(message)

	// Broadcast message to everyone in chat room
	for _, chatConnection := range chatRoom.websocketConnections {
		// @todo
		// Either we have to do this, or we have to not write it in frontend and
		// wait for it to appear back
		if chatConnection == selfWebsocketConnection {
			continue
		}

		err := sendMessageByte(chatConnection, chatMessageAsByteSlice)
		if err != nil {
			log.Println("Write error:", err)
		}
	}

	return nil
}
