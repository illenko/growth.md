package storage

import (
	"testing"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPathRepository(t *testing.T) {
	t.Run("creates repository successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		repo, err := NewPathRepository(tmpDir)

		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("fails with empty path", func(t *testing.T) {
		_, err := NewPathRepository("")

		assert.Error(t, err)
	})
}

func TestPathRepository_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewPathRepository(tmpDir)

	t.Run("creates and retrieves path", func(t *testing.T) {
		path, _ := core.NewLearningPath("path-001", "Backend Development", core.PathTypeManual)
		path.Body = "Comprehensive backend learning path"

		err := repo.Create(path)
		require.NoError(t, err)

		retrieved, err := repo.GetByIDWithBody("path-001")
		require.NoError(t, err)
		assert.Equal(t, "Backend Development", retrieved.Title)
		assert.Contains(t, retrieved.Body, "backend learning")
	})

	t.Run("updates path", func(t *testing.T) {
		path, _ := repo.GetByID("path-001")
		path.UpdateStatus(core.StatusCompleted)
		path.AddTag("engineering")

		err := repo.Update(path)
		require.NoError(t, err)

		updated, _ := repo.GetByID("path-001")
		assert.Equal(t, core.StatusCompleted, updated.Status)
		assert.Contains(t, updated.Tags, "engineering")
	})

	t.Run("deletes path", func(t *testing.T) {
		err := repo.Delete("path-001")
		require.NoError(t, err)

		exists, _ := repo.Exists("path-001")
		assert.False(t, exists)
	})
}

func TestPathRepository_FindByType(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewPathRepository(tmpDir)

	manual1, _ := core.NewLearningPath("path-001", "Manual Path 1", core.PathTypeManual)
	manual2, _ := core.NewLearningPath("path-002", "Manual Path 2", core.PathTypeManual)
	ai, _ := core.NewLearningPath("path-003", "AI Generated Path", core.PathTypeAIGenerated)
	ai.SetGenerationInfo("gpt-4", "backend engineer")

	repo.Create(manual1)
	repo.Create(manual2)
	repo.Create(ai)

	t.Run("finds manual paths", func(t *testing.T) {
		results, err := repo.FindByType(core.PathTypeManual)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, path := range results {
			assert.Equal(t, core.PathTypeManual, path.Type)
		}
	})

	t.Run("finds AI generated paths", func(t *testing.T) {
		results, err := repo.FindByType(core.PathTypeAIGenerated)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "AI Generated Path", results[0].Title)
	})
}

func TestPathRepository_FindByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewPathRepository(tmpDir)

	active1, _ := core.NewLearningPath("path-001", "Active Path 1", core.PathTypeManual)
	active2, _ := core.NewLearningPath("path-002", "Active Path 2", core.PathTypeManual)
	completed, _ := core.NewLearningPath("path-003", "Completed Path", core.PathTypeManual)
	completed.UpdateStatus(core.StatusCompleted)

	repo.Create(active1)
	repo.Create(active2)
	repo.Create(completed)

	t.Run("finds active paths", func(t *testing.T) {
		results, err := repo.FindByStatus(core.StatusActive)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, path := range results {
			assert.Equal(t, core.StatusActive, path.Status)
		}
	})

	t.Run("finds completed paths", func(t *testing.T) {
		results, err := repo.FindByStatus(core.StatusCompleted)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Completed Path", results[0].Title)
	})
}

func TestPathRepository_FindActive(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewPathRepository(tmpDir)

	active, _ := core.NewLearningPath("path-001", "Active Path", core.PathTypeManual)
	completed, _ := core.NewLearningPath("path-002", "Completed Path", core.PathTypeManual)
	completed.UpdateStatus(core.StatusCompleted)

	repo.Create(active)
	repo.Create(completed)

	t.Run("finds only active paths", func(t *testing.T) {
		results, err := repo.FindActive()

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Active Path", results[0].Title)
		assert.Equal(t, core.StatusActive, results[0].Status)
	})
}

func TestPathRepository_FindAIGenerated(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewPathRepository(tmpDir)

	manual, _ := core.NewLearningPath("path-001", "Manual Path", core.PathTypeManual)
	ai1, _ := core.NewLearningPath("path-002", "AI Path 1", core.PathTypeAIGenerated)
	ai2, _ := core.NewLearningPath("path-003", "AI Path 2", core.PathTypeAIGenerated)

	repo.Create(manual)
	repo.Create(ai1)
	repo.Create(ai2)

	t.Run("finds only AI generated paths", func(t *testing.T) {
		results, err := repo.FindAIGenerated()

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, path := range results {
			assert.Equal(t, core.PathTypeAIGenerated, path.Type)
		}
	})
}
