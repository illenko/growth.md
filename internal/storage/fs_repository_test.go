package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFilesystemRepository(t *testing.T) {
	t.Run("creates repository with valid parameters", func(t *testing.T) {
		tmpDir := t.TempDir()

		repo, err := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		require.NoError(t, err)
		assert.NotNil(t, repo)
		assert.Equal(t, tmpDir, repo.basePath)
		assert.Equal(t, "skill", repo.entityType)

		// Verify directory was created
		_, err = os.Stat(tmpDir)
		assert.NoError(t, err)
	})

	t.Run("fails with empty basePath", func(t *testing.T) {
		_, err := NewFilesystemRepository[core.Skill]("", "skill")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "basePath cannot be empty")
	})

	t.Run("fails with empty entityType", func(t *testing.T) {
		tmpDir := t.TempDir()

		_, err := NewFilesystemRepository[core.Skill](tmpDir, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "entityType cannot be empty")
	})

	t.Run("creates nested directory if it doesn't exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		nestedDir := filepath.Join(tmpDir, "nested", "path")

		repo, err := NewFilesystemRepository[core.Skill](nestedDir, "skill")

		require.NoError(t, err)
		assert.NotNil(t, repo)

		// Verify nested directory was created
		_, err = os.Stat(nestedDir)
		assert.NoError(t, err)
	})
}

func TestFilesystemRepository_Create(t *testing.T) {
	t.Run("creates entity successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		skill, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		skill.Body = "Python is a versatile programming language."

		err := repo.Create(skill)

		require.NoError(t, err)

		// Verify file was created
		files, _ := os.ReadDir(tmpDir)
		assert.Len(t, files, 1)
		assert.Contains(t, files[0].Name(), "skill-001")
		assert.Contains(t, files[0].Name(), "python")
	})

	t.Run("fails with nil entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		err := repo.Create(nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("fails when entity already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		skill, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)

		err := repo.Create(skill)
		require.NoError(t, err)

		err = repo.Create(skill)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestFilesystemRepository_GetByID(t *testing.T) {
	t.Run("retrieves existing entity without body", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		original, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		original.Body = "This is the body content."
		repo.Create(original)

		retrieved, err := repo.GetByID("skill-001")

		require.NoError(t, err)
		assert.Equal(t, original.ID, retrieved.ID)
		assert.Equal(t, original.Title, retrieved.Title)
		assert.Equal(t, original.Category, retrieved.Category)
		assert.Equal(t, original.Level, retrieved.Level)
		assert.Empty(t, retrieved.Body) // Body should not be included
	})

	t.Run("fails with non-existent entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		_, err := repo.GetByID("skill-999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		_, err := repo.GetByID("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestFilesystemRepository_GetByIDWithBody(t *testing.T) {
	t.Run("retrieves entity with body", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		original, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		original.Body = "This is the body content."
		repo.Create(original)

		retrieved, err := repo.GetByIDWithBody("skill-001")

		require.NoError(t, err)
		assert.Equal(t, original.ID, retrieved.ID)
		assert.Equal(t, original.Title, retrieved.Title)
		assert.Contains(t, retrieved.Body, "This is the body content") // Body should be included
	})
}

func TestFilesystemRepository_GetAll(t *testing.T) {
	t.Run("retrieves all entities", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		skill1, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		skill2, _ := core.NewSkill("skill-002", "Go", "programming", core.LevelAdvanced)
		skill3, _ := core.NewSkill("skill-003", "Docker", "devops", core.LevelIntermediate)

		repo.Create(skill1)
		repo.Create(skill2)
		repo.Create(skill3)

		entities, err := repo.GetAll()

		require.NoError(t, err)
		assert.Len(t, entities, 3)
	})

	t.Run("returns empty slice when no entities exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		entities, err := repo.GetAll()

		require.NoError(t, err)
		assert.Empty(t, entities)
	})
}

