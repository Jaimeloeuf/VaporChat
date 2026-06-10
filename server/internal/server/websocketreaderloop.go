package server

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsRequestEnvelope struct {
	ID        string          `json:"id"`
	Timestamp string          `json:"timestamp"`
	Username  string          `json:"author"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type WsRequestPayloadRoomCreate struct {
	ChatConfig ChatConfig `json:"chatConfig"`
}

type WsRequestPayloadRoomCreated struct {
	RoomID string `json:"roomID"`
}

type WsRequestPayloadMessageNew struct {
	Message string `json:"message"`
}

type WsRequestPayloadMessageDelete struct {
	MessageID string `json:"messageID"`
}

func NewWsRequest(username string, payload interface{}) (error, *WsRequestEnvelope) {
	var payloadType string

	// Determine the type string dynamically based on the input struct
	switch payload.(type) {
	case WsRequestPayloadRoomCreated, *WsRequestPayloadRoomCreated:
		payloadType = "room-created"
	default:
		return fmt.Errorf("Unknown payload type: %T", payload), nil
	}

	// Dynamic marshaling of the specific payload struct into []byte
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err), nil
	}

	return nil, &WsRequestEnvelope{
		ID:        uuid.New().String(),
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Username:  username,
		Type:      payloadType,
		Payload:   rawPayload,
	}
}

func CreateNewWsRequestAndSendIt(websocketConnection *websocket.Conn, username string, payload interface{}) error {
	err, wsRequestEnvelope := NewWsRequest("system", payload)
	if err != nil {
		return err
	}

	wsRequestEnvelopeBytes, err := json.Marshal(wsRequestEnvelope)
	if err != nil {
		return err
	}

	err = sendMessageByte(websocketConnection, wsRequestEnvelopeBytes)
	if err != nil {
		return err
	}

	return nil
}

// These are other WsRequest types with no payload value
// WsRequestPayloadNewStatusJoinRoom
// WsRequestPayloadNewStatusLeaveRoom
// WsRequestPayloadTyping

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

			err = CreateNewWsRequestAndSendIt(websocketConnection, "system", WsRequestPayloadRoomCreated{
				RoomID: newChatRoom.ID,
			})

		case "room-join":
			log.Printf("[Status] User %s joined the room at %s\n", wsRequestEnvelope.Username, wsRequestEnvelope.Timestamp)

		case "room-leave":
			log.Printf("[Status] User %s left the room at %s\n", wsRequestEnvelope.Username, wsRequestEnvelope.Timestamp)

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
