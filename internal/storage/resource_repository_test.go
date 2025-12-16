package storage

import (
	"testing"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResourceRepository(t *testing.T) {
	t.Run("creates repository successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		repo, err := NewResourceRepository(tmpDir)

		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("fails with empty path", func(t *testing.T) {
		_, err := NewResourceRepository("")

		assert.Error(t, err)
	})
}

func TestResourceRepository_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewResourceRepository(tmpDir)

	t.Run("creates and retrieves resource", func(t *testing.T) {
		resource, _ := core.NewResource("resource-001", "Clean Code", core.ResourceBook, "skill-001")
		resource.SetAuthor("Robert C. Martin")
		resource.Body = "Classic book on software craftsmanship"

		err := repo.Create(resource)
		require.NoError(t, err)

		retrieved, err := repo.GetByIDWithBody("resource-001")
		require.NoError(t, err)
		assert.Equal(t, "Clean Code", retrieved.Title)
		assert.Equal(t, "Robert C. Martin", retrieved.Author)
		assert.Contains(t, retrieved.Body, "software craftsmanship")
	})

	t.Run("updates resource", func(t *testing.T) {
		resource, _ := repo.GetByID("resource-001")
		resource.Start()
		resource.AddTag("fundamentals")
		resource.SetEstimatedHours(20)

		err := repo.Update(resource)
		require.NoError(t, err)

		updated, _ := repo.GetByID("resource-001")
		assert.Equal(t, core.ResourceInProgress, updated.Status)
		assert.Contains(t, updated.Tags, "fundamentals")
		assert.Equal(t, 20.0, updated.EstimatedHours)
	})

	t.Run("deletes resource", func(t *testing.T) {
		err := repo.Delete("resource-001")
		require.NoError(t, err)

		exists, _ := repo.Exists("resource-001")
		assert.False(t, exists)
	})
}

func TestResourceRepository_FindByType(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewResourceRepository(tmpDir)

	book1, _ := core.NewResource("resource-001", "Book 1", core.ResourceBook, "skill-001")
	book2, _ := core.NewResource("resource-002", "Book 2", core.ResourceBook, "skill-001")
	course, _ := core.NewResource("resource-003", "Course 1", core.ResourceCourse, "skill-001")
	video, _ := core.NewResource("resource-004", "Video 1", core.ResourceVideo, "skill-002")

	repo.Create(book1)
	repo.Create(book2)
	repo.Create(course)
	repo.Create(video)

	t.Run("finds books", func(t *testing.T) {
		results, err := repo.FindByType(core.ResourceBook)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, resource := range results {
			assert.Equal(t, core.ResourceBook, resource.Type)
		}
	})

	t.Run("finds courses", func(t *testing.T) {
		results, err := repo.FindByType(core.ResourceCourse)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Course 1", results[0].Title)
	})
}

func TestResourceRepository_FindBySkillID(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewResourceRepository(tmpDir)

	resource1, _ := core.NewResource("resource-001", "Resource 1", core.ResourceBook, "skill-001")
	resource2, _ := core.NewResource("resource-002", "Resource 2", core.ResourceCourse, "skill-001")
	resource3, _ := core.NewResource("resource-003", "Resource 3", core.ResourceVideo, "skill-002")

	repo.Create(resource1)
	repo.Create(resource2)
	repo.Create(resource3)

	t.Run("finds resources by skill ID", func(t *testing.T) {
		results, err := repo.FindBySkillID("skill-001")

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, resource := range results {
			assert.Equal(t, "skill-001", string(resource.SkillID))
		}
	})

	t.Run("finds resources for different skill", func(t *testing.T) {
		results, err := repo.FindBySkillID("skill-002")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Resource 3", results[0].Title)
	})

	t.Run("returns empty for non-existent skill", func(t *testing.T) {
		results, err := repo.FindBySkillID("skill-999")

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestResourceRepository_FindByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewResourceRepository(tmpDir)

	notStarted, _ := core.NewResource("resource-001", "Not Started", core.ResourceBook, "skill-001")
	inProgress, _ := core.NewResource("resource-002", "In Progress", core.ResourceCourse, "skill-001")
	inProgress.Start()
	completed, _ := core.NewResource("resource-003", "Completed", core.ResourceVideo, "skill-001")
	completed.Complete()

	repo.Create(notStarted)
	repo.Create(inProgress)
	repo.Create(completed)

	t.Run("finds not started resources", func(t *testing.T) {
		results, err := repo.FindByStatus(core.ResourceNotStarted)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Not Started", results[0].Title)
	})

	t.Run("finds in progress resources", func(t *testing.T) {
		results, err := repo.FindByStatus(core.ResourceInProgress)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "In Progress", results[0].Title)
	})

	t.Run("finds completed resources", func(t *testing.T) {
		results, err := repo.FindByStatus(core.ResourceCompleted)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Completed", results[0].Title)
	})
}
