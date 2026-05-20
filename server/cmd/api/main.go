package main

import (
	"time"

	"github.com/Jaimeloeuf/VaporChat/internal/server"
)

func main() {
	server.StartBackgroundChatStorageCleanupWorker(1 * time.Second)
	server.StartWebServer()
}
