package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("creates valid default config", func(t *testing.T) {
		config := DefaultConfig()

		require.NotNil(t, config)
		assert.Equal(t, "1.0", config.Version)
		assert.Equal(t, "openai", config.AI.Provider)
		assert.Equal(t, "gpt-4", config.AI.Model)
		assert.True(t, config.Git.AutoCommit)
		assert.Equal(t, "week", config.Progress.DefaultView)
		assert.Equal(t, "table", config.Display.OutputFormat)
		assert.False(t, config.MCP.Enabled)
	})

	t.Run("default config passes validation", func(t *testing.T) {
		config := DefaultConfig()

		err := config.Validate()

		assert.NoError(t, err)
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("loads valid config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

		original := DefaultConfig()
		original.User.Name = "Test User"
		original.User.Email = "test@example.com"

		err := SaveConfig(original, configPath)
		require.NoError(t, err)

		loaded, err := LoadConfig(configPath)

		require.NoError(t, err)
		assert.Equal(t, "Test User", loaded.User.Name)
		assert.Equal(t, "test@example.com", loaded.User.Email)
	})

	t.Run("fails with empty path", func(t *testing.T) {
		_, err := LoadConfig("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path cannot be empty")
	})

	t.Run("fails when file does not exist", func(t *testing.T) {
		_, err := LoadConfig("/nonexistent/config.yml")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("fails with invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

		err := os.WriteFile(configPath, []byte("invalid: yaml: content:\n  - broken"), 0644)
		require.NoError(t, err)

		_, err = LoadConfig(configPath)

		assert.Error(t, err)
	})

	t.Run("fails validation with invalid config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

		invalidYAML := `version: "1.0"
user:
  name: "Test"
ai:
  provider: "invalid-provider"
  model: "gpt-4"
git:
  autoCommit: true
  commitOnUpdate: true
  commitMessageTemplate: "{{.Action}}"
progress:
  defaultView: "week"
  weekStartDay: "monday"
display:
  outputFormat: "table"
  theme: "default"
  dateFormat: "2006-01-02"
mcp:
  enabled: false
`
		err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
		require.NoError(t, err)

		_, err = LoadConfig(configPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid AI provider")
	})
}

func TestSaveConfig(t *testing.T) {
	t.Run("saves config successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

		config := DefaultConfig()
		config.User.Name = "Test User"

		err := SaveConfig(config, configPath)
		require.NoError(t, err)

		assert.FileExists(t, configPath)
	})

	t.Run("creates directory if it does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "nested", "dir", "config.yml")

		config := DefaultConfig()

		err := SaveConfig(config, configPath)
		require.NoError(t, err)

		assert.FileExists(t, configPath)
	})

	t.Run("fails with nil config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

		err := SaveConfig(nil, configPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("fails with empty path", func(t *testing.T) {
		config := DefaultConfig()

		err := SaveConfig(config, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path cannot be empty")
	})

	t.Run("fails validation before saving", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

		config := DefaultConfig()
		config.AI.Provider = "invalid-provider"

		err := SaveConfig(config, configPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid AI provider")
	})

	t.Run("saved file can be loaded back", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

		original := DefaultConfig()
		original.User.Name = "John Doe"
		original.AI.Model = "gpt-3.5-turbo"
		original.Git.AutoCommit = false

		err := SaveConfig(original, configPath)
		require.NoError(t, err)

		loaded, err := LoadConfig(configPath)
		require.NoError(t, err)

		assert.Equal(t, original.User.Name, loaded.User.Name)
		assert.Equal(t, original.AI.Model, loaded.AI.Model)
		assert.Equal(t, original.Git.AutoCommit, loaded.Git.AutoCommit)
	})
}

func TestConfigValidate(t *testing.T) {
	t.Run("validates version is required", func(t *testing.T) {
		config := DefaultConfig()
		config.Version = ""

		err := config.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "version is required")
	})

	t.Run("validates AI provider", func(t *testing.T) {
		tests := []struct {
			provider string
			valid    bool
		}{
			{"openai", true},
			{"anthropic", true},
			{"local", true},
			{"invalid", false},
			{"", true}, // empty is allowed (optional)
		}

		for _, tt := range tests {
			config := DefaultConfig()
			config.AI.Provider = tt.provider

			err := config.Validate()

			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid AI provider")
			}
		}
	})

	t.Run("validates week start day", func(t *testing.T) {
		tests := []struct {
			day   string
			valid bool
		}{
			{"monday", true},
			{"sunday", true},
			{"saturday", true},
			{"tuesday", false},
			{"", true}, // empty is allowed (optional)
		}

		for _, tt := range tests {
			config := DefaultConfig()
			config.Progress.WeekStartDay = tt.day

			err := config.Validate()

			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid week start day")
			}
		}
	})

	t.Run("validates output format", func(t *testing.T) {
		tests := []struct {
			format string
			valid  bool
		}{
			{"table", true},
			{"json", true},
			{"yaml", true},
			{"xml", false},
			{"", true}, // empty is allowed (optional)
		}

		for _, tt := range tests {
			config := DefaultConfig()
			config.Display.OutputFormat = tt.format

			err := config.Validate()

			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid output format")
			}
		}
	})
}

