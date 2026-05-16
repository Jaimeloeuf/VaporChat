package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	ID string

	// Unix seconds
	createdAt uint64

	// Unix seconds
	expiresOn uint64

	// Slice to track all current participants (1 per websocket connection)
	websocketConnections []*websocket.Conn

	chatConfig ChatConfig
}

func NewChatRoom(chatConfig ChatConfig) *ChatRoom {
	currentTime := time.Now().Unix()

	return &ChatRoom{
		ID:         uuid.New().String(),
		createdAt:  uint64(currentTime),
		expiresOn:  uint64(currentTime + int64(chatConfig.ChatRoomTTL)),
		chatConfig: chatConfig,
	}
}
