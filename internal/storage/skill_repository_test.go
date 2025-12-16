package storage

import (
	"testing"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSkillRepository(t *testing.T) {
	t.Run("creates repository successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		repo, err := NewSkillRepository(tmpDir)

		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("fails with empty path", func(t *testing.T) {
		_, err := NewSkillRepository("")

		assert.Error(t, err)
	})
}

func TestSkillRepository_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewSkillRepository(tmpDir)

	t.Run("creates and retrieves skill", func(t *testing.T) {
		skill, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
		skill.Body = "Python programming language"

		err := repo.Create(skill)
		require.NoError(t, err)

		retrieved, err := repo.GetByIDWithBody("skill-001")
		require.NoError(t, err)
		assert.Equal(t, "Python", retrieved.Title)
		assert.Equal(t, "programming", retrieved.Category)
		assert.Contains(t, retrieved.Body, "Python programming")
	})

	t.Run("updates skill", func(t *testing.T) {
		skill, _ := repo.GetByID("skill-001")
		skill.UpdateLevel(core.LevelAdvanced)
		skill.AddTag("backend")

		err := repo.Update(skill)
		require.NoError(t, err)

		updated, _ := repo.GetByID("skill-001")
		assert.Equal(t, core.LevelAdvanced, updated.Level)
		assert.Contains(t, updated.Tags, "backend")
	})

	t.Run("deletes skill", func(t *testing.T) {
		err := repo.Delete("skill-001")
		require.NoError(t, err)

		exists, _ := repo.Exists("skill-001")
		assert.False(t, exists)
	})
}

func TestSkillRepository_FindByCategory(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewSkillRepository(tmpDir)

	python, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
	go1, _ := core.NewSkill("skill-002", "Go", "programming", core.LevelAdvanced)
	docker, _ := core.NewSkill("skill-003", "Docker", "devops", core.LevelIntermediate)

	repo.Create(python)
	repo.Create(go1)
	repo.Create(docker)

	t.Run("finds skills by category", func(t *testing.T) {
		results, err := repo.FindByCategory("programming")

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "programming", results[0].Category)
		assert.Equal(t, "programming", results[1].Category)
	})

	t.Run("returns empty for non-existent category", func(t *testing.T) {
		results, err := repo.FindByCategory("nonexistent")

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestSkillRepository_FindByLevel(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewSkillRepository(tmpDir)

	python, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
	go1, _ := core.NewSkill("skill-002", "Go", "programming", core.LevelAdvanced)
	docker, _ := core.NewSkill("skill-003", "Docker", "devops", core.LevelIntermediate)

	repo.Create(python)
	repo.Create(go1)
	repo.Create(docker)

	t.Run("finds skills by level", func(t *testing.T) {
		results, err := repo.FindByLevel(core.LevelIntermediate)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		for _, skill := range results {
			assert.Equal(t, core.LevelIntermediate, skill.Level)
		}
	})

	t.Run("finds advanced skills", func(t *testing.T) {
		results, err := repo.FindByLevel(core.LevelAdvanced)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Go", results[0].Title)
	})
}

func TestSkillRepository_FindByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewSkillRepository(tmpDir)

	python, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
	python.UpdateStatus(core.SkillLearning)

	go1, _ := core.NewSkill("skill-002", "Go", "programming", core.LevelAdvanced)
	go1.UpdateStatus(core.SkillMastered)

	docker, _ := core.NewSkill("skill-003", "Docker", "devops", core.LevelIntermediate)

	repo.Create(python)
	repo.Create(go1)
	repo.Create(docker)

	t.Run("finds skills by status", func(t *testing.T) {
		results, err := repo.FindByStatus(core.SkillLearning)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Python", results[0].Title)
	})

	t.Run("finds not started skills", func(t *testing.T) {
		results, err := repo.FindByStatus(core.SkillNotStarted)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Docker", results[0].Title)
	})

	t.Run("finds mastered skills", func(t *testing.T) {
		results, err := repo.FindByStatus(core.SkillMastered)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Go", results[0].Title)
	})
}

func TestSkillRepository_FindByCategoryAndLevel(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewSkillRepository(tmpDir)

	python, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
	go1, _ := core.NewSkill("skill-002", "Go", "programming", core.LevelAdvanced)
	docker, _ := core.NewSkill("skill-003", "Docker", "devops", core.LevelIntermediate)
	k8s, _ := core.NewSkill("skill-004", "Kubernetes", "devops", core.LevelBeginner)

	repo.Create(python)
	repo.Create(go1)
	repo.Create(docker)
	repo.Create(k8s)

	t.Run("finds skills by category and level", func(t *testing.T) {
		results, err := repo.FindByCategoryAndLevel("programming", core.LevelIntermediate)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Python", results[0].Title)
	})

	t.Run("finds devops intermediate skills", func(t *testing.T) {
		results, err := repo.FindByCategoryAndLevel("devops", core.LevelIntermediate)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Docker", results[0].Title)
	})

	t.Run("returns empty for non-matching combination", func(t *testing.T) {
		results, err := repo.FindByCategoryAndLevel("programming", core.LevelBeginner)

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}

func TestSkillRepository_Search(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewSkillRepository(tmpDir)

	python, _ := core.NewSkill("skill-001", "Python", "programming", core.LevelIntermediate)
	python.AddTag("backend")
	python.AddTag("ml")

	go1, _ := core.NewSkill("skill-002", "Go", "programming", core.LevelAdvanced)
	go1.AddTag("backend")
	go1.AddTag("systems")

	repo.Create(python)
	repo.Create(go1)

	t.Run("searches by title", func(t *testing.T) {
		results, err := repo.Search("python")

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Python", results[0].Title)
	})

	t.Run("searches by tag", func(t *testing.T) {
		results, err := repo.Search("backend")

		require.NoError(t, err)
		assert.Len(t, results, 2)
	})
}
