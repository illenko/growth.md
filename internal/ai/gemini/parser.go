package gemini

import (
	"encoding/json"
	"fmt"

	"github.com/illenko/growth.md/internal/ai"
	"github.com/illenko/growth.md/internal/core"
)

// PathGenerationOutput matches the JSON schema from the prompt
type PathGenerationOutput struct {
	Path      PathOutput    `json:"path"`
	Phases    []PhaseOutput `json:"phases"`
	Reasoning string        `json:"reasoning"`
}

type PathOutput struct {
	Title                  string `json:"title"`
	Description            string `json:"description"`
	EstimatedDurationWeeks int    `json:"estimated_duration_weeks"`
}

type PhaseOutput struct {
	Title             string                   `json:"title"`
	Description       string                   `json:"description"`
	DurationWeeks     int                      `json:"duration_weeks"`
	SkillRequirements []SkillRequirementOutput `json:"skill_requirements"`
	Milestones        []MilestoneOutput        `json:"milestones"`
	Resources         []ResourceOutput         `json:"resources"`
}

type SkillRequirementOutput struct {
	SkillTitle    string `json:"skill_title"`
	Category      string `json:"category"`
	RequiredLevel string `json:"required_level"`
}

type MilestoneOutput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type ResourceOutput struct {
	Title          string  `json:"title"`
	Type           string  `json:"type"`
	Author         string  `json:"author"`
	URL            string  `json:"url"`
	EstimatedHours float64 `json:"estimated_hours"`
	Description    string  `json:"description"`
	WhyRecommended string  `json:"why_recommended,omitempty"`
	Cost           string  `json:"cost,omitempty"`
}

// ResourceSuggestionOutput matches the resource suggestion JSON schema
type ResourceSuggestionOutput struct {
	Resources []ResourceOutput `json:"resources"`
	Reasoning string           `json:"reasoning"`
}

// ProgressAnalysisOutput matches the progress analysis JSON schema
type ProgressAnalysisOutput struct {
	Summary         string   `json:"summary"`
	Insights        []string `json:"insights"`
	Recommendations []string `json:"recommendations"`
	IsOnTrack       bool     `json:"is_on_track"`
	SuggestedFocus  []string `json:"suggested_focus"`
}

