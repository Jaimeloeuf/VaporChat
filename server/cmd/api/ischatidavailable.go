package main

// @todo This should be a method on chatRoom after chatRoom have its own independent lock
func (chatStorage *ChatStorage) isChatIDAvailable(chatID string) bool {
	chatStorage.RLock()
	defer chatStorage.RUnlock()

	_, chatRoomExists := chatStorage.chatRooms[chatID]
	return !chatRoomExists
}
