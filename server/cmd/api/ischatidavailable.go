package main

func (chatStorage *ChatStorage) isChatIDAvailable(chatID string) bool {
	chatStorage.RLock()
	defer chatStorage.RUnlock()

	chatConnections, chatExists := chatStorage.chats[chatID]

	if !chatExists {
		return true
	}

	// Loop to find a single nil
	for _, chatConnection := range chatConnections {
		if chatConnection == nil {
			return true
		}
	}

	return false
}
