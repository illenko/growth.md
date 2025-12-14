package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMilestone(t *testing.T) {
	t.Run("creates valid milestone", func(t *testing.T) {
		milestone, err := NewMilestone(
			"milestone-001",
			"First ML Model",
			MilestoneGoalLevel,
			ReferenceGoal,
			"goal-001",
		)

		require.NoError(t, err)
		assert.Equal(t, EntityID("milestone-001"), milestone.ID)
		assert.Equal(t, "First ML Model", milestone.Title)
		assert.Equal(t, MilestoneGoalLevel, milestone.Type)
		assert.Equal(t, ReferenceGoal, milestone.ReferenceType)
		assert.Equal(t, EntityID("goal-001"), milestone.ReferenceID)
		assert.Equal(t, StatusActive, milestone.Status)
		assert.Nil(t, milestone.AchievedDate)
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		_, err := NewMilestone("", "First ML Model", MilestoneGoalLevel, ReferenceGoal, "goal-001")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("fails with empty title", func(t *testing.T) {
		_, err := NewMilestone("milestone-001", "", MilestoneGoalLevel, ReferenceGoal, "goal-001")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("fails with invalid milestone type", func(t *testing.T) {
		_, err := NewMilestone("milestone-001", "First ML Model", MilestoneType("invalid"), ReferenceGoal, "goal-001")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type is invalid")
	})

	t.Run("fails with invalid reference type", func(t *testing.T) {
		_, err := NewMilestone("milestone-001", "First ML Model", MilestoneGoalLevel, ReferenceType("invalid"), "goal-001")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "referenceType is invalid")
	})

	t.Run("fails with empty referenceId", func(t *testing.T) {
		_, err := NewMilestone("milestone-001", "First ML Model", MilestoneGoalLevel, ReferenceGoal, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "referenceId is required")
	})
}

func TestMilestone_Achieve(t *testing.T) {
	milestone, _ := NewMilestone("milestone-001", "First ML Model", MilestoneGoalLevel, ReferenceGoal, "goal-001")

	t.Run("achieves milestone without proof", func(t *testing.T) {
		milestone.Achieve("")

		assert.Equal(t, StatusCompleted, milestone.Status)
		require.NotNil(t, milestone.AchievedDate)
		assert.True(t, milestone.IsAchieved())
		assert.Empty(t, milestone.Proof)
	})

	t.Run("achieves milestone with proof", func(t *testing.T) {
		milestone2, _ := NewMilestone("milestone-002", "Second ML Model", MilestoneGoalLevel, ReferenceGoal, "goal-001")

		milestone2.Achieve("https://github.com/user/project")

		assert.Equal(t, StatusCompleted, milestone2.Status)
		assert.Equal(t, "https://github.com/user/project", milestone2.Proof)
		assert.True(t, milestone2.IsAchieved())
	})
}

func TestMilestone_TargetDate(t *testing.T) {
	milestone, _ := NewMilestone("milestone-001", "First ML Model", MilestoneGoalLevel, ReferenceGoal, "goal-001")

	t.Run("sets target date", func(t *testing.T) {
		targetDate := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)
		milestone.SetTargetDate(targetDate)

		require.NotNil(t, milestone.TargetDate)
		assert.Equal(t, targetDate, *milestone.TargetDate)
	})

	t.Run("clears target date", func(t *testing.T) {
		milestone.ClearTargetDate()
		assert.Nil(t, milestone.TargetDate)
	})
}

func TestMilestone_SetProof(t *testing.T) {
	milestone, _ := NewMilestone("milestone-001", "First ML Model", MilestoneGoalLevel, ReferenceGoal, "goal-001")

	milestone.SetProof("https://github.com/user/project")

	assert.Equal(t, "https://github.com/user/project", milestone.Proof)
}

func TestMilestone_IsAchieved(t *testing.T) {
	milestone, _ := NewMilestone("milestone-001", "First ML Model", MilestoneGoalLevel, ReferenceGoal, "goal-001")

	t.Run("not achieved initially", func(t *testing.T) {
		assert.False(t, milestone.IsAchieved())
	})

	t.Run("achieved after calling Achieve", func(t *testing.T) {
		milestone.Achieve("")
		assert.True(t, milestone.IsAchieved())
	})
}
