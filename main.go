package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/c00/botman/botman"
	"github.com/c00/botman/cli"
	"github.com/c00/botman/config"
	"github.com/c00/botman/history"
	"github.com/c00/botman/models"
)

const version = "1.1.5"

var messages []models.ChatMessage = []models.ChatMessage{
	{
		Role:    models.ChatMessageRoleSystem,
		Content: "",
	},
}

var date = time.Now()
var appConfig config.BotmanConfig

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

	//Set System Prompt
	messages[0].Content = fmt.Sprintf("The current date and time is %v. %v", time.Now().Format(time.RFC1123Z), appConfig.SystemPrompt)

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

func getResponse(content string) {
	messages = append(messages, models.ChatMessage{
		Role:    models.ChatMessageRoleUser,
		Content: content,
	})

	//Instantiate chatter
	chatter := botman.GetChatter(appConfig)
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
