package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSkill(t *testing.T) {
	t.Run("creates valid skill", func(t *testing.T) {
		skill, err := NewSkill("skill-001", "Python", "programming", LevelIntermediate)

		require.NoError(t, err)
		assert.Equal(t, EntityID("skill-001"), skill.ID)
		assert.Equal(t, "Python", skill.Title)
		assert.Equal(t, "programming", skill.Category)
		assert.Equal(t, LevelIntermediate, skill.Level)
		assert.Equal(t, SkillNotStarted, skill.Status)
		assert.Empty(t, skill.Resources)
		assert.Empty(t, skill.Tags)
		assert.False(t, skill.Created.IsZero())
		assert.False(t, skill.Updated.IsZero())
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		_, err := NewSkill("", "Python", "programming", LevelIntermediate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("fails with empty title", func(t *testing.T) {
		_, err := NewSkill("skill-001", "", "programming", LevelIntermediate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("fails with empty category", func(t *testing.T) {
		_, err := NewSkill("skill-001", "Python", "", LevelIntermediate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category is required")
	})

	t.Run("fails with invalid level", func(t *testing.T) {
		_, err := NewSkill("skill-001", "Python", "programming", ProficiencyLevel("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid skill level")
	})
}

func TestSkill_Validate(t *testing.T) {
	tests := []struct {
		name    string
		skill   *Skill
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid skill",
			skill: &Skill{
				ID:         "skill-001",
				Title:      "Python",
				Category:   "programming",
				Level:      LevelIntermediate,
				Status:     SkillLearning,
				Timestamps: NewTimestamps(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			skill: &Skill{
				Title:      "Python",
				Category:   "programming",
				Level:      LevelIntermediate,
				Status:     SkillLearning,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "ID is required",
		},
		{
			name: "missing title",
			skill: &Skill{
				ID:         "skill-001",
				Category:   "programming",
				Level:      LevelIntermediate,
				Status:     SkillLearning,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "invalid level",
			skill: &Skill{
				ID:         "skill-001",
				Title:      "Python",
				Category:   "programming",
				Level:      ProficiencyLevel("invalid"),
				Status:     SkillLearning,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "invalid skill level",
		},
		{
			name: "invalid status",
			skill: &Skill{
				ID:         "skill-001",
				Title:      "Python",
				Category:   "programming",
				Level:      LevelIntermediate,
				Status:     SkillStatus("invalid"),
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "invalid skill status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.skill.Validate()
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

func TestSkill_AddResource(t *testing.T) {
	skill, _ := NewSkill("skill-001", "Python", "programming", LevelIntermediate)

	t.Run("adds resource", func(t *testing.T) {
		skill.AddResource("resource-001")
		assert.Len(t, skill.Resources, 1)
		assert.Contains(t, skill.Resources, EntityID("resource-001"))
	})

	t.Run("does not add duplicate resource", func(t *testing.T) {
		skill.AddResource("resource-001")
		assert.Len(t, skill.Resources, 1)
	})

	t.Run("adds multiple resources", func(t *testing.T) {
		skill.AddResource("resource-002")
		skill.AddResource("resource-003")
		assert.Len(t, skill.Resources, 3)
	})
}

func TestSkill_RemoveResource(t *testing.T) {
	skill, _ := NewSkill("skill-001", "Python", "programming", LevelIntermediate)
	skill.AddResource("resource-001")
	skill.AddResource("resource-002")
	skill.AddResource("resource-003")

	t.Run("removes existing resource", func(t *testing.T) {
		skill.RemoveResource("resource-002")
		assert.Len(t, skill.Resources, 2)
		assert.NotContains(t, skill.Resources, EntityID("resource-002"))
	})

	t.Run("does nothing for non-existent resource", func(t *testing.T) {
		skill.RemoveResource("resource-999")
		assert.Len(t, skill.Resources, 2)
	})
}

func TestSkill_AddTag(t *testing.T) {
	skill, _ := NewSkill("skill-001", "Python", "programming", LevelIntermediate)

	t.Run("adds tag", func(t *testing.T) {
		skill.AddTag("ml")
		assert.Len(t, skill.Tags, 1)
		assert.Contains(t, skill.Tags, "ml")
	})

	t.Run("normalizes tag to lowercase", func(t *testing.T) {
		skill.AddTag("BACKEND")
		assert.Contains(t, skill.Tags, "backend")
	})

	t.Run("trims whitespace", func(t *testing.T) {
		skill.AddTag("  data-science  ")
		assert.Contains(t, skill.Tags, "data-science")
	})

	t.Run("does not add duplicate tag", func(t *testing.T) {
		skill.AddTag("ml")
		skill.AddTag("ML")
		assert.Len(t, skill.Tags, 3) // ml, backend, data-science
	})

	t.Run("ignores empty tag", func(t *testing.T) {
		skill.AddTag("")
		skill.AddTag("   ")
		assert.Len(t, skill.Tags, 3)
	})
}

func TestSkill_UpdateLevel(t *testing.T) {
	skill, _ := NewSkill("skill-001", "Python", "programming", LevelBeginner)
	originalUpdated := skill.Updated

	t.Run("updates level", func(t *testing.T) {
		err := skill.UpdateLevel(LevelIntermediate)
		assert.NoError(t, err)
		assert.Equal(t, LevelIntermediate, skill.Level)
		assert.True(t, skill.Updated.After(originalUpdated) || skill.Updated.Equal(originalUpdated))
	})

	t.Run("fails with invalid level", func(t *testing.T) {
		err := skill.UpdateLevel(ProficiencyLevel("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid proficiency level")
	})
}

func TestSkill_UpdateStatus(t *testing.T) {
	skill, _ := NewSkill("skill-001", "Python", "programming", LevelIntermediate)

	t.Run("updates status", func(t *testing.T) {
		err := skill.UpdateStatus(SkillLearning)
		assert.NoError(t, err)
		assert.Equal(t, SkillLearning, skill.Status)
	})

	t.Run("fails with invalid status", func(t *testing.T) {
		err := skill.UpdateStatus(SkillStatus("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid skill status")
	})
}