func TestFilesystemRepository_Update(t *testing.T) {
	t.Run("updates existing entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		skill, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		repo.Create(skill)

		// Update skill
		skill.Level = core.LevelAdvanced
		skill.Body = "Updated body content."
		err := repo.Update(skill)

		require.NoError(t, err)

		// Verify update
		retrieved, _ := repo.GetByIDWithBody("skill-001")
		assert.Equal(t, core.LevelAdvanced, retrieved.Level)
		assert.Contains(t, retrieved.Body, "Updated body content")
	})

	t.Run("renames file when title changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		skill, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		repo.Create(skill)

		// Update title
		skill.Title = "Python 3"
		err := repo.Update(skill)

		require.NoError(t, err)

		// Verify old file is gone and new file exists
		files, _ := os.ReadDir(tmpDir)
		assert.Len(t, files, 1)
		assert.Contains(t, files[0].Name(), "python-3")
		assert.NotContains(t, files[0].Name(), "python.md")
	})

	t.Run("fails with non-existent entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		skill, _ := core.NewSkill("skill-999", "Python", "programming", core.LevelIntermediate)
		err := repo.Update(skill)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("fails with nil entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		err := repo.Update(nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

func TestFilesystemRepository_Delete(t *testing.T) {
	t.Run("deletes existing entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		skill, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		repo.Create(skill)

		err := repo.Delete("skill-001")

		require.NoError(t, err)

		// Verify file is gone
		files, _ := os.ReadDir(tmpDir)
		assert.Empty(t, files)
	})

	t.Run("fails with non-existent entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		err := repo.Delete("skill-999")

		assert.Error(t, err)
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		err := repo.Delete("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestFilesystemRepository_Search(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

	// Create test data
	python, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
	python.AddTag("ml")
	python.AddTag("backend")

	go1, _ := core.NewSkill("skill-002", "Go", "programming", core.LevelAdvanced)
	go1.AddTag("backend")
	go1.AddTag("systems")

	docker, _ := core.NewSkill("skill-003", "Docker", "devops", core.LevelIntermediate)
	docker.AddTag("containers")

	repo.Create(python)
	repo.Create(go1)
	repo.Create(docker)

	t.Run("searches by title", func(t *testing.T) {
		results, err := repo.Search("python")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Python", results[0].Title)
	})

	t.Run("searches by tag", func(t *testing.T) {
		results, err := repo.Search("backend")

		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("returns all entities with empty query", func(t *testing.T) {
		results, err := repo.Search("")

		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("returns empty results for non-matching query", func(t *testing.T) {
		results, err := repo.Search("nonexistent")

		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("is case insensitive", func(t *testing.T) {
		results, err := repo.Search("PYTHON")

		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

func TestFilesystemRepository_Exists(t *testing.T) {
	t.Run("returns true for existing entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		skill, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		repo.Create(skill)

		exists, err := repo.Exists("skill-001")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false for non-existent entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		exists, err := repo.Exists("skill-999")

		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, _ := NewFilesystemRepository[core.Skill](tmpDir, "skill")

		_, err := repo.Exists("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "converts to lowercase",
			input:    "Python Programming",
			expected: "python-programming",
		},
		{
			name:     "replaces spaces with hyphens",
			input:    "machine learning",
			expected: "machine-learning",
		},
		{
			name:     "removes special characters",
			input:    "C++ Programming!",
			expected: "c-programming",
		},
		{
			name:     "handles underscores",
			input:    "python_programming",
			expected: "python-programming",
		},
		{
			name:     "removes duplicate hyphens",
			input:    "python---programming",
			expected: "python-programming",
		},
		{
			name:     "trims hyphens from edges",
			input:    "-python-",
			expected: "python",
		},
		{
			name:     "handles long strings",
			input:    "This is a very long title that should be truncated to fifty characters maximum",
			expected: "this-is-a-very-long-title-that-should-be-truncated",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "handles numbers",
			input:    "Python 3.11",
			expected: "python-311",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slugify(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilesystemRepository_WithGoal(t *testing.T) {
	t.Run("works with Goal entity", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo, err := NewFilesystemRepository[core.Goal](tmpDir, "goal")
		require.NoError(t, err)

		targetDate := time.Now().Add(90 * 24 * time.Hour)
		goal, _ := core.NewGoal("goal-001", "Become ML Engineer", core.PriorityHigh)
		goal.SetTargetDate(targetDate)
		goal.Body = "I want to transition into machine learning engineering."

		err = repo.Create(goal)
		require.NoError(t, err)

		retrieved, err := repo.GetByIDWithBody("goal-001")
		require.NoError(t, err)
		assert.Equal(t, "Become ML Engineer", retrieved.Title)
		assert.Contains(t, retrieved.Body, "machine learning")
	})
}
