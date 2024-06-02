package config

type AppConfig struct {
	// Should be used to identify clients for creating locks.
	OpenAiKey   string `yaml:"openAiKey"`
	SaveHistory bool   `yaml:"saveHistory"`
}
