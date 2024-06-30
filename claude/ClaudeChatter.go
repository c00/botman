package claude

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

const apiUrl = "https://api.anthropic.com/v1/messages"

func NewChatBot(cfg config.ClaudeConfig) Chatbot {
	return Chatbot{
		cfg: cfg,
	}
}

type Chatbot struct {
	cfg config.ClaudeConfig
}

type claudePostBody struct {
	Model     string          `json:"model"`
	Messages  []claudeMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
	Stream    bool            `json:"stream,omitempty"`
	System    string          `json:"system"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c Chatbot) GetResponse(messages []models.ChatMessage, streamChan chan<- string) string {
	defer close(streamChan)

	//Set system prompt if there's only 2 messages, and the first one is a system prompt
	if len(messages) == 2 && messages[0].Role == models.ChatMessageRoleSystem && c.cfg.SystemPrompt != "" {
		messages[0].Content += " " + c.cfg.SystemPrompt
	}

	systemMessage, claudeMessages := messagesToClaudeMessages(messages)

	body := claudePostBody{
		Model:     c.cfg.Model,
		Messages:  claudeMessages,
		System:    systemMessage,
		Stream:    true,
		MaxTokens: c.cfg.MaxTokens,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Sprintf("cannot marshall JSON body: %v", err))
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01") //https://docs.anthropic.com/en/api/versioning
	req.Header.Set("x-api-key", c.cfg.ApiKey)
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
	event := ""
	dataPrefix := []byte("data: ")
	eventPrefix := []byte("event: ")
	for {
		line, err := reader.ReadBytes(byte('\n'))
		if errors.Is(err, io.EOF) {
			responseMessage := strings.Join(responseContent, "")
			return responseMessage
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting Claude Chat Completion:", err)
			panic(fmt.Sprintf("error getting Claude Chat Completion: %v", err))
		}

		if len(line) < 6 {
			//Ignore it
			continue
		} else if bytes.Equal(line[:7], eventPrefix) {
			event = string(line[7 : len(line)-1])
			continue
		} else if bytes.Equal(line[:6], dataPrefix) {
			chunk := parseChunk(line[6:], event)
			if chunk.LastMessage || chunk.Empty {
				continue
			}

			streamChan <- chunk.Delta
			responseContent = append(responseContent, chunk.Delta)
		} else {
			panic("Unexpected line chunk: " + string(line))
		}

	}
}

func convertMessage(m models.ChatMessage) claudeMessage {
	return claudeMessage{
		Role:    m.Role,
		Content: m.Content,
	}
}

func messagesToClaudeMessages(m []models.ChatMessage) (string, []claudeMessage) {
	systemMessage := ""
	result := make([]claudeMessage, 0, len(m))
	for _, message := range m {
		if message.Role == "system" {
			systemMessage = message.Content
		} else {
			result = append(result, convertMessage(message))
		}
	}
	return systemMessage, result
}
