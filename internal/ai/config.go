package ai

import (
	"fmt"
	"os"
)

// Config holds AI provider configuration
type Config struct {
	Provider    string  // "gemini", "openai", "anthropic", "local"
	APIKey      string  // API key or loaded from env
	Model       string  // Model name
	Temperature float32 // Temperature for generation (0.0 - 1.0)
	MaxTokens   int     // Maximum output tokens
	BaseURL     string  // For custom endpoints (optional)
}

func (c *Config) Validate() error {
	if c.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	if c.APIKey == "" {
		c.APIKey = c.loadAPIKeyFromEnv()
	}

	if c.APIKey == "" && c.Provider != "local" {
		return fmt.Errorf("API key is required for provider %s (set in config or use env var)", c.Provider)
	}

	if c.Temperature == 0 {
		c.Temperature = 0.7
	}

	if c.MaxTokens == 0 {
		c.MaxTokens = 8000
	}

	return nil
}

func (c *Config) loadAPIKeyFromEnv() string {
	switch c.Provider {
	case "gemini":
		return os.Getenv("GEMINI_API_KEY")
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	case "anthropic":
		return os.Getenv("ANTHROPIC_API_KEY")
	default:
		return ""
	}
}

func DefaultConfig() Config {
	return Config{
		Provider:    "gemini",
		Model:       "gemini-3-flash-preview",
		Temperature: 0.7,
		MaxTokens:   8000,
	}
}
