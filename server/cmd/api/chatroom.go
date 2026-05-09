package main

import "github.com/gorilla/websocket"

type ChatRoom struct {
	ID string

	createdAt  int64
	chatConfig ChatConfig

	// This is also the number of websocket connections there are currently
	currentNumberOfParticipants uint64

	// Slice to track all current participants (1 per websocket connection)
	websocketConnections []*websocket.Conn
}
