package main

func (chatStorage *ChatStorage) isChatIDAvailable(chatID string) bool {
	chatStorage.Lock()
	defer chatStorage.Unlock()

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
