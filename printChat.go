package main

import (
	"fmt"
	"os"

	"github.com/c00/botman/history"
)

func printChat(lookback int) {
	chat, err := history.LoadChat(lookback)
	if err != nil {
		fmt.Println("cannot load chat:", err)
		os.Exit(1)
	}

	chat.Print()
}

func printLastResponse() {
	chat, err := history.LoadChat(0)
	if err != nil {
		fmt.Println("cannot load chat:", err)
		os.Exit(1)
	}

	chat.PrintLastMessage()
}
