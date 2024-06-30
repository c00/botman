package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/c00/botman/cli"
	"github.com/c00/botman/config"
	"github.com/c00/botman/history"
	"github.com/c00/botman/models"
	"github.com/c00/botman/providers/claude"
	"github.com/c00/botman/providers/fireworks"
	"github.com/c00/botman/providers/openai"
)

const version = "1.1.2"

var messages []models.ChatMessage = []models.ChatMessage{
	{
		Role:    models.ChatMessageRoleSystem,
		Content: fmt.Sprintf("The current date and time is %v. Be concise. If code or a cli command is asked, only return the code or command. Do not add code block backticks. Output in plain text", time.Now().Format(time.RFC1123Z)),
	},
}

var date = time.Now()
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
	initFlag := flag.Bool("init", false, "Initialise or update the configuration and set API keys")

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

			//Save history after every response so we dont' miss out if we ctrl+c out on the prompt.
			if appConfig.SaveHistory {
				_, err := history.SaveChat(date, messages)
				if err != nil {
					fmt.Println("could not save chat history:", err)
				}
			}

			if *interactiveFlag {
				fmt.Print("\n\n")
			}
		}

		if *interactiveFlag {
			content = cli.GetInput("You")
			if content == "" {
				break
			}
		} else {
			break
		}
	}

	fmt.Println()
}

func getChatter() models.Chatter {
	if appConfig.LlmProvider == config.LlmProviderOpenAi {
		return openai.NewChatBot(appConfig.OpenAi)
	} else if appConfig.LlmProvider == config.LlmProviderFireworksAi {
		return fireworks.NewChatBot(appConfig.FireworksAi)
	} else if appConfig.LlmProvider == config.LlmProviderClaude {
		return claude.NewChatBot(appConfig.Claude)
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	//Let the channel stream to stdout
	go func(ch chan string) {
		for content := range ch {
			fmt.Print(content)
		}
		wg.Done()
	}(ch)

	//call GetResponse
	response := chatter.GetResponse(messages, ch)

	//Add the response message to the message slice
	messages = append(messages, models.ChatMessage{
		Role:    models.ChatMessageRoleAssistant,
		Content: response,
	})

	wg.Wait()
}
