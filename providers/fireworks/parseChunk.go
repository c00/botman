package fireworks

import (
	"bytes"
	"encoding/json"
)

type parsedChunk struct {
	Empty       bool
	LastMessage bool
	Delta       string
}

type chatCompletionChunk struct {
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
}

type choice struct {
	Index        int     `json:"index"`
	Delta        delta   `json:"delta"`
	FinishReason *string `json:"finish_reason,omitempty"`
}

type delta struct {
	Content string `json:"content"`
}

func parseChunk(data []byte) parsedChunk {
	if len(data) == 0 {
		return parsedChunk{Empty: true}
	}
	newLine := []byte("\n")
	doneLine := []byte("data: [DONE]\n")
	dataPrefix := []byte("data: ")

	if bytes.Equal(data, newLine) {
		return parsedChunk{Empty: true}
	} else if bytes.Equal(data, doneLine) {
		return parsedChunk{LastMessage: true}
	} else if len(data) < 6 || !bytes.Equal(data[:6], dataPrefix) {
		panic("weird data chunk: " + string(data))
	}

	//Parse JSON
	chunk := chatCompletionChunk{}
	err := json.Unmarshal(data[5:], &chunk)
	if err != nil {
		panic(err)
	}

	if len(chunk.Choices) == 0 {
		return parsedChunk{Empty: true}
	}

	return parsedChunk{
		Delta: chunk.Choices[0].Delta.Content,
		Empty: chunk.Choices[0].Delta.Content == "",
	}
}
