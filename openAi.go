package main

import (
	"github.com/c00/botman/models"
	openai "github.com/sashabaranov/go-openai"
)

func convertMessage(m models.ChatMessage) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    m.Role,
		Content: m.Content,
	}
}

func MessagesToOpenAiMessages(m []models.ChatMessage) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, 0, len(m))
	for _, message := range m {
		result = append(result, convertMessage(message))
	}
	return result
}
