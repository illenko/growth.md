package storage

import (
	"sort"

	"github.com/illenko/growth.md/internal/core"
)

type PhaseRepository struct {
	repo Repository[core.Phase]
}

func NewPhaseRepository(basePath string) (*PhaseRepository, error) {
	repo, err := NewFilesystemRepository[core.Phase](basePath, "phase")
	if err != nil {
		return nil, err
	}

	return &PhaseRepository{
		repo: repo,
	}, nil
}

func (r *PhaseRepository) Create(phase *core.Phase) error {
	return r.repo.Create(phase)
}

func (r *PhaseRepository) GetByID(id core.EntityID) (*core.Phase, error) {
	return r.repo.GetByID(id)
}

func (r *PhaseRepository) GetByIDWithBody(id core.EntityID) (*core.Phase, error) {
	return r.repo.GetByIDWithBody(id)
}

func (r *PhaseRepository) GetAll() ([]*core.Phase, error) {
	return r.repo.GetAll()
}

func (r *PhaseRepository) Update(phase *core.Phase) error {
	return r.repo.Update(phase)
}

func (r *PhaseRepository) Delete(id core.EntityID) error {
	return r.repo.Delete(id)
}

func (r *PhaseRepository) Search(query string) ([]*core.Phase, error) {
	return r.repo.Search(query)
}

func (r *PhaseRepository) Exists(id core.EntityID) (bool, error) {
	return r.repo.Exists(id)
}

func (r *PhaseRepository) FindByPathID(pathID core.EntityID) ([]*core.Phase, error) {
	allPhases, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Phase
	for _, phase := range allPhases {
		if phase.PathID == pathID {
			results = append(results, phase)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Order < results[j].Order
	})

	return results, nil
}
