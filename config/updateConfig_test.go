package config

import "testing"

func TestUpdateConfig(t *testing.T) {
	config := BotmanConfig{
		OpenAiKey: "12345",
	}

	config, needsSave := updateConfig(config)

	if !needsSave {
		t.Fatal("config should need saving")
	}

	if config.Version != currentVersion {
		t.Fatal("Expected version to be updated")
	}

	if config.OpenAiKey != "" {
		t.Fatal("Expected deprecated openAi key to be empty string")
	}

	if config.OpenAi.ApiKey != "12345" {
		t.Fatal("Expected new openAi key to be set")
	}

	_, saveAgain := updateConfig(config)
	if saveAgain {
		t.Fatal("We shouldn't need to save again")
	}
}
