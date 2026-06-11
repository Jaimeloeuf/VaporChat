package server

import (
	"encoding/json"
)

// "WebSocket Request Envelope" is the struct for client to server.
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

type WsRequestPayloadMessageNew struct {
	Message string `json:"message"`
}

type WsRequestPayloadMessageDelete struct {
	MessageID string `json:"messageID"`
}

// These are other WsRequest types with no payload value
// WsRequestPayloadNewStatusJoinRoom
// WsRequestPayloadNewStatusLeaveRoom
// WsRequestPayloadTyping
