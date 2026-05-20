package server

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

func broadcastMessageToChatRoom(chatID string, message string, selfWebsocketConnection *websocket.Conn) error {
	chatStorage.RLock()
	defer chatStorage.RUnlock()

	chatRoom, chatRoomExists := chatStorage.chatRooms[chatID]
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

		err := chatConnection.WriteMessage(websocket.TextMessage, chatMessageAsByteSlice)
		if err != nil {
			log.Println("Write error:", err)
		}
	}

	return nil
}
