package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func NewAppConfig() AppConfig {
	openAiKey := os.Getenv("OPENAI_API_KEY")
	return AppConfig{
		Version:     0,
		SaveHistory: true,
		LlmProvider: LlmProviderOpenAi,
		OpenAi: OpenAiConfig{
			ApiKey: openAiKey,
		},
	}
}

func SaveForUser(config AppConfig) error {
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

func LoadFromUser() AppConfig {
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

func Load(path string) (AppConfig, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{}, err
	}

	config := NewAppConfig()
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}
