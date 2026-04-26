package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// Configure the Upgrader with buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// @todo Allow connections from any origin for testing
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Create a thread-safe map structure (Go maps are not thread-safe by default)
type ChatStorage struct {
	sync.RWMutex
	// Maps a string UUID to an array of exactly 2 WebSocket connections
	chats map[string][2]*websocket.Conn
}

// @todo Have a timer to clear this regularly
var chatStorage = ChatStorage{
	chats: make(map[string][2]*websocket.Conn),
}

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

func (chatStorage *ChatStorage) setConnectionIfSpaceAvailable(chatID string, newConnection *websocket.Conn) error {
	chatStorage.Lock()
	defer chatStorage.Unlock()

	chatConnections := chatStorage.chats[chatID]

	// Assign connection to first empty slot
	for i, chatConnection := range chatConnections {
		if chatConnection == nil {
			chatConnections[i] = newConnection
			chatStorage.chats[chatID] = chatConnections
			return nil
		}
	}

	return errors.New("Chat room is full")

}

func main() {
	serverMux := http.NewServeMux()

	// "{$}" enforces an exact match for the root path only instead of making
	// this route act as the fallback path
	serverMux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, Gopher!")
	})

	serverMux.HandleFunc("POST /api/chat/new", func(w http.ResponseWriter, r *http.Request) {
		chatID := uuid.New().String()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"chatID": chatID})
	})

	serverMux.HandleFunc("POST /api/chat/join/{chatID}", func(w http.ResponseWriter, r *http.Request) {
		chatID, err := uuid.Parse(r.PathValue("chatID"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid UUID format"})
			return
		}

		fmt.Println("Received:", chatID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "joined"})
	})

	serverMux.HandleFunc("/api/chat/join/{chatID}/websocket", func(w http.ResponseWriter, r *http.Request) {
		chatID := r.PathValue("chatID")
		fmt.Println("Client joining:", chatID)

		if !chatStorage.isChatIDAvailable(chatID) {
			http.Error(w, "Chat ID is taken", http.StatusForbidden)
			return
		}

		// Upgrade HTTP server connection to WebSocket protocol
		websocketConnection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket upgrade error:", err)
			http.Error(w, "Could not upgrade to websocket connection", http.StatusBadRequest)
			return
		}
		defer websocketConnection.Close()
		fmt.Println("Client successfully connected!")

		if err := chatStorage.setConnectionIfSpaceAvailable(chatID, websocketConnection); err != nil {
			websocketConnection.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Chat room full"),
			)
			return
		}

		for {
			// Read message from browser
			messageType, msg, err := websocketConnection.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			// Print incoming message
			fmt.Printf("Received: %s\n", msg)

			// Echo message back to browser
			err = websocketConnection.WriteMessage(messageType, msg)
			if err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST)
	server := cors.Default().Handler(serverMux)

	fmt.Println("Server starting on :3000...")
	if err := http.ListenAndServe(":3000", server); err != nil {
		panic(err)
	}
}
