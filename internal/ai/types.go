package ai

import (
	"time"

	"github.com/illenko/growth.md/internal/core"
)

type PathGenerationRequest struct {
	Goal           *core.Goal
	CurrentSkills  []*core.Skill
	Background     string
	LearningStyle  string // e.g., "top-down", "bottom-up", "project-based"
	TimeCommitment string // e.g., "10 hours/week"
	TargetDate     *time.Time
}

type PathGenerationResponse struct {
	Path       *core.LearningPath
	Phases     []*core.Phase
	Resources  []*core.Resource
	Milestones []*core.Milestone
	Reasoning  string
}

type ResourceSuggestionRequest struct {
	Skill         *core.Skill
	CurrentLevel  core.ProficiencyLevel
	TargetLevel   core.ProficiencyLevel
	LearningStyle string
	Budget        string // e.g., "free", "paid", "any"
}

type ResourceSuggestionResponse struct {
	Resources []*core.Resource
	Reasoning string
}

type ProgressAnalysisRequest struct {
	Goal          *core.Goal
	Path          *core.LearningPath
	ProgressLogs  []*core.ProgressLog
	CurrentSkills []*core.Skill
}

type ProgressAnalysisResponse struct {
	Summary         string
	Insights        []string
	Recommendations []string
	IsOnTrack       bool
	SuggestedFocus  []string
}
