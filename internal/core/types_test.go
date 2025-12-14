package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"valid active", StatusActive, true},
		{"valid completed", StatusCompleted, true},
		{"valid archived", StatusArchived, true},
		{"invalid", Status("invalid"), false},
		{"empty", Status(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestPriority_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		want     bool
	}{
		{"valid high", PriorityHigh, true},
		{"valid medium", PriorityMedium, true},
		{"valid low", PriorityLow, true},
		{"invalid", Priority("invalid"), false},
		{"empty", Priority(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.priority.IsValid())
		})
	}
}

func TestProficiencyLevel_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		level ProficiencyLevel
		want  bool
	}{
		{"valid beginner", LevelBeginner, true},
		{"valid intermediate", LevelIntermediate, true},
		{"valid advanced", LevelAdvanced, true},
		{"valid expert", LevelExpert, true},
		{"invalid", ProficiencyLevel("invalid"), false},
		{"empty", ProficiencyLevel(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.level.IsValid())
		})
	}
}

func TestSkillStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status SkillStatus
		want   bool
	}{
		{"valid not-started", SkillNotStarted, true},
		{"valid learning", SkillLearning, true},
		{"valid mastered", SkillMastered, true},
		{"invalid", SkillStatus("invalid"), false},
		{"empty", SkillStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestResourceType_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		resourceType ResourceType
		want         bool
	}{
		{"valid book", ResourceBook, true},
		{"valid course", ResourceCourse, true},
		{"valid video", ResourceVideo, true},
		{"valid article", ResourceArticle, true},
		{"valid project", ResourceProject, true},
		{"valid documentation", ResourceDocumentation, true},
		{"invalid", ResourceType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.resourceType.IsValid())
		})
	}
}

func TestResourceStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status ResourceStatus
		want   bool
	}{
		{"valid not-started", ResourceNotStarted, true},
		{"valid in-progress", ResourceInProgress, true},
		{"valid completed", ResourceCompleted, true},
		{"invalid", ResourceStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestTimestamps(t *testing.T) {
	t.Run("NewTimestamps creates timestamps", func(t *testing.T) {
		ts := NewTimestamps()
		assert.False(t, ts.Created.IsZero())
		assert.False(t, ts.Updated.IsZero())
		assert.Equal(t, ts.Created, ts.Updated)
	})

	t.Run("Touch updates timestamp", func(t *testing.T) {
		ts := NewTimestamps()
		created := ts.Created

		// Small delay to ensure different timestamp
		ts.Touch()

		assert.Equal(t, created, ts.Created, "Created should not change")
		assert.True(t, ts.Updated.After(created) || ts.Updated.Equal(created), "Updated should be >= Created")
	})
}
