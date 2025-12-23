package core

import (
	"errors"
	"strings"
)

// Skill represents a technical or professional competency
type Skill struct {
	ID        EntityID         `yaml:"id"`
	Title     string           `yaml:"title"`
	Category  string           `yaml:"category"`
	Level     ProficiencyLevel `yaml:"level"`
	Status    SkillStatus      `yaml:"status"`
	Resources []EntityID       `yaml:"resources,omitempty"`
	Tags      []string         `yaml:"tags,omitempty"`
	Timestamps

	// Free-form notes, learning goals, projects, etc.
	Body string `yaml:"-"`
}

func NewSkill(id EntityID, title, category string, level ProficiencyLevel) (*Skill, error) {
	skill := &Skill{
		ID:         id,
		Title:      title,
		Category:   category,
		Level:      level,
		Status:     SkillNotStarted,
		Resources:  []EntityID{},
		Tags:       []string{},
		Timestamps: NewTimestamps(),
	}

	if err := skill.Validate(); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *Skill) Validate() error {
	if s.ID == "" {
		return errors.New("skill ID is required")
	}

	if strings.TrimSpace(s.Title) == "" {
		return errors.New("skill title is required and cannot be empty")
	}

	if strings.TrimSpace(s.Category) == "" {
		return errors.New("skill category is required (e.g., 'Programming Languages', 'DevOps', 'Data Science')")
	}

	if !s.Level.IsValid() {
		return errors.New("invalid skill level: must be one of: beginner, intermediate, advanced, expert")
	}

	if !s.Status.IsValid() {
		return errors.New("invalid skill status: must be one of: not-started, learning, mastered")
	}

	if s.Created.IsZero() {
		return errors.New("skill created timestamp is required")
	}

	if s.Updated.IsZero() {
		return errors.New("skill updated timestamp is required")
	}

	return nil
}

func (s *Skill) AddResource(resourceID EntityID) {
	for _, id := range s.Resources {
		if id == resourceID {
			return
		}
	}
	s.Resources = append(s.Resources, resourceID)
	s.Touch()
}

func (s *Skill) RemoveResource(resourceID EntityID) {
	for i, id := range s.Resources {
		if id == resourceID {
			s.Resources = append(s.Resources[:i], s.Resources[i+1:]...)
			s.Touch()
			return
		}
	}
}

func (s *Skill) AddTag(tag string) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return
	}

	for _, t := range s.Tags {
		if t == tag {
			return
		}
	}
	s.Tags = append(s.Tags, tag)
	s.Touch()
}

func (s *Skill) UpdateLevel(level ProficiencyLevel) error {
	if !level.IsValid() {
		return errors.New("invalid proficiency level: must be one of: beginner, intermediate, advanced, expert")
	}
	s.Level = level
	s.Touch()
	return nil
}

func (s *Skill) UpdateStatus(status SkillStatus) error {
	if !status.IsValid() {
		return errors.New("invalid skill status: must be one of: not-started, learning, mastered")
	}
	s.Status = status
	s.Touch()
	return nil
}
