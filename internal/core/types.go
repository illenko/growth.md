package core

import "time"

// EntityID is a unique identifier for entities
type EntityID string

// Status represents the status of goals and paths
type Status string

const (
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusCompleted, StatusArchived:
		return true
	}
	return false
}

// Priority represents the priority level of goals
type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

// IsValid checks if the priority is valid
func (p Priority) IsValid() bool {
	switch p {
	case PriorityHigh, PriorityMedium, PriorityLow:
		return true
	}
	return false
}

// ProficiencyLevel represents skill proficiency levels
type ProficiencyLevel string

const (
	LevelBeginner     ProficiencyLevel = "beginner"
	LevelIntermediate ProficiencyLevel = "intermediate"
	LevelAdvanced     ProficiencyLevel = "advanced"
	LevelExpert       ProficiencyLevel = "expert"
)

// IsValid checks if the proficiency level is valid
func (l ProficiencyLevel) IsValid() bool {
	switch l {
	case LevelBeginner, LevelIntermediate, LevelAdvanced, LevelExpert:
		return true
	}
	return false
}

// SkillStatus represents the status of a skill
type SkillStatus string

const (
	SkillNotStarted SkillStatus = "not-started"
	SkillLearning   SkillStatus = "learning"
	SkillMastered   SkillStatus = "mastered"
)

// IsValid checks if the skill status is valid
func (s SkillStatus) IsValid() bool {
	switch s {
	case SkillNotStarted, SkillLearning, SkillMastered:
		return true
	}
	return false
}

// ResourceType represents the type of learning resource
type ResourceType string

const (
	ResourceBook          ResourceType = "book"
	ResourceCourse        ResourceType = "course"
	ResourceVideo         ResourceType = "video"
	ResourceArticle       ResourceType = "article"
	ResourceProject       ResourceType = "project"
	ResourceDocumentation ResourceType = "documentation"
)

// IsValid checks if the resource type is valid
func (r ResourceType) IsValid() bool {
	switch r {
	case ResourceBook, ResourceCourse, ResourceVideo, ResourceArticle, ResourceProject, ResourceDocumentation:
		return true
	}
	return false
}

// ResourceStatus represents the status of a resource
type ResourceStatus string

const (
	ResourceNotStarted ResourceStatus = "not-started"
	ResourceInProgress ResourceStatus = "in-progress"
	ResourceCompleted  ResourceStatus = "completed"
)

// IsValid checks if the resource status is valid
func (r ResourceStatus) IsValid() bool {
	switch r {
	case ResourceNotStarted, ResourceInProgress, ResourceCompleted:
		return true
	}
	return false
}

// PathType represents whether a learning path is AI-generated or manual
type PathType string

const (
	PathTypeAIGenerated PathType = "ai-generated"
	PathTypeManual      PathType = "manual"
)

// IsValid checks if the path type is valid
func (p PathType) IsValid() bool {
	switch p {
	case PathTypeAIGenerated, PathTypeManual:
		return true
	}
	return false
}

// MilestoneType represents what type of entity the milestone is associated with
type MilestoneType string

const (
	MilestoneGoalLevel  MilestoneType = "goal-level"
	MilestonePathLevel  MilestoneType = "path-level"
	MilestoneSkillLevel MilestoneType = "skill-level"
)

// IsValid checks if the milestone type is valid
func (m MilestoneType) IsValid() bool {
	switch m {
	case MilestoneGoalLevel, MilestonePathLevel, MilestoneSkillLevel:
		return true
	}
	return false
}

// ReferenceType represents the type of entity being referenced
type ReferenceType string

const (
	ReferenceGoal  ReferenceType = "goal"
	ReferencePath  ReferenceType = "path"
	ReferenceSkill ReferenceType = "skill"
)

// IsValid checks if the reference type is valid
func (r ReferenceType) IsValid() bool {
	switch r {
	case ReferenceGoal, ReferencePath, ReferenceSkill:
		return true
	}
	return false
}

// Timestamps is a common struct for created/updated times
type Timestamps struct {
	Created time.Time `yaml:"created"`
	Updated time.Time `yaml:"updated"`
}

// NewTimestamps creates a new Timestamps with current time
func NewTimestamps() Timestamps {
	now := time.Now()
	return Timestamps{
		Created: now,
		Updated: now,
	}
}

// Touch updates the Updated timestamp to now
func (t *Timestamps) Touch() {
	t.Updated = time.Now()
}
