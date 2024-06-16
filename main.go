package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/c00/botman/config"
	"github.com/c00/botman/history"
	"github.com/c00/botman/models"
	"github.com/c00/botman/openai"
)

const version = "1.1.0"

var messages []models.ChatMessage = []models.ChatMessage{
	{
		Role:    models.ChatMessageRoleSystem,
		Content: fmt.Sprintf("The current date and time is %v. Be concise. If code or a cli command is asked, only return the code or command. Do not add code block backticks. Output in plain text", time.Now().Format(time.RFC1123Z)),
	},
}

var appConfig config.AppConfig

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
	return strings.Join(flag.Args(), " ")
}

func main() {
	helpFlag := flag.Bool("h", false, "Display help information")
	interactiveFlag := flag.Bool("i", false, "Interactive mode")
	historyFlag := flag.Int("history", -1, "Show historical chat, looking back x chats")
	printLast := flag.Bool("l", false, "Print the last response")
	initFlag := flag.Bool("init", false, "Initialise the configuration and set the OpenAI API Key")

	appConfig = config.LoadFromUser()

	flag.Parse()
	//Print help
	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	if *initFlag {
		setupConfig()
		os.Exit(0)
	}

	//Print chat from history
	if *historyFlag >= 0 {
		printChat(*historyFlag)
		os.Exit(0)
	}

	if *printLast {
		printLastResponse()
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

	if appConfig.SaveHistory {
		_, err := history.SaveChat(messages)
		if err != nil {
			fmt.Println("could not save chat history:", err)
		}
	}
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

func getChatter() models.Chatter {
	if appConfig.LlmProvider == config.LlmProviderOpenAi {
		return openai.NewChatBot(appConfig.OpenAiKey)
	}

	panic(fmt.Sprintf("chatter '%v' not implemented", appConfig.LlmProvider))
}

func getResponse(content string) {
	messages = append(messages, models.ChatMessage{
		Role:    models.ChatMessageRoleUser,
		Content: content,
	})

	//Instantiate chatter
	chatter := getChatter()
	ch := make(chan string)

	//Let the channel stream to stdout
	go func(ch chan string) {
		for content := range ch {
			fmt.Print(content)
		}
	}(ch)

	//call GetResponse
	response := chatter.GetResponse(messages, ch)

	//Add the response message to the message slice
	messages = append(messages, models.ChatMessage{
		Role:    models.ChatMessageRoleAssistant,
		Content: response,
	})
}
