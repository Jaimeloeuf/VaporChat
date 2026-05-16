package main

import "sync"

// Create a thread-safe map structure (Go maps are not thread-safe by default)
type ChatStorage struct {
	sync.RWMutex

	chatRooms map[string]*ChatRoom
}

func (chatStorage *ChatStorage) AddNewChatRoom(newChatRoom *ChatRoom) {
	chatStorage.Lock()
	defer chatStorage.Unlock()
	chatStorage.chatRooms[newChatRoom.ID] = newChatRoom
}
