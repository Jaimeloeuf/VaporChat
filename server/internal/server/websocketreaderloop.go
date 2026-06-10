package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

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

		log.Printf("[WsRequest] incoming")

		// Exiting loop will hit the defer and clean up websocket connection
		if err != nil {
			log.Printf("Client disconnected or ws read message error: %v", err)
			// @todo
			// broadcastMessageToRoom("chatID", "Other user has left", websocketConnection)
			break
		}

		// Ensure it's a TextMessage frame (skip binary frames if they happen)
		// @todo handle the other types
		if messageType != websocket.TextMessage {
			continue
		}

		// Unmarshal JSON into top level struct
		var wsRequestEnvelope WsRequestEnvelope
		if err := json.Unmarshal(rawBytes, &wsRequestEnvelope); err != nil {
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
		switch wsRequestEnvelope.Type {

		case "room-create":
			var wsRequestPayload WsRequestPayloadRoomCreate
			if err := json.Unmarshal(wsRequestEnvelope.Payload, &wsRequestPayload); err != nil {
				log.Printf("Malformed %s payload: %v", wsRequestEnvelope.Type, err)
				continue
			}

			newChatRoom := NewChatRoom(wsRequestPayload.ChatConfig)
			chatStorage.AddNewChatRoom(newChatRoom)

			log.Printf("[Status] Created room at %s\n", newChatRoom.ID)

			err = CreateNewWsResponseAndSendIt(websocketConnection, "system", WsRequestPayloadRoomCreated{
				RoomID: newChatRoom.ID,
			})

		case "room-join":
			log.Printf("[Status] User %s joined the room at %s\n", wsRequestEnvelope.UserID, wsRequestEnvelope.Timestamp)

		case "room-leave":
			log.Printf("[Status] User %s left the room at %s\n", wsRequestEnvelope.UserID, wsRequestEnvelope.Timestamp)

		case "message-new":
			var wsRequestPayload WsRequestPayloadMessageNew
			if err := json.Unmarshal(wsRequestEnvelope.Payload, &wsRequestPayload); err != nil {
				log.Printf("Malformed %s payload: %v", wsRequestEnvelope.Type, err)
				continue
			}

		case "message-delete":
			var wsRequestPayload WsRequestPayloadMessageDelete
			if err := json.Unmarshal(wsRequestEnvelope.Payload, &wsRequestPayload); err != nil {
				log.Printf("Malformed %s payload: %v", wsRequestEnvelope.Type, err)
				continue
			}

		case "typing":

		default:
			log.Printf("Received unhandled or unknown WsRequest type: %s\n", wsRequestEnvelope.Type)
		}

	}
}
