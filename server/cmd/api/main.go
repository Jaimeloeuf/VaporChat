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

// @todo Have a timer to clear this regularly
var chatStorage = ChatStorage{
	chats: make(map[string][2]*websocket.Conn),
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

		log.Println("Joining:", chatID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "joined"})
	})

	serverMux.HandleFunc("/api/chat/join/{chatID}/websocket", func(w http.ResponseWriter, r *http.Request) {
		chatID := r.PathValue("chatID")
		log.Println("Client connecting to:", chatID)

		if !chatStorage.isChatIDAvailable(chatID) {
			http.Error(w, "Chat ID is taken", http.StatusForbidden)
			return
		}

		// Upgrade HTTP server connection to WebSocket protocol
		websocketConnection, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket upgrade error:", err)
			http.Error(w, "Could not upgrade to websocket connection", http.StatusBadRequest)
			return
		}
		defer websocketConnection.Close()
		log.Println("Client connected")

		if err := chatStorage.setConnectionIfSpaceAvailable(chatID, websocketConnection); err != nil {
			websocketConnection.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Chat room full"),
			)
			return
		}

		for {
			_, msg, err := websocketConnection.ReadMessage()

			// Exiting loop will hit the defer and clean up websocket connection
			if err != nil {
				log.Printf("Client disconnected or ws read message error: %v", err)
				break
			}

			// @todo Debug only, leave no logs in server
			// Print incoming message
			log.Printf("Received: %s\n", msg)

			chatConnections := chatStorage.chats[chatID]

			// @todo Do nothing until other party joined the chat
			// Broadcast message to everyone in chat room
			for _, chatConnection := range chatConnections {
				if chatConnection == nil {
					continue
				}

				err = chatConnection.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Write error:", err)
				}
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
