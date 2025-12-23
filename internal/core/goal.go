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

func (g *Goal) Validate() error {
	if g.ID == "" {
		return errors.New("goal ID is required")
	}

	if strings.TrimSpace(g.Title) == "" {
		return errors.New("goal title is required and cannot be empty")
	}

	if !g.Status.IsValid() {
		return errors.New("invalid goal status: must be one of: active, completed, archived")
	}

	if !g.Priority.IsValid() {
		return errors.New("invalid goal priority: must be one of: high, medium, low")
	}

	if g.Created.IsZero() {
		return errors.New("goal created timestamp is required")
	}

	if g.Updated.IsZero() {
		return errors.New("goal updated timestamp is required")
	}

	return nil
}

func (g *Goal) AddLearningPath(pathID EntityID) {
	for _, id := range g.LearningPaths {
		if id == pathID {
			return
		}
	}
	g.LearningPaths = append(g.LearningPaths, pathID)
	g.Touch()
}

func (g *Goal) RemoveLearningPath(pathID EntityID) {
	for i, id := range g.LearningPaths {
		if id == pathID {
			g.LearningPaths = append(g.LearningPaths[:i], g.LearningPaths[i+1:]...)
			g.Touch()
			return
		}
	}
}

func (g *Goal) AddMilestone(milestoneID EntityID) {
	for _, id := range g.Milestones {
		if id == milestoneID {
			return
		}
	}
	g.Milestones = append(g.Milestones, milestoneID)
	g.Touch()
}

func (g *Goal) RemoveMilestone(milestoneID EntityID) {
	for i, id := range g.Milestones {
		if id == milestoneID {
			g.Milestones = append(g.Milestones[:i], g.Milestones[i+1:]...)
			g.Touch()
			return
		}
	}
}

func (g *Goal) AddTag(tag string) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return
	}

	for _, t := range g.Tags {
		if t == tag {
			return
		}
	}
	g.Tags = append(g.Tags, tag)
	g.Touch()
}

func (g *Goal) UpdateStatus(status Status) error {
	if !status.IsValid() {
		return errors.New("invalid goal status: must be one of: active, completed, archived")
	}
	g.Status = status
	g.Touch()
	return nil
}

func (g *Goal) UpdatePriority(priority Priority) error {
	if !priority.IsValid() {
		return errors.New("invalid goal priority: must be one of: high, medium, low")
	}
	g.Priority = priority
	g.Touch()
	return nil
}

func (g *Goal) SetTargetDate(date time.Time) {
	g.TargetDate = &date
	g.Touch()
}

func (g *Goal) ClearTargetDate() {
	g.TargetDate = nil
	g.Touch()
}