// ParsePathGeneration converts AI response to domain types
func ParsePathGeneration(responseText string, pathID, goalID core.EntityID) (*ai.PathGenerationResponse, error) {
	var output PathGenerationOutput

	if err := json.Unmarshal([]byte(responseText), &output); err != nil {
		return nil, &ai.ParseError{
			Provider: "gemini",
			Message:  "failed to parse path generation response",
			Err:      err,
		}
	}

	// Validate required fields
	if output.Path.Title == "" {
		return nil, &ai.ParseError{
			Provider: "gemini",
			Message:  "path title is missing from response",
		}
	}

	// Convert path
	path := &core.LearningPath{
		ID:          pathID,
		Title:       output.Path.Title,
		Body:        output.Path.Description,
		Type:        core.PathTypeAIGenerated,
		Status:      core.StatusActive,
		GeneratedBy: "gemini-2.0-flash-exp",
		Phases:      []core.EntityID{},
		Tags:        []string{},
		Timestamps:  core.NewTimestamps(),
	}

	// Convert phases
	phases := make([]*core.Phase, 0, len(output.Phases))
	resources := make([]*core.Resource, 0)
	milestones := make([]*core.Milestone, 0)

	for i, phaseOut := range output.Phases {
		phaseID := core.EntityID(fmt.Sprintf("phase-%03d", i+1))

		phase := &core.Phase{
			ID:                phaseID,
			PathID:            pathID,
			Title:             phaseOut.Title,
			Body:              phaseOut.Description,
			Order:             i + 1,
			RequiredSkills:    []core.SkillRequirement{},
			Milestones:        []core.EntityID{},
			EstimatedDuration: fmt.Sprintf("%d weeks", phaseOut.DurationWeeks),
			Timestamps:        core.NewTimestamps(),
		}

		// Add skill requirements
		for _, skillReq := range phaseOut.SkillRequirements {
			level := core.ProficiencyLevel(skillReq.RequiredLevel)
			if level.IsValid() {
				phase.RequiredSkills = append(phase.RequiredSkills, core.SkillRequirement{
					SkillID:     "", // Will be matched/created later
					TargetLevel: level,
				})
			}
		}

		// Convert phase milestones
		for j, milestoneOut := range phaseOut.Milestones {
			milestoneID := core.EntityID(fmt.Sprintf("milestone-%03d", len(milestones)+j+1))

			milestoneType := core.MilestoneType(milestoneOut.Type)
			if !milestoneType.IsValid() {
				milestoneType = core.MilestonePathLevel // Default
			}

			milestone := &core.Milestone{
				ID:            milestoneID,
				Title:         milestoneOut.Title,
				Body:          milestoneOut.Description,
				Type:          milestoneType,
				ReferenceType: core.ReferencePath,
				ReferenceID:   pathID,
				Status:        core.StatusActive,
				Timestamps:    core.NewTimestamps(),
			}

			milestones = append(milestones, milestone)
			phase.Milestones = append(phase.Milestones, milestoneID)
		}

		// Convert phase resources
		for k, resourceOut := range phaseOut.Resources {
			resourceID := core.EntityID(fmt.Sprintf("resource-%03d", len(resources)+k+1))

			resourceType := core.ResourceType(resourceOut.Type)
			if !resourceType.IsValid() {
				resourceType = core.ResourceCourse // Default
			}

			resource := &core.Resource{
				ID:             resourceID,
				Title:          resourceOut.Title,
				Type:           resourceType,
				SkillID:        "", // Will be linked later
				Body:           resourceOut.Description,
				Author:         resourceOut.Author,
				URL:            resourceOut.URL,
				EstimatedHours: resourceOut.EstimatedHours,
				Status:         core.ResourceNotStarted,
				Tags:           []string{},
				Timestamps:     core.NewTimestamps(),
			}

			resources = append(resources, resource)
		}

		path.Phases = append(path.Phases, phaseID)
		phases = append(phases, phase)
	}

	return &ai.PathGenerationResponse{
		Path:       path,
		Phases:     phases,
		Resources:  resources,
		Milestones: milestones,
		Reasoning:  output.Reasoning,
	}, nil
}

// ParseResourceSuggestion converts AI response to resource suggestions
func ParseResourceSuggestion(responseText string, skillID core.EntityID) (*ai.ResourceSuggestionResponse, error) {
	var output ResourceSuggestionOutput

	if err := json.Unmarshal([]byte(responseText), &output); err != nil {
		return nil, &ai.ParseError{
			Provider: "gemini",
			Message:  "failed to parse resource suggestion response",
			Err:      err,
		}
	}

	resources := make([]*core.Resource, 0, len(output.Resources))

	for i, resourceOut := range output.Resources {
		resourceID := core.EntityID(fmt.Sprintf("resource-%03d", i+1))

		resourceType := core.ResourceType(resourceOut.Type)
		if !resourceType.IsValid() {
			resourceType = core.ResourceCourse // Default
		}

		resource := &core.Resource{
			ID:             resourceID,
			Title:          resourceOut.Title,
			Type:           resourceType,
			SkillID:        skillID,
			Body:           resourceOut.Description,
			Author:         resourceOut.Author,
			URL:            resourceOut.URL,
			EstimatedHours: resourceOut.EstimatedHours,
			Status:         core.ResourceNotStarted,
			Tags:           []string{},
			Timestamps:     core.NewTimestamps(),
		}

		resources = append(resources, resource)
	}

	return &ai.ResourceSuggestionResponse{
		Resources: resources,
		Reasoning: output.Reasoning,
	}, nil
}

// ParseProgressAnalysis converts AI response to progress analysis
func ParseProgressAnalysis(responseText string) (*ai.ProgressAnalysisResponse, error) {
	var output ProgressAnalysisOutput

	if err := json.Unmarshal([]byte(responseText), &output); err != nil {
		return nil, &ai.ParseError{
			Provider: "gemini",
			Message:  "failed to parse progress analysis response",
			Err:      err,
		}
	}

	return &ai.ProgressAnalysisResponse{
		Summary:         output.Summary,
		Insights:        output.Insights,
		Recommendations: output.Recommendations,
		IsOnTrack:       output.IsOnTrack,
		SuggestedFocus:  output.SuggestedFocus,
	}, nil
}
