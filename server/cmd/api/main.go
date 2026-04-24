package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	serverMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP server connection to the WebSocket protocol
		websocketConnection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Upgrade error:", err)
			return
		}
		defer websocketConnection.Close()
		fmt.Println("Client successfully connected!")

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
