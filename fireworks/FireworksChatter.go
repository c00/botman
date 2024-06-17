package fireworks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
)

const apiUrl = "https://api.fireworks.ai/inference/v1/chat/completions"

func NewChatBot(cfg config.FireworksConfig) Chatbot {
	return Chatbot{
		cfg: cfg,
	}
}

type Chatbot struct {
	cfg config.FireworksConfig
}

type fireworksPostBody struct {
	Model            string             `json:"model"`
	Messages         []fireworksMessage `json:"messages"`
	MaxTokens        string             `json:"max_tokens,omitempty"`
	TopP             int                `json:"top_p,omitempty"`
	TopK             int                `json:"top_k,omitempty"`
	PresencePenalty  int                `json:"presence_penalty,omitempty"`
	FrequencyPenalty int                `json:"frequency_penalty,omitempty"`
	Temperature      float32            `json:"temperature,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	N                int                `json:"n,omitempty"`
}

type fireworksMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c Chatbot) GetResponse(messages []models.ChatMessage, streamChan chan<- string) string {
	defer close(streamChan)

	//Set system prompt if there's only 2 messages, and the first one is a system prompt
	if len(messages) == 2 && messages[0].Role == models.ChatMessageRoleSystem && c.cfg.SystemPrompt != "" {
		messages[0].Content += " " + c.cfg.SystemPrompt
	}

	body := fireworksPostBody{
		Model:    c.cfg.Model,
		Messages: messagesToFireworksMessages(messages),
		Stream:   true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Sprintf("cannot marshall JSON body: %v", err))
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.cfg.ApiKey))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	if resp.StatusCode != 200 {
		bodyContent, err := io.ReadAll(reader)
		if err != nil {
			panic(fmt.Sprintf("got status: %v, error: %v", resp.StatusCode, err))
		}
		panic(fmt.Sprintf("got status: %v, %v", resp.StatusCode, string(bodyContent)))
	}

	responseContent := make([]string, 0, 50)

	//Read the streaming response.
	//I'm assuming each message will have a newline at the end.
	for {
		line, err := reader.ReadBytes(byte('\n'))
		chunk := parseChunk(line)

		if errors.Is(err, io.EOF) {
			if !chunk.Empty {
				streamChan <- chunk.Delta
				responseContent = append(responseContent, chunk.Delta)
			}

			responseMessage := strings.Join(responseContent, "")
			return responseMessage
		}

		if chunk.Empty {
			continue
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting FireworksAI Chat Completion:", err)
			panic(fmt.Sprintf("error getting FireworksAI Chat Completion: %v", err))
		}

		if chunk.LastMessage {
			continue
		}

		streamChan <- chunk.Delta
		responseContent = append(responseContent, chunk.Delta)
	}
}

func convertMessage(m models.ChatMessage) fireworksMessage {
	return fireworksMessage{
		Role:    m.Role,
		Content: m.Content,
	}
}

func messagesToFireworksMessages(m []models.ChatMessage) []fireworksMessage {
	result := make([]fireworksMessage, 0, len(m))
	for _, message := range m {
		result = append(result, convertMessage(message))
	}
	return result
}
