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
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	APIKey   string `yaml:"apiKey,omitempty"`
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
			Provider: "openai",
			Model:    "gpt-4",
			APIKey:   "",
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
			"openai":    true,
			"anthropic": true,
			"google":    true,
			"local":     true,
		}
		if !validProviders[c.AI.Provider] {
			return errors.New("invalid AI provider")
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
