package claude

import (
	"encoding/json"
)

type parsedChunk struct {
	Empty       bool
	LastMessage bool
	Delta       string
}

type contentBlockStartChunk struct {
	Type         string    `json:"type"`
	Index        int       `json:"index"`
	ContentBlock textChunk `json:"content_block"`
}

type contentBlockDeltaChunk struct {
	Type  string    `json:"type"`
	Index int       `json:"index"`
	Delta textChunk `json:"delta"`
}

type textChunk struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func getStartText(data []byte) parsedChunk {
	//Parse JSON
	chunk := contentBlockStartChunk{}
	err := json.Unmarshal(data, &chunk)
	if err != nil {
		panic(err)
	}

	return parsedChunk{
		Delta: chunk.ContentBlock.Text,
		Empty: chunk.ContentBlock.Text == "",
	}
}

func getDeltaText(data []byte) parsedChunk {
	//Parse JSON
	chunk := contentBlockDeltaChunk{}
	err := json.Unmarshal(data, &chunk)
	if err != nil {
		panic(err)
	}

	return parsedChunk{
		Delta: chunk.Delta.Text,
		Empty: chunk.Delta.Text == "",
	}
}

func parseChunk(data []byte, event string) parsedChunk {
	switch event {
	case "message_stop":
		return parsedChunk{LastMessage: true, Empty: true}
	case "content_block_stop":
		fallthrough
	case "message_start":
		fallthrough
	case "message_delta":
		fallthrough
	case "ping":
		return parsedChunk{Empty: true}
	case "content_block_start":
		return getStartText(data)
	case "content_block_delta":
		return getDeltaText(data)
	}

	panic("unexpected event: " + event)
}
