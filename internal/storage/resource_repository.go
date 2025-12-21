package storage

import (
	"github.com/illenko/growth.md/internal/core"
)

type ResourceRepository struct {
	repo Repository[core.Resource]
}

func NewResourceRepository(basePath string) (*ResourceRepository, error) {
	repo, err := NewFilesystemRepository[core.Resource](basePath, "resource")
	if err != nil {
		return nil, err
	}

	return &ResourceRepository{
		repo: repo,
	}, nil
}

// SetConfig sets the configuration for git auto-commit.
func (r *ResourceRepository) SetConfig(config *Config) {
	if fsRepo, ok := r.repo.(*FilesystemRepository[core.Resource]); ok {
		fsRepo.SetConfig(config)
	}
}

func (r *ResourceRepository) Create(resource *core.Resource) error {
	return r.repo.Create(resource)
}

func (r *ResourceRepository) GetByID(id core.EntityID) (*core.Resource, error) {
	return r.repo.GetByID(id)
}

func (r *ResourceRepository) GetByIDWithBody(id core.EntityID) (*core.Resource, error) {
	return r.repo.GetByIDWithBody(id)
}

func (r *ResourceRepository) GetAll() ([]*core.Resource, error) {
	return r.repo.GetAll()
}

func (r *ResourceRepository) Update(resource *core.Resource) error {
	return r.repo.Update(resource)
}

func (r *ResourceRepository) Delete(id core.EntityID) error {
	return r.repo.Delete(id)
}

func (r *ResourceRepository) Search(query string) ([]*core.Resource, error) {
	return r.repo.Search(query)
}

func (r *ResourceRepository) Exists(id core.EntityID) (bool, error) {
	return r.repo.Exists(id)
}

func (r *ResourceRepository) FindByType(resourceType core.ResourceType) ([]*core.Resource, error) {
	allResources, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Resource
	for _, resource := range allResources {
		if resource.Type == resourceType {
			results = append(results, resource)
		}
	}

	return results, nil
}

func (r *ResourceRepository) FindBySkillID(skillID core.EntityID) ([]*core.Resource, error) {
	allResources, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Resource
	for _, resource := range allResources {
		if resource.SkillID == skillID {
			results = append(results, resource)
		}
	}

	return results, nil
}

func (r *ResourceRepository) FindByStatus(status core.ResourceStatus) ([]*core.Resource, error) {
	allResources, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Resource
	for _, resource := range allResources {
		if resource.Status == status {
			results = append(results, resource)
		}
	}

	return results, nil
}
