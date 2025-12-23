package storage

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  string         `yaml:"version"`
	User     UserConfig     `yaml:"user"`
	AI       AIConfig       `yaml:"ai"`
	Git      GitConfig      `yaml:"git"`
	Progress ProgressConfig `yaml:"progress"`
	Display  DisplayConfig  `yaml:"display"`
	MCP      MCPConfig      `yaml:"mcp"`
}

type UserConfig struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email,omitempty"`
}

type AIConfig struct {
	Provider      string  `yaml:"provider"`         // gemini, openai, anthropic, local
	Model         string  `yaml:"model"`            // model name (uses provider default if empty)
	APIKey        string  `yaml:"apiKey,omitempty"` // optional, prefers env var
	Temperature   float32 `yaml:"temperature"`      // 0.0 - 1.0, controls randomness
	MaxTokens     int     `yaml:"maxTokens"`        // max output tokens
	DefaultStyle  string  `yaml:"defaultStyle"`     // learning style preference
	DefaultBudget string  `yaml:"defaultBudget"`    // resource budget preference
}

type GitConfig struct {
	AutoCommit            bool   `yaml:"autoCommit"`
	CommitOnUpdate        bool   `yaml:"commitOnUpdate"`
	CommitMessageTemplate string `yaml:"commitMessageTemplate"`
}

type ProgressConfig struct {
	DefaultView  string `yaml:"defaultView"`
	WeekStartDay string `yaml:"weekStartDay"`
}

type DisplayConfig struct {
	OutputFormat string `yaml:"outputFormat"`
	Theme        string `yaml:"theme"`
	DateFormat   string `yaml:"dateFormat"`
}

type MCPConfig struct {
	Enabled    bool   `yaml:"enabled"`
	ServerPath string `yaml:"serverPath,omitempty"`
	Port       int    `yaml:"port,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		User: UserConfig{
			Name:  "",
			Email: "",
		},
		AI: AIConfig{
			Provider:      "gemini",
			Model:         "gemini-3-flash-preview",
			APIKey:        "",
			Temperature:   0.7,
			MaxTokens:     8000,
			DefaultStyle:  "project-based",
			DefaultBudget: "any",
		},
		Git: GitConfig{
			AutoCommit:            true,
			CommitOnUpdate:        true,
			CommitMessageTemplate: "{{.Action}} {{.EntityType}}: {{.Title}}",
		},
		Progress: ProgressConfig{
			DefaultView:  "week",
			WeekStartDay: "monday",
		},
		Display: DisplayConfig{
			OutputFormat: "table",
			Theme:        "default",
			DateFormat:   "2006-01-02",
		},
		MCP: MCPConfig{
			Enabled:    false,
			ServerPath: "",
			Port:       3000,
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path cannot be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *Config, path string) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}

	if path == "" {
		return errors.New("config path cannot be empty")
	}

	if err := config.Validate(); err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Version == "" {
		return errors.New("config version is required")
	}

	if c.AI.Provider != "" {
		validProviders := map[string]bool{
			"gemini":    true,
			"openai":    true,
			"anthropic": true,
			"local":     true,
		}
		if !validProviders[c.AI.Provider] {
			return errors.New("invalid AI provider (must be: gemini, openai, anthropic, or local)")
		}
	}

	// Validate AI temperature
	if c.AI.Temperature < 0 || c.AI.Temperature > 1 {
		return errors.New("AI temperature must be between 0.0 and 1.0")
	}

	// Validate AI max tokens
	if c.AI.MaxTokens < 100 || c.AI.MaxTokens > 100000 {
		return errors.New("AI max tokens must be between 100 and 100000")
	}

	// Validate learning style
	if c.AI.DefaultStyle != "" {
		validStyles := map[string]bool{
			"top-down":      true,
			"bottom-up":     true,
			"project-based": true,
		}
		if !validStyles[c.AI.DefaultStyle] {
			return errors.New("invalid learning style (must be: top-down, bottom-up, or project-based)")
		}
	}

	// Validate budget
	if c.AI.DefaultBudget != "" {
		validBudgets := map[string]bool{
			"free": true,
			"paid": true,
			"any":  true,
		}
		if !validBudgets[c.AI.DefaultBudget] {
			return errors.New("invalid budget (must be: free, paid, or any)")
		}
	}

	if c.Progress.WeekStartDay != "" {
		validDays := map[string]bool{
			"monday":   true,
			"sunday":   true,
			"saturday": true,
		}
		if !validDays[c.Progress.WeekStartDay] {
			return errors.New("invalid week start day")
		}
	}

	if c.Display.OutputFormat != "" {
		validFormats := map[string]bool{
			"table": true,
			"json":  true,
			"yaml":  true,
		}
		if !validFormats[c.Display.OutputFormat] {
			return errors.New("invalid output format")
		}
	}

	return nil
}
