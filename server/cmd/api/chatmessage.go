package main

type ChatMessage struct {
	Timestamp string `json:"timestamp"`
	Username  string `json:"author"`
	Message   string `json:"message"`
}
