package core

import (
	"errors"
	"strings"
	"time"
)

// Milestone represents a significant achievement
type Milestone struct {
	ID            EntityID      `yaml:"id"`
	Title         string        `yaml:"title"`
	Type          MilestoneType `yaml:"type"`
	ReferenceType ReferenceType `yaml:"referenceType"`
	ReferenceID   EntityID      `yaml:"referenceId"`
	Status        Status        `yaml:"status"` // pending or completed
	AchievedDate  *time.Time    `yaml:"achievedDate,omitempty"`
	TargetDate    *time.Time    `yaml:"targetDate,omitempty"`
	Proof         string        `yaml:"proof,omitempty"` // URL to evidence
	Timestamps

	// Body contains the markdown content (definition of done, success metrics, importance, notes)
	Body string `yaml:"-"`
}

// NewMilestone creates a new Milestone
func NewMilestone(id EntityID, title string, milestoneType MilestoneType, refType ReferenceType, refID EntityID) (*Milestone, error) {
	milestone := &Milestone{
		ID:            id,
		Title:         title,
		Type:          milestoneType,
		ReferenceType: refType,
		ReferenceID:   refID,
		Status:        StatusActive, // Using "active" for pending milestones
		Timestamps:    NewTimestamps(),
	}

	if err := milestone.Validate(); err != nil {
		return nil, err
	}

	return milestone, nil
}

func (m *Milestone) Validate() error {
	if m.ID == "" {
		return errors.New("milestone ID is required")
	}

	if strings.TrimSpace(m.Title) == "" {
		return errors.New("milestone title is required and cannot be empty")
	}

	if !m.Type.IsValid() {
		return errors.New("invalid milestone type: must be one of: goal-level, path-level, skill-level")
	}

	if !m.ReferenceType.IsValid() {
		return errors.New("invalid reference type: must be one of: goal, path, skill (use --ref-type flag)")
	}

	if m.ReferenceID == "" {
		return errors.New("milestone reference ID is required (use --ref-id flag)")
	}

	if !m.Status.IsValid() {
		return errors.New("invalid milestone status: must be one of: active, completed, archived")
	}

	if m.Created.IsZero() {
		return errors.New("milestone created timestamp is required")
	}

	if m.Updated.IsZero() {
		return errors.New("milestone updated timestamp is required")
	}

	return nil
}

func (m *Milestone) Achieve(proof string) {
	now := time.Now()
	m.Status = StatusCompleted
	m.AchievedDate = &now
	if proof != "" {
		m.Proof = proof
	}
	m.Touch()
}

// SetTargetDate sets when the milestone should be achieved
func (m *Milestone) SetTargetDate(date time.Time) {
	m.TargetDate = &date
	m.Touch()
}

// ClearTargetDate removes the target date
func (m *Milestone) ClearTargetDate() {
	m.TargetDate = nil
	m.Touch()
}

// SetProof sets evidence/proof URL
func (m *Milestone) SetProof(proof string) {
	m.Proof = proof
	m.Touch()
}

func (m *Milestone) IsAchieved() bool {
	return m.Status == StatusCompleted && m.AchievedDate != nil
}
