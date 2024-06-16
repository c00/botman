package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
	openai "github.com/sashabaranov/go-openai"
)

func NewChatBot(cfg config.OpenAiConfig) Chatbot {
	return Chatbot{
		client: openai.NewClient(cfg.ApiKey),
		cfg:    cfg,
	}
}

type Chatbot struct {
	client *openai.Client
	cfg    config.OpenAiConfig
}

func (c Chatbot) GetResponse(messages []models.ChatMessage, streamChan chan<- string) string {
	defer close(streamChan)

	//Set system prompt if there's only 2 messages, and the first one is a system prompt
	if len(messages) == 2 && messages[0].Role == models.ChatMessageRoleSystem && c.cfg.SystemPrompt != "" {
		messages[0].Content += " " + c.cfg.SystemPrompt
	}

	stream, err := c.client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: messagesToOpenAiMessages(messages),
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting OpenAi Chat Completion:", err)
		panic(fmt.Sprintf("error getting OpenAi Chat Completion: %v", err))
	}
	defer stream.Close()

	responseContent := make([]string, 0, 50)

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			responseMessage := strings.Join(responseContent, "")
			return responseMessage
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Stream error:", err)
			panic(fmt.Sprintf("stream error: %v", err))
		}

		streamChan <- response.Choices[0].Delta.Content

		responseContent = append(responseContent, response.Choices[0].Delta.Content)
	}
}

func convertMessage(m models.ChatMessage) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    m.Role,
		Content: m.Content,
	}
}

func messagesToOpenAiMessages(m []models.ChatMessage) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, 0, len(m))
	for _, message := range m {
		result = append(result, convertMessage(message))
	}
	return result
}
