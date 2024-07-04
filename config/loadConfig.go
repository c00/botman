package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

const defaultPrompt = "Be concise. If code or a cli command is asked, only return the code or command. Do not add code block backticks. Output in plain text"

func NewAppConfig() BotmanConfig {
	openAiKey := os.Getenv("OPENAI_API_KEY")
	return BotmanConfig{
		Version:      0,
		SaveHistory:  true,
		LlmProvider:  LlmProviderOpenAi,
		SystemPrompt: defaultPrompt,
		OpenAi: OpenAiConfig{
			ApiKey: openAiKey,
			Model:  "gpt-4o",
		},
		FireworksAi: FireworksConfig{
			Model: "accounts/fireworks/models/mixtral-8x22b-instruct",
		},
		Claude: ClaudeConfig{
			Model:     "claude-3-5-sonnet-20240620",
			MaxTokens: 1024,
		},
	}
}

func SaveForUser(config BotmanConfig) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	configPath := filepath.Join(u.HomeDir, APP_FOLDER, "config.yaml")

	dir := filepath.Dir(configPath)
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
			if err != nil {
				return err
			}
		}
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("%v should be a directory", dir)
	}

	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, bytes, 0700)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromUser() BotmanConfig {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Could not load config from user:", err)
		os.Exit(1)
	}

	configPath := filepath.Join(u.HomeDir, APP_FOLDER, "config.yaml")

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		//Create it instead
		config := NewAppConfig()
		SaveForUser(config)
		return config
	}

	config, err := Load(configPath)
	if err != nil {
		fmt.Println("Could not read config from user:", err)
		os.Exit(1)
	}

	//Update to latest version
	config, shouldSave := updateConfig(config)
	if shouldSave {
		SaveForUser(config)
	}

	return config
}

func LoadFromEnv() BotmanConfig {
	def := NewAppConfig()
	return BotmanConfig{
		Version:      currentVersion,
		SaveHistory:  boolFromEnv("BOTMAN_SAVE_HISTORY", false),
		LlmProvider:  stringFromEnv("BOTMAN_LLM", def.LlmProvider),
		SystemPrompt: stringFromEnv("BOTMAN_PROMPT", def.SystemPrompt),
		OpenAi: OpenAiConfig{
			ApiKey:       stringFromEnv("BOTMAN_OPENAI_API_KEY", def.OpenAi.ApiKey),
			Model:        stringFromEnv("BOTMAN_OPENAI_MODEL", def.OpenAi.Model),
			SystemPrompt: stringFromEnv("BOTMAN_OPENAI_PROMPT", def.OpenAi.SystemPrompt),
		},
		FireworksAi: FireworksConfig{
			ApiKey:       stringFromEnv("BOTMAN_FIREWORKS_API_KEY", def.FireworksAi.ApiKey),
			Model:        stringFromEnv("BOTMAN_FIREWORKS_MODEL", def.FireworksAi.Model),
			SystemPrompt: stringFromEnv("BOTMAN_FIREWORKS_PROMPT", def.FireworksAi.SystemPrompt),
		},
		Claude: ClaudeConfig{
			ApiKey:       stringFromEnv("BOTMAN_CLAUDE_API_KEY", def.Claude.ApiKey),
			Model:        stringFromEnv("BOTMAN_CLAUDE_MODEL", def.Claude.Model),
			SystemPrompt: stringFromEnv("BOTMAN_CLAUDE_PROMPT", def.Claude.SystemPrompt),
			MaxTokens:    intFromEnv("BOTMAN_CLAUDE_MAX_TOKENS", def.Claude.MaxTokens),
		},
	}
}

func stringFromEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func boolFromEnv(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	if val == "1" || val == "true" {
		return true
	}

	return false
}

func intFromEnv(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return intVal
}

func Load(path string) (BotmanConfig, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return BotmanConfig{}, err
	}

	config := NewAppConfig()
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return BotmanConfig{}, err
	}

	return config, nil
}
