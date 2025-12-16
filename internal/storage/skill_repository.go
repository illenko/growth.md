package storage

import (
	"github.com/illenko/growth.md/internal/core"
)

type SkillRepository struct {
	repo Repository[core.Skill]
}

func NewSkillRepository(basePath string) (*SkillRepository, error) {
	repo, err := NewFilesystemRepository[core.Skill](basePath, "skill")
	if err != nil {
		return nil, err
	}

	return &SkillRepository{
		repo: repo,
	}, nil
}

func (r *SkillRepository) Create(skill *core.Skill) error {
	return r.repo.Create(skill)
}

func (r *SkillRepository) GetByID(id core.EntityID) (*core.Skill, error) {
	return r.repo.GetByID(id)
}

func (r *SkillRepository) GetByIDWithBody(id core.EntityID) (*core.Skill, error) {
	return r.repo.GetByIDWithBody(id)
}

func (r *SkillRepository) GetAll() ([]*core.Skill, error) {
	return r.repo.GetAll()
}

func (r *SkillRepository) Update(skill *core.Skill) error {
	return r.repo.Update(skill)
}

func (r *SkillRepository) Delete(id core.EntityID) error {
	return r.repo.Delete(id)
}

func (r *SkillRepository) Search(query string) ([]*core.Skill, error) {
	return r.repo.Search(query)
}

func (r *SkillRepository) Exists(id core.EntityID) (bool, error) {
	return r.repo.Exists(id)
}

func (r *SkillRepository) FindByCategory(category string) ([]*core.Skill, error) {
	allSkills, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Skill
	for _, skill := range allSkills {
		if skill.Category == category {
			results = append(results, skill)
		}
	}

	return results, nil
}

func (r *SkillRepository) FindByLevel(level core.ProficiencyLevel) ([]*core.Skill, error) {
	allSkills, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Skill
	for _, skill := range allSkills {
		if skill.Level == level {
			results = append(results, skill)
		}
	}

	return results, nil
}

func (r *SkillRepository) FindByStatus(status core.SkillStatus) ([]*core.Skill, error) {
	allSkills, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Skill
	for _, skill := range allSkills {
		if skill.Status == status {
			results = append(results, skill)
		}
	}

	return results, nil
}

func (r *SkillRepository) FindByCategoryAndLevel(category string, level core.ProficiencyLevel) ([]*core.Skill, error) {
	allSkills, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*core.Skill
	for _, skill := range allSkills {
		if skill.Category == category && skill.Level == level {
			results = append(results, skill)
		}
	}

	return results, nil
}
