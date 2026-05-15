package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	ID string

	// Unix seconds
	createdAt int64

	// Unix seconds
	expiresOn int64

	// This is also the number of websocket connections there are currently
	currentNumberOfParticipants uint64

	// Slice to track all current participants (1 per websocket connection)
	websocketConnections []*websocket.Conn

	chatConfig ChatConfig
}

func NewChatRoom(chatConfig ChatConfig) *ChatRoom {
	currentTime := time.Now().Unix()

	return &ChatRoom{
		ID:                          uuid.New().String(),
		createdAt:                   currentTime,
		expiresOn:                   currentTime + int64(chatConfig.ChatRoomTTL),
		currentNumberOfParticipants: 0,
		chatConfig:                  chatConfig,
	}
}
