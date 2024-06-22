package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetInput(label string) string {
	//Wait for an enter
	reader := bufio.NewReader(os.Stdin)
	if label != "" {
		fmt.Printf("%v: ", label)
	}
	text, _ := reader.ReadString('\n')

	text = strings.TrimSuffix(text, "\n")

	if text == "" {
		return ""
	}

	fmt.Println()

	return text
}

// Returns the index of the chosen option.
func GetChoice(choices []string, current int) int {
	if len(choices) == 0 {
		panic("No choices to choose from")
	}

	//Print choices
	for i, choice := range choices {
		fmt.Printf("%v. %v", i+1, choice)
		if i == current {
			fmt.Print(" (current)\n")
		} else {
			fmt.Println()
		}
	}
	fmt.Println()

	fmt.Printf("Choose [1-%v]: ", len(choices))

	for {
		var err error = nil
		choiceStr := GetInput("")
		if choiceStr == "" {
			return current
		}

		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Please enter a number")
		} else if choice < 1 {
			fmt.Println("Number should be greater than 0")
		} else if choice > len(choices) {
			fmt.Printf("Number should not be larger than %v\n", len(choices))
		} else {
			return choice - 1
		}
	}
}
