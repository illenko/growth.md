package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Python Programming", "python-programming"},
		{"Go & Rust", "go-rust"},
		{"Machine Learning (ML)", "machine-learning-ml"},
		{"  Spaces   Everywhere  ", "spaces-everywhere"},
		{"Hello___World", "hello-world"},
		{"a", "a"},
		{"", "untitled"},
		{"Very Long Title That Should Be Truncated Because It Exceeds Maximum Length", "very-long-title-that-should-be-truncated-because-i"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateFileName(t *testing.T) {
	tests := []struct {
		id       core.EntityID
		title    string
		expected string
	}{
		{"skill-001", "Python", "skill-001-python.md"},
		{"goal-042", "Senior Engineer", "goal-042-senior-engineer.md"},
		{"resource-099", "Clean Code (Book)", "resource-099-clean-code-book.md"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result := GenerateFileName(tt.id, tt.title)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateNextID(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("generates first ID when no files exist", func(t *testing.T) {
		skillsDir := filepath.Join(tmpDir, "skills")
		require.NoError(t, os.MkdirAll(skillsDir, 0755))

		id, err := GenerateNextIDInPath("skill", tmpDir)

		require.NoError(t, err)
		assert.Equal(t, "skill-001", string(id))
	})

	t.Run("generates next ID based on existing files", func(t *testing.T) {
		skillsDir := filepath.Join(tmpDir, "skills")
		require.NoError(t, os.WriteFile(filepath.Join(skillsDir, "skill-001-python.md"), []byte(""), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(skillsDir, "skill-002-go.md"), []byte(""), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(skillsDir, "skill-005-rust.md"), []byte(""), 0644))

		id, err := GenerateNextIDInPath("skill", tmpDir)

		require.NoError(t, err)
		assert.Equal(t, "skill-006", string(id))
	})

	t.Run("handles gaps in numbering", func(t *testing.T) {
		goalsDir := filepath.Join(tmpDir, "goals")
		require.NoError(t, os.MkdirAll(goalsDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(goalsDir, "goal-001-first.md"), []byte(""), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(goalsDir, "goal-010-tenth.md"), []byte(""), 0644))

		id, err := GenerateNextIDInPath("goal", tmpDir)

		require.NoError(t, err)
		assert.Equal(t, "goal-011", string(id))
	})

	t.Run("returns error for unknown entity type", func(t *testing.T) {
		_, err := GenerateNextIDInPath("unknown", tmpDir)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown entity type")
	})
}
