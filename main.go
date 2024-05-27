package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const version = "1.0.1"

func getPipedIn() string {
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting stdin info:", err)
		os.Exit(1)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		//No stdin
		return ""
	}

	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting stdin info:", err)
		os.Exit(1)
	}

	return string(content)
}

func getArg() string {
	first := true
	relevant := []string{}
	for _, a := range os.Args {
		if first {
			first = false
			continue
		}

		if strings.HasPrefix("-", a) {
			continue
		}

		relevant = append(relevant, a)
	}

	return strings.Join(relevant, " ")
}

func main() {
	helpFlag := flag.Bool("h", false, "Display help information")
	flag.Parse()
	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	in := getPipedIn()
	arg := getArg()

	getResponse(strings.TrimSpace(fmt.Sprintf("%v %v", in, arg)))
}

func printHelp() {
	fmt.Printf(`Command Name: botman 

Usage: botman [OPTIONS] PROMPT

Version: %v

Description:
Botman lets you talk to an LLM. It is optimized for use in the terminal. It accepts both stdin and arguments.

Options:
	-h                      Show this help message and exit
	
PROMPT: Any text prompt to ask the LLM.

Examples:
	1. Basic usage: botman "tell me a joke about the golang gopher"
	2. using stdin: echo Quote a Bob Kelso joke | botman
`, version)
}

func getResponse(content string) {
	if content == "" {
		fmt.Print("No input in stdin, nor as an argument.\n\n")
		printHelp()
		os.Exit(0)
	}
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Be concise. If code or a cli command is asked, only return the code or command. Do not add code block backticks. Output in plain text",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting Chat Completion:", err)
		os.Exit(1)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println()
			return
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Stream error:", err)
			os.Exit(1)
		}

		fmt.Print(response.Choices[0].Delta.Content)
	}

}
