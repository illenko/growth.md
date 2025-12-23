package aifactory

import (
	"fmt"

	"github.com/illenko/growth.md/internal/ai"
	"github.com/illenko/growth.md/internal/ai/gemini"
	"github.com/illenko/growth.md/internal/ai/openai"
)

// NewClient creates an AI client based on config
func NewClient(cfg ai.Config) (ai.AIClient, error) {
	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Create client based on provider
	switch cfg.Provider {
	case "gemini":
		return gemini.NewClient(cfg)
	case "openai":
		return openai.NewClient(cfg)
	case "anthropic":
		return nil, fmt.Errorf("anthropic provider: %w (coming soon)", ai.ErrProviderNotSupported)
	case "local":
		return nil, fmt.Errorf("local provider: %w (coming soon)", ai.ErrProviderNotSupported)
	default:
		return nil, fmt.Errorf("unknown provider '%s': %w", cfg.Provider, ai.ErrProviderNotSupported)
	}
}
