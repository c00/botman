package claude

import (
	"os"
	"testing"

	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
)

func TestGetResponse(t *testing.T) {
	chatter := NewChatBot(config.ClaudeConfig{
		ApiKey:    os.Getenv("CLAUDE_API_KEY"),
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 50,
	})

	ch := make(chan string)
	chunks := []string{}

	go func(ch chan string) {
		for content := range ch {
			chunks = append(chunks, content)
		}
	}(ch)

	messages := []models.ChatMessage{
		{
			Role:    models.ChatMessageRoleSystem,
			Content: "You give short and concise answers",
		},
		{
			Role:    models.ChatMessageRoleUser,
			Content: "A haiku about the word Hello",
		},
	}

	response := chatter.GetResponse(messages, ch)
	if len(chunks) < 3 {
		t.Fatal("Not chunky enough")
	}

	if response == "" {
		t.Fatal("Response is too small")
	}
}
