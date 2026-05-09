package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Create a thread-safe map structure (Go maps are not thread-safe by default)
type ChatStorage struct {
	sync.RWMutex

	// Maps a string UUID to an array of exactly 2 WebSocket connections
	chats map[string][2]*websocket.Conn

	chatRooms map[string]ChatRoom
}
