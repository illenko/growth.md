package core

import (
	"errors"
	"strings"
)

// SkillRequirement defines a skill needed for a phase with target level
type SkillRequirement struct {
	SkillID     EntityID         `yaml:"skillId"`
	TargetLevel ProficiencyLevel `yaml:"targetLevel"`
}

// Phase represents a stage within a learning path
type Phase struct {
	ID                EntityID           `yaml:"id"`
	PathID            EntityID           `yaml:"pathId"`
	Title             string             `yaml:"title"`
	Order             int                `yaml:"order"`
	EstimatedDuration string             `yaml:"estimatedDuration,omitempty"` // e.g., "2 months"
	RequiredSkills    []SkillRequirement `yaml:"requiredSkills,omitempty"`
	Milestones        []EntityID         `yaml:"milestones,omitempty"`
	Timestamps

	// Body contains the markdown content (goal, projects, timeline)
	Body string `yaml:"-"`
}

// NewPhase creates a new Phase
func NewPhase(id, pathID EntityID, title string, order int) (*Phase, error) {
	phase := &Phase{
		ID:             id,
		PathID:         pathID,
		Title:          title,
		Order:          order,
		RequiredSkills: []SkillRequirement{},
		Milestones:     []EntityID{},
		Timestamps:     NewTimestamps(),
	}

	if err := phase.Validate(); err != nil {
		return nil, err
	}

	return phase, nil
}

// Validate checks if the phase is valid
func (p *Phase) Validate() error {
	if p.ID == "" {
		return errors.New("phase ID is required")
	}

	if p.PathID == "" {
		return errors.New("phase pathID is required")
	}

	if strings.TrimSpace(p.Title) == "" {
		return errors.New("phase title is required")
	}

	if p.Order < 0 {
		return errors.New("phase order must be non-negative")
	}

	// Validate skill requirements
	for _, req := range p.RequiredSkills {
		if req.SkillID == "" {
			return errors.New("skill requirement must have a skillId")
		}
		if !req.TargetLevel.IsValid() {
			return errors.New("skill requirement has invalid target level")
		}
	}

	if p.Created.IsZero() {
		return errors.New("phase created timestamp is required")
	}

	if p.Updated.IsZero() {
		return errors.New("phase updated timestamp is required")
	}

	return nil
}

// AddSkillRequirement adds a required skill for this phase
func (p *Phase) AddSkillRequirement(skillID EntityID, targetLevel ProficiencyLevel) error {
	if !targetLevel.IsValid() {
		return errors.New("invalid target level")
	}

	// Check if skill already required
	for _, req := range p.RequiredSkills {
		if req.SkillID == skillID {
			// Update target level if different
			req.TargetLevel = targetLevel
			p.Touch()
			return nil
		}
	}

	p.RequiredSkills = append(p.RequiredSkills, SkillRequirement{
		SkillID:     skillID,
		TargetLevel: targetLevel,
	})
	p.Touch()
	return nil
}

// RemoveSkillRequirement removes a skill requirement
func (p *Phase) RemoveSkillRequirement(skillID EntityID) {
	for i, req := range p.RequiredSkills {
		if req.SkillID == skillID {
			p.RequiredSkills = append(p.RequiredSkills[:i], p.RequiredSkills[i+1:]...)
			p.Touch()
			return
		}
	}
}

// AddMilestone adds a milestone to the phase
func (p *Phase) AddMilestone(milestoneID EntityID) {
	for _, id := range p.Milestones {
		if id == milestoneID {
			return
		}
	}
	p.Milestones = append(p.Milestones, milestoneID)
	p.Touch()
}
