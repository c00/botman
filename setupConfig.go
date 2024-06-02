package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/c00/botman/config"
)

func setupConfig() {
	fmt.Print("Botman initialization\n\n")

	//Wait for an enter
	reader := bufio.NewReader(os.Stdin)
	if appConfig.OpenAiKey == "" {
		fmt.Println("There is currently no OpenAI API key set")
	} else {
		fmt.Println("Current OpenAI API key:", appConfig.OpenAiKey)
	}

	fmt.Print("Enter your (new) OpenAi API key: ")
	text, _ := reader.ReadString('\n')

	fmt.Println()

	if text == "\n" {
		fmt.Println("New key is empty. No changes are made.")
		return
	}

	//trim the latest newline
	regex := regexp.MustCompile(`^\s+|\s+$`)
	replaced := regex.ReplaceAll([]byte(text), []byte{})

	newKey := string(replaced)
	if newKey == "" {
		fmt.Println("New key is empty. No changes are made.")
		return
	}

	appConfig.OpenAiKey = newKey
	err := config.SaveForUser(appConfig)
	if err != nil {
		fmt.Println("could not update the configuration:", err)
		os.Exit(1)
	}

	fmt.Println("Configuration has been updated")
}
