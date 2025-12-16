package storage

import (
	"testing"
	"time"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGoalRepository(t *testing.T) {
	t.Run("creates repository successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		repo, err := NewGoalRepository(tmpDir)

		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("fails with empty path", func(t *testing.T) {
		_, err := NewGoalRepository("")

		assert.Error(t, err)
	})
}

func TestGoalRepository_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewGoalRepository(tmpDir)

	t.Run("creates and retrieves goal", func(t *testing.T) {
		goal, _ := core.NewGoal("goal-001", "Become Senior Engineer", core.PriorityHigh)
		goal.Body = "Career advancement goal"

		err := repo.Create(goal)
		require.NoError(t, err)

		retrieved, err := repo.GetByIDWithBody("goal-001")
		require.NoError(t, err)
		assert.Equal(t, "Become Senior Engineer", retrieved.Title)
		assert.Contains(t, retrieved.Body, "Career advancement")
	})

	t.Run("updates goal", func(t *testing.T) {
		goal, _ := repo.GetByID("goal-001")
		goal.UpdateStatus(core.StatusCompleted)
		goal.AddTag("career")

		err := repo.Update(goal)
		require.NoError(t, err)

		updated, _ := repo.GetByID("goal-001")
		assert.Equal(t, core.StatusCompleted, updated.Status)
		assert.Contains(t, updated.Tags, "career")
	})

	t.Run("deletes goal", func(t *testing.T) {
		err := repo.Delete("goal-001")
		require.NoError(t, err)

		exists, _ := repo.Exists("goal-001")
		assert.False(t, exists)
	})
}

func TestGoalRepository_FindByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewGoalRepository(tmpDir)

	goal1, _ := core.NewGoal("goal-001", "Active Goal", core.PriorityHigh)
	goal2, _ := core.NewGoal("goal-002", "Completed Goal", core.PriorityMedium)
	goal2.UpdateStatus(core.StatusCompleted)
	goal3, _ := core.NewGoal("goal-003", "Another Active", core.PriorityLow)

	repo.Create(goal1)
	repo.Create(goal2)
	repo.Create(goal3)

	t.Run("finds active goals", func(t *testing.T) {
		results, err := repo.FindByStatus(core.StatusActive)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, goal := range results {
			assert.Equal(t, core.StatusActive, goal.Status)
		}
	})

	t.Run("finds completed goals", func(t *testing.T) {
		results, err := repo.FindByStatus(core.StatusCompleted)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Completed Goal", results[0].Title)
	})
}

func TestGoalRepository_FindByPriority(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewGoalRepository(tmpDir)

	high1, _ := core.NewGoal("goal-001", "High Priority 1", core.PriorityHigh)
	high2, _ := core.NewGoal("goal-002", "High Priority 2", core.PriorityHigh)
	medium, _ := core.NewGoal("goal-003", "Medium Priority", core.PriorityMedium)

	repo.Create(high1)
	repo.Create(high2)
	repo.Create(medium)

	t.Run("finds high priority goals", func(t *testing.T) {
		results, err := repo.FindByPriority(core.PriorityHigh)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, goal := range results {
			assert.Equal(t, core.PriorityHigh, goal.Priority)
		}
	})

	t.Run("finds medium priority goals", func(t *testing.T) {
		results, err := repo.FindByPriority(core.PriorityMedium)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Medium Priority", results[0].Title)
	})
}

func TestGoalRepository_FindActive(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewGoalRepository(tmpDir)

	active, _ := core.NewGoal("goal-001", "Active Goal", core.PriorityHigh)
	completed, _ := core.NewGoal("goal-002", "Completed Goal", core.PriorityMedium)
	completed.UpdateStatus(core.StatusCompleted)

	repo.Create(active)
	repo.Create(completed)

	t.Run("finds only active goals", func(t *testing.T) {
		results, err := repo.FindActive()

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Active Goal", results[0].Title)
		assert.Equal(t, core.StatusActive, results[0].Status)
	})
}

func TestGoalRepository_FindByTargetDateRange(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewGoalRepository(tmpDir)

	jan15 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	feb15 := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
	mar15 := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

	goal1, _ := core.NewGoal("goal-001", "Q1 Goal", core.PriorityHigh)
	goal1.SetTargetDate(jan15)

	goal2, _ := core.NewGoal("goal-002", "Q1-Q2 Goal", core.PriorityMedium)
	goal2.SetTargetDate(feb15)

	goal3, _ := core.NewGoal("goal-003", "Q2 Goal", core.PriorityLow)
	goal3.SetTargetDate(mar15)

	goal4, _ := core.NewGoal("goal-004", "No Date Goal", core.PriorityMedium)

	repo.Create(goal1)
	repo.Create(goal2)
	repo.Create(goal3)
	repo.Create(goal4)

	t.Run("finds goals in date range", func(t *testing.T) {
		start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)

		results, err := repo.FindByTargetDateRange(start, end)

		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("excludes goals outside range", func(t *testing.T) {
		start := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)

		results, err := repo.FindByTargetDateRange(start, end)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Q2 Goal", results[0].Title)
	})
}
