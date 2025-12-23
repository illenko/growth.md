package ai

import (
	"context"

	"github.com/illenko/growth.md/internal/core"
)

// MockClient is a mock AI client for testing
type MockClient struct {
	GenerateLearningPathFunc func(ctx context.Context, req PathGenerationRequest) (*PathGenerationResponse, error)
	SuggestResourcesFunc     func(ctx context.Context, req ResourceSuggestionRequest) (*ResourceSuggestionResponse, error)
	AnalyzeProgressFunc      func(ctx context.Context, req ProgressAnalysisRequest) (*ProgressAnalysisResponse, error)
	ProviderName             string
}

// GenerateLearningPath calls the mock function or returns a default response
func (m *MockClient) GenerateLearningPath(ctx context.Context, req PathGenerationRequest) (*PathGenerationResponse, error) {
	if m.GenerateLearningPathFunc != nil {
		return m.GenerateLearningPathFunc(ctx, req)
	}

	// Return default mock response
	return &PathGenerationResponse{
		Path: &core.LearningPath{
			ID:     "path-001",
			Title:  "Mock Learning Path",
			Body:   "This is a mock learning path for testing",
			Type:   core.PathTypeAIGenerated,
			Status: core.StatusActive,
		},
		Phases:     []*core.Phase{},
		Resources:  []*core.Resource{},
		Milestones: []*core.Milestone{},
		Reasoning:  "Mock reasoning for testing",
	}, nil
}

// SuggestResources calls the mock function or returns a default response
func (m *MockClient) SuggestResources(ctx context.Context, req ResourceSuggestionRequest) (*ResourceSuggestionResponse, error) {
	if m.SuggestResourcesFunc != nil {
		return m.SuggestResourcesFunc(ctx, req)
	}

	// Return default mock response
	return &ResourceSuggestionResponse{
		Resources: []*core.Resource{
			{
				ID:      "resource-001",
				Title:   "Mock Resource",
				Type:    core.ResourceBook,
				SkillID: req.Skill.ID,
				Body:    "This is a mock resource for testing",
				Status:  core.ResourceNotStarted,
			},
		},
		Reasoning: "Mock resource reasoning",
	}, nil
}

// AnalyzeProgress calls the mock function or returns a default response
func (m *MockClient) AnalyzeProgress(ctx context.Context, req ProgressAnalysisRequest) (*ProgressAnalysisResponse, error) {
	if m.AnalyzeProgressFunc != nil {
		return m.AnalyzeProgressFunc(ctx, req)
	}

	// Return default mock response
	return &ProgressAnalysisResponse{
		Summary:         "Mock progress summary",
		Insights:        []string{"Mock insight 1", "Mock insight 2"},
		Recommendations: []string{"Mock recommendation 1", "Mock recommendation 2"},
		IsOnTrack:       true,
		SuggestedFocus:  []string{"Mock focus area"},
	}, nil
}

// Provider returns the provider name
func (m *MockClient) Provider() string {
	if m.ProviderName != "" {
		return m.ProviderName
	}
	return "mock"
}

// NewMockClient creates a new mock client
func NewMockClient() *MockClient {
	return &MockClient{
		ProviderName: "mock",
	}
}
