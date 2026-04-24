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

func (chatStorage *ChatStorage) setConnectionIfSpaceAvailable(chatID string, newConnection *websocket.Conn) error {
	chatStorage.Lock()
	defer chatStorage.Unlock()

	chatConnections := chatStorage.chats[chatID]

	conn0 := chatConnections[0]
	conn1 := chatConnections[1]

	if conn0 != nil && conn1 != nil {
		return errors.New("Chat room is full")
	}

	if conn0 == nil {
		chatConnections[0] = newConnection
	} else if conn1 == nil {
		chatConnections[1] = newConnection
	}

	chatStorage.chats[chatID] = chatConnections

	return nil
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

		// Upgrade the HTTP server connection to the WebSocket protocol
		websocketConnection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Upgrade error:", err)
			return
		}
		defer websocketConnection.Close()
		fmt.Println("Client successfully connected!")

		if err := chatStorage.setConnectionIfSpaceAvailable(chatID, websocketConnection); err != nil {
			websocketConnection.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "Chat room full"),
			)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
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
