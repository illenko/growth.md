package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLearningPath(t *testing.T) {
	t.Run("creates valid AI-generated path", func(t *testing.T) {
		path, err := NewLearningPath("path-001", "ML Engineer Track", PathTypeAIGenerated)

		require.NoError(t, err)
		assert.Equal(t, EntityID("path-001"), path.ID)
		assert.Equal(t, "ML Engineer Track", path.Title)
		assert.Equal(t, PathTypeAIGenerated, path.Type)
		assert.Equal(t, StatusActive, path.Status)
		assert.Empty(t, path.Phases)
		assert.Empty(t, path.Tags)
		assert.False(t, path.Created.IsZero())
	})

	t.Run("creates valid manual path", func(t *testing.T) {
		path, err := NewLearningPath("path-002", "Custom Path", PathTypeManual)

		require.NoError(t, err)
		assert.Equal(t, PathTypeManual, path.Type)
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		_, err := NewLearningPath("", "ML Engineer Track", PathTypeAIGenerated)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("fails with empty title", func(t *testing.T) {
		_, err := NewLearningPath("path-001", "", PathTypeAIGenerated)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("fails with invalid type", func(t *testing.T) {
		_, err := NewLearningPath("path-001", "ML Track", PathType("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid path type")
	})
}

func TestLearningPath_Validate(t *testing.T) {
	tests := []struct {
		name    string
		path    *LearningPath
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid path",
			path: &LearningPath{
				ID:         "path-001",
				Title:      "ML Engineer Track",
				Type:       PathTypeAIGenerated,
				Status:     StatusActive,
				Timestamps: NewTimestamps(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			path: &LearningPath{
				Title:      "ML Engineer Track",
				Type:       PathTypeAIGenerated,
				Status:     StatusActive,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "ID is required",
		},
		{
			name: "invalid type",
			path: &LearningPath{
				ID:         "path-001",
				Title:      "ML Engineer Track",
				Type:       PathType("invalid"),
				Status:     StatusActive,
				Timestamps: NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "invalid path type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.path.Validate()
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

func TestLearningPath_AddPhase(t *testing.T) {
	path, _ := NewLearningPath("path-001", "ML Engineer Track", PathTypeAIGenerated)

	t.Run("adds phase", func(t *testing.T) {
		path.AddPhase("phase-001")
		assert.Len(t, path.Phases, 1)
		assert.Contains(t, path.Phases, EntityID("phase-001"))
	})

	t.Run("maintains order", func(t *testing.T) {
		path.AddPhase("phase-002")
		path.AddPhase("phase-003")
		assert.Len(t, path.Phases, 3)
		assert.Equal(t, EntityID("phase-001"), path.Phases[0])
		assert.Equal(t, EntityID("phase-002"), path.Phases[1])
		assert.Equal(t, EntityID("phase-003"), path.Phases[2])
	})

	t.Run("does not add duplicate", func(t *testing.T) {
		path.AddPhase("phase-002")
		assert.Len(t, path.Phases, 3)
	})
}

func TestLearningPath_RemovePhase(t *testing.T) {
	path, _ := NewLearningPath("path-001", "ML Engineer Track", PathTypeAIGenerated)
	path.AddPhase("phase-001")
	path.AddPhase("phase-002")
	path.AddPhase("phase-003")

	t.Run("removes phase", func(t *testing.T) {
		path.RemovePhase("phase-002")
		assert.Len(t, path.Phases, 2)
		assert.NotContains(t, path.Phases, EntityID("phase-002"))
	})
}

func TestLearningPath_SetGenerationInfo(t *testing.T) {
	path, _ := NewLearningPath("path-001", "ML Engineer Track", PathTypeAIGenerated)

	path.SetGenerationInfo("claude-opus-4-5", "Background: 5 years backend")

	assert.Equal(t, "claude-opus-4-5", path.GeneratedBy)
	assert.Equal(t, "Background: 5 years backend", path.GenerationContext)
}

func TestLearningPath_UpdateStatus(t *testing.T) {
	path, _ := NewLearningPath("path-001", "ML Engineer Track", PathTypeAIGenerated)

	t.Run("updates status", func(t *testing.T) {
		err := path.UpdateStatus(StatusCompleted)
		assert.NoError(t, err)
		assert.Equal(t, StatusCompleted, path.Status)
	})

	t.Run("fails with invalid status", func(t *testing.T) {
		err := path.UpdateStatus(Status("invalid"))
		assert.Error(t, err)
	})
}
