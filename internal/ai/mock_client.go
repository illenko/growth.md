package ai

import (
	"context"

	"github.com/illenko/growth.md/internal/core"
)

type MockClient struct {
	GenerateLearningPathFunc func(ctx context.Context, req PathGenerationRequest) (*PathGenerationResponse, error)
	SuggestResourcesFunc     func(ctx context.Context, req ResourceSuggestionRequest) (*ResourceSuggestionResponse, error)
	AnalyzeProgressFunc      func(ctx context.Context, req ProgressAnalysisRequest) (*ProgressAnalysisResponse, error)
	ProviderName             string
}

func (m *MockClient) GenerateLearningPath(ctx context.Context, req PathGenerationRequest) (*PathGenerationResponse, error) {
	if m.GenerateLearningPathFunc != nil {
		return m.GenerateLearningPathFunc(ctx, req)
	}

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

func (m *MockClient) SuggestResources(ctx context.Context, req ResourceSuggestionRequest) (*ResourceSuggestionResponse, error) {
	if m.SuggestResourcesFunc != nil {
		return m.SuggestResourcesFunc(ctx, req)
	}

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

func (m *MockClient) AnalyzeProgress(ctx context.Context, req ProgressAnalysisRequest) (*ProgressAnalysisResponse, error) {
	if m.AnalyzeProgressFunc != nil {
		return m.AnalyzeProgressFunc(ctx, req)
	}

	return &ProgressAnalysisResponse{
		Summary:         "Mock progress summary",
		Insights:        []string{"Mock insight 1", "Mock insight 2"},
		Recommendations: []string{"Mock recommendation 1", "Mock recommendation 2"},
		IsOnTrack:       true,
		SuggestedFocus:  []string{"Mock focus area"},
	}, nil
}

func (m *MockClient) Provider() string {
	if m.ProviderName != "" {
		return m.ProviderName
	}
	return "mock"
}

func NewMockClient() *MockClient {
	return &MockClient{
		ProviderName: "mock",
	}
}
