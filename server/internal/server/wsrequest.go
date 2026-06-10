package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsRequestEnvelope struct {
	ID        string          `json:"id"`
	Timestamp string          `json:"timestamp"`
	UserID    string          `json:"userID"`
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

func NewWsRequest(userID string, payload interface{}) (error, *WsRequestEnvelope) {
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
		UserID:    userID,
		Type:      payloadType,
		Payload:   rawPayload,
	}
}

func CreateNewWsRequestAndSendIt(websocketConnection *websocket.Conn, userID string, payload interface{}) error {
	err, wsRequestEnvelope := NewWsRequest(userID, payload)
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
