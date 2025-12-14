package core

import (
	"errors"
	"strings"
	"time"
)

// Goal represents a high-level career objective
type Goal struct {
	ID            EntityID   `yaml:"id"`
	Title         string     `yaml:"title"`
	Status        Status     `yaml:"status"`
	Priority      Priority   `yaml:"priority"`
	TargetDate    *time.Time `yaml:"targetDate,omitempty"`
	LearningPaths []EntityID `yaml:"learningPaths,omitempty"`
	Milestones    []EntityID `yaml:"milestones,omitempty"`
	Tags          []string   `yaml:"tags,omitempty"`
	Timestamps

	// Body contains the markdown content (motivation, success criteria, timeline, notes)
	Body string `yaml:"-"`
}

// NewGoal creates a new Goal with the given title and priority
func NewGoal(id EntityID, title string, priority Priority) (*Goal, error) {
	goal := &Goal{
		ID:            id,
		Title:         title,
		Status:        StatusActive,
		Priority:      priority,
		LearningPaths: []EntityID{},
		Milestones:    []EntityID{},
		Tags:          []string{},
		Timestamps:    NewTimestamps(),
	}

	if err := goal.Validate(); err != nil {
		return nil, err
	}

	return goal, nil
}

// Validate checks if the goal is valid
func (g *Goal) Validate() error {
	if g.ID == "" {
		return errors.New("goal ID is required")
	}

	if strings.TrimSpace(g.Title) == "" {
		return errors.New("goal title is required")
	}

	if !g.Status.IsValid() {
		return errors.New("goal status is invalid")
	}

	if !g.Priority.IsValid() {
		return errors.New("goal priority is invalid")
	}

	if g.Created.IsZero() {
		return errors.New("goal created timestamp is required")
	}

	if g.Updated.IsZero() {
		return errors.New("goal updated timestamp is required")
	}

	return nil
}

// AddLearningPath adds a learning path to the goal
func (g *Goal) AddLearningPath(pathID EntityID) {
	// Check if path already exists
	for _, id := range g.LearningPaths {
		if id == pathID {
			return
		}
	}
	g.LearningPaths = append(g.LearningPaths, pathID)
	g.Touch()
}

// RemoveLearningPath removes a learning path from the goal
func (g *Goal) RemoveLearningPath(pathID EntityID) {
	for i, id := range g.LearningPaths {
		if id == pathID {
			g.LearningPaths = append(g.LearningPaths[:i], g.LearningPaths[i+1:]...)
			g.Touch()
			return
		}
	}
}

// AddMilestone adds a milestone to the goal
func (g *Goal) AddMilestone(milestoneID EntityID) {
	// Check if milestone already exists
	for _, id := range g.Milestones {
		if id == milestoneID {
			return
		}
	}
	g.Milestones = append(g.Milestones, milestoneID)
	g.Touch()
}

// RemoveMilestone removes a milestone from the goal
func (g *Goal) RemoveMilestone(milestoneID EntityID) {
	for i, id := range g.Milestones {
		if id == milestoneID {
			g.Milestones = append(g.Milestones[:i], g.Milestones[i+1:]...)
			g.Touch()
			return
		}
	}
}

// AddTag adds a tag to the goal
func (g *Goal) AddTag(tag string) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return
	}

	// Check if tag already exists
	for _, t := range g.Tags {
		if t == tag {
			return
		}
	}
	g.Tags = append(g.Tags, tag)
	g.Touch()
}

// UpdateStatus updates the goal's status
func (g *Goal) UpdateStatus(status Status) error {
	if !status.IsValid() {
		return errors.New("invalid goal status")
	}
	g.Status = status
	g.Touch()
	return nil
}

// UpdatePriority updates the goal's priority
func (g *Goal) UpdatePriority(priority Priority) error {
	if !priority.IsValid() {
		return errors.New("invalid goal priority")
	}
	g.Priority = priority
	g.Touch()
	return nil
}

// SetTargetDate sets the target completion date
func (g *Goal) SetTargetDate(date time.Time) {
	g.TargetDate = &date
	g.Touch()
}

// ClearTargetDate removes the target date
func (g *Goal) ClearTargetDate() {
	g.TargetDate = nil
	g.Touch()
}
