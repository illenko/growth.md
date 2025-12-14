package core

import (
	"errors"
	"strings"
)

// LearningPath represents a structured plan for achieving a goal
type LearningPath struct {
	ID                EntityID   `yaml:"id"`
	Title             string     `yaml:"title"`
	Type              PathType   `yaml:"type"`
	Status            Status     `yaml:"status"`
	GeneratedBy       string     `yaml:"generatedBy,omitempty"`       // AI model used for generation
	GenerationContext string     `yaml:"generationContext,omitempty"` // Context used for AI generation
	Phases            []EntityID `yaml:"phases,omitempty"`
	Tags              []string   `yaml:"tags,omitempty"`
	Timestamps

	// Body contains the markdown content (philosophy, starting point, phases detail)
	Body string `yaml:"-"`
}

// NewLearningPath creates a new Learning Path
func NewLearningPath(id EntityID, title string, pathType PathType) (*LearningPath, error) {
	path := &LearningPath{
		ID:         id,
		Title:      title,
		Type:       pathType,
		Status:     StatusActive,
		Phases:     []EntityID{},
		Tags:       []string{},
		Timestamps: NewTimestamps(),
	}

	if err := path.Validate(); err != nil {
		return nil, err
	}

	return path, nil
}

// Validate checks if the learning path is valid
func (p *LearningPath) Validate() error {
	if p.ID == "" {
		return errors.New("path ID is required")
	}

	if strings.TrimSpace(p.Title) == "" {
		return errors.New("path title is required")
	}

	if !p.Type.IsValid() {
		return errors.New("path type is invalid")
	}

	if !p.Status.IsValid() {
		return errors.New("path status is invalid")
	}

	if p.Created.IsZero() {
		return errors.New("path created timestamp is required")
	}

	if p.Updated.IsZero() {
		return errors.New("path updated timestamp is required")
	}

	return nil
}

// AddPhase adds a phase to the learning path (maintains order)
func (p *LearningPath) AddPhase(phaseID EntityID) {
	// Check if phase already exists
	for _, id := range p.Phases {
		if id == phaseID {
			return
		}
	}
	p.Phases = append(p.Phases, phaseID)
	p.Touch()
}

// RemovePhase removes a phase from the learning path
func (p *LearningPath) RemovePhase(phaseID EntityID) {
	for i, id := range p.Phases {
		if id == phaseID {
			p.Phases = append(p.Phases[:i], p.Phases[i+1:]...)
			p.Touch()
			return
		}
	}
}

// AddTag adds a tag to the path
func (p *LearningPath) AddTag(tag string) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return
	}

	for _, t := range p.Tags {
		if t == tag {
			return
		}
	}
	p.Tags = append(p.Tags, tag)
	p.Touch()
}

// UpdateStatus updates the path's status
func (p *LearningPath) UpdateStatus(status Status) error {
	if !status.IsValid() {
		return errors.New("invalid path status")
	}
	p.Status = status
	p.Touch()
	return nil
}

// SetGenerationInfo sets AI generation metadata
func (p *LearningPath) SetGenerationInfo(model, context string) {
	p.GeneratedBy = model
	p.GenerationContext = context
	p.Touch()
}
