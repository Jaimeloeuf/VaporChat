package main

import (
	"encoding/json"
	"time"

	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

var chatStorage = ChatStorage{
	chats:     make(map[string][2]*websocket.Conn),
	chatRooms: make(map[string]*ChatRoom),
}

func broadcastMessage(chatID string, message string, selfWebsocketConnection *websocket.Conn) {
	chatConnections := chatStorage.chats[chatID]

	chatMessageAsByteSlice := []byte(message)

	// Broadcast message to everyone in chat room
	for _, chatConnection := range chatConnections {
		if chatConnection == nil {
			continue
		}

		// @todo
		// Either we have to do this, or we have to not write it in frontend and
		// wait for it to appear back
		if chatConnection == selfWebsocketConnection {
			continue
		}

		err := chatConnection.WriteMessage(websocket.TextMessage, chatMessageAsByteSlice)
		if err != nil {
			log.Println("Write error:", err)
		}
	}
}

func handleWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("chatID")
	log.Println("Client connecting to:", chatID)

	if !chatStorage.isChatIDAvailable(chatID) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(JSendError("Chat ID is taken"))
		return
	}

	// Upgrade HTTP server connection to WebSocket protocol
	websocketConnection, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JSendError("Could not upgrade to websocket connection"))
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
		var chatMessage ChatMessage
		err := websocketConnection.ReadJSON(&chatMessage)

		// @todo If JSON parsing failed because of field issues, we might not want to break the connection?
		// Exiting loop will hit the defer and clean up websocket connection
		if err != nil {
			log.Printf("Client disconnected or ws read message error: %v", err)
			broadcastMessage(chatID, "Other user has left", websocketConnection)
			break
		}

		// @todo Debug only, leave no logs in server
		// Print incoming message
		log.Printf("Received: %s\n", chatMessage.Message)

		broadcastMessage(chatID, chatMessage.Message, websocketConnection)
	}
}

type ChatRequest struct {
	UserID     string     `json:"userID"`
	ChatConfig ChatConfig `json:"chatConfig"`
}

// Spawns the clean up function in the background with a goroutine
func startBackgroundChatStorageCleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)

	// This goroutine runs forever in the background
	go func() {
		for range ticker.C {
			log.Println("Running background ChatStorage cleanup job...")

			chatStorage.Lock()

			// Loop through your rooms to find empty or expired ones
			for chatID, connections := range chatStorage.chats {

				// Both slots in the array are empty/nil
				if connections[0] == nil && connections[1] == nil {
					log.Printf("Cleaning up empty chat room: %s", chatID)

					// Delete from both maps safely under the lock
					delete(chatStorage.chats, chatID)
					// delete(chatStorage.chatRooms, chatID)
				}
			}

			// Loop through chat rooms to delete expired ones
			for chatID, chatRoom := range chatStorage.chatRooms {
				if uint64(time.Now().Unix()) >= uint64(chatRoom.expiresOn) {
					log.Printf("Cleaning up empty chat room: %s", chatID)
					delete(chatStorage.chatRooms, chatID)
				}
			}

			chatStorage.Unlock()
		}
	}()
}

func main() {
	startBackgroundChatStorageCleanupWorker(1 * time.Second)

	serverMux := http.NewServeMux()

	// "{$}" enforces an exact match for the root path only instead of making
	// this route act as the fallback path
	serverMux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ok")
	})

	serverMux.HandleFunc("POST /api/chat/new", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var requestBody ChatRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JSendError("Invalid JSON format"))
			return
		}

		newChatRoom := NewChatRoom(requestBody.ChatConfig)
		chatStorage.AddNewChatRoom(newChatRoom)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSendSuccess(map[string]string{"chatID": newChatRoom.ID}))
	})

	serverMux.HandleFunc("POST /api/chat/join/{chatID}", func(w http.ResponseWriter, r *http.Request) {
		chatID, err := uuid.Parse(r.PathValue("chatID"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JSendError("invalid UUID format"))
			return
		}

		log.Println("Joining:", chatID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSendSuccess(map[string]string{"status": "joined"}))
	})

	serverMux.HandleFunc("/api/chat/join/{chatID}/websocket", handleWebsocketConnection)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST)
	server := cors.Default().Handler(serverMux)

	serverPort := ":3000"

	log.Printf("Server starting on %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, server); err != nil {
		panic(err)
	}
}
