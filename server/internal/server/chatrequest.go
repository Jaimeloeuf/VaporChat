package server

type ChatRequest struct {
	UserID     string     `json:"userID"`
	ChatConfig ChatConfig `json:"chatConfig"`
}
