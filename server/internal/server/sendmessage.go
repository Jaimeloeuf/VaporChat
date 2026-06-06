package server

import (
	"log"

	"github.com/gorilla/websocket"
)

func sendMessageString(websocketConnection *websocket.Conn, message string) error {
	chatMessageAsByteSlice := []byte(message)
	return sendMessageByte(websocketConnection, chatMessageAsByteSlice)
}

func sendMessageByte(websocketConnection *websocket.Conn, message []byte) error {
	err := websocketConnection.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Websocket write message error:", err)
		return err
	}

	return nil
}
