package botman

import (
	"fmt"

	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
	"github.com/c00/botman/providers/claude"
	"github.com/c00/botman/providers/fireworks"
	"github.com/c00/botman/providers/openai"
)

func GetChatter(cfg config.BotmanConfig) models.Chatter {
	if cfg.LlmProvider == config.LlmProviderOpenAi {
		return openai.NewChatBot(cfg.OpenAi)
	} else if cfg.LlmProvider == config.LlmProviderFireworksAi {
		return fireworks.NewChatBot(cfg.FireworksAi)
	} else if cfg.LlmProvider == config.LlmProviderClaude {
		return claude.NewChatBot(cfg.Claude)
	}

	panic(fmt.Sprintf("chatter '%v' not implemented", cfg.LlmProvider))
}
