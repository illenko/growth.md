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

	// Body contains the markdown content (not in YAML frontmatter)
	// This is free-form notes, learning goals, projects, etc.
	Body string `yaml:"-"` // yaml:"-" means don't include in frontmatter
}

// NewSkill creates a new Skill with the given title
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

// Validate checks if the skill is valid
func (s *Skill) Validate() error {
	if s.ID == "" {
		return errors.New("skill ID is required")
	}

	if strings.TrimSpace(s.Title) == "" {
		return errors.New("skill title is required")
	}

	if strings.TrimSpace(s.Category) == "" {
		return errors.New("skill category is required")
	}

	if !s.Level.IsValid() {
		return errors.New("skill level is invalid")
	}

	if !s.Status.IsValid() {
		return errors.New("skill status is invalid")
	}

	if s.Created.IsZero() {
		return errors.New("skill created timestamp is required")
	}

	if s.Updated.IsZero() {
		return errors.New("skill updated timestamp is required")
	}

	return nil
}

// AddResource adds a resource to the skill
func (s *Skill) AddResource(resourceID EntityID) {
	// Check if resource already exists
	for _, id := range s.Resources {
		if id == resourceID {
			return
		}
	}
	s.Resources = append(s.Resources, resourceID)
	s.Touch()
}

// RemoveResource removes a resource from the skill
func (s *Skill) RemoveResource(resourceID EntityID) {
	for i, id := range s.Resources {
		if id == resourceID {
			s.Resources = append(s.Resources[:i], s.Resources[i+1:]...)
			s.Touch()
			return
		}
	}
}

// AddTag adds a tag to the skill
func (s *Skill) AddTag(tag string) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return
	}

	// Check if tag already exists
	for _, t := range s.Tags {
		if t == tag {
			return
		}
	}
	s.Tags = append(s.Tags, tag)
	s.Touch()
}

// UpdateLevel updates the skill's proficiency level
func (s *Skill) UpdateLevel(level ProficiencyLevel) error {
	if !level.IsValid() {
		return errors.New("invalid proficiency level")
	}
	s.Level = level
	s.Touch()
	return nil
}

// UpdateStatus updates the skill's status
func (s *Skill) UpdateStatus(status SkillStatus) error {
	if !status.IsValid() {
		return errors.New("invalid skill status")
	}
	s.Status = status
	s.Touch()
	return nil
}
