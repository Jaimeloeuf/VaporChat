package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

func handler(w http.ResponseWriter, r *http.Request) {
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
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, Gopher!")
	})

	http.HandleFunc("/api/chat/new", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			chatID := uuid.New().String()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"chatID": chatID})

		default:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	})

	http.HandleFunc("/ws", handler)

	fmt.Println("Server starting on :3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
