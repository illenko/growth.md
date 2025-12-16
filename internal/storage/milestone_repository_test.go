package storage

import (
	"testing"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMilestoneRepository(t *testing.T) {
	t.Run("creates repository successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		repo, err := NewMilestoneRepository(tmpDir)

		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("fails with empty path", func(t *testing.T) {
		_, err := NewMilestoneRepository("")

		assert.Error(t, err)
	})
}

func TestMilestoneRepository_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewMilestoneRepository(tmpDir)

	t.Run("creates and retrieves milestone", func(t *testing.T) {
		milestone, _ := core.NewMilestone("milestone-001", "Complete Course", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-001")
		milestone.Body = "Finish the Go programming course"

		err := repo.Create(milestone)
		require.NoError(t, err)

		retrieved, err := repo.GetByIDWithBody("milestone-001")
		require.NoError(t, err)
		assert.Equal(t, "Complete Course", retrieved.Title)
		assert.Contains(t, retrieved.Body, "Go programming")
	})

	t.Run("updates milestone", func(t *testing.T) {
		milestone, _ := repo.GetByID("milestone-001")
		milestone.Achieve("https://cert.example.com/12345")

		err := repo.Update(milestone)
		require.NoError(t, err)

		updated, _ := repo.GetByID("milestone-001")
		assert.Equal(t, core.StatusCompleted, updated.Status)
		assert.NotNil(t, updated.AchievedDate)
		assert.Equal(t, "https://cert.example.com/12345", updated.Proof)
	})

	t.Run("deletes milestone", func(t *testing.T) {
		err := repo.Delete("milestone-001")
		require.NoError(t, err)

		exists, _ := repo.Exists("milestone-001")
		assert.False(t, exists)
	})
}

func TestMilestoneRepository_FindByReferenceID(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewMilestoneRepository(tmpDir)

	milestone1, _ := core.NewMilestone("milestone-001", "Goal Milestone 1", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-001")
	milestone2, _ := core.NewMilestone("milestone-002", "Goal Milestone 2", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-001")
	milestone3, _ := core.NewMilestone("milestone-003", "Path Milestone", core.MilestonePathLevel, core.ReferencePath, "path-001")
	milestone4, _ := core.NewMilestone("milestone-004", "Other Goal Milestone", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-002")

	repo.Create(milestone1)
	repo.Create(milestone2)
	repo.Create(milestone3)
	repo.Create(milestone4)

	t.Run("finds milestones by goal reference", func(t *testing.T) {
		results, err := repo.FindByReferenceID(core.ReferenceGoal, "goal-001")

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, milestone := range results {
			assert.Equal(t, core.ReferenceGoal, milestone.ReferenceType)
			assert.Equal(t, "goal-001", string(milestone.ReferenceID))
		}
	})

	t.Run("finds milestones by path reference", func(t *testing.T) {
		results, err := repo.FindByReferenceID(core.ReferencePath, "path-001")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Path Milestone", results[0].Title)
	})

	t.Run("returns empty for non-existent reference", func(t *testing.T) {
		results, err := repo.FindByReferenceID(core.ReferenceGoal, "goal-999")

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestMilestoneRepository_FindByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewMilestoneRepository(tmpDir)

	active1, _ := core.NewMilestone("milestone-001", "Active 1", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-001")
	active2, _ := core.NewMilestone("milestone-002", "Active 2", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-001")
	completed, _ := core.NewMilestone("milestone-003", "Completed", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-001")
	completed.Achieve("")

	repo.Create(active1)
	repo.Create(active2)
	repo.Create(completed)

	t.Run("finds active milestones", func(t *testing.T) {
		results, err := repo.FindByStatus(core.StatusActive)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, milestone := range results {
			assert.Equal(t, core.StatusActive, milestone.Status)
		}
	})

	t.Run("finds completed milestones", func(t *testing.T) {
		results, err := repo.FindByStatus(core.StatusCompleted)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Completed", results[0].Title)
		assert.Equal(t, core.StatusCompleted, results[0].Status)
	})
}

func TestMilestoneRepository_FindByType(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewMilestoneRepository(tmpDir)

	goal1, _ := core.NewMilestone("milestone-001", "Goal Level 1", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-001")
	goal2, _ := core.NewMilestone("milestone-002", "Goal Level 2", core.MilestoneGoalLevel, core.ReferenceGoal, "goal-001")
	path, _ := core.NewMilestone("milestone-003", "Path Level", core.MilestonePathLevel, core.ReferencePath, "path-001")
	skill, _ := core.NewMilestone("milestone-004", "Skill Level", core.MilestoneSkillLevel, core.ReferenceSkill, "skill-001")

	repo.Create(goal1)
	repo.Create(goal2)
	repo.Create(path)
	repo.Create(skill)

	t.Run("finds goal-level milestones", func(t *testing.T) {
		results, err := repo.FindByType(core.MilestoneGoalLevel)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, milestone := range results {
			assert.Equal(t, core.MilestoneGoalLevel, milestone.Type)
		}
	})

	t.Run("finds path-level milestones", func(t *testing.T) {
		results, err := repo.FindByType(core.MilestonePathLevel)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Path Level", results[0].Title)
	})

	t.Run("finds skill-level milestones", func(t *testing.T) {
		results, err := repo.FindByType(core.MilestoneSkillLevel)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Skill Level", results[0].Title)
	})
}
