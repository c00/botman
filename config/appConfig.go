package config

const LlmProviderOpenAi = "openai"
const LlmProviderFireworksAi = "fireworksai"

// To keep track of breaking changes in the config file
const currentVersion = 1

type AppConfig struct {
	Version     int             `yaml:"version"`
	OpenAiKey   string          `yaml:"openAiKey,omitempty"` //Deprecated, use OpenAi.ApiKey instead
	SaveHistory bool            `yaml:"saveHistory"`
	LlmProvider string          `yaml:"llmProvider"`
	OpenAi      OpenAiConfig    `yaml:"openAi"`
	FireworksAi FireworksConfig `yaml:"fireworksAi"`
}

type OpenAiConfig struct {
	ApiKey       string `yaml:"apiKey"`
	Model        string `yaml:"model"`
	SystemPrompt string `yaml:"systemPrompt"`
}

type FireworksConfig struct {
	ApiKey       string `yaml:"apiKey"`
	Model        string `yaml:"model"`
	SystemPrompt string `yaml:"systemPrompt"`
}
