package main

type ChatConfig struct {
	ChatRoomTTL                 uint64 `json:"chatRoomTTL"`
	MaxNumberOfParticipants     uint64 `json:"maxNumberOfParticipants"`
	MaxHistoryDurationInSeconds uint64 `json:"maxHistoryDurationInSeconds"`
	MaxMessagesLength           uint64 `json:"maxMessagesLength"`
}