func TestConfigRoundTrip(t *testing.T) {
	t.Run("full config round trip", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yml")

		original := &Config{
			Version: "1.0",
			User: UserConfig{
				Name:  "Alice Smith",
				Email: "alice@example.com",
			},
			AI: AIConfig{
				Provider: "anthropic",
				Model:    "claude-3-opus",
				APIKey:   "sk-test-key",
			},
			Git: GitConfig{
				AutoCommit:            false,
				CommitOnUpdate:        true,
				CommitMessageTemplate: "custom: {{.Title}}",
			},
			Progress: ProgressConfig{
				DefaultView:  "month",
				WeekStartDay: "sunday",
			},
			Display: DisplayConfig{
				OutputFormat: "json",
				Theme:        "dark",
				DateFormat:   "2006/01/02",
			},
			MCP: MCPConfig{
				Enabled:    true,
				ServerPath: "/usr/local/bin/mcp-server",
				Port:       8080,
			},
		}

		err := SaveConfig(original, configPath)
		require.NoError(t, err)

		loaded, err := LoadConfig(configPath)
		require.NoError(t, err)

		assert.Equal(t, original.Version, loaded.Version)
		assert.Equal(t, original.User.Name, loaded.User.Name)
		assert.Equal(t, original.User.Email, loaded.User.Email)
		assert.Equal(t, original.AI.Provider, loaded.AI.Provider)
		assert.Equal(t, original.AI.Model, loaded.AI.Model)
		assert.Equal(t, original.AI.APIKey, loaded.AI.APIKey)
		assert.Equal(t, original.Git.AutoCommit, loaded.Git.AutoCommit)
		assert.Equal(t, original.Git.CommitOnUpdate, loaded.Git.CommitOnUpdate)
		assert.Equal(t, original.Git.CommitMessageTemplate, loaded.Git.CommitMessageTemplate)
		assert.Equal(t, original.Progress.DefaultView, loaded.Progress.DefaultView)
		assert.Equal(t, original.Progress.WeekStartDay, loaded.Progress.WeekStartDay)
		assert.Equal(t, original.Display.OutputFormat, loaded.Display.OutputFormat)
		assert.Equal(t, original.Display.Theme, loaded.Display.Theme)
		assert.Equal(t, original.Display.DateFormat, loaded.Display.DateFormat)
		assert.Equal(t, original.MCP.Enabled, loaded.MCP.Enabled)
		assert.Equal(t, original.MCP.ServerPath, loaded.MCP.ServerPath)
		assert.Equal(t, original.MCP.Port, loaded.MCP.Port)
	})
}
