package server

type ChatConfig struct {
	ChatRoomTTL                 uint32 `json:"chatRoomTTL"`
	MaxNumberOfParticipants     uint32 `json:"maxNumberOfParticipants"`
	MaxHistoryDurationInSeconds uint32 `json:"maxHistoryDurationInSeconds"`
	MaxMessagesLength           uint32 `json:"maxMessagesLength"`
}
