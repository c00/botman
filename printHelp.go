package main

import "fmt"

func printHelp() {
	fmt.Printf(`Command Name: botman 

Usage: botman [OPTIONS] PROMPT

Version: %v

Description:
Botman lets you talk to an LLM. It is optimized for use in the terminal. It accepts both stdin and arguments.

Options:
	-h            Show this help message and exit
	-i            Interactive mode. Keep interacting to continue the conversation.
	
PROMPT: Any text prompt to ask the LLM.

Examples:
	1. Basic usage: botman "tell me a joke about the golang gopher"
	2. Using stdin: echo Quote a Bob Kelso joke | botman
	3. Interactive mode: botman -i
`, version)
}
