package server

import (
	"log"
	"time"
)

// Spawns the clean up function in the background with a goroutine
func StartBackgroundChatStorageCleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)

	// This goroutine runs forever in the background
	go func() {
		for range ticker.C {
			log.Println("Running background ChatStorage cleanup job...")

			chatStorage.Lock()

			// Loop through chat rooms to delete expired ones
			for chatID, chatRoom := range chatStorage.chatRooms {
				if uint64(time.Now().Unix()) >= chatRoom.expiresOn {
					log.Printf("Cleaning up empty chat room: %s", chatID)
					delete(chatStorage.chatRooms, chatID)
				}
			}

			chatStorage.Unlock()
		}
	}()
}
