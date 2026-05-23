package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Jaimeloeuf/VaporChat/internal/jsend"
	"github.com/rs/cors"
)

func StartWebServer() {
	serverMux := http.NewServeMux()

	// "{$}" enforces an exact match for the root path only instead of making
	// this route act as the fallback path
	serverMux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ok")
	})

	serverMux.HandleFunc("/api/websocket", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Client connecting via websocket")

		// Upgrade HTTP server connection to WebSocket protocol
		websocketConnection, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket upgrade error:", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(jsend.Error("Could not upgrade to websocket connection"))
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
