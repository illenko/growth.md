package storage

import (
	"sort"
	"time"

	"github.com/illenko/growth.md/internal/core"
)

type ProgressLogRepository struct {
	repo Repository[core.ProgressLog]
}

func NewProgressLogRepository(basePath string) (*ProgressLogRepository, error) {
	repo, err := NewFilesystemRepository[core.ProgressLog](basePath, "progress")
	if err != nil {
		return nil, err
	}

	return &ProgressLogRepository{
		repo: repo,
	}, nil
}

// SetConfig sets the configuration for git auto-commit.
func (r *ProgressLogRepository) SetConfig(config *Config) {
	if fsRepo, ok := r.repo.(*FilesystemRepository[core.ProgressLog]); ok {
		fsRepo.SetConfig(config)
	}
}

func (r *ProgressLogRepository) Create(log *core.ProgressLog) error {
	return r.repo.Create(log)
}

func (r *ProgressLogRepository) GetByID(id core.EntityID) (*core.ProgressLog, error) {
	return r.repo.GetByID(id)
}

func (r *ProgressLogRepository) GetByIDWithBody(id core.EntityID) (*core.ProgressLog, error) {
	return r.repo.GetByIDWithBody(id)
}

func (r *ProgressLogRepository) GetAll() ([]*core.ProgressLog, error) {
	return r.repo.GetAll()
}

func (r *ProgressLogRepository) Update(log *core.ProgressLog) error {
	return r.repo.Update(log)
}

func (r *ProgressLogRepository) Delete(id core.EntityID) error {
	return r.repo.Delete(id)
}

func (r *ProgressLogRepository) Search(query string) ([]*core.ProgressLog, error) {
	return r.repo.Search(query)
}

func (r *ProgressLogRepository) Exists(id core.EntityID) (bool, error) {
	return r.repo.Exists(id)
}

func (r *ProgressLogRepository) FindByDateRange(start, end time.Time) ([]*core.ProgressLog, error) {
	allLogs, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.ProgressLog
	for _, log := range allLogs {
		if (log.Date.Equal(start) || log.Date.After(start)) &&
			(log.Date.Equal(end) || log.Date.Before(end)) {
			results = append(results, log)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.After(results[j].Date)
	})

	return results, nil
}

func (r *ProgressLogRepository) FindRecent(limit int) ([]*core.ProgressLog, error) {
	allLogs, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].Date.After(allLogs[j].Date)
	})

	if limit > 0 && len(allLogs) > limit {
		return allLogs[:limit], nil
	}

	return allLogs, nil
}

func (r *ProgressLogRepository) FindBySkillID(skillID core.EntityID) ([]*core.ProgressLog, error) {
	allLogs, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.ProgressLog
	for _, log := range allLogs {
		for _, id := range log.SkillsWorked {
			if id == skillID {
				results = append(results, log)
				break
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.After(results[j].Date)
	})

	return results, nil
}

func (r *ProgressLogRepository) FindByResourceID(resourceID core.EntityID) ([]*core.ProgressLog, error) {
	allLogs, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.ProgressLog
	for _, log := range allLogs {
		for _, id := range log.ResourcesUsed {
			if id == resourceID {
				results = append(results, log)
				break
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.After(results[j].Date)
	})

	return results, nil
}
