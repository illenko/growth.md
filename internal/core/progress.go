package core

import (
	"errors"
	"time"
)

// ProgressLog represents a time-based journal entry
type ProgressLog struct {
	ID                 EntityID   `yaml:"id"`
	Date               time.Time  `yaml:"date"`
	HoursInvested      float64    `yaml:"hoursInvested,omitempty"`
	SkillsWorked       []EntityID `yaml:"skillsWorked,omitempty"`
	ResourcesUsed      []EntityID `yaml:"resourcesUsed,omitempty"`
	MilestonesAchieved []EntityID `yaml:"milestonesAchieved,omitempty"`
	Mood               string     `yaml:"mood,omitempty"` // e.g., "motivated", "frustrated", "focused"
	Timestamps

	// Body contains the markdown content (summary, accomplishments, challenges,
	// time breakdown, what I learned, reflections, energy level)
	Body string `yaml:"-"`
}

// NewProgressLog creates a new Progress Log for a given date
func NewProgressLog(id EntityID, date time.Time) (*ProgressLog, error) {
	// Normalize to midnight
	dateNormalized := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	log := &ProgressLog{
		ID:                 id,
		Date:               dateNormalized,
		SkillsWorked:       []EntityID{},
		ResourcesUsed:      []EntityID{},
		MilestonesAchieved: []EntityID{},
		Timestamps:         NewTimestamps(),
	}

	if err := log.Validate(); err != nil {
		return nil, err
	}

	return log, nil
}

func (p *ProgressLog) Validate() error {
	if p.ID == "" {
		return errors.New("progress log ID is required")
	}

	if p.Date.IsZero() {
		return errors.New("progress log date is required (use --date flag in YYYY-MM-DD format)")
	}

	if p.HoursInvested < 0 {
		return errors.New("progress log hours invested cannot be negative (must be >= 0)")
	}

	if p.Created.IsZero() {
		return errors.New("progress log created timestamp is required")
	}

	if p.Updated.IsZero() {
		return errors.New("progress log updated timestamp is required")
	}

	return nil
}

// AddSkillWorked adds a skill to the skills worked list
func (p *ProgressLog) AddSkillWorked(skillID EntityID) {
	for _, id := range p.SkillsWorked {
		if id == skillID {
			return
		}
	}
	p.SkillsWorked = append(p.SkillsWorked, skillID)
	p.Touch()
}

// AddResourceUsed adds a resource to the resources used list
func (p *ProgressLog) AddResourceUsed(resourceID EntityID) {
	for _, id := range p.ResourcesUsed {
		if id == resourceID {
			return
		}
	}
	p.ResourcesUsed = append(p.ResourcesUsed, resourceID)
	p.Touch()
}

// AddMilestoneAchieved adds a milestone to the achievements list
func (p *ProgressLog) AddMilestoneAchieved(milestoneID EntityID) {
	for _, id := range p.MilestonesAchieved {
		if id == milestoneID {
			return
		}
	}
	p.MilestonesAchieved = append(p.MilestonesAchieved, milestoneID)
	p.Touch()
}

// SetHoursInvested sets the total hours invested
func (p *ProgressLog) SetHoursInvested(hours float64) error {
	if hours < 0 {
		return errors.New("hours invested cannot be negative (must be >= 0)")
	}
	p.HoursInvested = hours
	p.Touch()
	return nil
}

// SetMood sets the mood for this period
func (p *ProgressLog) SetMood(mood string) {
	p.Mood = mood
	p.Touch()
}
