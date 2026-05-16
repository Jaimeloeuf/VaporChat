package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func websocketLoop(websocketConnection *websocket.Conn, chatID string) {
	defer websocketConnection.Close()

	for {
		var chatMessage ChatMessage
		err := websocketConnection.ReadJSON(&chatMessage)

		// @todo If JSON parsing failed because of field issues, we might not want to break the connection?
		// Exiting loop will hit the defer and clean up websocket connection
		if err != nil {
			log.Printf("Client disconnected or ws read message error: %v", err)
			broadcastMessageToChatRoom(chatID, "Other user has left", websocketConnection)
			break
		}

		// @todo Debug only, leave no logs in server
		// Print incoming message
		log.Printf("Received: %s\n", chatMessage.Message)

		broadcastMessageToChatRoom(chatID, chatMessage.Message, websocketConnection)
	}
}
