package fireworks

import (
	"os"
	"testing"

	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
)

func TestGetResponse(t *testing.T) {
	chatter := NewChatBot(config.FireworksConfig{
		ApiKey: os.Getenv("FIREWORKS_API_KEY"),
		Model:  "accounts/fireworks/models/mixtral-8x7b-instruct",
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
