package storage

import (
	"testing"
	"time"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProgressLogRepository(t *testing.T) {
	t.Run("creates repository successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		repo, err := NewProgressLogRepository(tmpDir)

		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("fails with empty path", func(t *testing.T) {
		_, err := NewProgressLogRepository("")

		assert.Error(t, err)
	})
}

func TestProgressLogRepository_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewProgressLogRepository(tmpDir)

	week1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("creates and retrieves progress log", func(t *testing.T) {
		log, _ := core.NewProgressLog("progress-001", week1)
		log.AddSkillWorked("skill-001")
		log.SetHoursInvested(15.5)
		log.Body = "Great week of learning"

		err := repo.Create(log)
		require.NoError(t, err)

		retrieved, err := repo.GetByIDWithBody("progress-001")
		require.NoError(t, err)
		assert.Equal(t, 15.5, retrieved.HoursInvested)
		assert.Contains(t, retrieved.SkillsWorked, core.EntityID("skill-001"))
		assert.Contains(t, retrieved.Body, "Great week")
	})

	t.Run("updates progress log", func(t *testing.T) {
		log, _ := repo.GetByID("progress-001")
		log.AddResourceUsed("resource-001")
		log.SetMood("motivated")

		err := repo.Update(log)
		require.NoError(t, err)

		updated, _ := repo.GetByID("progress-001")
		assert.Contains(t, updated.ResourcesUsed, core.EntityID("resource-001"))
		assert.Equal(t, "motivated", updated.Mood)
	})

	t.Run("deletes progress log", func(t *testing.T) {
		err := repo.Delete("progress-001")
		require.NoError(t, err)

		exists, _ := repo.Exists("progress-001")
		assert.False(t, exists)
	})
}

func TestProgressLogRepository_FindByDateRange(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewProgressLogRepository(tmpDir)

	week1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	week2 := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	week3 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	week4 := time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC)

	log1, _ := core.NewProgressLog("progress-001", week1)
	log2, _ := core.NewProgressLog("progress-002", week2)
	log3, _ := core.NewProgressLog("progress-003", week3)
	log4, _ := core.NewProgressLog("progress-004", week4)

	repo.Create(log1)
	repo.Create(log2)
	repo.Create(log3)
	repo.Create(log4)

	t.Run("finds logs in date range", func(t *testing.T) {
		start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		results, err := repo.FindByDateRange(start, end)

		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("returns logs in descending order", func(t *testing.T) {
		start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC)

		results, err := repo.FindByDateRange(start, end)

		require.NoError(t, err)
		assert.Len(t, results, 4)
		assert.Equal(t, "progress-004", string(results[0].ID))
		assert.Equal(t, "progress-001", string(results[3].ID))
	})
}

func TestProgressLogRepository_FindRecent(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewProgressLogRepository(tmpDir)

	week1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	week2 := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	week3 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	log1, _ := core.NewProgressLog("progress-001", week1)
	log2, _ := core.NewProgressLog("progress-002", week2)
	log3, _ := core.NewProgressLog("progress-003", week3)

	repo.Create(log1)
	repo.Create(log2)
	repo.Create(log3)

	t.Run("finds recent logs with limit", func(t *testing.T) {
		results, err := repo.FindRecent(2)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "progress-003", string(results[0].ID))
		assert.Equal(t, "progress-002", string(results[1].ID))
	})

	t.Run("returns all logs when limit is larger", func(t *testing.T) {
		results, err := repo.FindRecent(10)

		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("returns all logs when limit is 0", func(t *testing.T) {
		results, err := repo.FindRecent(0)

		require.NoError(t, err)
		assert.Len(t, results, 3)
	})
}

func TestProgressLogRepository_FindBySkillID(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewProgressLogRepository(tmpDir)

	week1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	week2 := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	week3 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	log1, _ := core.NewProgressLog("progress-001", week1)
	log1.AddSkillWorked("skill-001")
	log1.AddSkillWorked("skill-002")

	log2, _ := core.NewProgressLog("progress-002", week2)
	log2.AddSkillWorked("skill-001")

	log3, _ := core.NewProgressLog("progress-003", week3)
	log3.AddSkillWorked("skill-003")

	repo.Create(log1)
	repo.Create(log2)
	repo.Create(log3)

	t.Run("finds logs by skill ID", func(t *testing.T) {
		results, err := repo.FindBySkillID("skill-001")

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, log := range results {
			assert.Contains(t, log.SkillsWorked, core.EntityID("skill-001"))
		}
	})

	t.Run("returns logs in descending order", func(t *testing.T) {
		results, err := repo.FindBySkillID("skill-001")

		require.NoError(t, err)
		assert.Equal(t, "progress-002", string(results[0].ID))
		assert.Equal(t, "progress-001", string(results[1].ID))
	})

	t.Run("returns empty for non-existent skill", func(t *testing.T) {
		results, err := repo.FindBySkillID("skill-999")

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestProgressLogRepository_FindByResourceID(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewProgressLogRepository(tmpDir)

	week1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	week2 := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	week3 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	log1, _ := core.NewProgressLog("progress-001", week1)
	log1.AddResourceUsed("resource-001")
	log1.AddResourceUsed("resource-002")

	log2, _ := core.NewProgressLog("progress-002", week2)
	log2.AddResourceUsed("resource-001")

	log3, _ := core.NewProgressLog("progress-003", week3)
	log3.AddResourceUsed("resource-003")

	repo.Create(log1)
	repo.Create(log2)
	repo.Create(log3)

	t.Run("finds logs by resource ID", func(t *testing.T) {
		results, err := repo.FindByResourceID("resource-001")

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, log := range results {
			assert.Contains(t, log.ResourcesUsed, core.EntityID("resource-001"))
		}
	})

	t.Run("returns logs in descending order", func(t *testing.T) {
		results, err := repo.FindByResourceID("resource-001")

		require.NoError(t, err)
		assert.Equal(t, "progress-002", string(results[0].ID))
		assert.Equal(t, "progress-001", string(results[1].ID))
	})

	t.Run("returns empty for non-existent resource", func(t *testing.T) {
		results, err := repo.FindByResourceID("resource-999")

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}
