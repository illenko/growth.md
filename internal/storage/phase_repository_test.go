package storage

import (
	"testing"

	"github.com/illenko/growth.md/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPhaseRepository(t *testing.T) {
	t.Run("creates repository successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		repo, err := NewPhaseRepository(tmpDir)

		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("fails with empty path", func(t *testing.T) {
		_, err := NewPhaseRepository("")

		assert.Error(t, err)
	})
}

func TestPhaseRepository_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewPhaseRepository(tmpDir)

	t.Run("creates and retrieves phase", func(t *testing.T) {
		phase, _ := core.NewPhase("phase-001", "path-001", "Foundation Phase", 1)
		phase.Body = "Learn the fundamentals"

		err := repo.Create(phase)
		require.NoError(t, err)

		retrieved, err := repo.GetByIDWithBody("phase-001")
		require.NoError(t, err)
		assert.Equal(t, "Foundation Phase", retrieved.Title)
		assert.Contains(t, retrieved.Body, "fundamentals")
		assert.Equal(t, "path-001", string(retrieved.PathID))
	})

	t.Run("updates phase", func(t *testing.T) {
		phase, _ := repo.GetByID("phase-001")
		phase.EstimatedDuration = "2 months"
		phase.AddSkillRequirement("skill-001", core.LevelIntermediate)

		err := repo.Update(phase)
		require.NoError(t, err)

		updated, _ := repo.GetByID("phase-001")
		assert.Equal(t, "2 months", updated.EstimatedDuration)
		assert.Len(t, updated.RequiredSkills, 1)
	})

	t.Run("deletes phase", func(t *testing.T) {
		err := repo.Delete("phase-001")
		require.NoError(t, err)

		exists, _ := repo.Exists("phase-001")
		assert.False(t, exists)
	})
}

func TestPhaseRepository_FindByPathID(t *testing.T) {
	tmpDir := t.TempDir()
	repo, _ := NewPhaseRepository(tmpDir)

	phase1, _ := core.NewPhase("phase-001", "path-001", "Phase 1", 1)
	phase2, _ := core.NewPhase("phase-002", "path-001", "Phase 2", 2)
	phase3, _ := core.NewPhase("phase-003", "path-001", "Phase 3", 3)
	phase4, _ := core.NewPhase("phase-004", "path-002", "Other Path Phase", 1)

	repo.Create(phase1)
	repo.Create(phase2)
	repo.Create(phase3)
	repo.Create(phase4)

	t.Run("finds phases by path ID", func(t *testing.T) {
		results, err := repo.FindByPathID("path-001")

		require.NoError(t, err)
		assert.Len(t, results, 3)
		for _, phase := range results {
			assert.Equal(t, "path-001", string(phase.PathID))
		}
	})

	t.Run("returns phases ordered by order field", func(t *testing.T) {
		results, err := repo.FindByPathID("path-001")

		require.NoError(t, err)
		assert.Equal(t, "Phase 1", results[0].Title)
		assert.Equal(t, "Phase 2", results[1].Title)
		assert.Equal(t, "Phase 3", results[2].Title)
		assert.Equal(t, 1, results[0].Order)
		assert.Equal(t, 2, results[1].Order)
		assert.Equal(t, 3, results[2].Order)
	})

	t.Run("handles out of order insertion", func(t *testing.T) {
		phase5, _ := core.NewPhase("phase-005", "path-003", "Phase C", 3)
		phase6, _ := core.NewPhase("phase-006", "path-003", "Phase A", 1)
		phase7, _ := core.NewPhase("phase-007", "path-003", "Phase B", 2)

		repo.Create(phase5)
		repo.Create(phase6)
		repo.Create(phase7)

		results, err := repo.FindByPathID("path-003")

		require.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Equal(t, "Phase A", results[0].Title)
		assert.Equal(t, "Phase B", results[1].Title)
		assert.Equal(t, "Phase C", results[2].Title)
	})

	t.Run("returns empty for non-existent path", func(t *testing.T) {
		results, err := repo.FindByPathID("path-999")

		require.NoError(t, err)
		assert.Empty(t, results)
	})
}
