package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPhase(t *testing.T) {
	t.Run("creates valid phase", func(t *testing.T) {
		phase, err := NewPhase("phase-001", "path-001", "Foundations", 1)

		require.NoError(t, err)
		assert.Equal(t, EntityID("phase-001"), phase.ID)
		assert.Equal(t, EntityID("path-001"), phase.PathID)
		assert.Equal(t, "Foundations", phase.Title)
		assert.Equal(t, 1, phase.Order)
		assert.Empty(t, phase.RequiredSkills)
		assert.Empty(t, phase.Milestones)
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		_, err := NewPhase("", "path-001", "Foundations", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("fails with empty pathID", func(t *testing.T) {
		_, err := NewPhase("phase-001", "", "Foundations", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pathID is required")
	})

	t.Run("fails with negative order", func(t *testing.T) {
		_, err := NewPhase("phase-001", "path-001", "Foundations", -1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order must be non-negative")
	})
}

func TestPhase_Validate(t *testing.T) {
	tests := []struct {
		name    string
		phase   *Phase
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid phase",
			phase: &Phase{
				ID:         "phase-001",
				PathID:     "path-001",
				Title:      "Foundations",
				Order:      1,
				Timestamps: NewTimestamps(),
			},
			wantErr: false,
		},
		{
			name: "missing title",
			phase: &Phase{
				ID:         "phase-001",
				PathID:     "path-001",
				Order:      1,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "invalid skill requirement - empty skillId",
			phase: &Phase{
				ID:     "phase-001",
				PathID: "path-001",
				Title:  "Foundations",
				Order:  1,
				RequiredSkills: []SkillRequirement{
					{SkillID: "", TargetLevel: LevelBeginner},
				},
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "must have a skillId",
		},
		{
			name: "invalid skill requirement - invalid level",
			phase: &Phase{
				ID:     "phase-001",
				PathID: "path-001",
				Title:  "Foundations",
				Order:  1,
				RequiredSkills: []SkillRequirement{
					{SkillID: "skill-001", TargetLevel: ProficiencyLevel("invalid")},
				},
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "invalid target level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.phase.Validate()
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

func TestPhase_AddSkillRequirement(t *testing.T) {
	phase, _ := NewPhase("phase-001", "path-001", "Foundations", 1)

	t.Run("adds skill requirement", func(t *testing.T) {
		err := phase.AddSkillRequirement("skill-001", LevelIntermediate)
		require.NoError(t, err)
		assert.Len(t, phase.RequiredSkills, 1)
		assert.Equal(t, EntityID("skill-001"), phase.RequiredSkills[0].SkillID)
		assert.Equal(t, LevelIntermediate, phase.RequiredSkills[0].TargetLevel)
	})

	t.Run("adds multiple requirements", func(t *testing.T) {
		phase.AddSkillRequirement("skill-002", LevelBeginner)
		phase.AddSkillRequirement("skill-003", LevelAdvanced)
		assert.Len(t, phase.RequiredSkills, 3)
	})

	t.Run("fails with invalid level", func(t *testing.T) {
		err := phase.AddSkillRequirement("skill-004", ProficiencyLevel("invalid"))
		assert.Error(t, err)
	})
}

func TestPhase_RemoveSkillRequirement(t *testing.T) {
	phase, _ := NewPhase("phase-001", "path-001", "Foundations", 1)
	phase.AddSkillRequirement("skill-001", LevelIntermediate)
	phase.AddSkillRequirement("skill-002", LevelBeginner)

	t.Run("removes skill requirement", func(t *testing.T) {
		phase.RemoveSkillRequirement("skill-001")
		assert.Len(t, phase.RequiredSkills, 1)
		assert.Equal(t, EntityID("skill-002"), phase.RequiredSkills[0].SkillID)
	})

	t.Run("does nothing for non-existent skill", func(t *testing.T) {
		phase.RemoveSkillRequirement("skill-999")
		assert.Len(t, phase.RequiredSkills, 1)
	})
}

func TestPhase_AddMilestone(t *testing.T) {
	phase, _ := NewPhase("phase-001", "path-001", "Foundations", 1)

	t.Run("adds milestone", func(t *testing.T) {
		phase.AddMilestone("milestone-001")
		assert.Len(t, phase.Milestones, 1)
		assert.Contains(t, phase.Milestones, EntityID("milestone-001"))
	})

	t.Run("does not add duplicate", func(t *testing.T) {
		phase.AddMilestone("milestone-001")
		assert.Len(t, phase.Milestones, 1)
	})
}
