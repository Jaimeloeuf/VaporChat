package server

import (
	"errors"

	"github.com/gorilla/websocket"
)

// @todo This should be a method on chatRoom after chatRoom have its own independent lock
func (chatStorage *ChatStorage) saveConnection(chatID string, newConnection *websocket.Conn) error {
	chatStorage.Lock()
	defer chatStorage.Unlock()

	chatRoom, chatRoomExists := chatStorage.chatRooms[chatID]
	if !chatRoomExists {
		return errors.New("Chat room not available")
	}

	if len(chatRoom.websocketConnections) >= int(chatRoom.chatConfig.MaxNumberOfParticipants) {
		return errors.New("Chat room not available")
	}

	chatRoom.websocketConnections = append(chatRoom.websocketConnections, newConnection)

	return nil
}
