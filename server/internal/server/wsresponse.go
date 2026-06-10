package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsResponseEnvelope struct {
	ID        string          `json:"id"`
	Timestamp string          `json:"timestamp"`
	UserID    string          `json:"userID"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type WsResponsePayloadRoomCreated struct {
	RoomID string `json:"roomID"`
}

// These are other WsResponse types with no payload value
// WsResponsePayloadNewStatusJoinRoom
// WsResponsePayloadNewStatusLeaveRoom
// WsResponsePayloadTyping

func NewWsResponse(userID string, payload interface{}) (error, *WsResponseEnvelope) {
	var payloadType string

	// Determine the type string dynamically based on the input struct
	switch payload.(type) {
	case WsResponsePayloadRoomCreated, *WsResponsePayloadRoomCreated:
		payloadType = "room-created"
	default:
		return fmt.Errorf("Unknown payload type: %T", payload), nil
	}

	// Dynamic marshaling of the specific payload struct into []byte
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err), nil
	}

	return nil, &WsResponseEnvelope{
		ID:        uuid.New().String(),
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		UserID:    userID,
		Type:      payloadType,
		Payload:   rawPayload,
	}
}

func CreateNewWsResponseAndSendIt(websocketConnection *websocket.Conn, userID string, payload interface{}) error {
	err, wsResponseEnvelope := NewWsResponse(userID, payload)
	if err != nil {
		return err
	}

	wsResponseEnvelopeBytes, err := json.Marshal(wsResponseEnvelope)
	if err != nil {
		return err
	}

	err = sendMessageByte(websocketConnection, wsResponseEnvelopeBytes)
	if err != nil {
		return err
	}

	return nil
}
