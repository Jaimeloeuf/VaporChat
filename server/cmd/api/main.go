package main

import (
	"time"
)

var chatStorage = ChatStorage{
	chatRooms: make(map[string]*ChatRoom),
}

type ChatRequest struct {
	UserID     string     `json:"userID"`
	ChatConfig ChatConfig `json:"chatConfig"`
}

func main() {
	startBackgroundChatStorageCleanupWorker(1 * time.Second)
	startWebServer()
}
