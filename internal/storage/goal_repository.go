package storage

import (
	"time"

	"github.com/illenko/growth.md/internal/core"
)

type GoalRepository struct {
	repo Repository[core.Goal]
}

func NewGoalRepository(basePath string) (*GoalRepository, error) {
	repo, err := NewFilesystemRepository[core.Goal](basePath, "goal")
	if err != nil {
		return nil, err
	}

	return &GoalRepository{
		repo: repo,
	}, nil
}

// SetConfig sets the configuration for git auto-commit.
func (r *GoalRepository) SetConfig(config *Config) {
	if fsRepo, ok := r.repo.(*FilesystemRepository[core.Goal]); ok {
		fsRepo.SetConfig(config)
	}
}

func (r *GoalRepository) Create(goal *core.Goal) error {
	return r.repo.Create(goal)
}

func (r *GoalRepository) GetByID(id core.EntityID) (*core.Goal, error) {
	return r.repo.GetByID(id)
}

func (r *GoalRepository) GetByIDWithBody(id core.EntityID) (*core.Goal, error) {
	return r.repo.GetByIDWithBody(id)
}

func (r *GoalRepository) GetAll() ([]*core.Goal, error) {
	return r.repo.GetAll()
}

func (r *GoalRepository) Update(goal *core.Goal) error {
	return r.repo.Update(goal)
}

func (r *GoalRepository) Delete(id core.EntityID) error {
	return r.repo.Delete(id)
}

func (r *GoalRepository) Search(query string) ([]*core.Goal, error) {
	return r.repo.Search(query)
}

func (r *GoalRepository) Exists(id core.EntityID) (bool, error) {
	return r.repo.Exists(id)
}

func (r *GoalRepository) FindByStatus(status core.Status) ([]*core.Goal, error) {
	allGoals, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Goal
	for _, goal := range allGoals {
		if goal.Status == status {
			results = append(results, goal)
		}
	}

	return results, nil
}

func (r *GoalRepository) FindByPriority(priority core.Priority) ([]*core.Goal, error) {
	allGoals, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Goal
	for _, goal := range allGoals {
		if goal.Priority == priority {
			results = append(results, goal)
		}
	}

	return results, nil
}

func (r *GoalRepository) FindActive() ([]*core.Goal, error) {
	return r.FindByStatus(core.StatusActive)
}

func (r *GoalRepository) FindByTargetDateRange(start, end time.Time) ([]*core.Goal, error) {
	allGoals, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Goal
	for _, goal := range allGoals {
		if goal.TargetDate != nil {
			if (goal.TargetDate.Equal(start) || goal.TargetDate.After(start)) &&
				(goal.TargetDate.Equal(end) || goal.TargetDate.Before(end)) {
				results = append(results, goal)
			}
		}
	}

	return results, nil
}
