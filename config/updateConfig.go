package config

func updateConfig(cfg AppConfig) (AppConfig, bool) {
	if cfg.Version == currentVersion {
		return cfg, false
	}

	//This is a slice of update functions to update the appConfig to the latest version.
	updaters := []func(AppConfig) AppConfig{
		//0 - openAi got its own config.
		func(cfg AppConfig) AppConfig {
			cfg.OpenAi = OpenAiConfig{
				ApiKey:       cfg.OpenAiKey,
				SystemPrompt: "",
			}
			cfg.OpenAiKey = ""
			cfg.Version = 1
			return cfg
		},
		//Add more when needed.
	}

	for i := cfg.Version; i < currentVersion; i++ {
		if updaters[i] == nil {
			//This panic will also occur if you forgot to add an updater function for a newer version.
			panic("Cannot upgrade app config. Current config is too old.")
		}
		cfg = updaters[i](cfg)
	}

	return cfg, true
}
