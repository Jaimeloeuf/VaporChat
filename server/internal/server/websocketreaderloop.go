package server

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type ChatUpdateEnvelope struct {
	Timestamp string          `json:"timestamp"`
	Username  string          `json:"author"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type ChatUpdatePayloadRoomCreate struct {
	ChatConfig ChatConfig `json:"chatConfig"`
}

type ChatUpdatePayloadMessageNew struct {
	Message string `json:"message"`
}

type ChatUpdatePayloadMessageDelete struct {
	MessageID string `json:"messageID"`
}

// These are other ChatUpdate types with no payload value
// ChatUpdatePayloadNewStatusJoinRoom
// ChatUpdatePayloadNewStatusLeaveRoom
// ChatUpdatePayloadTyping

func websocketReaderLoop(websocketConnection *websocket.Conn) {
	// Ensure connection is cleaned up from memory and closed when loop ends
	defer func() {
		websocketConnection.Close()

		// @todo Delete connection from mapping
		// websocketConnections.mu.Lock()
		// delete(websocketConnections, websocketConnection)
		// websocketConnections.mu.Unlock()

		log.Printf("Closed connection ...$connID...\n")
	}()

	for {
		// Read raw message bytes from WebSocket connection
		// We use ReadMessage instead of ReadJSON to preserve raw bytes for
		// `json.RawMessage` parsing later
		messageType, rawBytes, err := websocketConnection.ReadMessage()

		log.Printf("[ChatUpdate] incoming")

		// Exiting loop will hit the defer and clean up websocket connection
		if err != nil {
			log.Printf("Client disconnected or ws read message error: %v", err)
			// @todo
			// broadcastChatUpdateToChatRoom(chatID, "Other user has left", websocketConnection)
			break
		}

		// Ensure it's a TextMessage frame (skip binary frames if they happen)
		// @todo handle the other types
		if messageType != websocket.TextMessage {
			continue
		}

		// Unmarshal JSON into top level ChatUpdateEnvelope struct
		var chatUpdateEnvelope ChatUpdateEnvelope
		if err := json.Unmarshal(rawBytes, &chatUpdateEnvelope); err != nil {
			// @todo If JSON parsing failed because of field issues, we might not want to break the connection?
			// CRITICAL FIX: If JSON is malformed, log it and 'continue' instead of 'break'
			// This prevents bad client inputs from crashing the entire user session.
			log.Printf("Failed to unmarshal outer envelope: %v. Raw data: %s", err, string(rawBytes))
			continue
		}

		// @todo Debug only, leave no logs in server
		// Print incoming message
		log.Printf("Received raw data: %s", string(rawBytes))

		// Switch based on the polymorphic 'Type' field
		switch chatUpdateEnvelope.Type {

		case "room-create":
			var chatUpdatePayload ChatUpdatePayloadRoomCreate
			if err := json.Unmarshal(chatUpdateEnvelope.Payload, &chatUpdatePayload); err != nil {
				log.Printf("Malformed %s payload: %v", chatUpdateEnvelope.Type, err)
				continue
			}

			newChatRoom := NewChatRoom(chatUpdatePayload.ChatConfig)
			chatStorage.AddNewChatRoom(newChatRoom)

		case "message-new":
			var chatUpdatePayload ChatUpdatePayloadMessageNew
			if err := json.Unmarshal(chatUpdateEnvelope.Payload, &chatUpdatePayload); err != nil {
				log.Printf("Malformed %s payload: %v", chatUpdateEnvelope.Type, err)
				continue
			}

		case "message-delete":
			var chatUpdatePayload ChatUpdatePayloadMessageDelete
			if err := json.Unmarshal(chatUpdateEnvelope.Payload, &chatUpdatePayload); err != nil {
				log.Printf("Malformed %s payload: %v", chatUpdateEnvelope.Type, err)
				continue
			}

		case "join_room":

		case "leave_room":

		case "typing":

		default:
			log.Printf("Received unhandled or unknown ChatUpdate type: %s\n", chatUpdateEnvelope.Type)
		}

	}
}
