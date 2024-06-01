package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const version = "1.0.3"

var messages []openai.ChatCompletionMessage = []openai.ChatCompletionMessage{
	{
		Role:    openai.ChatMessageRoleSystem,
		Content: "Be concise. If code or a cli command is asked, only return the code or command. Do not add code block backticks. Output in plain text",
	},
}

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

		if strings.HasPrefix(a, "-") {
			continue
		}

		relevant = append(relevant, a)
	}

	return strings.Join(relevant, " ")
}

func main() {
	helpFlag := flag.Bool("h", false, "Display help information")
	interactiveFlag := flag.Bool("i", false, "Interactive mode")

	flag.Parse()
	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	in := getPipedIn()
	arg := getArg()

	content := strings.TrimSpace(fmt.Sprintf("%v %v", in, arg))

	if content == "" {
		*interactiveFlag = true
	}

	//Main program loop
	for {
		if content != "" {
			getResponse(content)

			if *interactiveFlag {
				fmt.Print("\n\n")
			}
		}

		if *interactiveFlag {
			content = getCliInput()
			if content == "" {
				break
			}
		} else {
			break
		}
	}

	fmt.Println()
}

func getCliInput() string {
	//Wait for an enter
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("You: ")
	text, _ := reader.ReadString('\n')

	if text == "\n" {
		return ""
	}

	fmt.Println()

	return text
}

func getResponse(content string) {
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: messages,
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting Chat Completion:", err)
		os.Exit(1)
	}
	defer stream.Close()

	responseContent := make([]string, 0, 50)

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			responseMessage := strings.Join(responseContent, "")
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: responseMessage,
			})
			return
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Stream error:", err)
			os.Exit(1)
		}

		fmt.Print(response.Choices[0].Delta.Content)
		responseContent = append(responseContent, response.Choices[0].Delta.Content)
	}

}
