package main

type ChatConfig struct {
	MaxNumberOfParticipants     uint64 `json:"maxNumberOfParticipants"`
	MaxHistoryDurationInSeconds uint64 `json:"maxHistoryDurationInSeconds"`
	MaxMessagesLength           uint64 `json:"maxMessagesLength"`
}
