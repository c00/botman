package config

const LlmProviderOpenAi = "openai"

type AppConfig struct {
	// Should be used to identify clients for creating locks.
	OpenAiKey   string `yaml:"openAiKey"`
	SaveHistory bool   `yaml:"saveHistory"`
	LlmProvider string `yaml:"llmProvider"`
}
