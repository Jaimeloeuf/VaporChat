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

func startWebServer() {
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

	serverMux.HandleFunc("/api/chat/join/{chatID}/websocket", func(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Client connected")

		if err := chatStorage.saveConnection(chatID, websocketConnection); err != nil {
			log.Printf("Connection refused for chat room ID %s: %v", chatID, err)
			websocketConnection.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Connection refused: chat room not available"),
			)
			return
		}

		go websocketLoop(websocketConnection, chatID)
	})

	serverMux.HandleFunc("/api/websocket", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Client connecting via websocket")

		// Upgrade HTTP server connection to WebSocket protocol
		websocketConnection, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket upgrade error:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JSendError("Could not upgrade to websocket connection"))
			return
		}
		log.Println("Client connected via websocket")

		go websocketReaderLoop(websocketConnection)
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST)
	server := cors.Default().Handler(serverMux)

	serverPort := ":3000"

	log.Printf("Server starting on %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, server); err != nil {
		panic(err)
	}
}
