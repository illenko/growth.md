package storage

import (
	"github.com/illenko/growth.md/internal/core"
)

type MilestoneRepository struct {
	repo Repository[core.Milestone]
}

func NewMilestoneRepository(basePath string) (*MilestoneRepository, error) {
	repo, err := NewFilesystemRepository[core.Milestone](basePath, "milestone")
	if err != nil {
		return nil, err
	}

	return &MilestoneRepository{
		repo: repo,
	}, nil
}

// SetConfig sets the configuration for git auto-commit.
func (r *MilestoneRepository) SetConfig(config *Config) {
	if fsRepo, ok := r.repo.(*FilesystemRepository[core.Milestone]); ok {
		fsRepo.SetConfig(config)
	}
}

func (r *MilestoneRepository) Create(milestone *core.Milestone) error {
	return r.repo.Create(milestone)
}

func (r *MilestoneRepository) GetByID(id core.EntityID) (*core.Milestone, error) {
	return r.repo.GetByID(id)
}

func (r *MilestoneRepository) GetByIDWithBody(id core.EntityID) (*core.Milestone, error) {
	return r.repo.GetByIDWithBody(id)
}

func (r *MilestoneRepository) GetAll() ([]*core.Milestone, error) {
	return r.repo.GetAll()
}

func (r *MilestoneRepository) Update(milestone *core.Milestone) error {
	return r.repo.Update(milestone)
}

func (r *MilestoneRepository) Delete(id core.EntityID) error {
	return r.repo.Delete(id)
}

func (r *MilestoneRepository) Search(query string) ([]*core.Milestone, error) {
	return r.repo.Search(query)
}

func (r *MilestoneRepository) Exists(id core.EntityID) (bool, error) {
	return r.repo.Exists(id)
}

func (r *MilestoneRepository) FindByReferenceID(refType core.ReferenceType, refID core.EntityID) ([]*core.Milestone, error) {
	allMilestones, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Milestone
	for _, milestone := range allMilestones {
		if milestone.ReferenceType == refType && milestone.ReferenceID == refID {
			results = append(results, milestone)
		}
	}

	return results, nil
}

func (r *MilestoneRepository) FindByStatus(status core.Status) ([]*core.Milestone, error) {
	allMilestones, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Milestone
	for _, milestone := range allMilestones {
		if milestone.Status == status {
			results = append(results, milestone)
		}
	}

	return results, nil
}

func (r *MilestoneRepository) FindByType(milestoneType core.MilestoneType) ([]*core.Milestone, error) {
	allMilestones, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Milestone
	for _, milestone := range allMilestones {
		if milestone.Type == milestoneType {
			results = append(results, milestone)
		}
	}

	return results, nil
}
