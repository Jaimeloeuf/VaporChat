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

	// This is also the number of websocket connections there are currently
	currentNumberOfParticipants uint64

	// Slice to track all current participants (1 per websocket connection)
	websocketConnections []*websocket.Conn

	chatConfig ChatConfig
}

func NewChatRoom(chatConfig ChatConfig) *ChatRoom {
	return &ChatRoom{
		ID:                          uuid.New().String(),
		createdAt:                   time.Now().Unix(),
		currentNumberOfParticipants: 0,
		chatConfig:                  chatConfig,
	}
}
