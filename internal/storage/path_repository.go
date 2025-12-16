package storage

import (
	"github.com/illenko/growth.md/internal/core"
)

type PathRepository struct {
	repo Repository[core.LearningPath]
}

func NewPathRepository(basePath string) (*PathRepository, error) {
	repo, err := NewFilesystemRepository[core.LearningPath](basePath, "path")
	if err != nil {
		return nil, err
	}

	return &PathRepository{
		repo: repo,
	}, nil
}

func (r *PathRepository) Create(path *core.LearningPath) error {
	return r.repo.Create(path)
}

func (r *PathRepository) GetByID(id core.EntityID) (*core.LearningPath, error) {
	return r.repo.GetByID(id)
}

func (r *PathRepository) GetByIDWithBody(id core.EntityID) (*core.LearningPath, error) {
	return r.repo.GetByIDWithBody(id)
}

func (r *PathRepository) GetAll() ([]*core.LearningPath, error) {
	return r.repo.GetAll()
}

func (r *PathRepository) Update(path *core.LearningPath) error {
	return r.repo.Update(path)
}

func (r *PathRepository) Delete(id core.EntityID) error {
	return r.repo.Delete(id)
}

func (r *PathRepository) Search(query string) ([]*core.LearningPath, error) {
	return r.repo.Search(query)
}

func (r *PathRepository) Exists(id core.EntityID) (bool, error) {
	return r.repo.Exists(id)
}

func (r *PathRepository) FindByType(pathType core.PathType) ([]*core.LearningPath, error) {
	allPaths, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.LearningPath
	for _, path := range allPaths {
		if path.Type == pathType {
			results = append(results, path)
		}
	}

	return results, nil
}

func (r *PathRepository) FindByStatus(status core.Status) ([]*core.LearningPath, error) {
	allPaths, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.LearningPath
	for _, path := range allPaths {
		if path.Status == status {
			results = append(results, path)
		}
	}

	return results, nil
}

func (r *PathRepository) FindActive() ([]*core.LearningPath, error) {
	return r.FindByStatus(core.StatusActive)
}

func (r *PathRepository) FindAIGenerated() ([]*core.LearningPath, error) {
	return r.FindByType(core.PathTypeAIGenerated)
}
