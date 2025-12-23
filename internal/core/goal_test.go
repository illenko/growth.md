package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGoal(t *testing.T) {
	t.Run("creates valid goal", func(t *testing.T) {
		goal, err := NewGoal("goal-001", "Become ML Engineer", PriorityHigh)

		require.NoError(t, err)
		assert.Equal(t, EntityID("goal-001"), goal.ID)
		assert.Equal(t, "Become ML Engineer", goal.Title)
		assert.Equal(t, StatusActive, goal.Status)
		assert.Equal(t, PriorityHigh, goal.Priority)
		assert.Nil(t, goal.TargetDate)
		assert.Empty(t, goal.LearningPaths)
		assert.Empty(t, goal.Milestones)
		assert.Empty(t, goal.Tags)
		assert.False(t, goal.Created.IsZero())
		assert.False(t, goal.Updated.IsZero())
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		_, err := NewGoal("", "Become ML Engineer", PriorityHigh)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("fails with empty title", func(t *testing.T) {
		_, err := NewGoal("goal-001", "", PriorityHigh)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("fails with invalid priority", func(t *testing.T) {
		_, err := NewGoal("goal-001", "Become ML Engineer", Priority("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid goal priority")
	})
}

func TestGoal_Validate(t *testing.T) {
	tests := []struct {
		name    string
		goal    *Goal
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid goal",
			goal: &Goal{
				ID:         "goal-001",
				Title:      "Become ML Engineer",
				Status:     StatusActive,
				Priority:   PriorityHigh,
				Timestamps: NewTimestamps(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			goal: &Goal{
				Title:      "Become ML Engineer",
				Status:     StatusActive,
				Priority:   PriorityHigh,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "ID is required",
		},
		{
			name: "missing title",
			goal: &Goal{
				ID:         "goal-001",
				Status:     StatusActive,
				Priority:   PriorityHigh,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "invalid status",
			goal: &Goal{
				ID:         "goal-001",
				Title:      "Become ML Engineer",
				Status:     Status("invalid"),
				Priority:   PriorityHigh,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "invalid goal status",
		},
		{
			name: "invalid priority",
			goal: &Goal{
				ID:         "goal-001",
				Title:      "Become ML Engineer",
				Status:     StatusActive,
				Priority:   Priority("invalid"),
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "invalid goal priority",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.goal.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGoal_AddLearningPath(t *testing.T) {
	goal, _ := NewGoal("goal-001", "Become ML Engineer", PriorityHigh)

	t.Run("adds learning path", func(t *testing.T) {
		goal.AddLearningPath("path-001")
		assert.Len(t, goal.LearningPaths, 1)
		assert.Contains(t, goal.LearningPaths, EntityID("path-001"))
	})

	t.Run("does not add duplicate path", func(t *testing.T) {
		goal.AddLearningPath("path-001")
		assert.Len(t, goal.LearningPaths, 1)
	})

	t.Run("adds multiple paths", func(t *testing.T) {
		goal.AddLearningPath("path-002")
		goal.AddLearningPath("path-003")
		assert.Len(t, goal.LearningPaths, 3)
	})
}

func TestGoal_RemoveLearningPath(t *testing.T) {
	goal, _ := NewGoal("goal-001", "Become ML Engineer", PriorityHigh)
	goal.AddLearningPath("path-001")
	goal.AddLearningPath("path-002")
	goal.AddLearningPath("path-003")

	t.Run("removes existing path", func(t *testing.T) {
		goal.RemoveLearningPath("path-002")
		assert.Len(t, goal.LearningPaths, 2)
		assert.NotContains(t, goal.LearningPaths, EntityID("path-002"))
	})

	t.Run("does nothing for non-existent path", func(t *testing.T) {
		goal.RemoveLearningPath("path-999")
		assert.Len(t, goal.LearningPaths, 2)
	})
}

func TestGoal_AddMilestone(t *testing.T) {
	goal, _ := NewGoal("goal-001", "Become ML Engineer", PriorityHigh)

	t.Run("adds milestone", func(t *testing.T) {
		goal.AddMilestone("milestone-001")
		assert.Len(t, goal.Milestones, 1)
		assert.Contains(t, goal.Milestones, EntityID("milestone-001"))
	})

	t.Run("does not add duplicate milestone", func(t *testing.T) {
		goal.AddMilestone("milestone-001")
		assert.Len(t, goal.Milestones, 1)
	})
}

func TestGoal_AddTag(t *testing.T) {
	goal, _ := NewGoal("goal-001", "Become ML Engineer", PriorityHigh)

	t.Run("adds tag", func(t *testing.T) {
		goal.AddTag("ml")
		assert.Len(t, goal.Tags, 1)
		assert.Contains(t, goal.Tags, "ml")
	})

	t.Run("normalizes tag to lowercase", func(t *testing.T) {
		goal.AddTag("CAREER-CHANGE")
		assert.Contains(t, goal.Tags, "career-change")
	})

	t.Run("does not add duplicate tag", func(t *testing.T) {
		goal.AddTag("ml")
		assert.Len(t, goal.Tags, 2) // ml, career-change
	})
}

func TestGoal_UpdateStatus(t *testing.T) {
	goal, _ := NewGoal("goal-001", "Become ML Engineer", PriorityHigh)

	t.Run("updates status", func(t *testing.T) {
		err := goal.UpdateStatus(StatusCompleted)
		assert.NoError(t, err)
		assert.Equal(t, StatusCompleted, goal.Status)
	})

	t.Run("fails with invalid status", func(t *testing.T) {
		err := goal.UpdateStatus(Status("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid goal status")
	})
}

func TestGoal_UpdatePriority(t *testing.T) {
	goal, _ := NewGoal("goal-001", "Become ML Engineer", PriorityHigh)

	t.Run("updates priority", func(t *testing.T) {
		err := goal.UpdatePriority(PriorityMedium)
		assert.NoError(t, err)
		assert.Equal(t, PriorityMedium, goal.Priority)
	})

	t.Run("fails with invalid priority", func(t *testing.T) {
		err := goal.UpdatePriority(Priority("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid goal priority")
	})
}

func TestGoal_TargetDate(t *testing.T) {
	goal, _ := NewGoal("goal-001", "Become ML Engineer", PriorityHigh)

	t.Run("sets target date", func(t *testing.T) {
		targetDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		goal.SetTargetDate(targetDate)
		require.NotNil(t, goal.TargetDate)
		assert.Equal(t, targetDate, *goal.TargetDate)
	})

	t.Run("clears target date", func(t *testing.T) {
		goal.ClearTargetDate()
		assert.Nil(t, goal.TargetDate)
	})
}
