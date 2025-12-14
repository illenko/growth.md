package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResource(t *testing.T) {
	t.Run("creates valid resource", func(t *testing.T) {
		resource, err := NewResource("resource-001", "Fluent Python", ResourceBook, "skill-002")

		require.NoError(t, err)
		assert.Equal(t, EntityID("resource-001"), resource.ID)
		assert.Equal(t, "Fluent Python", resource.Title)
		assert.Equal(t, ResourceBook, resource.Type)
		assert.Equal(t, EntityID("skill-002"), resource.SkillID)
		assert.Equal(t, ResourceNotStarted, resource.Status)
		assert.Empty(t, resource.Tags)
	})

	t.Run("fails with empty ID", func(t *testing.T) {
		_, err := NewResource("", "Fluent Python", ResourceBook, "skill-002")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("fails with empty title", func(t *testing.T) {
		_, err := NewResource("resource-001", "", ResourceBook, "skill-002")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("fails with invalid type", func(t *testing.T) {
		_, err := NewResource("resource-001", "Fluent Python", ResourceType("invalid"), "skill-002")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type is invalid")
	})

	t.Run("fails with empty skillId", func(t *testing.T) {
		_, err := NewResource("resource-001", "Fluent Python", ResourceBook, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "skillId is required")
	})
}

func TestResource_Validate(t *testing.T) {
	tests := []struct {
		name     string
		resource *Resource
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid resource",
			resource: &Resource{
				ID:             "resource-001",
				Title:          "Fluent Python",
				Type:           ResourceBook,
				SkillID:        "skill-002",
				Status:         ResourceNotStarted,
				EstimatedHours: 40,
				Timestamps:     NewTimestamps(),
			},
			wantErr: false,
		},
		{
			name: "negative estimated hours",
			resource: &Resource{
				ID:             "resource-001",
				Title:          "Fluent Python",
				Type:           ResourceBook,
				SkillID:        "skill-002",
				Status:         ResourceNotStarted,
				EstimatedHours: -10,
				Timestamps:     NewTimestamps(),
			},
			wantErr: true,
			errMsg:  "cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.resource.Validate()
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

func TestResource_UpdateStatus(t *testing.T) {
	resource, _ := NewResource("resource-001", "Fluent Python", ResourceBook, "skill-002")

	t.Run("updates status", func(t *testing.T) {
		err := resource.UpdateStatus(ResourceInProgress)
		assert.NoError(t, err)
		assert.Equal(t, ResourceInProgress, resource.Status)
	})

	t.Run("fails with invalid status", func(t *testing.T) {
		err := resource.UpdateStatus(ResourceStatus("invalid"))
		assert.Error(t, err)
	})
}

func TestResource_Start(t *testing.T) {
	resource, _ := NewResource("resource-001", "Fluent Python", ResourceBook, "skill-002")

	resource.Start()

	assert.Equal(t, ResourceInProgress, resource.Status)
}

func TestResource_Complete(t *testing.T) {
	resource, _ := NewResource("resource-001", "Fluent Python", ResourceBook, "skill-002")

	resource.Complete()

	assert.Equal(t, ResourceCompleted, resource.Status)
}

func TestResource_AddTag(t *testing.T) {
	resource, _ := NewResource("resource-001", "Fluent Python", ResourceBook, "skill-002")

	t.Run("adds tag", func(t *testing.T) {
		resource.AddTag("python")
		assert.Len(t, resource.Tags, 1)
		assert.Contains(t, resource.Tags, "python")
	})

	t.Run("normalizes to lowercase", func(t *testing.T) {
		resource.AddTag("ADVANCED")
		assert.Contains(t, resource.Tags, "advanced")
	})

	t.Run("does not add duplicate", func(t *testing.T) {
		resource.AddTag("python")
		assert.Len(t, resource.Tags, 2) // python, advanced
	})
}

func TestResource_SetURL(t *testing.T) {
	resource, _ := NewResource("resource-001", "Fluent Python", ResourceBook, "skill-002")

	resource.SetURL("https://example.com/book")

	assert.Equal(t, "https://example.com/book", resource.URL)
}

func TestResource_SetAuthor(t *testing.T) {
	resource, _ := NewResource("resource-001", "Fluent Python", ResourceBook, "skill-002")

	resource.SetAuthor("Luciano Ramalho")

	assert.Equal(t, "Luciano Ramalho", resource.Author)
}

func TestResource_SetEstimatedHours(t *testing.T) {
	resource, _ := NewResource("resource-001", "Fluent Python", ResourceBook, "skill-002")

	t.Run("sets valid hours", func(t *testing.T) {
		err := resource.SetEstimatedHours(40)
		assert.NoError(t, err)
		assert.Equal(t, 40.0, resource.EstimatedHours)
	})

	t.Run("fails with negative hours", func(t *testing.T) {
		err := resource.SetEstimatedHours(-10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")
	})
}
