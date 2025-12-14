package core

import (
	"errors"
	"strings"
)

// Resource represents a learning material
type Resource struct {
	ID             EntityID       `yaml:"id"`
	Title          string         `yaml:"title"`
	Type           ResourceType   `yaml:"type"`
	SkillID        EntityID       `yaml:"skillId"`
	Status         ResourceStatus `yaml:"status"`
	URL            string         `yaml:"url,omitempty"`
	Author         string         `yaml:"author,omitempty"`
	EstimatedHours float64        `yaml:"estimatedHours,omitempty"`
	Tags           []string       `yaml:"tags,omitempty"`
	Timestamps

	// Body contains the markdown content (overview, progress, key takeaways, application, rating)
	Body string `yaml:"-"`
}

// NewResource creates a new Resource
func NewResource(id EntityID, title string, resourceType ResourceType, skillID EntityID) (*Resource, error) {
	resource := &Resource{
		ID:         id,
		Title:      title,
		Type:       resourceType,
		SkillID:    skillID,
		Status:     ResourceNotStarted,
		Tags:       []string{},
		Timestamps: NewTimestamps(),
	}

	if err := resource.Validate(); err != nil {
		return nil, err
	}

	return resource, nil
}

// Validate checks if the resource is valid
func (r *Resource) Validate() error {
	if r.ID == "" {
		return errors.New("resource ID is required")
	}

	if strings.TrimSpace(r.Title) == "" {
		return errors.New("resource title is required")
	}

	if !r.Type.IsValid() {
		return errors.New("resource type is invalid")
	}

	if r.SkillID == "" {
		return errors.New("resource skillId is required")
	}

	if !r.Status.IsValid() {
		return errors.New("resource status is invalid")
	}

	if r.EstimatedHours < 0 {
		return errors.New("resource estimatedHours cannot be negative")
	}

	if r.Created.IsZero() {
		return errors.New("resource created timestamp is required")
	}

	if r.Updated.IsZero() {
		return errors.New("resource updated timestamp is required")
	}

	return nil
}

// UpdateStatus updates the resource's status
func (r *Resource) UpdateStatus(status ResourceStatus) error {
	if !status.IsValid() {
		return errors.New("invalid resource status")
	}
	r.Status = status
	r.Touch()
	return nil
}

// Start marks the resource as in-progress
func (r *Resource) Start() {
	r.Status = ResourceInProgress
	r.Touch()
}

// Complete marks the resource as completed
func (r *Resource) Complete() {
	r.Status = ResourceCompleted
	r.Touch()
}

// AddTag adds a tag to the resource
func (r *Resource) AddTag(tag string) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return
	}

	for _, t := range r.Tags {
		if t == tag {
			return
		}
	}
	r.Tags = append(r.Tags, tag)
	r.Touch()
}

// SetURL sets the resource URL
func (r *Resource) SetURL(url string) {
	r.URL = url
	r.Touch()
}

// SetAuthor sets the resource author
func (r *Resource) SetAuthor(author string) {
	r.Author = author
	r.Touch()
}

// SetEstimatedHours sets the estimated time investment
func (r *Resource) SetEstimatedHours(hours float64) error {
	if hours < 0 {
		return errors.New("estimated hours cannot be negative")
	}
	r.EstimatedHours = hours
	r.Touch()
	return nil
}
