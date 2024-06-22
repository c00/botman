package main

import (
	"fmt"
	"os"

	"github.com/c00/botman/cli"
	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
)

func setupConfig() {
	fmt.Print("Botman initialization\n\n")

	//Preferred provider
	fmt.Println("Select your preferred LLM Provider")
	currentChoiceIndex := 0
	if appConfig.LlmProvider == config.LlmProviderFireworksAi {
		currentChoiceIndex = 1
	}
	choice := cli.GetChoice([]string{"Open AI", "Fireworks AI"}, currentChoiceIndex)
	if choice == 0 {
		appConfig.LlmProvider = config.LlmProviderOpenAi
		setupApiKey(&appConfig.OpenAi.ApiKey, "OpenAI")
		chooseModel(&appConfig.OpenAi.Model, models.OpenAiModels)
	} else if choice == 1 {
		appConfig.LlmProvider = config.LlmProviderFireworksAi
		setupApiKey(&appConfig.FireworksAi.ApiKey, "Fireworks AI")
		chooseModel(&appConfig.FireworksAi.Model, models.FireworksAIModels)
	}

	fmt.Println()
	err := config.SaveForUser(appConfig)
	if err != nil {
		fmt.Println("could not update the configuration:", err)
		os.Exit(1)
	}

	fmt.Println("Configuration has been updated")
}

// Setup API key and return true if changes were made.
func setupApiKey(key *string, name string) bool {
	if *key == "" {
		input := cli.GetInput(fmt.Sprintf("Enter your %v API key", name))
		*key = input
		return true
	} else {
		fmt.Printf("Current %v API key: %v\n", name, *key)
		input := cli.GetInput(fmt.Sprintf("Enter your new %v API key, or press [enter] to keep the current one", name))
		if input != "" {
			*key = input
			return true
		}
	}

	return false
}

// Select a model
func chooseModel(currentModel *string, models []string) bool {

	fmt.Println("\nChoose a model:")

	index := indexOf(models, *currentModel)
	chosen := cli.GetChoice(models, index)
	if chosen == -1 {
		return false
	}

	newModel := models[chosen]
	if newModel == *currentModel {
		return false
	} else {
		*currentModel = newModel
		return true
	}
}

func indexOf(s []string, searchString string) int {
	for i, v := range s {
		if v == searchString {
			return i
		}
	}
	return -1
}
