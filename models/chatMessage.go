package models

type ChatMessage struct {
	Role    string
	Content string
}

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
	ChatMessageRoleFunction  = "function"
	ChatMessageRoleTool      = "tool"
)

type Chatter interface {
	GetResponse(messages []ChatMessage, streamChan chan<- string) string
}
